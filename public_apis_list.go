package apitools

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/OpenUdon/apitools/internal/sourceguard"
)

const (
	publicAPIsMaxEntries  = 10000
	publicAPIsMaxRowBytes = 64 * 1024
)

// parsePublicAPIsCatalog accepts legacy JSON mirrors and the public-apis
// Markdown list. Limits fail closed: a partial list is never searchable.
func parsePublicAPIsCatalog(ctx context.Context, content []byte, sourceURL string) ([]publicAPIEntry, error) {
	fail := func(reason string) ([]publicAPIEntry, error) {
		return nil, fmt.Errorf("parse public-apis catalog %q: %s", sourceURL, reason)
	}
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("parse public-apis catalog %q: %w", sourceURL, err)
	}
	if len(content) > DefaultMaxBytes {
		return fail(fmt.Sprintf("response exceeds %d bytes", DefaultMaxBytes))
	}
	content = bytes.TrimSpace(content)
	if len(content) > 0 && content[0] == '{' {
		if err := sourceguard.CheckJSONWithLimits(ctx, "public-apis", content, sourceguard.DefaultLimits()); err != nil {
			return nil, fmt.Errorf("parse public-apis catalog %q: %w", sourceURL, err)
		}
		var raw struct {
			Entries []json.RawMessage `json:"entries"`
		}
		if err := json.Unmarshal(content, &raw); err != nil {
			return fail(err.Error())
		}
		if raw.Entries == nil {
			return fail("JSON entries must be an array")
		}
		if len(raw.Entries) > publicAPIsMaxEntries {
			return fail(fmt.Sprintf("entry count exceeds %d", publicAPIsMaxEntries))
		}
		entries := make([]publicAPIEntry, 0, len(raw.Entries))
		for _, data := range raw.Entries {
			if err := ctx.Err(); err != nil {
				return nil, fmt.Errorf("parse public-apis catalog %q: %w", sourceURL, err)
			}
			if len(data) > publicAPIsMaxRowBytes {
				return fail(fmt.Sprintf("entry exceeds %d bytes", publicAPIsMaxRowBytes))
			}
			var entry publicAPIEntry
			if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
				return fail("JSON entry must be an object")
			}
			if err := json.Unmarshal(data, &entry); err != nil {
				return fail(err.Error())
			}
			entries = append(entries, entry)
		}
		return entries, nil
	}

	var entries []publicAPIEntry
	category := ""
	inTable := false
	scanner := bufio.NewScanner(bytes.NewReader(content))
	scanner.Buffer(make([]byte, 4096), publicAPIsMaxRowBytes+2)
	for scanner.Scan() {
		if err := ctx.Err(); err != nil {
			return nil, fmt.Errorf("parse public-apis catalog %q: %w", sourceURL, err)
		}
		if len(scanner.Bytes()) > publicAPIsMaxRowBytes {
			return fail(fmt.Sprintf("row exceeds %d bytes", publicAPIsMaxRowBytes))
		}
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "### ") {
			category = strings.TrimSpace(strings.TrimPrefix(line, "### "))
			inTable = false
			continue
		}
		if !strings.HasPrefix(line, "|") {
			inTable = false
			continue
		}
		columns := publicAPIsColumns(line)
		if len(columns) == 5 && strings.EqualFold(columns[0], "API") && strings.EqualFold(columns[1], "Description") && strings.EqualFold(columns[2], "Auth") && strings.EqualFold(columns[3], "HTTPS") && strings.EqualFold(columns[4], "CORS") {
			inTable = true
			continue
		}
		if !inTable {
			continue
		}
		if len(columns) != 5 {
			return fail("malformed Markdown table row: expected five columns")
		}
		separator := true
		for _, column := range columns {
			if !strings.Contains(column, "-") || strings.Trim(column, " -:") != "" {
				separator = false
				break
			}
		}
		if separator {
			continue
		}
		name, link, ok := publicAPIsMarkdownLink(columns[0])
		if !ok || category == "" {
			return fail("malformed Markdown API link or missing category")
		}
		if len(entries) == publicAPIsMaxEntries {
			return fail(fmt.Sprintf("entry count exceeds %d", publicAPIsMaxEntries))
		}
		entries = append(entries, publicAPIEntry{API: name, Description: columns[1], Link: link, Category: category})
	}
	if err := scanner.Err(); err != nil {
		return fail(fmt.Sprintf("row exceeds %d bytes or cannot be read: %v", publicAPIsMaxRowBytes, err))
	}
	if len(entries) == 0 {
		return fail("unrecognized or empty Markdown API list")
	}
	return entries, nil
}

func publicAPIsColumns(line string) []string {
	line = strings.TrimSuffix(strings.TrimPrefix(line, "|"), "|")
	var columns []string
	start := 0
	backslashes := 0
	for i := 0; i < len(line); i++ {
		if line[i] == '|' && backslashes%2 == 0 {
			columns = append(columns, strings.ReplaceAll(strings.TrimSpace(line[start:i]), `\|`, "|"))
			start = i + 1
		}
		if line[i] == '\\' {
			backslashes++
		} else {
			backslashes = 0
		}
	}
	return append(columns, strings.ReplaceAll(strings.TrimSpace(line[start:]), `\|`, "|"))
}

func publicAPIsMarkdownLink(text string) (string, string, bool) {
	if !strings.HasPrefix(text, "[") || !strings.HasSuffix(text, ")") {
		return "", "", false
	}
	end := strings.Index(text, "](")
	if end < 2 {
		return "", "", false
	}
	name := strings.TrimSpace(text[1:end])
	link := strings.TrimSpace(text[end+2 : len(text)-1])
	// URL validation, including unsafe-host protection, belongs to the probe
	// downloader. Parsing a list never resolves or follows its links.
	if strings.HasPrefix(link, "<") && strings.HasSuffix(link, ">") {
		link = link[1 : len(link)-1]
	}
	return name, link, name != "" && link != ""
}
