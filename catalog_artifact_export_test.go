package apitools_test

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

func catalogExportReferences(t *testing.T, options apitools.CatalogIndexOptions) []apitools.CatalogArtifactReference {
	t.Helper()
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), options)
	if err != nil {
		t.Fatal(err)
	}
	artifact := index.Artifacts[0]
	var refs []apitools.CatalogArtifactReference
	for _, link := range artifact.Links {
		refs = append(refs, apitools.CatalogArtifactReference{CatalogSHA256: index.CatalogSHA256, ProviderID: link.ProviderID, SpecRefID: link.SpecRefID, ArtifactID: link.ArtifactID, Kind: artifact.Kind, SHA256: artifact.SHA256, Bytes: artifact.Bytes, Selector: artifact.Operations[0].Source.Selector})
	}
	return refs
}

func TestCatalogArtifactExportScopedSharedAtomic(t *testing.T) {
	options := syntheticCatalogIndex(t)
	other := options.Catalog.Providers[0].SpecReferences[0]
	other.ID = "other-notes-spec"
	options.Catalog.Providers[0].SpecReferences = append(options.Catalog.Providers[0].SpecReferences, other)
	provider := options.Catalog.Providers[0].ID
	spec := options.Catalog.Providers[0].SpecReferences[0].ID
	for _, row := range []struct{ id, provider, spec string }{{"provider-auth", provider, ""}, {"selected-auth", provider, spec}, {"other-spec-auth", provider, other.ID}, {"other-provider-auth", options.Catalog.Providers[1].ID, ""}} {
		options.Catalog.SecurityOverlays = append(options.Catalog.SecurityOverlays, catalog.SecurityOverlay{ID: row.id, ProviderID: row.provider, SpecRefID: row.spec, Status: catalog.AuthStatusUnknown, SourceNote: "Synthetic auth metadata only.", SourceRefs: []string{"https://example.invalid/auth"}})
	}
	refs := catalogExportReferences(t, options)
	reader := options.ReadRegistrations
	options.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		for j := range rows {
			rows[j].SourceURL = "https://user:PROVENANCE_SECRET@example.invalid/spec?token=PROVENANCE_SECRET"
			rows[j].Metadata = map[string]string{"private": "REGISTRY_SECRET"}
		}
		extra := rows[0]
		extra.ProviderID = provider
		extra.SpecRefID = other.ID
		extra.ArtifactID = "unselected-artifact"
		return append(rows, extra), err
	}
	// Select only one of the provider links, repeating it to prove exact dedup.
	var selected apitools.CatalogArtifactReference
	for _, ref := range refs {
		if ref.ProviderID == provider {
			selected = ref
		}
	}
	workflow := t.TempDir()
	export := apitools.CatalogArtifactExportOptions{Index: options, References: []apitools.CatalogArtifactReference{selected, selected}, WorkflowDir: workflow}
	report, err := apitools.ExportCatalogArtifacts(context.Background(), export)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Artifacts) != 1 || len(report.SecurityOverlays) != 2 {
		t.Fatalf("scope leaked: %#v", report)
	}
	target := filepath.Join(workflow, "api-artifacts")
	source, _ := os.ReadFile(filepath.Join(options.Root.Directory, "openapi/notes.json"))
	saved, _ := os.ReadFile(filepath.Join(target, filepath.FromSlash(report.Artifacts[0].Path)))
	if !bytes.Equal(source, saved) {
		t.Fatal("changed raw source")
	}
	manifest, _ := os.ReadFile(filepath.Join(target, "provenance.json"))
	if bytes.Contains(manifest, []byte("PROVENANCE_SECRET")) || bytes.Contains(manifest, []byte("REGISTRY_SECRET")) {
		t.Fatal("private registry metadata in provenance")
	}
	if bytes.Contains(manifest, []byte(options.Root.Directory)) || bytes.Contains(manifest, []byte(workflow)) {
		t.Fatal("absolute path in provenance")
	}
	again, err := apitools.ExportCatalogArtifacts(context.Background(), export)
	if err != nil || !again.Reused {
		t.Fatalf("identical reuse: %v", err)
	}
	// Change selected provenance: collision is explicit, Force replaces atomically.
	export.References = refs
	if _, err := apitools.ExportCatalogArtifacts(context.Background(), export); err == nil {
		t.Fatal("silently replaced a different selection")
	}
	before, _ := os.ReadFile(filepath.Join(target, "provenance.json"))
	if !bytes.Equal(before, manifest) {
		t.Fatal("collision changed target")
	}
	export.Force = true
	report, err = apitools.ExportCatalogArtifacts(context.Background(), export)
	if err != nil || len(report.Artifacts) != 2 || report.Artifacts[0].Path != report.Artifacts[1].Path {
		t.Fatalf("shared copy: %v %#v", err, report)
	}
	entries, err := os.ReadDir(filepath.Join(target, "openapi"))
	if err != nil || len(entries) != 1 {
		t.Fatalf("shared source duplicated: %v", err)
	}
	// One bad selection refuses the whole tree, preserving the old tree under Force.
	export.References[1].SHA256 = strings.Repeat("0", 64)
	before, _ = os.ReadFile(filepath.Join(target, "provenance.json"))
	if _, err := apitools.ExportCatalogArtifacts(context.Background(), export); err == nil {
		t.Fatal("accepted mismatched reference")
	}
	after, _ := os.ReadFile(filepath.Join(target, "provenance.json"))
	if !bytes.Equal(before, after) {
		t.Fatal("partial replacement")
	}
	children, _ := os.ReadDir(workflow)
	if len(children) != 1 {
		t.Fatal("left transaction scratch tree")
	}
}

