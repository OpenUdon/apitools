package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/OpenUdon/apitools/catalog"
)

// ReviewedEmbeddingConfig is an opt-in OpenAI-compatible embeddings endpoint
// and its per-search request, input, and UTF-8 text byte budgets.
type ReviewedEmbeddingConfig struct {
	Endpoint    string `json:"endpoint"`
	Model       string `json:"model"`
	MaxRequests int    `json:"max_requests"`
	MaxInputs   int    `json:"max_inputs"`
	MaxBytes    int    `json:"max_bytes"`
}

// ReviewedEmbeddingUsage reports status and consumption for a single lookup.
type ReviewedEmbeddingUsage struct {
	Status   string `json:"status"`
	Requests int    `json:"requests"`
	Inputs   int    `json:"inputs"`
	Bytes    int    `json:"bytes"`
}

const (
	reviewedEmbeddingBatch             = 32
	reviewedEmbeddingTimeout           = 12 * time.Second
	reviewedEmbeddingMaxResponse       = 8 << 20
	reviewedEmbeddingMaxDimensions     = 4096
	reviewedEmbeddingMinimumSimilarity = 0.8
)

var errReviewedEmbeddingResponse = errors.New("embedding service response is unavailable or invalid")

type reviewedEmbeddingInput struct {
	digest string
	text   string
}

func (c *Client) rankReviewedEmbeddings(ctx context.Context, db *sql.DB, options ReviewedProviderSearchOptions, profiles []catalog.ProviderSearchProfile, ranked []rankedReviewedProvider, query string) ([]rankedReviewedProvider, *ReviewedEmbeddingUsage) {
	usage := &ReviewedEmbeddingUsage{}
	cfg := *options.Embeddings
	switch {
	case options.Offline:
		usage.Status = "offline"
		return ranked, usage
	case len(ranked) > 0 && ranked[0].result.Match == ReviewedProviderExactMatch:
		usage.Status = "exact-provider"
		return ranked, usage
	case strings.TrimSpace(cfg.Endpoint) == "" || strings.TrimSpace(cfg.Model) == "":
		usage.Status = "unconfigured"
		return ranked, usage
	}
	c = c.effective()
	endpoint, err := c.versionURLSyntax(cfg.Endpoint)
	if err != nil || endpoint.Scheme != "https" || endpoint.Fragment != "" || strings.Contains(cfg.Endpoint, "#") ||
		endpoint.RawQuery != "" || endpoint.ForceQuery || !utf8.ValidString(cfg.Model) ||
		len(cfg.Model) > 128 || cfg.Model != strings.TrimSpace(cfg.Model) || strings.IndexFunc(cfg.Model, unicode.IsControl) != -1 ||
		cfg.MaxRequests <= 0 || cfg.MaxRequests > 128 || cfg.MaxInputs <= 0 || cfg.MaxInputs > 1024 ||
		cfg.MaxBytes <= 0 || cfg.MaxBytes > 4<<20 {
		usage.Status = "invalid-configuration"
		return ranked, usage
	}
	if !utf8.ValidString(query) {
		usage.Status = "invalid-query-text"
		return ranked, usage
	}
	if len(query) > 1024 || len(query) > cfg.MaxBytes {
		usage.Status = "budget-exhausted"
		return ranked, usage
	}
	ctx, cancel := context.WithTimeout(ctx, reviewedEmbeddingTimeout)
	defer cancel()

	namespace := fmt.Sprintf("%x", sha256.Sum256([]byte(endpoint.String()+"\x00"+cfg.Model)))
	inputs := make([]reviewedEmbeddingInput, 0, len(profiles))
	digests := make(map[string]string, len(profiles))
	for _, profile := range profiles {
		publicText := reviewedEmbeddingText(profile)
		if !utf8.ValidString(publicText) {
			usage.Status = "invalid-profile-text"
			return ranked, usage
		}
		digest := fmt.Sprintf("%x", sha256.Sum256([]byte(publicText)))
		inputs = append(inputs, reviewedEmbeddingInput{digest: digest, text: publicText})
		digests[profile.ProviderID] = digest
	}
	vectors, err := readReviewedEmbeddingCache(ctx, db, namespace, inputs)
	if err != nil {
		usage.Status = "cache-unavailable"
		return ranked, usage
	}

	var missing []reviewedEmbeddingInput
	plannedBytes := len(query)
	for _, input := range inputs {
		if _, exists := vectors[input.digest]; !exists {
			missing = append(missing, input)
			plannedBytes += len(input.text)
		}
	}
	requests := (len(missing)+reviewedEmbeddingBatch-1)/reviewedEmbeddingBatch + 1
	if requests > cfg.MaxRequests || len(missing)+1 > cfg.MaxInputs || plannedBytes > cfg.MaxBytes {
		usage.Status = "budget-exhausted"
		return ranked, usage
	}
	client, err := c.reviewedEmbeddingHTTPClient(ctx, endpoint)
	if err != nil {
		usage.Status = "endpoint-unavailable"
		return ranked, usage
	}
	defer client.CloseIdleConnections()

	newVectors := make([][]float32, 0, len(missing))
	for first := 0; first < len(missing); first += reviewedEmbeddingBatch {
		batch := missing[first:min(first+reviewedEmbeddingBatch, len(missing))]
		text := make([]string, 0, len(batch))
		for _, input := range batch {
			text = append(text, input.text)
		}
		got, err := requestReviewedEmbeddings(ctx, client, endpoint, cfg.Model, text, usage)
		if err != nil {
			usage.Status = "endpoint-unavailable"
			return ranked, usage
		}
		newVectors = append(newVectors, got...)
		for i, input := range batch {
			vectors[input.digest] = got[i]
		}
	}
	queryResult, err := requestReviewedEmbeddings(ctx, client, endpoint, cfg.Model, []string{query}, usage)
	if err != nil {
		usage.Status = "endpoint-unavailable"
		return ranked, usage
	}
	queryVector := queryResult[0]
	for _, vector := range vectors {
		if len(vector) != len(queryVector) {
			usage.Status = "inconsistent-dimensions"
			return ranked, usage
		}
	}
	if len(missing) > 0 {
		if err := storeReviewedEmbeddingBatch(ctx, db, namespace, missing, newVectors); err != nil {
			usage.Status = "cache-unavailable"
			return ranked, usage
		}
	}
	usage.Status = "used"
	return mergeReviewedEmbeddingRanks(profiles, ranked, digests, vectors, queryVector), usage
}

