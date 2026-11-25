package apitools

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func TestC07PinnedPublisherOperations(t *testing.T) {
	qualificationDir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if qualificationDir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	for _, fixture := range []struct {
		provider, reference, file, method, path, operationID string
	}{
		{"hubspot", "hubspot-crm-contacts-2026-09-openapi", "hubspot-contacts-2026-09-424.json", "GET", "/crm/objects/2026-09/{objectType}", "get-/crm/objects/2026-09/{objectType}_/crm/objects/2026-03/{objectType}"},
		{"hubspot", "hubspot-crm-deals-2026-09-openapi", "hubspot-deals-2026-09-424.json", "GET", "/crm/objects/2026-09/{objectType}", "get-/crm/objects/2026-09/{objectType}_/crm/objects/2026-03/{objectType}"},
		{"kubernetes", "kubernetes-v1-37-2-core-openapi-v3", "kubernetes-commit-core-v3.json", "GET", "/api/v1/namespaces/{namespace}/pods", "listCoreV1NamespacedPod"},
	} {
		t.Run(fixture.reference, func(t *testing.T) {
			provider, ok := catalog.FindBuiltInProvider(fixture.provider)
			if !ok {
				t.Fatalf("missing provider %s", fixture.provider)
			}
			var source catalogUpdateSource
			for _, candidate := range catalogUpdateReferenceSources(provider) {
				if candidate.ReferenceID == fixture.reference {
					source = candidate
					break
				}
			}
			content, err := os.ReadFile(filepath.Join(qualificationDir, fixture.file))
			if err != nil {
				t.Fatal(err)
			}
			digest := fmt.Sprintf("%x", sha256.Sum256(content))
			if catalogUpdateReviewedDigest(source.HeldRevision) != digest {
				t.Fatalf("catalog pin did not identify the retained artifact: %s", digest)
			}
			if validation := catalogUpdateNativeOpenAPI(t.Context(), content, source.URL, source.refreshReference(), digest); !validation.accepted() {
				t.Fatalf("reviewed source failed updater validation: %+v", validation)
			}
			table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
				ID: fixture.reference, OperationSourceInput: OperationSourceInput{Kind: OperationSourceOpenAPI, URL: source.URL, Content: content},
			}}})
			if err != nil || len(table.Sources) != 1 || table.Sources[0].SHA256 != digest {
				t.Fatalf("native source identity not established: %v", err)
			}
			for _, shape := range table.Operations {
				if shape.Method == fixture.method && shape.Path == fixture.path {
					for _, alias := range shape.Aliases {
						if alias.Value == fixture.operationID {
							return
						}
					}
					t.Fatalf("operation %s %s lacks exact alias %s", fixture.method, fixture.path, fixture.operationID)
				}
			}
			t.Fatalf("native source omitted %s %s", fixture.method, fixture.path)
		})
	}
}

func TestC07HubSpotReleaseCompatibility(t *testing.T) {
	dir := os.Getenv("APITOOLS_C07_QUALIFICATION_DIR")
	if dir == "" {
		t.Skip("set APITOOLS_C07_QUALIFICATION_DIR to the pinned C07 source directory")
	}
	const collection = "https://raw.githubusercontent.com/HubSpot/HubSpot-public-api-spec-collection/c6f46a91db551f1f4fcbe61050ea002050e27d55/PublicApiSpecs/CRM/"
	const selectedKey = "#/paths/~1crm~1objects~12026-09~1{objectType}/get"
	for _, service := range []struct {
		name, group, legacySelector, legacyDigest, selectedDigest, rolloutDigest string
	}{
		{"contacts", "Contacts", "#/paths/~1crm~1v3~1objects~1contacts/get", "f144b194fc68e7d6ea338b7f89882d247479a80437c0f581456c10f1cc90bd5c", "f67bf5699a481ef85735c1846884f3c7c72ea85f858341de3499a90d9039f344", "0f2fb5042f2fd9c18d70e053ed8503978bb9cd641cc82c40e9d841d7f63911e9"},
		{"deals", "Deals", "#/paths/~1crm~1v3~1objects~10-3/get", "3b9f0d5a0f6b0684aa4210c661ccadb13df513238c01b811f13902a5a4cd2cf0", "cafe9ee81824c7994fdf3298ef483d5c7fddcebd869ad02ea8b23b09968f1aa1", "ef6fa7c6149224af71c811b1e86c1b6b0904fda46019525781a66d7cccb84d23"},
	} {
		t.Run(service.name, func(t *testing.T) {
			read := func(file, url, expectedDigest string) OperationSourceInput {
				t.Helper()
				content, err := os.ReadFile(filepath.Join(dir, file))
				if err != nil {
					t.Fatal(err)
				}
				if fmt.Sprintf("%x", sha256.Sum256(content)) != expectedDigest {
					t.Fatalf("qualification artifact changed: %s", file)
				}
				return OperationSourceInput{URL: collection + service.group + "/Rollouts/" + url + "/" + service.name + ".json", Content: content, Kind: OperationSourceOpenAPI}
			}
			prefix := "hubspot-" + service.name
			legacy := read(prefix+"-v3-424.json", "424/v3", service.legacyDigest)
			selected := read(prefix+"-2026-09-424.json", "424/2026-09", service.selectedDigest)
			alternate := read(prefix+"-2026-09.json", "321895/2026-09", service.rolloutDigest)
			for _, comparison := range []struct {
				baseline, candidate OperationSourceInput
				key                 string
			}{
				{legacy, selected, service.legacySelector},
				{selected, alternate, selectedKey},
			} {
				table, err := BuildOperationShapeTable(t.Context(), OperationShapeOptions{Sources: []ShapeSourceInput{{
					ID: "baseline", OperationSourceInput: comparison.baseline,
				}}})
				if _, ok := indexUpdateShapes(table.Operations)[comparison.key]; err != nil || !ok {
					t.Fatalf("baseline does not reproduce the reviewed key %s: %v", comparison.key, err)
				}
				got := compareCatalogUpdate(t.Context(), comparison.baseline, comparison.candidate, []string{comparison.key}, catalog.RefreshableSpecReference{}, "")
				if got != catalogUpdateChanged && got != catalogUpdateUnexamined {
					t.Fatalf("unproven HubSpot migration was marked %s for %s", got, comparison.key)
				}
				t.Logf("%s: %s", comparison.key, got)
			}
		})
	}
}

func TestC07HubSpotSecuritySchemesHaveSpecScope(t *testing.T) {
	view, err := catalog.BuiltInSecurityInspectionView("hubspot")
	if err != nil {
		t.Fatal(err)
	}
	for _, conflict := range view.Conflicts {
		if conflict.Type == catalog.SecurityInspectionConflictDuplicateScheme {
			t.Fatalf("HubSpot qualification introduced a conflicting scheme: %+v", conflict)
		}
	}
}
