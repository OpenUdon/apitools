package apitools_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/sqlitecache"
)

func syntheticCatalogIndex(t *testing.T) apitools.CatalogIndexOptions {
	t.Helper()
	dir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dir, "openapi"), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"catalog.json", "registrations.json", "openapi/notes.json"} {
		data, err := os.ReadFile(filepath.Join("testdata/catalog-root", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	var cat catalog.Catalog
	data, _ := os.ReadFile(filepath.Join(dir, "catalog.json"))
	if err := json.Unmarshal(data, &cat); err != nil {
		t.Fatal(err)
	}
	var rows []sqlitecache.CatalogArtifact
	data, _ = os.ReadFile(filepath.Join(dir, "registrations.json"))
	if err := json.Unmarshal(data, &rows); err != nil {
		t.Fatal(err)
	}
	cache, err := sqlitecache.Open(filepath.Join(dir, "cache.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range rows {
		if err := cache.StoreCatalogArtifact(context.Background(), row); err != nil {
			t.Fatal(err)
		}
	}
	if err := cache.Close(); err != nil {
		t.Fatal(err)
	}
	return apitools.CatalogIndexOptions{Root: catalog.RootOptions{Directory: dir}, Catalog: &cat, ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts}
}

func TestCatalogIndexSharedDeterministicRelocation(t *testing.T) {
	first, second := syntheticCatalogIndex(t), syntheticCatalogIndex(t)
	ctx := context.Background()
	index, err := apitools.BuildCatalogOperationIndex(ctx, first)
	if err != nil {
		t.Fatal(err)
	}
	if len(index.Artifacts) != 1 || len(index.Artifacts[0].Links) != 2 || len(index.Artifacts[0].Operations) != 2 || index.Artifacts[0].State != "indexed" {
		t.Fatalf("shared artifact: %#v", index.Artifacts)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, first, index); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(first.Root.Directory, catalog.DefaultOperationIndexPath))
	if bytes.Contains(before, []byte(first.Root.Directory)) {
		t.Fatal("absolute path in index")
	}
	// Permute catalog and registration input order; timestamp differences in
	// the two prepared caches must not affect deterministic metadata.
	second.Catalog.Providers[0], second.Catalog.Providers[1] = second.Catalog.Providers[1], second.Catalog.Providers[0]
	reader := second.ReadRegistrations
	second.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		rows[0], rows[1] = rows[1], rows[0]
		return rows, err
	}
	other, err := apitools.BuildCatalogOperationIndex(ctx, second)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, second, other); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(filepath.Join(second.Root.Directory, catalog.DefaultOperationIndexPath))
	if !bytes.Equal(before, after) {
		t.Fatal("relocation/permutation changed index bytes")
	}
	if err := os.Remove(filepath.Join(first.Root.Directory, "openapi/notes.json")); err != nil {
		t.Fatal(err)
	}
	loaded, err := apitools.ReadCatalogOperationIndex(ctx, first)
	if err != nil || len(loaded.Artifacts[0].Operations) != 2 {
		t.Fatalf("index reader reparsed source: %v", err)
	}
	missing, err := apitools.BuildCatalogOperationIndex(ctx, first)
	if err != nil || missing.Artifacts[0].State != "missing" {
		t.Fatalf("missing coverage: %v %#v", err, missing.Artifacts)
	}
}

func TestCatalogIndexRefusesInputOverwriteAndRegistryDrift(t *testing.T) {
	options := syntheticCatalogIndex(t)
	ctx := context.Background()
	index, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, options, index); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(filepath.Join(options.Root.Directory, catalog.DefaultOperationIndexPath))
	unsafe := options
	unsafe.Root.IndexPath = "openapi/notes.json"
	if _, err := apitools.BuildCatalogOperationIndex(ctx, unsafe); err == nil {
		t.Fatal("accepted source overwrite destination")
	}
	unsafe = options
	unsafe.Root.IndexPath = "cache.sqlite-wal"
	if err := apitools.WriteCatalogOperationIndex(ctx, unsafe, index); err == nil {
		t.Fatal("accepted registry sidecar overwrite")
	}
	reader := options.ReadRegistrations
	options.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		rows[0].SHA256 = strings.Repeat("0", 64)
		return rows, err
	}
	stale, err := apitools.ReadCatalogOperationIndex(ctx, options)
	if err != nil || stale.Artifacts[0].State != "stale" || len(stale.Artifacts[0].Operations) != 0 {
		t.Fatalf("stale scope: %v %#v", err, stale.Artifacts)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, options, index); err == nil {
		t.Fatal("published obsolete generation")
	}
	after, _ := os.ReadFile(filepath.Join(options.Root.Directory, catalog.DefaultOperationIndexPath))
	if !bytes.Equal(before, after) {
		t.Fatal("failed publication changed old index")
	}
	options.ReadRegistrations = reader
	calls := 0
	options.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		calls++
		if calls > 1 {
			rows[0].SHA256 = strings.Repeat("0", 64)
		}
		return rows, err
	}
	if _, err := apitools.BuildCatalogOperationIndex(ctx, options); err == nil {
		t.Fatal("accepted registry drift while building")
	}
}