// Keep the embedding payload projection independent of sources, locators,
// revisions, auth, persistent database contents, and the public search chain.
func reviewedEmbeddingText(profile catalog.ProviderSearchProfile) string {
	parts := []string{profile.DisplayName, profile.Category, profile.Description}
	parts = append(parts, profile.Aliases...)
	parts = append(parts, profile.Keywords...)
	for _, source := range profile.Coverage.Sources {
		for _, operation := range source.Operations {
			parts = append(parts, operation.Label)
		}
	}
	return strings.Join(parts, "\n")
}

func readReviewedEmbeddingCache(ctx context.Context, db *sql.DB, namespace string, inputs []reviewedEmbeddingInput) (map[string][]float32, error) {
	if _, err := db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS search_embedding (
    namespace TEXT NOT NULL, digest TEXT NOT NULL, vector BLOB NOT NULL,
    PRIMARY KEY (namespace, digest)
)`); err != nil {
		return nil, err
	}
	selectVector, err := db.PrepareContext(ctx, "SELECT vector FROM search_embedding WHERE namespace=? AND digest=? AND length(vector)<=?")
	if err != nil {
		return nil, err
	}
	defer selectVector.Close()
	vectors := make(map[string][]float32, len(inputs))
	for _, input := range inputs {
		var data []byte
		err := selectVector.QueryRowContext(ctx, namespace, input.digest, 128<<10).Scan(&data)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return nil, err
		}
		var vector []float32
		if json.Unmarshal(data, &vector) == nil && validReviewedEmbeddingVector(vector) {
			vectors[input.digest] = vector
		}
	}
	return vectors, nil
}

func storeReviewedEmbeddingBatch(ctx context.Context, db *sql.DB, namespace string, inputs []reviewedEmbeddingInput, vectors [][]float32) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, input := range inputs {
		raw, err := json.Marshal(vectors[i])
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT OR REPLACE INTO search_embedding VALUES (?, ?, ?)", namespace, input.digest, raw); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (c *Client) reviewedEmbeddingHTTPClient(ctx context.Context, endpoint *url.URL) (*http.Client, error) {
	if _, err := c.validateHTTPURL(ctx, endpoint.String()); err != nil {
		return nil, err
	}
	transport := &http.Transport{
		DialContext: c.safeDialContext, ForceAttemptHTTP2: true,
		TLSHandshakeTimeout: 8 * time.Second,
	}
	if c.AllowUnsafeHosts {
		transport.DialContext = (&net.Dialer{}).DialContext
		// An explicitly selected local fixture may use its own trust roots.
		// The new transport owns all dialing, routing and client identity.
		if ambient, ok := c.client().Transport.(*http.Transport); ok && ambient != nil && ambient.TLSClientConfig != nil &&
			ambient.TLSClientConfig.RootCAs != nil {
			transport.TLSClientConfig = &tls.Config{RootCAs: ambient.TLSClientConfig.RootCAs.Clone()}
		}
	}
	return &http.Client{
		Transport: transport, Timeout: reviewedEmbeddingTimeout,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}, nil
}

func requestReviewedEmbeddings(ctx context.Context, client *http.Client, endpoint *url.URL, model string, inputs []string, usage *ReviewedEmbeddingUsage) ([][]float32, error) {
	request := struct {
		Model string   `json:"model"`
		Input []string `json:"input"`
	}{Model: model, Input: inputs}
	raw, err := json.Marshal(request)
	if err != nil {
		return nil, errReviewedEmbeddingResponse
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), bytes.NewReader(raw))
	if err != nil {
		return nil, errReviewedEmbeddingResponse
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	usage.Requests++
	usage.Inputs += len(inputs)
	for _, input := range inputs {
		usage.Bytes += len(input)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errReviewedEmbeddingResponse
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || resp.ContentLength > reviewedEmbeddingMaxResponse {
		return nil, errReviewedEmbeddingResponse
	}
	data, err := io.ReadAll(io.LimitReader(resp.Body, reviewedEmbeddingMaxResponse+1))
	if err != nil || len(data) > reviewedEmbeddingMaxResponse {
		return nil, errReviewedEmbeddingResponse
	}
	var decoded struct {
		Data []struct {
			Index     *int      `json:"index"`
			Embedding []float32 `json:"embedding"`
		} `json:"data"`
	}
	if json.Unmarshal(data, &decoded) != nil || len(decoded.Data) != len(inputs) {
		return nil, errReviewedEmbeddingResponse
	}
	vectors := make([][]float32, len(inputs))
	for _, row := range decoded.Data {
		if row.Index == nil || *row.Index < 0 || *row.Index >= len(inputs) || vectors[*row.Index] != nil ||
			!validReviewedEmbeddingVector(row.Embedding) || len(row.Embedding) != len(decoded.Data[0].Embedding) {
			return nil, errReviewedEmbeddingResponse
		}
		vectors[*row.Index] = row.Embedding
	}
	return vectors, nil
}

func validReviewedEmbeddingVector(vector []float32) bool {
	if len(vector) == 0 || len(vector) > reviewedEmbeddingMaxDimensions {
		return false
	}
	var nonzero bool
	for _, item := range vector {
		if math.IsNaN(float64(item)) || math.IsInf(float64(item), 0) {
			return false
		}
		nonzero = nonzero || item != 0
	}
	return nonzero
}

func mergeReviewedEmbeddingRanks(profiles []catalog.ProviderSearchProfile, ranked []rankedReviewedProvider, digests map[string]string, vectors map[string][]float32, query []float32) []rankedReviewedProvider {
	positions := make(map[string]int, len(ranked))
	for i, entry := range ranked {
		positions[entry.result.Profile.ProviderID] = i
	}
	for _, profile := range profiles {
		similarity := reviewedCosine(vectors[digests[profile.ProviderID]], query)
		i, found := positions[profile.ProviderID]
		switch {
		case found && ranked[i].rank == 1:
			lexical := ranked[i].score
			ranked[i].score = lexical/(1+math.Abs(lexical)) - similarity/2
		case similarity >= reviewedEmbeddingMinimumSimilarity && profile.Coverage.Status == catalog.ProviderSearchSelectedSubset:
			semantic := rankedReviewedProvider{
				result: ReviewedProviderSearchResult{Profile: profile, Match: ReviewedProviderSemanticMatch},
				rank:   2, score: -similarity,
			}
			if found {
				ranked[i] = semantic
			} else {
				ranked = append(ranked, semantic)
			}
		}
	}
	sortRankedReviewedProviders(ranked)
	return ranked
}

func reviewedCosine(vector, query []float32) float64 {
	var dot, a, b float64
	for i, value := range vector {
		x, y := float64(value), float64(query[i])
		dot += x * y
		a += x * x
		b += y * y
	}
	return dot / math.Sqrt(a*b)
}
