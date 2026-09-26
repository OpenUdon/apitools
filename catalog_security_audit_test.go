package apitools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
)

func TestBuiltInCatalogSecurityAuditClassifiesEveryProvider(t *testing.T) {
	report, err := BuiltInCatalogSecurityAuditReport(CatalogSecurityAuditOptions{})
	if err != nil {
		t.Fatal(err)
	}
	wantProviderCount := catalog.BuiltInCatalogManifest().ProviderCount
	if report.Summary.ProviderCount != wantProviderCount {
		t.Fatalf("provider count = %d, want %d", report.Summary.ProviderCount, wantProviderCount)
	}
	if countAuditDisposition(report, SecurityAuditDispositionQueuedSourceReReview) != 0 {
		t.Fatalf("queued source re-review count = %d, want 0", countAuditDisposition(report, SecurityAuditDispositionQueuedSourceReReview))
	}
	if countAuditDisposition(report, SecurityAuditDispositionPresentIncompleteReviewed) != 37 {
		t.Fatalf("present-incomplete-reviewed count = %d, want 37", countAuditDisposition(report, SecurityAuditDispositionPresentIncompleteReviewed))
	}
	for _, disposition := range []string{
		SecurityAuditDispositionCompleteUpstream,
		SecurityAuditDispositionCompleteViaOverlay,
		SecurityAuditDispositionPresentIncompleteReviewed,
		SecurityAuditDispositionIntentionallyAnonymous,
	} {
		if countAuditDisposition(report, disposition) == 0 {
			t.Fatalf("missing disposition %q in %#v", disposition, report.Summary.Dispositions)
		}
	}
	dockerRegistry := findAuditProvider(t, report, "docker-registry")
	if dockerRegistry.AuthStatus != catalog.AuthStatusPresentIncomplete {
		t.Fatalf("docker-registry auth status = %q, want %q", dockerRegistry.AuthStatus, catalog.AuthStatusPresentIncomplete)
	}
	if dockerRegistry.Disposition != SecurityAuditDispositionPresentIncompleteReviewed {
		t.Fatalf("docker-registry disposition = %q, want %q", dockerRegistry.Disposition, SecurityAuditDispositionPresentIncompleteReviewed)
	}
	if len(dockerRegistry.OverlayIDs) == 0 {
		t.Fatalf("docker-registry should retain overlay IDs")
	}
}

func TestCatalogSecurityAuditInspectsOpenAPIArtifactSecurity(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "example.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {"/widgets": {"get": {"operationId": "listWidgets"}}},
  "components": {"securitySchemes": {"BearerAuth": {"type": "http", "scheme": "bearer"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	options := CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
		},
		SecurityClassifications: []catalog.SecurityClassification{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Status:     catalog.AuthStatusPresentIncomplete,
			SourceRefs: []string{"https://example.com/docs"},
			SourceNote: "test classification",
		}},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			ArtifactID: "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	}
	report, err := BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Providers) != 1 {
		t.Fatalf("providers = %d, want 1", len(report.Providers))
	}
	artifacts := report.Providers[0].ArtifactSecurity
	if len(artifacts) != 1 {
		t.Fatalf("artifact rows = %d, want 1", len(artifacts))
	}
	row := artifacts[0]
	if row.Status != SecurityAuditArtifactMissingSecurityRequirements {
		t.Fatalf("artifact status = %q, want %q", row.Status, SecurityAuditArtifactMissingSecurityRequirements)
	}
	if row.SecuritySchemeCount != 1 || row.OperationCount != 1 || row.OperationSecurityCount != 0 {
		t.Fatalf("artifact counts = schemes:%d ops:%d secured:%d, want 1/1/0", row.SecuritySchemeCount, row.OperationCount, row.OperationSecurityCount)
	}
	if !row.SecurityDefinitionsInspected || !row.SecurityRequirementsInspected {
		t.Fatalf("artifact security inspection flags not set: %#v", row)
	}
}