func TestCatalogIndexDigestFailureAndCoverage(t *testing.T) {
	options := syntheticCatalogIndex(t)
	ctx := context.Background()
	path := filepath.Join(options.Root.Directory, "openapi/notes.json")
	data, _ := os.ReadFile(path)
	data[0] = ' '
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	index, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil || index.Artifacts[0].State != "digest_mismatch" || len(index.Artifacts[0].Operations) != 0 {
		t.Fatalf("digest failure: %v %#v", err, index.Artifacts)
	}
	reader := options.ReadRegistrations
	options.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		for i := range rows {
			rows[i].Bytes = 129 << 20
		}
		return rows, err
	}
	index, err = apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil || index.Artifacts[0].State != "oversize" {
		t.Fatalf("oversize coverage: %v", err)
	}
	options.ReadRegistrations = func(context.Context, catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) { return nil, nil }
	index, err = apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil || len(index.Artifacts) != 0 || len(index.Coverage) != 2 {
		t.Fatalf("missing registrations: %v", err)
	}
	for _, row := range index.Coverage {
		if row.State != "missing" {
			t.Fatal("missing scope misrepresented")
		}
	}
}

func TestCatalogIndexRejectsCorruptMetadata(t *testing.T) {
	options := syntheticCatalogIndex(t)
	ctx := context.Background()
	original, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(original)
	mutations := map[string]func(*apitools.CatalogOperationIndex){
		"version":            func(i *apitools.CatalogOperationIndex) { i.MetadataVersion = "future" },
		"digest":             func(i *apitools.CatalogOperationIndex) { i.RegistrationSHA256 = strings.Repeat("z", 64) },
		"missing coverage":   func(i *apitools.CatalogOperationIndex) { i.Coverage = i.Coverage[:1] },
		"false coverage":     func(i *apitools.CatalogOperationIndex) { i.Coverage[0].State = "missing" },
		"duplicate coverage": func(i *apitools.CatalogOperationIndex) { i.Coverage[1] = i.Coverage[0] },
		"unregistered link":  func(i *apitools.CatalogOperationIndex) { i.Artifacts[0].Links[0].SpecRefID = "invented" },
		"unsafe path":        func(i *apitools.CatalogOperationIndex) { i.Artifacts[0].Path = "openapi/../notes.json" },
		"wrong count":        func(i *apitools.CatalogOperationIndex) { i.Artifacts[0].OperationCount++ },
		"wrong selector identity": func(i *apitools.CatalogOperationIndex) {
			i.Artifacts[0].Operations[0].Source.SHA256 = strings.Repeat("0", 64)
		},
		"duplicate operation": func(i *apitools.CatalogOperationIndex) { i.Artifacts[0].Operations[1] = i.Artifacts[0].Operations[0] },
		"secret provenance": func(i *apitools.CatalogOperationIndex) {
			i.Artifacts[0].Operations[0].Source.DocumentURL = "https://example.invalid/spec?token=secret"
		},
		"false authority": func(i *apitools.CatalogOperationIndex) { i.Artifacts[0].Links[0].Advisory = true },
	}
	for name, mutate := range mutations {
		t.Run(name, func(t *testing.T) {
			var index apitools.CatalogOperationIndex
			if err := json.Unmarshal(encoded, &index); err != nil {
				t.Fatal(err)
			}
			mutate(&index)
			if err := apitools.WriteCatalogOperationIndex(ctx, options, index); err == nil {
				t.Fatal("published corrupt metadata")
			}
			data, _ := json.Marshal(index)
			if err := os.WriteFile(filepath.Join(options.Root.Directory, catalog.DefaultOperationIndexPath), data, 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := apitools.ReadCatalogOperationIndex(ctx, options); err == nil {
				t.Fatal("read corrupt metadata")
			}
		})
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, options, original); err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := apitools.ReadCatalogOperationIndex(cancelled, options); err == nil {
		t.Fatal("ignored cancellation")
	}
	if err := apitools.WriteCatalogOperationIndex(cancelled, options, original); err == nil {
		t.Fatal("published after cancellation")
	}
	// An invalid source records unexamined coverage and retains no candidates.
	reader := options.ReadRegistrations
	data := []byte(`{"openapi":"3.0.3","info":{"title":"invalid","version":"1"}}`)
	if err := os.WriteFile(filepath.Join(options.Root.Directory, "openapi/notes.json"), data, 0o600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	options.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		for j := range rows {
			rows[j].SHA256 = hex.EncodeToString(digest[:])
			rows[j].Bytes = int64(len(data))
		}
		return rows, err
	}
	index, err := apitools.BuildCatalogOperationIndex(ctx, options)
	if err != nil || index.Artifacts[0].State != "parse_failure" || len(index.Artifacts[0].Operations) != 0 {
		t.Fatalf("invalid source treated as examined: %v", err)
	}
	if err := apitools.WriteCatalogOperationIndex(ctx, options, index); err != nil {
		t.Fatal(err)
	}
}
