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
