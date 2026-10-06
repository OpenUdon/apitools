package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/internal/artifactio"
)

type apiVersionState struct {
	SchemaVersion string                    `json:"schema_version"`
	Identity      string                    `json:"identity"`
	Report        APIVersionDiscoveryReport `json:"report"`
}

type apiVersionListCache struct {
	SchemaVersion string          `json:"schema_version"`
	URL           string          `json:"url"`
	ETag          string          `json:"etag,omitempty"`
	LastModified  string          `json:"last_modified,omitempty"`
	CheckedAt     time.Time       `json:"checked_at"`
	Data          json.RawMessage `json:"data"`
}

func versionRequestIdentity(req APIVersionDiscoveryRequest, opts APIVersionDiscoveryOptions) string {
	data, _ := json.Marshal(struct {
		Request         APIVersionDiscoveryRequest
		Inventory       *OperationInventory
		InventoryDigest string
	}{req, opts.BaselineInventory, opts.BaselineInventorySHA256})
	d := sha256.Sum256(data)
	return hex.EncodeToString(d[:])
}

func readVersionMetadata(file string, max int64) ([]byte, error) {
	abs, err := filepath.Abs(file)
	if err != nil {
		return nil, err
	}
	data, err := artifactio.ReadFile(filepath.Dir(abs), filepath.Base(abs), artifactio.ReadOptions{MaxBytes: max})
	return data.Data, err
}

func writeVersionMetadata(ctx context.Context, file string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	// A caller opt-in does not authorize placing freshness/cache data inside a
	// source checkout. Walk existing parents without following target symlinks.
	for p := filepath.Dir(abs); p != filepath.Dir(p); p = filepath.Dir(p) {
		if _, err := os.Lstat(filepath.Join(p, ".git")); err == nil {
			return fmt.Errorf("version metadata must be outside a repository")
		}
	}
	_, err = artifactio.WriteFile(filepath.Dir(abs), filepath.Base(abs), data, artifactio.WriteOptions{Force: true, Mode: 0600})
	return err
}

func loadVersionState(file, identity string) (APIVersionDiscoveryReport, bool, error) {
	data, err := readVersionMetadata(file, apiVersionMaxReportBytes+4096)
	if err != nil {
		return APIVersionDiscoveryReport{}, false, err
	}
	var state apiVersionState
	if err = decodeVersionJSON(data, &state, apiVersionMaxReportBytes+4096); err != nil {
		return state.Report, false, err
	}
	if state.SchemaVersion != APIVersionStateSchemaVersion {
		return state.Report, false, fmt.Errorf("unsupported API version state schema")
	}
	return state.Report, state.Identity == identity && freshVersionReport(state.Report), nil
}

func freshVersionReport(report APIVersionDiscoveryReport) bool {
	age := time.Since(report.CheckedAt)
	return report.SchemaVersion == APIVersionDiscoverySchemaVersion && age >= 0 && age < 24*time.Hour && !report.Truncated && (report.Status == "none_found_in_scope" || report.Status == "newer_found" || report.Status == "conflicting")
}

func saveVersionSource(ctx context.Context, dir string, record APIVersionRecord, data []byte) (*ImportedSpec, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	for p := abs; p != filepath.Dir(p); p = filepath.Dir(p) {
		if _, e := os.Lstat(filepath.Join(p, ".git")); e == nil {
			return nil, fmt.Errorf("version source save must be outside a repository")
		}
	}
	u, _ := url.Parse(record.FinalURL)
	metadata, ok := downloadedSpecMetadata(ctx, data, record.FinalURL)
	if !ok {
		return nil, fmt.Errorf("version source is not Import-valid OpenAPI/Swagger")
	}
	name := fileNameForImport("", u, data)
	stem := strings.TrimSuffix(name, filepath.Ext(name))
	ext := filepath.Ext(name)
	for n := 1; n <= 100; n++ {
		candidate := name
		if n > 1 {
			candidate = fmt.Sprintf("%s-%d%s", stem, n, ext)
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		result, e := artifactio.WriteFile(abs, candidate, data, artifactio.WriteOptions{Mode: 0644})
		if e != nil {
			if strings.Contains(e.Error(), "collision") {
				continue
			}
			return nil, e
		}
		return &ImportedSpec{Name: candidate, Path: result.Path, URL: record.FinalURL, SHA256: result.SHA256, Bytes: result.Bytes, Metadata: metadata, Title: metadata.Title, Description: metadata.Description}, nil
	}
	return nil, fmt.Errorf("version save collision limit reached")
}
