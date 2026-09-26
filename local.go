package apitools

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// ErrLocalScanTruncated wraps the error LocalFiles and DiscoverOpenAPI return
// when a local scan reaches its visit or candidate bound before completing.
// Results returned alongside this error are valid and partial, not invalid:
// callers may keep and use them, and multi-source discovery treats this as
// recoverable rather than aborting. Use errors.Is to distinguish it from a
// scan failure such as a missing or unreadable root.
var ErrLocalScanTruncated = errors.New("local scan reached a bound before completing")

type LocalOptions struct {
	Dir      string
	BaseDir  string
	Query    string
	MaxBytes int64 // per-file limit; zero uses DefaultMaxBytes
}

type LocalResult struct {
	Path         string       `json:"path"`
	RelativePath string       `json:"relative_path"`
	Title        string       `json:"title,omitempty"`
	Description  string       `json:"description,omitempty"`
	Score        int          `json:"score"`
	Metadata     SpecMetadata `json:"metadata"`
}

// LocalFiles discovers OpenAPI or Swagger documents already present on disk.
// Invalid or unreadable entries are ignored so callers can scan mixed project
// directories. Traversal, candidate, and per-file byte limits are bounded;
// reaching a traversal or candidate limit returns an error with partial results.
// Use DiscoverLocalSources when structured per-entry rejection details are needed.
func LocalFiles(ctx context.Context, opts LocalOptions) ([]LocalResult, error) {
	results, _, err := localFilesWithDigests(ctx, opts)
	return results, err
}

func localFilesWithDigests(ctx context.Context, opts LocalOptions) ([]LocalResult, []string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	dir := strings.TrimSpace(opts.Dir)
	if dir == "" {
		return nil, nil, fmt.Errorf("local OpenAPI directory is required")
	}
	dir, err := resolveLocalScanRoot(dir)
	if err != nil {
		return nil, nil, err
	}
	baseDir := strings.TrimSpace(opts.BaseDir)
	if baseDir == "" {
		baseDir = dir
	} else {
		baseDir, err = filepath.Abs(baseDir)
		if err != nil {
			return nil, nil, err
		}
		baseDir = filepath.Clean(baseDir)
		if resolved, resolveErr := filepath.EvalSymlinks(baseDir); resolveErr == nil {
			baseDir = resolved
		}
	}
	state, err := discoverLocalSources(ctx, LocalSourceDiscoveryOptions{
		Roots:    []string{dir},
		Query:    opts.Query,
		MaxBytes: opts.MaxBytes,
	}, true, baseDir)
	if err != nil {
		return nil, nil, err
	}
	digests := make([]string, len(state.localResults))
	for i, result := range state.localResults {
		digests[i] = state.localResultDigests[result.Path]
	}
	if state.report.Truncated && len(state.report.Diagnostics) > 0 {
		return state.localResults, digests, fmt.Errorf("local OpenAPI scan is incomplete: %s: %w", state.report.Diagnostics[0].Message, ErrLocalScanTruncated)
	}
	return state.localResults, digests, nil
}

func localSpecMetadata(ctx context.Context, content []byte) (SpecMetadata, bool) {
	if metadata, ok := specMetadata(ctx, content); ok {
		return metadata, true
	}
	return looseSpecMetadata(content)
}

func looseSpecMetadata(content []byte) (SpecMetadata, bool) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return SpecMetadata{}, false
	}
	var root struct {
		OpenAPI string `json:"openapi" yaml:"openapi"`
		Swagger string `json:"swagger" yaml:"swagger"`
		Info    struct {
			Title       string `json:"title" yaml:"title"`
			Description string `json:"description" yaml:"description"`
		} `json:"info" yaml:"info"`
	}
	if trimmed[0] == '{' {
		if err := json.Unmarshal(trimmed, &root); err != nil {
			return SpecMetadata{}, false
		}
	} else if err := yaml.Unmarshal(trimmed, &root); err != nil {
		return SpecMetadata{}, false
	}
	if strings.TrimSpace(root.OpenAPI) == "" && strings.TrimSpace(root.Swagger) == "" {
		return SpecMetadata{}, false
	}
	return SpecMetadata{
		Title:       strings.TrimSpace(root.Info.Title),
		Description: strings.TrimSpace(root.Info.Description),
		OpenAPI:     strings.TrimSpace(root.OpenAPI),
		Swagger:     strings.TrimSpace(root.Swagger),
	}, true
}

func hasOpenAPIFileExt(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json", ".yaml", ".yml":
		return true
	default:
		return false
	}
}