func TestCatalogSecurityAuditDetectsPartialOperationSecurity(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "partial.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {
    "/widgets": {"get": {"operationId": "listWidgets", "security": [{"BearerAuth": []}]}},
    "/widgets/{id}": {"get": {"operationId": "getWidget"}}
  },
  "components": {"securitySchemes": {"BearerAuth": {"type": "http", "scheme": "bearer"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := BuildCatalogSecurityAuditReport(CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
		},
		SecurityClassifications: []catalog.SecurityClassification{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Status:     catalog.AuthStatusPresentIncomplete,
			SourceRefs: []string{"https://example.com/auth"},
			SourceNote: "reviewed example auth documentation",
		}},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	provider := findAuditProvider(t, report, "example")
	if len(provider.ArtifactSecurity) != 1 {
		t.Fatalf("artifact audit rows = %d, want 1", len(provider.ArtifactSecurity))
	}
	artifact := provider.ArtifactSecurity[0]
	if artifact.Status != SecurityAuditArtifactPartialOperationSecurity {
		t.Fatalf("artifact status = %q, want %q", artifact.Status, SecurityAuditArtifactPartialOperationSecurity)
	}
	if artifact.RootSecurityRequirementCount != 0 || artifact.OperationCount != 2 || artifact.OperationSecurityCount != 1 || artifact.OperationSecurityDeclarationCount != 1 {
		t.Fatalf("artifact security counts = %#v, want root=0 operations=2 named=1 declared=1", artifact)
	}
	const followUp = "Confirm operations without a declared security requirement are intentionally anonymous, or add reviewed operation security requirements."
	if !strings.Contains(strings.Join(artifact.ManualFollowUps, " "), followUp) {
		t.Fatalf("artifact follow-ups = %#v, want %q", artifact.ManualFollowUps, followUp)
	}
	if !strings.Contains(strings.Join(provider.ManualFollowUps, " "), followUp) {
		t.Fatalf("provider follow-ups = %#v, want partial coverage follow-up", provider.ManualFollowUps)
	}
}

func TestCatalogSecurityAuditCountsAnonymousRootSecurityAlternative(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "anonymous-root.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "security": [{}],
  "paths": {
    "/widgets": {"get": {"operationId": "listWidgets", "security": [{"BearerAuth": []}]}},
    "/widgets/{id}": {"get": {"operationId": "getWidget"}}
  },
  "components": {"securitySchemes": {"BearerAuth": {"type": "http", "scheme": "bearer"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	options := CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
		},
		SecurityClassifications: []catalog.SecurityClassification{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Status:     catalog.AuthStatusPresentIncomplete,
			SourceRefs: []string{"https://example.com/auth"},
			SourceNote: "reviewed example auth documentation",
		}},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	}
	report, err := BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	artifact := findAuditProvider(t, report, "example").ArtifactSecurity[0]
	if !artifact.RootSecurityDeclared || artifact.RootSecurityRequirementCount != 1 {
		t.Fatalf("root security declaration/count = %t/%d, want true/1 anonymous alternative", artifact.RootSecurityDeclared, artifact.RootSecurityRequirementCount)
	}
	if artifact.Status != SecurityAuditArtifactHasSecurityMetadata {
		t.Fatalf("artifact status = %q, want %q", artifact.Status, SecurityAuditArtifactHasSecurityMetadata)
	}
	if strings.Contains(strings.Join(artifact.ManualFollowUps, " "), partialOperationSecurityFollowUp) {
		t.Fatalf("artifact follow-ups spuriously report partial operation security: %#v", artifact.ManualFollowUps)
	}

	content = []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "security": [{}],
  "paths": {"/widgets": {"get": {"operationId": "listWidgets"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err = BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	artifact = findAuditProvider(t, report, "example").ArtifactSecurity[0]
	if artifact.Status != SecurityAuditArtifactHasSecurityMetadata || artifact.SecuritySchemeCount != 0 {
		t.Fatalf("scheme-less anonymous root requirement = %#v, want explicit security metadata without schemes", artifact)
	}

	content = []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "security": [],
  "paths": {
    "/widgets": {"get": {"operationId": "listWidgets", "security": [{"BearerAuth": []}]}},
    "/widgets/{id}": {"get": {"operationId": "getWidget"}}
  },
  "components": {"securitySchemes": {"BearerAuth": {"type": "http", "scheme": "bearer"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err = BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	artifact = findAuditProvider(t, report, "example").ArtifactSecurity[0]
	if !artifact.RootSecurityDeclared || artifact.RootSecurityRequirementCount != 0 || artifact.Status != SecurityAuditArtifactHasSecurityMetadata {
		t.Fatalf("explicit empty root security declaration = %#v, want declared anonymous policy without partial status", artifact)
	}
}

func TestCatalogSecurityAuditCountsAnonymousOperationRequirement(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "anonymous-operation.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {
    "/widgets": {"get": {"operationId": "listWidgets", "security": [{"BearerAuth": []}]}},
    "/widgets/{id}": {"get": {"operationId": "getWidget", "security": [{}]}},
    "/widgets/{other}": {"get": {"operationId": "getOtherWidget", "security": []}}
  },
  "components": {"securitySchemes": {"BearerAuth": {"type": "http", "scheme": "bearer"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	options := CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
		},
		SecurityClassifications: []catalog.SecurityClassification{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Status:     catalog.AuthStatusPresentIncomplete,
			SourceRefs: []string{"https://example.com/auth"},
			SourceNote: "reviewed example auth documentation",
		}},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	}
	report, err := BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	artifact := findAuditProvider(t, report, "example").ArtifactSecurity[0]
	if artifact.OperationSecurityCount != 1 || artifact.OperationSecurityDeclarationCount != 3 || artifact.Status != SecurityAuditArtifactHasSecurityMetadata {
		t.Fatalf("anonymous operation requirement audit = %#v, want 1 named/3 declared operations without partial status", artifact)
	}
	if strings.Contains(strings.Join(artifact.ManualFollowUps, " "), partialOperationSecurityFollowUp) {
		t.Fatalf("artifact follow-ups spuriously report anonymous operation as missing: %#v", artifact.ManualFollowUps)
	}

	content = []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {"/widgets": {"get": {"operationId": "listWidgets", "security": [{}]}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err = BuildCatalogSecurityAuditReport(options)
	if err != nil {
		t.Fatal(err)
	}
	artifact = findAuditProvider(t, report, "example").ArtifactSecurity[0]
	if artifact.Status != SecurityAuditArtifactHasSecurityMetadata || artifact.OperationSecurityDeclarationCount != 1 || artifact.SecuritySchemeCount != 0 {
		t.Fatalf("scheme-less anonymous operation requirement = %#v, want explicit security metadata without schemes", artifact)
	}
}

func TestCatalogSecurityAuditDetectsMissingOpenAPISecurityMetadata(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "public.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {"/status": {"get": {"operationId": "getStatus"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := BuildCatalogSecurityAuditReport(CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
		},
		SecurityClassifications: []catalog.SecurityClassification{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Status:     catalog.AuthStatusPresentIncomplete,
			SourceRefs: []string{"https://example.com/docs"},
			SourceNote: "test classification",
		}},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	row := report.Providers[0].ArtifactSecurity[0]
	if row.Status != SecurityAuditArtifactMissingSecurityMetadata {
		t.Fatalf("artifact status = %q, want %q", row.Status, SecurityAuditArtifactMissingSecurityMetadata)
	}
}

func TestCatalogSecurityAuditNamesCoveringOverlayForArtifactFinding(t *testing.T) {
	dir := t.TempDir()
	openAPIDir := filepath.Join(dir, "openapi")
	if err := os.MkdirAll(openAPIDir, 0o755); err != nil {
		t.Fatal(err)
	}
	specPath := filepath.Join("openapi", "public.json")
	content := []byte(`{
  "openapi": "3.0.3",
  "info": {"title": "Example", "version": "1.0.0"},
  "paths": {"/status": {"get": {"operationId": "getStatus"}}}
}`)
	if err := os.WriteFile(filepath.Join(dir, filepath.FromSlash(specPath)), content, 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := BuildCatalogSecurityAuditReport(CatalogSecurityAuditOptions{
		Catalog: catalog.Catalog{
			Providers: []catalog.Provider{testSecurityAuditProvider()},
			SecurityOverlays: []catalog.SecurityOverlay{{
				ID:         "example-auth-overlay",
				ProviderID: "example",
				SpecRefID:  "example-openapi",
				Status:     catalog.AuthStatusOverlayRequired,
				SecuritySchemes: []catalog.SecurityScheme{{
					Name:         "exampleBearer",
					Type:         catalog.SecuritySchemeHTTP,
					Scheme:       "bearer",
					BearerFormat: "example token",
				}},
				RootSecurity: []catalog.SecurityRequirement{{Scheme: "exampleBearer"}},
				SourceRefs:   []string{"https://example.com/openapi.json", "https://example.com/docs"},
				SourceNote:   "test overlay",
			}},
		},
		Artifacts: []catalog.CatalogSpecArtifact{{
			ProviderID: "example",
			SpecRefID:  "example-openapi",
			Kind:       "openapi",
			Path:       specPath,
		}},
		CacheDir: dir,
	})
	if err != nil {
		t.Fatal(err)
	}
	row := report.Providers[0].ArtifactSecurity[0]
	if row.Status != SecurityAuditArtifactMissingSecurityMetadata {
		t.Fatalf("artifact status = %q, want %q", row.Status, SecurityAuditArtifactMissingSecurityMetadata)
	}
	if len(row.ManualFollowUps) != 1 || !strings.Contains(row.ManualFollowUps[0], "example-auth-overlay") {
		t.Fatalf("artifact follow-ups = %#v, want covering overlay id", row.ManualFollowUps)
	}
	if strings.Contains(row.ManualFollowUps[0], "add a source-backed security overlay") {
		t.Fatalf("artifact follow-up still asks for another overlay: %#v", row.ManualFollowUps)
	}
}

func TestCoveringSecurityOverlayIDsMatchSpecURL(t *testing.T) {
	ref := catalog.SpecReference{
		ID:  "example-openapi",
		URL: "https://example.com/openapi.json",
	}
	ids := coveringSecurityOverlayIDs(ref, []catalog.SecurityOverlay{
		{
			ID:         "example-source-overlay",
			ProviderID: "example",
			SpecRefID:  "example-auth-docs",
			SourceRefs: []string{"https://example.com/openapi.json"},
		},
	})
	if len(ids) != 1 || ids[0] != "example-source-overlay" {
		t.Fatalf("covering overlay ids = %#v, want source-ref match", ids)
	}
}

func countAuditDisposition(report CatalogSecurityAuditReport, disposition string) int {
	for _, bucket := range report.Summary.Dispositions {
		if bucket.Name == disposition {
			return bucket.Count
		}
	}
	return 0
}

func findAuditProvider(t *testing.T, report CatalogSecurityAuditReport, providerID string) CatalogSecurityAuditRow {
	t.Helper()
	for _, provider := range report.Providers {
		if provider.ProviderID == providerID {
			return provider
		}
	}
	t.Fatalf("missing provider %q", providerID)
	return CatalogSecurityAuditRow{}
}

func testSecurityAuditProvider() catalog.Provider {
	return catalog.Provider{
		ID:                              "example",
		DisplayName:                     "Example",
		Category:                        "test",
		WorkflowRelevance:               "test",
		SourceHints:                     []string{"test"},
		ReviewState:                     catalog.ProviderReviewedCatalogEntry,
		CandidateID:                     "example",
		OfficialOpenAPIAvailability:     catalog.SpecAvailabilityKnown,
		OfficialMachineSpecAvailability: catalog.SpecAvailabilityUnknown,
		UserOpenAPINeed:                 catalog.UserOpenAPINeedNotExpected,
		SpecReferences: []catalog.SpecReference{{
			ID:              "example-openapi",
			Kind:            catalog.SpecKindOpenAPI,
			URL:             "https://example.com/openapi.json",
			SourceAuthority: catalog.SourceAuthorityOfficialProvider,
			SourceNote:      "test source",
			VerifiedAt:      "2026-05-21",
		}},
	}
}
