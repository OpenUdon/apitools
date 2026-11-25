package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestReviewedSearchCLITextAndJSON(t *testing.T) {
	dir := t.TempDir()
	config := filepath.Join(dir, "search.json")
	if err := os.WriteFile(config, []byte(`{"index_db":"reviewed.sqlite","state_db":"missing.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	for _, jsonMode := range []bool{false, true} {
		args := []string{"--query", "send slack message", "--source", "reviewed", "--config", config, "--limit", "1"}
		if jsonMode {
			args = append(args, "--json")
		}
		out.Reset()
		errOut.Reset()
		code := runSearchWithClient(args, &out, &errOut, func(string) (*apitools.Client, func(), error) {
			t.Fatal("explicit reviewed search opened the public directory client")
			return nil, func() {}, nil
		})
		if code != exitSuccess {
			t.Fatalf("exit %d: %s", code, errOut.String())
		}
		if jsonMode {
			var report apitools.ReviewedProviderSearchReport
			if err := json.Unmarshal(out.Bytes(), &report); err != nil {
				t.Fatal(err)
			}
			if report.Source != "reviewed-catalog" || report.PublisherState != "unavailable" ||
				len(report.Results) != 1 || report.Results[0].Profile.ProviderID != "slack" ||
				len(report.Results[0].MatchedCapabilities) != 1 {
				t.Fatalf("unexpected reviewed search: %+v", report)
			}
		} else {
			for _, want := range []string{"Slack", "covered-operation", "selected-subset", "Send Slack message",
				"chat.postMessage", "advisory-overlay", "Catalog review date",
				"Operation source slack-web-api-docs (slack-web-message-history-overlay) sha256:"} {
				if !strings.Contains(out.String(), want) {
					t.Fatalf("text result missing %q:\n%s", want, out.String())
				}
			}
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "reviewed.sqlite")); err != nil {
		t.Fatal(err)
	}
}

func TestReviewedSearchConfigRejectsUnknownFields(t *testing.T) {
	p := filepath.Join(t.TempDir(), "search.json")
	for _, invalid := range []string{`{"metadata_endpoint":"https://example.test"}`, `null`, `{}{}`} {
		if err := os.WriteFile(p, []byte(invalid), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := readReviewedSearchConfig(p); err == nil {
			t.Fatalf("accepted malformed config %s", invalid)
		}
	}
}

func TestReviewedSearchOfflineConfigDisablesEmbeddingRequests(t *testing.T) {
	path := filepath.Join(t.TempDir(), "search.json")
	data := []byte(`{
	  "embeddings": {
	    "endpoint": "https://127.0.0.1/v1/embeddings",
	    "model": "fixture",
	    "max_requests": 16,
	    "max_inputs": 350,
	    "max_bytes": 1048576
	  }
	}`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	for _, mode := range [][]string{{"--offline"}, {"--cache-mode", "offline"}} {
		var out, errOut bytes.Buffer
		args := append([]string{"search", "--query", "send slack message", "--source", "reviewed", "--config", path, "--json"}, mode...)
		if code := run(args, &out, &errOut); code != exitSuccess {
			t.Fatalf("exit %d: %s", code, errOut.String())
		}
		var report apitools.ReviewedProviderSearchReport
		if err := json.Unmarshal(out.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		if report.Embeddings == nil || report.Embeddings.Status != "offline" || report.Embeddings.Requests != 0 ||
			len(report.Results) == 0 || report.Results[0].Profile.ProviderID != "slack" {
			t.Fatalf("offline search did not stay lexical: %+v", report.Embeddings)
		}
	}
}

func TestReviewedSearchRejectsAnInvalidCacheMode(t *testing.T) {
	var out, errOut bytes.Buffer
	code := run([]string{"search", "--source", "reviewed", "--query", "send slack message",
		"--cache-mode", "offine", "--config", "missing-config.json"}, &out, &errOut)
	if code != exitRuntime || !strings.Contains(errOut.String(), "unknown cache mode") || out.Len() != 0 {
		t.Fatalf("invalid mode allowed lookup: %d, %q, %q", code, out.String(), errOut.String())
	}
}
