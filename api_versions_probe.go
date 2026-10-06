package apitools

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type apiVersionBudget struct {
	mu       sync.Mutex
	requests int
	bodies   int
	sem      chan struct{}
}

func (b *apiVersionBudget) charge() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.requests >= 12 {
		return fmt.Errorf("API version request budget exhausted")
	}
	b.requests++
	return nil
}
func (b *apiVersionBudget) reserve() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.bodies >= 3 {
		return false
	}
	b.bodies++
	return true
}
func (b *apiVersionBudget) release() { b.mu.Lock(); b.bodies--; b.mu.Unlock() }
func (b *apiVersionBudget) counts() (int, int) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.requests, b.bodies
}

type apiVersionProbe struct {
	content     []byte
	finalURL    string
	status      int
	bytes       int64
	declared    *int64
	digest      string
	etag        string
	modified    string
	retryAfter  string
	contentType string
	nonSpec     bool
	err         error
}

func (c *Client) probeAPIVersion(ctx context.Context, budget *apiVersionBudget, raw string, scopes []APIVersionOfficialScope, source bool, etag, modified string) apiVersionProbe {
	var result apiVersionProbe
	select {
	case budget.sem <- struct{}{}:
		defer func() { <-budget.sem }()
	case <-ctx.Done():
		result.err = ctx.Err()
		return result
	}
	if source && !budget.reserve() {
		result.err = fmt.Errorf("API version document budget exhausted")
		return result
	}
	release := false
	defer func() {
		if source && release {
			budget.release()
		}
	}()
	check := func(raw string) error {
		if _, err := c.versionURLSyntax(raw); err != nil {
			return err
		}
		if source {
			if err := c.versionScopeAllows(raw, scopes); err != nil {
				return err
			}
		}
		_, err := c.validateHTTPURL(ctx, raw)
		return err
	}
	if err := check(raw); err != nil {
		result.err = err
		return result
	}
	if err := budget.charge(); err != nil {
		result.err = err
		return result
	}
	copyClient := *c.effective()
	copyClient.Cache = nil
	base := *c.client()
	base.Jar = nil
	base.CheckRedirect = nil
	copyClient.HTTPClient = &base
	client, err := copyClient.redirectSafeClient()
	if err != nil {
		result.err = err
		return result
	}
	initial, _ := url.Parse(raw)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		if len(via) > 0 && via[0].URL.Scheme == "https" && req.URL.Scheme != "https" {
			return fmt.Errorf("HTTPS downgrade is not permitted")
		}
		if !source && versionOrigin(req.URL) != versionOrigin(initial) {
			return fmt.Errorf("metadata redirect changes origin")
		}
		if err := check(req.URL.String()); err != nil {
			return err
		}
		return budget.charge()
	}
	client.Timeout = 3 * time.Second
	reqCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, raw, nil)
	if err != nil {
		result.err = err
		return result
	}
	req.Header.Set("User-Agent", "apitools/api-version-discovery-v1")
	req.Header.Set("Accept-Encoding", "gzip")
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	} else if modified != "" {
		req.Header.Set("If-Modified-Since", modified)
	}
	resp, err := client.Do(req)
	if err != nil {
		result.err = err
		return result
	}
	defer resp.Body.Close()
	result.status = resp.StatusCode
	result.finalURL = raw
	if resp.Request != nil && resp.Request.URL != nil {
		result.finalURL = resp.Request.URL.String()
	}
	result.contentType = resp.Header.Get("Content-Type")
	result.etag = boundedVersionHeader(resp.Header.Get("ETag"))
	result.modified = boundedVersionHeader(resp.Header.Get("Last-Modified"))
	result.retryAfter = boundedVersionHeader(resp.Header.Get("Retry-After"))
	if resp.ContentLength >= 0 {
		n := resp.ContentLength
		result.declared = &n
	}
	if result.status == 304 || result.status == 404 {
		release = true
		return result
	}
	if result.status < 200 || result.status >= 300 {
		result.err = HTTPStatusError{Code: result.status, Status: resp.Status}
		return result
	}
	max := int64(1 << 20)
	if source {
		max = apiVersionMaxSourceBytes
	} else if raw == c.APIsGuruListURL || strings.HasPrefix(raw, "https://api.apis.guru/v2/") {
		max = DefaultMaxBytes
	}
	if c.MaxBytes > 0 && c.MaxBytes < max {
		max = c.MaxBytes
	}
	if resp.ContentLength > max {
		result.err = fmt.Errorf("source declared length exceeds %d bytes", max)
		return result
	}
	wire := &io.LimitedReader{R: resp.Body, N: max + 1}
	var decoded io.Reader = wire
	switch strings.ToLower(strings.TrimSpace(resp.Header.Get("Content-Encoding"))) {
	case "", "identity":
	case "gzip":
		g, e := gzip.NewReader(wire)
		if e != nil {
			result.err = e
			return result
		}
		defer g.Close()
		decoded = g
	default:
		result.err = fmt.Errorf("unsupported content encoding")
		return result
	}
	prefix, err := io.ReadAll(io.LimitReader(decoded, 4096))
	result.bytes = int64(len(prefix))
	if err != nil {
		result.err = err
		return result
	}
	trim := bytes.ToLower(bytes.TrimSpace(prefix))
	if source && (bytes.HasPrefix(trim, []byte("<!doctype html")) || bytes.HasPrefix(trim, []byte("<html"))) {
		result.nonSpec = true
		release = true
		return result
	}
	data, err := io.ReadAll(io.LimitReader(io.MultiReader(bytes.NewReader(prefix), decoded), max+1))
	result.bytes = int64(len(data))
	if err != nil {
		result.err = err
		return result
	}
	if int64(len(data)) > max {
		result.err = fmt.Errorf("decoded source exceeds %d bytes", max)
		return result
	}
	if _, err = io.Copy(io.Discard, wire); err != nil {
		result.err = err
		return result
	}
	if wire.N == 0 {
		result.err = fmt.Errorf("wire source exceeds %d bytes", max)
		return result
	}
	if err = ctx.Err(); err != nil {
		result.err = err
		return result
	}
	result.content = data
	digest := sha256.Sum256(data)
	result.digest = hex.EncodeToString(digest[:])
	return result
}

func boundedVersionHeader(value string) string {
	if len(value) > 256 || strings.ContainsAny(value, "\r\n") {
		return ""
	}
	return value
}
