package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const catalogUpdateMaxDownloadBytes = int64(1_000_000_000)

type catalogUpdateArtifact struct {
	URL       string
	Bytes     int64
	ETag      string
	SHA256    string
	CheckedAt time.Time
	Status    string
}

const (
	catalogUpdateSizeUnknown = "size-unknown"
	catalogUpdateOversized   = "over-download-limit"
	catalogUpdateParserLimit = "over-parser-limit"
	catalogUpdateSizeChanged = "size-changed"
	catalogUpdateDownloaded  = "downloaded"
)

// acquireCatalogUpdateArtifact uses HEAD to establish a bound before issuing
// a body request. dst is private staging. The caller discards it on any
// outcome other than catalogUpdateDownloaded.
func (c *Client) acquireCatalogUpdateArtifact(ctx context.Context, rawURL string, parserLimit int64, dst io.Writer) (catalogUpdateArtifact, error) {
	c = c.effective()
	if parserLimit <= 0 {
		parserLimit = DefaultMaxBytes
	}
	if c.MaxBytes > 0 && c.MaxBytes < parserLimit {
		parserLimit = c.MaxBytes
	}
	timeout := c.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	parsed, err := c.versionURLSyntax(rawURL)
	if err != nil {
		return catalogUpdateArtifact{}, err
	}
	if parsed.Path == "" {
		parsed.Path = "/"
	}
	scope := []APIVersionOfficialScope{scopeForVersionURL(parsed.String())}
	if err := c.versionScopeAllows(parsed.String(), scope); err != nil {
		return catalogUpdateArtifact{}, err
	}
	if _, err := c.validateHTTPURL(ctx, parsed.String()); err != nil {
		return catalogUpdateArtifact{}, err
	}
	copyClient := *c
	base := *c.client()
	base.Jar = nil
	base.CheckRedirect = nil
	copyClient.HTTPClient = &base
	client, err := copyClient.redirectSafeClient()
	if err != nil {
		return catalogUpdateArtifact{}, err
	}
	defer client.CloseIdleConnections()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many artifact redirects")
		}
		if via[len(via)-1].URL.Scheme == "https" && req.URL.Scheme != "https" {
			return fmt.Errorf("artifact redirect downgrades HTTPS")
		}
		if err := c.versionScopeAllows(req.URL.String(), scope); err != nil {
			return err
		}
		_, err := c.validateHTTPURL(req.Context(), req.URL.String())
		return err
	}

	head, err := requestCatalogUpdateArtifact(ctx, client, http.MethodHead, parsed.String(), "")
	if err != nil {
		return catalogUpdateArtifact{}, err
	}
	err = head.Body.Close()
	probe := catalogUpdateArtifact{URL: head.Request.URL.String(), CheckedAt: time.Now().UTC()}
	if err != nil {
		return probe, err
	}
	switch head.StatusCode {
	case http.StatusMethodNotAllowed, http.StatusNotImplemented:
		probe.Status = catalogUpdateSizeUnknown
		return probe, nil
	case http.StatusOK:
	default:
		return probe, HTTPStatusError{Code: head.StatusCode, Status: head.Status}
	}
	probe.Bytes = catalogUpdateContentLength(head)
	switch {
	case probe.Bytes <= 0:
		probe.Status = catalogUpdateSizeUnknown
	case probe.Bytes > catalogUpdateMaxDownloadBytes:
		probe.Status = catalogUpdateOversized
	case probe.Bytes > parserLimit:
		probe.Status = catalogUpdateParserLimit
	}
	if probe.Status != "" {
		return probe, nil
	}
	probe.ETag = boundedVersionHeader(head.Header.Get("ETag"))

	// The exact final HEAD target is already checked. A GET redirect may
	// represent a different artifact. It requires its own preflight.
	client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}
	get, err := requestCatalogUpdateArtifact(ctx, client, http.MethodGet, probe.URL, probe.ETag)
	if err != nil {
		return probe, err
	}
	defer get.Body.Close()
	if get.StatusCode != http.StatusOK {
		return probe, HTTPStatusError{Code: get.StatusCode, Status: get.Status}
	}
	if probe.Bytes != catalogUpdateContentLength(get) ||
		(probe.ETag != "" && probe.ETag != get.Header.Get("ETag")) {
		probe.Status = catalogUpdateSizeChanged
		return probe, nil
	}
	digest := sha256.New()
	if _, err := io.CopyN(io.MultiWriter(dst, digest), get.Body, probe.Bytes); err != nil {
		return probe, fmt.Errorf("stream artifact: %w", err)
	}
	var extra [1]byte
	if n, err := io.ReadFull(get.Body, extra[:]); n != 0 || err != io.EOF {
		return probe, fmt.Errorf("artifact body does not match declared length")
	}
	probe.SHA256 = hex.EncodeToString(digest.Sum(nil))
	probe.Status = catalogUpdateDownloaded
	return probe, nil
}

func requestCatalogUpdateArtifact(ctx context.Context, client *http.Client, method, rawURL, etag string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("User-Agent", "apitools/catalog-update")
	if len(etag) >= 2 && strings.HasPrefix(etag, "\"") && strings.HasSuffix(etag, "\"") {
		req.Header.Set("If-Match", etag)
	}
	return client.Do(req)
}

func catalogUpdateContentLength(resp *http.Response) int64 {
	if resp.Uncompressed || len(resp.TransferEncoding) != 0 {
		return -1
	}
	if encoding := resp.Header.Get("Content-Encoding"); encoding != "" && !strings.EqualFold(encoding, "identity") {
		return -1
	}
	n, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil || n != resp.ContentLength {
		return -1
	}
	return n
}