func TestCatalogArtifactExportIntegrityAndConfinement(t *testing.T) {
	for _, name := range []string{"raw-digest", "catalog", "provider", "spec", "kind", "bytes", "selector", "escape", "source-overlap", "symlink", "cancel", "registry-drift"} {
		t.Run(name, func(t *testing.T) {
			options := syntheticCatalogIndex(t)
			refs := catalogExportReferences(t, options)
			workflow := t.TempDir()
			export := apitools.CatalogArtifactExportOptions{Index: options, References: refs, WorkflowDir: workflow, Force: true}
			ctx := context.Background()
			switch name {
			case "raw-digest":
				os.WriteFile(filepath.Join(options.Root.Directory, "openapi/notes.json"), []byte("changed"), 0o600)
			case "catalog":
				export.References[0].CatalogSHA256 = strings.Repeat("0", 64)
			case "provider":
				export.References[0].ProviderID = "unknown"
			case "spec":
				export.References[0].SpecRefID = "unknown"
			case "kind":
				export.References[0].Kind = apitools.OperationSourceKind("graphql")
			case "bytes":
				export.References[0].Bytes++
			case "selector":
				export.References[0].Selector = "bad\nselector"
			case "escape":
				export.ArtifactDir = "../escaped"
			case "source-overlap":
				export.WorkflowDir = options.Root.Directory
			case "symlink":
				if err := os.Symlink(options.Root.Directory, filepath.Join(workflow, "linked")); err != nil {
					t.Fatal(err)
				}
				export.ArtifactDir = "linked/nested"
			case "cancel":
				cancelled, cancel := context.WithCancel(ctx)
				cancel()
				ctx = cancelled
			case "registry-drift":
				reader := options.ReadRegistrations
				calls := 0
				export.Index.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
					rows, err := reader(ctx, root)
					calls++
					if calls > 1 {
						rows[0].SHA256 = strings.Repeat("0", 64)
					}
					return rows, err
				}
			}
			if _, err := apitools.ExportCatalogArtifacts(ctx, export); err == nil {
				t.Fatal("accepted invalid export")
			}
			if _, err := os.Stat(filepath.Join(workflow, "api-artifacts")); !os.IsNotExist(err) {
				t.Fatal("published partial output")
			}
		})
	}
}

func TestCatalogArtifactExportPreservesOverlayAlternatives(t *testing.T) {
	options := syntheticCatalogIndex(t)
	var overlay catalog.SecurityOverlay
	data := []byte(`{"id":"alternatives","provider_id":"example-notes","status":"present-incomplete","security_schemes":[{"name":"a","type":"apiKey","parameter_name":"X-A","in":"header"},{"name":"b","type":"apiKey","parameter_name":"X-B","in":"header"}],"root_security_sets":[{"requirements":[{"scheme":"a"}]},{"requirements":[{"scheme":"b"}]}],"source_note":"Synthetic OR alternatives","source_refs":["https://example.invalid/auth"]}`)
	if err := json.Unmarshal(data, &overlay); err != nil {
		t.Fatal(err)
	}
	options.Catalog.SecurityOverlays = []catalog.SecurityOverlay{overlay}
	refs := catalogExportReferences(t, options)
	workflow := t.TempDir()
	report, err := apitools.ExportCatalogArtifacts(context.Background(), apitools.CatalogArtifactExportOptions{Index: options, References: refs, WorkflowDir: workflow})
	if err != nil {
		t.Fatal(err)
	}
	saved, err := os.ReadFile(filepath.Join(workflow, "api-artifacts", filepath.FromSlash(report.SecurityOverlays[0].Path)))
	if err != nil {
		t.Fatal(err)
	}
	var actual catalog.SecurityOverlay
	if err := json.Unmarshal(saved, &actual); err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(actual)
	b, _ := json.Marshal(overlay)
	if !bytes.Equal(a, b) {
		t.Fatal("changed auth alternatives")
	}
}
