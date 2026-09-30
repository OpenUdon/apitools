package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogIndexDispatchHelpAndRootRequired(t *testing.T) {
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "index", "--help"}, &out, &errOut); code != exitSuccess || !strings.Contains(out.String(), "--root DIRECTORY") {
		t.Fatalf("help: %d %s %s", code, out.String(), errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := run([]string{"catalog", "index"}, &out, &errOut); code != exitUsage {
		t.Fatalf("missing root: %d", code)
	}
	root := filepath.Join(t.TempDir(), "missing")
	if code := run([]string{"catalog", "index", "--root", root}, &out, &errOut); code != exitRuntime {
		t.Fatalf("missing directory: %d", code)
	}
	if _, err := os.Stat(root); !os.IsNotExist(err) {
		t.Fatal("created missing root")
	}
	if code := run([]string{"catalog", "index", "--root", t.TempDir()}, &out, &errOut); code != exitRuntime {
		t.Fatalf("missing registry: %d", code)
	}
}

func TestCatalogIndexRefusesWhitespaceInputOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "catalog.json")
	before := []byte(`{"providers":[]}`)
	if err := os.WriteFile(path, before, 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := run([]string{"catalog", "index", "--root", dir, "--catalog", "catalog.json", "--index", " catalog.json "}, &out, &errOut); code != exitUsage || !strings.Contains(errOut.String(), "overlap") {
		t.Fatalf("overwrite guard: %d %s", code, errOut.String())
	}
	after, err := os.ReadFile(path)
	if err != nil || !bytes.Equal(before, after) {
		t.Fatal("changed input metadata")
	}
}
