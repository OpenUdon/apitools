//go:build ignore

// prepare creates a new disposable synthetic catalog root. It never refreshes
// remote sources and refuses an existing destination.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/sqlitecache"
)

func main() {
	if err := prepare(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func prepare() error {
	if len(os.Args) != 2 || !filepath.IsAbs(os.Args[1]) {
		return fmt.Errorf("usage: go run testdata/catalog-root/prepare.go /absolute/new/directory")
	}
	_, source, _, _ := runtime.Caller(0)
	fixture := filepath.Dir(source)
	var registrations []sqlitecache.CatalogArtifact
	data, err := os.ReadFile(filepath.Join(fixture, "registrations.json"))
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &registrations); err != nil {
		return err
	}
	dest := os.Args[1]
	if err := os.Mkdir(dest, 0o700); err != nil {
		return err
	}
	if err := os.Mkdir(filepath.Join(dest, "openapi"), 0o700); err != nil {
		return err
	}
	for _, name := range []string{"catalog.json", "registrations.json", "openapi/notes.json"} {
		data, err := os.ReadFile(filepath.Join(fixture, name))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dest, name), data, 0o600); err != nil {
			return err
		}
	}
	cache, err := sqlitecache.Open(filepath.Join(dest, catalog.DefaultRegistryPath))
	if err != nil {
		return err
	}
	for _, artifact := range registrations {
		if err := cache.StoreCatalogArtifact(context.Background(), artifact); err != nil {
			_ = cache.Close()
			return err
		}
	}
	return cache.Close()
}
