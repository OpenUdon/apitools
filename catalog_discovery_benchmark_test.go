package apitools_test

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/sqlitecache"
)

// Explicit operator-selected roots only. This benchmark does not prepare,
// refresh, fetch, parse raw artifacts, or alter the selected root.
func BenchmarkCatalogDiscoveryLocalRoot(b *testing.B) {
	root := os.Getenv("APITOOLS_CATALOG_BENCHMARK_ROOT")
	if root == "" {
		b.Skip("set APITOOLS_CATALOG_BENCHMARK_ROOT to an existing populated/indexed local root")
	}
	options := apitools.CatalogDiscoveryOptions{Index: apitools.CatalogIndexOptions{Root: catalog.RootOptions{Directory: root}, ReadRegistrations: sqlitecache.ReadCatalogSpecArtifacts}, Request: apitools.CatalogDiscoveryRequest{SchemaVersion: apitools.CatalogDiscoverySchemaVersion, Contract: apitools.StepContract{Purpose: "list notes records", Effect: apitools.OperationEffectRead}}}
	if name := os.Getenv("APITOOLS_CATALOG_BENCHMARK_CATALOG"); name != "" {
		if filepath.IsAbs(name) || filepath.Clean(name) != name || filepath.Base(name) != name {
			b.Fatal("benchmark catalog must be one confined root-relative filename")
		}
		paths, err := catalog.ResolveRoot(options.Index.Root)
		if err != nil {
			b.Fatal(err)
		}
		path := filepath.Join(paths.Directory, name)
		info, err := os.Lstat(path)
		if err != nil || !info.Mode().IsRegular() {
			b.Fatal("benchmark catalog must be an existing regular file")
		}
		file, err := os.Open(path)
		if err != nil {
			b.Fatal(err)
		}
		data, err := io.ReadAll(io.LimitReader(file, 4<<20+1))
		_ = file.Close()
		if err != nil || len(data) > 4<<20 {
			b.Fatal("benchmark catalog exceeds 4 MiB")
		}
		var cat catalog.Catalog
		if err := json.Unmarshal(data, &cat); err != nil {
			b.Fatal(err)
		}
		options.Index.Catalog = &cat
	}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
		if err != nil || report.Outcome == apitools.CatalogDiscoveryBlocked {
			b.Fatalf("configured benchmark root failed: %v %s", err, report.Outcome)
		}
		if report.ExaminedOperations == 0 {
			b.Fatal("benchmark root has no examined operations; explicitly prepare its index first")
		}
		b.ReportMetric(float64(report.ExaminedOperations), "operations/query")
		b.ReportMetric(float64(report.QualifiedOperations), "qualified/query")
		encoded, err := json.Marshal(report)
		if err != nil {
			b.Fatal(err)
		}
		b.ReportMetric(float64(len(encoded)), "report-bytes/query")
	}
}
