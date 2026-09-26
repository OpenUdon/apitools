package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (c *Client) Import(ctx context.Context, opts ImportOptions) (ImportedSpec, error) {
	report, err := c.ImportWithReport(ctx, opts)
	if err != nil {
		return ImportedSpec{}, err
	}
	return report.Imported, nil
}

// ImportWithReport imports a document and reports any cache timestamp and age
// without changing the legacy Import result shape.
func (c *Client) ImportWithReport(ctx context.Context, opts ImportOptions) (ImportReport, error) {
	c = c.effective()
	dir := strings.TrimSpace(opts.Dir)
	if dir == "" {
		return ImportReport{}, fmt.Errorf("target directory is required")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return ImportReport{}, err
	}
	mode, err := normalizeCacheMode(opts.CacheMode)
	if err != nil {
		return ImportReport{}, err
	}
	content, finalURL, metadata, cacheStoredAt, err := c.downloadSpecWithCache(ctx, opts.URL, mode, normalizeCacheMaxAge(opts.CacheMaxAge))
	if err != nil {
		return ImportReport{}, err
	}
	name, path, err := writeUniqueFile(dir, fileNameForImport(opts.Name, finalURL, content), content)
	if err != nil {
		return ImportReport{}, err
	}
	digest := sha256.Sum256(content)
	report := ImportReport{Imported: ImportedSpec{
		Name:        name,
		Path:        path,
		Title:       metadata.Title,
		Description: metadata.Description,
		URL:         finalURL.String(),
		SHA256:      hex.EncodeToString(digest[:]),
		Bytes:       int64(len(content)),
		Metadata:    metadata,
	}}
	if !cacheStoredAt.IsZero() {
		report.CacheStoredAt = &cacheStoredAt
		age := time.Since(cacheStoredAt)
		if age < 0 {
			age = 0
		}
		report.CacheAge = age.Round(time.Second).String()
	}
	return report, nil
}

func fileNameForImport(suggestedName string, parsed *url.URL, content []byte) string {
	suggestedName = strings.TrimSpace(suggestedName)
	suggestedExt := strings.ToLower(filepath.Ext(suggestedName))
	stem := sanitizeName(strings.TrimSuffix(suggestedName, suggestedExt))
	if stem == "" && parsed != nil {
		stem = sanitizeName(strings.TrimSuffix(filepath.Base(parsed.Path), filepath.Ext(parsed.Path)))
	}
	if stem == "" {
		digest := sha256.Sum256(content)
		stem = "openapi-" + hex.EncodeToString(digest[:])[:12]
	}
	ext := ".yaml"
	switch suggestedExt {
	case ".json", ".yaml", ".yml":
		ext = suggestedExt
	}
	if parsed != nil {
		switch strings.ToLower(filepath.Ext(parsed.Path)) {
		case ".json", ".yaml", ".yml":
			if suggestedExt == "" {
				ext = strings.ToLower(filepath.Ext(parsed.Path))
			}
		}
	}
	if suggestedExt == "" && len(bytes.TrimSpace(content)) > 0 && bytes.TrimSpace(content)[0] == '{' {
		ext = ".json"
	}
	return stem + ext
}

func writeUniqueFile(dir, name string, content []byte) (string, string, error) {
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	for i := 1; ; i++ {
		candidate := name
		if i > 1 {
			candidate = fmt.Sprintf("%s-%d%s", stem, i, ext)
		}
		path := filepath.Join(dir, candidate)
		file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if os.IsExist(err) {
			reused, compareErr := existingFileMatches(path, content)
			if compareErr != nil {
				return "", "", compareErr
			}
			if reused {
				return candidate, path, nil
			}
			continue
		}
		if err != nil {
			return "", "", err
		}
		_, writeErr := file.Write(content)
		closeErr := file.Close()
		if writeErr != nil || closeErr != nil {
			_ = os.Remove(path)
			if writeErr != nil {
				return "", "", writeErr
			}
			return "", "", closeErr
		}
		return candidate, path, nil
	}
}

func existingFileMatches(path string, content []byte) (bool, error) {
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() != int64(len(content)) {
		return false, nil
	}
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	openedInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	if !openedInfo.Mode().IsRegular() || !os.SameFile(info, openedInfo) {
		return false, nil
	}
	existing, err := io.ReadAll(io.LimitReader(file, int64(len(content))+1))
	if err != nil {
		return false, err
	}
	return bytes.Equal(existing, content), nil
}

func sanitizeName(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastDash := false
	for _, r := range value {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
