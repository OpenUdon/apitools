package apitools_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/catalog"
)

// Rebuild from changed synthetic source bytes and exact registrations, rather
// than inserting ranking expectations into the derived index.
func alterDiscoveryNotes(t *testing.T, options *apitools.CatalogDiscoveryOptions, change func(map[string]any)) {
	t.Helper()
	path := filepath.Join(options.Index.Root.Directory, "openapi/notes.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var doc map[string]any
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatal(err)
	}
	change(doc)
	data, err = json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(data)
	reader := options.Index.ReadRegistrations
	options.Index.ReadRegistrations = func(ctx context.Context, root catalog.RootOptions) ([]catalog.CatalogSpecArtifact, error) {
		rows, err := reader(ctx, root)
		for i := range rows {
			rows[i].SHA256 = hex.EncodeToString(digest[:])
			rows[i].Bytes = int64(len(data))
		}
		return rows, err
	}
	index, err := apitools.BuildCatalogOperationIndex(context.Background(), options.Index)
	if err != nil {
		t.Fatal(err)
	}
	if err := apitools.WriteCatalogOperationIndex(context.Background(), options.Index, index); err != nil {
		t.Fatal(err)
	}
}

func TestCatalogDiscoveryQualifiedOutcomes(t *testing.T) {
	for _, name := range []string{"match", "one-term", "weak-overlap", "identity-only", "effect-conflict", "output-conflict", "output-unknown", "unknown-effect", "ambiguous", "ambiguous-limit", "unused-input", "missing-auth", "nullable-output", "field-loss", "actual-effect-unknown"} {
		t.Run(name, func(t *testing.T) {
			options := preparedCatalogDiscovery(t)
			want, count := apitools.CatalogDiscoveryMatch, 1
			switch name {
			case "unused-input":
				required := true
				options.Request.Contract.Inputs = map[string]apitools.ContractValue{"invoice": {Type: "string", Required: &required}}
				want, count = apitools.CatalogDiscoveryInsufficientEvidence, 0
			case "missing-auth", "nullable-output", "field-loss", "actual-effect-unknown":
				alterDiscoveryNotes(t, &options, func(doc map[string]any) {
					op := doc["paths"].(map[string]any)["/notes"].(map[string]any)["get"].(map[string]any)
					switch name {
					case "missing-auth":
						op["security"] = []any{map[string]any{"missingBearer": []any{}}}
					case "actual-effect-unknown":
						op["operationId"] = "unknownNotes"
						op["summary"] = "Notes records available"
					default:
						schema := op["responses"].(map[string]any)["200"].(map[string]any)["content"].(map[string]any)["application/json"].(map[string]any)["schema"].(map[string]any)
						properties := schema["properties"].(map[string]any)
						if name == "nullable-output" {
							properties["notes"].(map[string]any)["nullable"] = true
						} else {
							properties["api_key"] = map[string]any{"type": "string"}
						}
					}
				})
				if name == "nullable-output" {
					required := false
					options.Request.Contract.Outputs = map[string]apitools.ContractValue{"notes": {Type: "array", Required: &required}}
				}
				want, count = apitools.CatalogDiscoveryInsufficientEvidence, 0
			case "one-term":
				options.Request.Contract.Purpose = "notes"
				want, count = apitools.CatalogDiscoveryInsufficientEvidence, 0
			case "weak-overlap":
				options.Request.Contract.Purpose = "notes invoice payment"
				want, count = apitools.CatalogDiscoveryInsufficientEvidence, 0
			case "identity-only":
				alterDiscoveryNotes(t, &options, func(doc map[string]any) {
					delete(doc["paths"].(map[string]any)["/notes"].(map[string]any)["get"].(map[string]any), "summary")
				})
				want, count = apitools.CatalogDiscoveryInsufficientEvidence, 0
			case "effect-conflict":
				options.Request.Contract.Effect = apitools.OperationEffectWrite
				want, count = apitools.CatalogDiscoveryNoQualifyingAPI, 0
			case "output-conflict", "output-unknown":
				required := true
				kind := "string"
				if name == "output-unknown" {
					kind = "array"
					want = apitools.CatalogDiscoveryInsufficientEvidence
				} else {
					want = apitools.CatalogDiscoveryNoQualifyingAPI
				}
				options.Request.Contract.Outputs = map[string]apitools.ContractValue{"notes": {Type: kind, Required: &required}}
				if name == "output-unknown" {
					options.Request.Contract.Outputs["notes"] = apitools.ContractValue{Type: kind}
				}
				count = 0
			case "unknown-effect":
				options.Request.Contract.Effect = apitools.OperationEffectUnknown
			case "ambiguous", "ambiguous-limit":
				alterDiscoveryNotes(t, &options, func(doc map[string]any) {
					paths := doc["paths"].(map[string]any)
					data, _ := json.Marshal(paths["/notes"].(map[string]any)["get"])
					var op map[string]any
					_ = json.Unmarshal(data, &op)
					op["operationId"] = "listArchivedNotes"
					paths["/archive"] = map[string]any{"get": op}
				})
				want, count = apitools.CatalogDiscoveryAmbiguous, 2
				if name == "ambiguous-limit" {
					options.Request.Limits.MaxResults = 1
				}
			}
			report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
			if err != nil || report.Outcome != want || report.QualifiedOperations != count {
				t.Fatalf("%s: outcome=%s count=%d error=%v candidates=%+v", name, report.Outcome, report.QualifiedOperations, err, report.Candidates)
			}
			if name == "ambiguous" && (len(report.Ties) != 1 || len(report.Ties[0].CandidateIndexes) != 2) {
				t.Fatal("rank tie lost")
			}
			if name == "ambiguous-limit" && (len(report.Candidates) != 1 || !report.Truncated || !report.Candidates[0].Qualified) {
				t.Fatal("display limit hid ambiguity")
			}
			if name == "match" {
				workflow := t.TempDir()
				export, err := apitools.ExportCatalogArtifacts(context.Background(), apitools.CatalogArtifactExportOptions{Index: options.Index, References: report.Candidates[0].References, WorkflowDir: workflow})
				if err != nil || len(export.Artifacts) != 2 {
					t.Fatalf("selected reference roundtrip: %v %+v", err, export)
				}
			}
		})
	}
}

func TestCatalogDiscoveryPositiveMatchKeepsCoverageGaps(t *testing.T) {
	options := preparedCatalogDiscovery(t)
	options.Request.Limits.MaxOperations = 1
	report, err := apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryMatch || !report.Incomplete || report.Scope.Complete {
		t.Fatalf("positive partial scope: %v %+v", err, report)
	}
	options.Request.Contract.Effect = apitools.OperationEffectWrite
	report, err = apitools.DiscoverCatalogOperations(context.Background(), options)
	if err != nil || report.Outcome != apitools.CatalogDiscoveryInsufficientEvidence {
		t.Fatal("work limit proved absence")
	}
}
