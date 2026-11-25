package apitools

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/googlediscovery"
)

func TestGoogleSlidesCatalogScopeEvidence(t *testing.T) {
	const (
		readScope  = "https://www.googleapis.com/auth/presentations.readonly"
		writeScope = "https://www.googleapis.com/auth/presentations"
	)
	methods := []struct {
		id, method, path, scope string
	}{
		{"slides.presentations.get", "GET", "/v1/presentations/{+presentationId}", readScope},
		{"slides.presentations.pages.get", "GET", "/v1/presentations/{presentationId}/pages/{pageObjectId}", readScope},
		{"slides.presentations.pages.getThumbnail", "GET", "/v1/presentations/{presentationId}/pages/{pageObjectId}/thumbnail", readScope},
		{"slides.presentations.create", "POST", "/v1/presentations", writeScope},
		{"slides.presentations.batchUpdate", "POST", "/v1/presentations/{presentationId}:batchUpdate", writeScope},
	}
	t.Run("advisory", func(t *testing.T) {
		overlays := catalog.SecurityOverlaysForProvider("google-slides")
		if len(overlays) != 1 {
			t.Fatalf("Slides security overlays = %d, want 1", len(overlays))
		}
		overlay := overlays[0]
		if overlay.SpecRefID != "slides-discovery-v1" || overlay.Status != catalog.AuthStatusOverlayRequired ||
			len(overlay.RootSecurity) != 0 || len(overlay.RootSecuritySets) != 0 || len(overlay.OperationSecurity) != len(methods) {
			t.Fatalf("Slides overlay must apply only to the reviewed methods: %+v", overlay)
		}
		for _, expected := range methods {
			index := slices.IndexFunc(overlay.OperationSecurity, func(operation catalog.OperationSecurity) bool {
				return operation.Match.Method == expected.method && operation.Match.Path == expected.path
			})
			if index < 0 {
				t.Errorf("%s is absent from the overlay", expected.id)
				continue
			}
			rule := overlay.OperationSecurity[index]
			if rule.Match.OperationID != "" || len(rule.Match.Tags) != 0 || len(rule.SecuritySets) != 0 ||
				len(rule.Security) != 1 || rule.Security[0].Scheme != "googleOAuth2" ||
				!slices.Equal(rule.Security[0].Scopes, []string{expected.scope}) {
				t.Errorf("%s has the wrong advisory match or scope: %+v", expected.id, rule)
			}
		}
	})
	t.Run("retained Discovery", func(t *testing.T) {
		dir := os.Getenv("APITOOLS_C13_SLIDES_DISCOVERY_DIR")
		if dir == "" {
			t.Skip("set APITOOLS_C13_SLIDES_DISCOVERY_DIR to the two retained SDK Discovery copies")
		}
		for _, fixture := range []struct{ file, digest string }{
			{"v0.221.0.json", "c56a87d562396319c1af8a2a2068ed50a31fe1e34437dd3ba593a57990d268d9"},
			{"v0.259.0.json", "f45d36026cc10cbb790a83405dc774c5815d9287604bdfd335475fd1cae0b022"},
		} {
			t.Run(fixture.file, func(t *testing.T) {
				data, err := readLocalSpecFile(filepath.Join(dir, fixture.file), DefaultMaxBytes)
				if err != nil {
					t.Fatal(err)
				}
				if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != fixture.digest {
					t.Fatalf("Discovery digest = %s, want %s", got, fixture.digest)
				}
				model, err := googlediscovery.Parse(data)
				if err != nil {
					t.Fatal(err)
				}
				for _, expected := range methods {
					index := slices.IndexFunc(model.Operations, func(operation *googlediscovery.Operation) bool {
						return operation.ID == expected.id
					})
					if index < 0 {
						t.Errorf("Discovery omits %s", expected.id)
						continue
					}
					native := model.Operations[index]
					if native.HTTPMethod != expected.method || native.Path != expected.path ||
						!slices.Contains(native.Scopes, expected.scope) {
						t.Errorf("%s disagrees with the reviewed method, path or scope: %s %s %v",
							expected.id, native.HTTPMethod, native.Path, native.Scopes)
					}
					if expected.scope == writeScope && slices.Contains(native.Scopes, readScope) {
						t.Errorf("%s explicitly lists the Slides read-only scope", expected.id)
					}
				}
			})
		}
	})
}
