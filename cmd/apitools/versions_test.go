package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/OpenUdon/apitools"
)

func TestVersionsCLIExplicitNetworkAndHelp(t *testing.T) {
	var requests atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"title":"Fixture","version":"2"},"paths":{}}`))
	}))
	defer server.Close()
	req := apitools.APIVersionDiscoveryRequest{SchemaVersion: apitools.APIVersionDiscoverySchemaVersion, Known: apitools.APIVersionBaseline{SourceKind: apitools.OperationSourceOpenAPI, SourceURL: server.URL + "/v2.json", Version: "1", SHA256: strings.Repeat("a", 64)}, OfficialScopes: []apitools.APIVersionOfficialScope{{Origin: server.URL, PathPrefix: "/"}}}
	data, _ := json.Marshal(req)
	file := filepath.Join(t.TempDir(), "request.json")
	if err := os.WriteFile(file, data, 0600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"--request", file, "--json"}, {"--request", file, "--no-network", "--json"}, {"--request", file, "--network", "--json"}} {
		var out, errOut bytes.Buffer
		before := requests.Load()
		code := runVersionsWithClient(args, &out, &errOut, &apitools.Client{AllowUnsafeHosts: true})
		if code != 0 || errOut.Len() != 0 {
			t.Fatalf("code=%d stderr=%s", code, errOut.String())
		}
		var report apitools.APIVersionDiscoveryReport
		if err := json.Unmarshal(out.Bytes(), &report); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.Join(args, " "), " --network ") {
			if requests.Load() != before+1 || report.Status != "newer_found" {
				t.Fatalf("network %#v", report)
			}
		} else if requests.Load() != before || report.Status != "unexamined" {
			t.Fatalf("default fetched %#v", report)
		}
	}
	for _, args := range [][]string{{"versions", "--help"}, {"versions", "--request", file, "--network", "--no-network"}} {
		var out, errOut bytes.Buffer
		code := run(args, &out, &errOut)
		if args[1] == "--help" {
			if code != 0 || out.Len() == 0 || errOut.Len() != 0 {
				t.Fatal("help policy")
			}
		} else if code != 2 || errOut.Len() == 0 {
			t.Fatal("usage policy")
		}
	}
}
