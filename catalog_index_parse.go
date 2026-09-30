package apitools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/artifactio"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

const (
	maxCatalogArtifactBytes           = 128 << 20
	maxCatalogArtifactStructuralItems = 4_000_000
	maxCatalogArtifactOperations      = 100_000
	maxCatalogArtifactMetadataBytes   = 256 << 20
	catalogArtifactDeadline           = 45 * time.Second
)

// parseRegisteredCatalogArtifact is private to the catalog index path. Its
// registration comes from the caller's explicit catalog registry snapshot;
// neither discovery requests nor existing direct parser APIs can select it.
// Non-OpenAPI family parsers retain their original 20-MiB limits.
func parseRegisteredCatalogArtifact(ctx context.Context, root string, registration catalog.CatalogSpecArtifact) (operationSourceAdapterResult, error) {
	ctx, cancel := context.WithTimeout(ctx, catalogArtifactDeadline)
	defer cancel()
	if err := ctx.Err(); err != nil {
		return operationSourceAdapterResult{}, err
	}
	if registration.ProviderID == "" || registration.SpecRefID == "" || registration.ArtifactID == "" || registration.SHA256 == "" || registration.Bytes <= 0 {
		return operationSourceAdapterResult{}, fmt.Errorf("catalog artifact requires registration identity and integrity evidence")
	}
	file, err := artifactio.ReadFile(root, registration.Path, artifactio.ReadOptions{MaxBytes: maxCatalogArtifactBytes, SHA256: registration.SHA256, Bytes: registration.Bytes})
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	if err := ctx.Err(); err != nil {
		return operationSourceAdapterResult{}, err
	}
	url, _, err := sanitizeOperationSourceURL(registration.SourceURL)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	source := loadedOperationSource{input: OperationSourceInput{Kind: OperationSourceKind(registration.Kind), Name: registration.ArtifactID, Path: registration.Path, URL: url}, content: file.Data, digest: file.SHA256}
	budget := DefaultPromptBudget()
	budget.MaxOperations = maxCatalogArtifactOperations
	budget.MaxContextBytes = maxCatalogArtifactMetadataBytes
	if err := validateOperationSourceIdentity(source.input, budget); err != nil {
		return operationSourceAdapterResult{}, err
	}
	if source.input.Kind != OperationSourceOpenAPI {
		if len(file.Data) > DefaultMaxBytes {
			return operationSourceAdapterResult{}, fmt.Errorf("large catalog artifacts require the OpenAPI index adapter")
		}
		return adaptOperationSource(ctx, source, budget)
	}
	rootValue, err := parseCatalogOpenAPI(ctx, file.Data)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	doc := InventoryDocument{Name: registration.ArtifactID, Path: registration.Path, URL: url}
	var inventory OperationInventory
	_, err = addDocumentInventoryBounded(ctx, &inventory, doc, 0, rootValue, "", budget.MaxOperations, maxCatalogArtifactMetadataBytes)
	sanitizeInventory(&inventory, budget)
	result, parseErr := openAPICandidatesFromInventory(ctx, source, budget, inventory, err, maxCatalogArtifactMetadataBytes)
	if err := ctx.Err(); err != nil {
		return operationSourceAdapterResult{}, err
	}
	// Count bytes incrementally, avoiding a second whole-result allocation.
	var total int
	for _, candidate := range result.candidates {
		if err := ctx.Err(); err != nil {
			return operationSourceAdapterResult{}, err
		}
		encoded, err := json.Marshal(candidate)
		if err != nil {
			return operationSourceAdapterResult{}, err
		}
		total += len(encoded)
		if total > maxCatalogArtifactMetadataBytes {
			return operationSourceAdapterResult{}, fmt.Errorf("catalog artifact metadata exceeds byte budget")
		}
	}
	return result, parseErr
}

func parseCatalogOpenAPI(ctx context.Context, content []byte) (map[string]any, error) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("catalog artifact is empty")
	}
	limits := sourceguard.Limits{MaxDocumentBytes: maxCatalogArtifactBytes, MaxStructuralItems: maxCatalogArtifactStructuralItems, MaxNestingDepth: sourceguard.MaxNestingDepth}
	var value any
	if trimmed[0] == '{' {
		if err := sourceguard.CheckJSONWithLimits(ctx, "catalog openapi", trimmed, limits); err != nil {
			return nil, err
		}
		decoder := json.NewDecoder(bytes.NewReader(trimmed))
		decoder.UseNumber()
		if err := decoder.Decode(&value); err != nil {
			return nil, err
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			return nil, fmt.Errorf("catalog artifact must contain exactly one JSON document")
		}
	} else {
		document, err := sourceguard.YAMLDocument(ctx, "catalog openapi", trimmed, limits)
		if err != nil {
			return nil, err
		}
		if err := document.Decode(&value); err != nil {
			return nil, err
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	root, ok := normalizeInventoryValue(value).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("catalog artifact root must be an object")
	}
	openapi := stringValue(root["openapi"])
	swagger := stringValue(root["swagger"])
	if swagger == "" && root["swagger"] == float64(2) {
		swagger = "2.0"
		root["swagger"] = swagger
	}
	if !strings.HasPrefix(openapi, "3.0.") && !strings.HasPrefix(openapi, "3.1.") && swagger != "2.0" {
		return nil, fmt.Errorf("unsupported catalog OpenAPI version")
	}
	return root, ctx.Err()
}
