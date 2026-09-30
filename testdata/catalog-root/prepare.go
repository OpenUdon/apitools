//go:build ignore

// prepare creates a new disposable synthetic catalog root. It never refreshes
// remote sources and refuses an existing destination.
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

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
	count := 0
	if len(os.Args) == 4 && os.Args[2] == "--operations" {
		var err error
		count, err = strconv.Atoi(os.Args[3])
		if err != nil || count < 1 || count > 100000 {
			return fmt.Errorf("synthetic operation count must be 1 through 100000")
		}
	} else if len(os.Args) != 2 {
		return fmt.Errorf("usage: go run testdata/catalog-root/prepare.go /absolute/new/directory [--operations COUNT]")
	}
	if !filepath.IsAbs(os.Args[1]) {
		return fmt.Errorf("destination must be an absolute new directory")
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
	if count > 0 {
		data, err := os.ReadFile(filepath.Join(fixture, "openapi/notes.json"))
		if err != nil {
			return err
		}
		var doc map[string]any
		if err := json.Unmarshal(data, &doc); err != nil {
			return err
		}
		template := doc["paths"].(map[string]any)["/notes"].(map[string]any)["get"].(map[string]any)
		paths := map[string]any{}
		for i := 0; i < count; i++ {
			op := map[string]any{}
			for key, value := range template {
				op[key] = value
			}
			op["operationId"] = fmt.Sprintf("listNotes%06d", i)
			paths[fmt.Sprintf("/notes/%06d", i)] = map[string]any{"get": op}
		}
		doc["paths"] = paths
		data, err = json.Marshal(doc)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dest, "openapi/notes.json"), data, 0600); err != nil {
			return err
		}
		digest := sha256.Sum256(data)
		for i := range registrations {
			registrations[i].SHA256 = hex.EncodeToString(digest[:])
			registrations[i].Bytes = int64(len(data))
		}
		data, err = json.Marshal(registrations)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dest, "registrations.json"), data, 0600); err != nil {
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
