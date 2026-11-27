package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

const nextTestCommit = "0123456789abcdef0123456789abcdef01234568"

func TestCompareReviewedCatalogs(t *testing.T) {
	before := siteCatalog{SchemaVersion: 1, SourceCommit: testCommit, Providers: []siteProvider{
		{ID: "agnostic", Name: "Same", URL: "https://previous.example.com/path"},
		{ID: "retired", Name: "Retired"},
		{ID: "updated", Coverage: coverageDocs},
	}}
	after := siteCatalog{SchemaVersion: 1, SourceCommit: nextTestCommit, Providers: []siteProvider{
		{ID: "added", Name: "New"},
		{ID: "agnostic", Name: "Same", URL: "https://next.example.com/path"},
		{ID: "updated", Coverage: coverageOverlay},
	}}
	files := make([]string, 0, 2)
	for i, fixture := range []siteCatalog{before, after} {
		contents, err := json.Marshal(fixture)
		if err != nil {
			t.Fatal(err)
		}
		file := filepath.Join(t.TempDir(), fmt.Sprintf("catalog-%d.json", i))
		if err := os.WriteFile(file, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		files = append(files, file)
	}
	var output bytes.Buffer
	err := run([]string{"compare", "-from", files[0], "-to", files[1], "-period", "2026-11", "-first-publication"}, &output, io.Discard)
	if err != nil {
		t.Fatal(err)
	}
	var record monthlyRecord
	if err := json.Unmarshal(output.Bytes(), &record); err != nil {
		t.Fatal(err)
	}
	want := monthlyRecord{Period: "2026-11", FromCommit: testCommit, ToCommit: nextTestCommit, FirstPublication: true,
		Added: []string{"added"}, Changed: []string{"updated"}, Removed: []string{"retired"}}
	if !reflect.DeepEqual(record, want) {
		t.Fatalf("compared record = %+v, want %+v", record, want)
	}
}

func TestCommittedChangeHistory(t *testing.T) {
	history, err := loadChangeHistory()
	if err != nil {
		t.Fatal(err)
	}
	var output = make(map[string][]byte)
	view := pageView{Title: "Monthly changes", Kind: "changes", History: history}
	if err := addPage(output, "index.html", view); err != nil {
		t.Fatal(err)
	}
	actual := string(output["index.html"])
	if strings.Count(actual, "First publication. There is no previous site publication.") != 1 {
		t.Error("history must identify the first publication once")
	}
	for _, record := range history {
		_, article, found := strings.Cut(actual, `<time datetime="`+record.Period+`">`)
		article, _, closed := strings.Cut(article, "</article>")
		if !found || !closed || !strings.Contains(article, record.FromCommit) || !strings.Contains(article, record.ToCommit) {
			t.Errorf("missing rendered record for %s", record.Period)
		}
		for index, list := range [][]string{record.Added, record.Changed, record.Removed} {
			label := []string{"added", "changed", "removed"}[index]
			if !strings.Contains(article, fmt.Sprintf("<strong>%d</strong> %s", len(list), label)) {
				t.Errorf("rendered %s count for %s differs from the record", label, record.Period)
			}
			for _, id := range list {
				if !strings.Contains(article, "<code>"+id+"</code>") {
					t.Errorf("rendered %s %s is missing for %s", label, id, record.Period)
				}
			}
		}
	}
}

func TestLoopbackPreviewRollback(t *testing.T) {
	base, err := parseBaseURL("https://apis.kinet.com/")
	if err != nil {
		t.Fatal(err)
	}
	data, err := builtInSiteCatalog(testCommit, base)
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	filesByCommit := make(map[string]map[string][]byte)
	for _, revision := range []string{testCommit, nextTestCommit} {
		data.SourceCommit = revision
		data.Providers[0].Description = "Reviewed preview fixture " + revision
		files, err := buildFiles(data)
		if err != nil {
			t.Fatal(err)
		}
		filesByCommit[revision] = files
		if err := writeOutput(filepath.Join(root, revision), files); err != nil {
			t.Fatal(err)
		}
	}
	address := "127.0.0.1:0"
	firstProvider := "providers/" + data.Providers[0].ID + "/index.html"
	for _, revision := range []string{testCommit, nextTestCommit, testCommit} {
		func() {
			artifact := filepath.Join(root, revision)
			manifest := verifyPreviewManifest(t, artifact, filesByCommit[revision])
			listener, err := net.Listen("tcp4", address)
			if err != nil {
				t.Fatal(err)
			}
			address = listener.Addr().String()
			preview := httptest.NewUnstartedServer(http.FileServer(http.Dir(artifact)))
			preview.Listener.Close()
			preview.Listener = listener
			preview.Start()
			defer preview.Close()
			client := preview.Client()
			client.Timeout = 3 * time.Second
			for _, name := range []string{"index.html", "catalog.json", "changes/index.html", "privacy/index.html", firstProvider, "SHA256SUMS"} {
				response, err := client.Get("http://" + address + "/" + name)
				if err != nil {
					t.Fatal(err)
				}
				contents, err := io.ReadAll(response.Body)
				response.Body.Close()
				if err != nil || response.StatusCode != http.StatusOK {
					t.Fatalf("preview of %s at %s returned %d: %v", name, revision, response.StatusCode, err)
				}
				expected := filesByCommit[revision][name]
				if name == "SHA256SUMS" {
					expected = manifest
				}
				if !bytes.Equal(expected, contents) {
					t.Fatalf("preview served incorrect %s for %s", name, revision)
				}
			}
		}()
	}
}

func verifyPreviewManifest(t *testing.T, artifact string, files map[string][]byte) []byte {
	t.Helper()
	manifest, err := os.ReadFile(filepath.Join(artifact, "SHA256SUMS"))
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(manifest)), "\n")
	seen := make(map[string]bool, len(lines))
	for _, line := range lines {
		digest, name, split := strings.Cut(line, "  ")
		if _, exists := files[name]; !split || !exists || seen[name] {
			t.Fatalf("invalid preview manifest line %q", line)
		}
		seen[name] = true
		contents, err := os.ReadFile(filepath.Join(artifact, filepath.FromSlash(name)))
		if err != nil || fmt.Sprintf("%x", sha256.Sum256(contents)) != digest {
			t.Fatalf("preview manifest verification failed for %s: %v", name, err)
		}
	}
	if len(seen) != len(files) {
		t.Fatalf("preview manifest has %d entries, want %d", len(seen), len(files))
	}
	return manifest
}
