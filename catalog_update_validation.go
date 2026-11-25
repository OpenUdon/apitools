package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

type catalogUpdateOpenAPIEvidence struct {
	RawStatus       string   `json:"raw_status"`
	CorrectedStatus string   `json:"corrected_status,omitempty"`
	CorrectionNotes []string `json:"correction_notes,omitempty"`
	ReviewedRoute   string   `json:"reviewed_route,omitempty"`
	Version         string   `json:"-"`
}

func (e catalogUpdateOpenAPIEvidence) accepted() bool {
	return e.RawStatus == CatalogRefreshValidOpenAPI || e.RawStatus == CatalogRefreshValidSwagger ||
		e.CorrectedStatus != "" || e.ReviewedRoute != ""
}

// These artifact/operation pairs have source-backed qualification in the C04
// evidence or the existing NVIDIA Air fixture and catalog review. A new digest
// requires its own review, even if the same operation name is still present.
var catalogUpdateReviewedOpenAPI = map[string]struct {
	digest, method, path string
}{
	"datadog/datadog-api-v2":                    {"37702169e22c9f10aaf840d83b7e621896b17463033ae5ff06a5b91af2802189", "get", "/api/v2/synthetics/settings/on_demand_concurrency_cap"},
	"digitalocean/digitalocean-api-v2":          {"120b1f5ca7ac46159c0428c318c4b1e88598d923081ba338ebe6838f8c1a08ba", "get", "/v2/account"},
	"hugging-face-hub/hugging-face-hub-openapi": {"da5be4f2ee566534ee0c5fd74e1c380cc4575f19ca7a90c7fed4c51cc360c01d", "get", "/api/models/{namespace}/{repo}/refs"},
	"miro/miro-rest-v2-openapi":                 {"f0e17b5b43c92ff2d85c7fbe9a84301e39df1d3c14e6738cd939a6d7142f44dd", "get", "/v2/boards/{board_id}"},
	"mongodb-atlas/mongodb-atlas-admin-openapi": {"8eaeb645ded2fb91b6c9c07d29b5315ebb9a8d5ae4397f32b512457d74df34c5", "get", "/api/atlas/v2/groups"},
	"nvidia-dsx-air/nvidia-dsx-air-openapi":     {"dff59377915e43f5d9b30fe4f315c01163cf55eb243b2b28a88566fb639ebb28", "get", "/api/v3/announcements/"},
	"plaid/plaid-openapi-2020-09-14":            {"59246af3f2aef2e55f11a19af1baf6dd1aa19497c3c49a6fc3fe24c68634717f", "post", "/processor/token/permissions/get"},
	"vercel/vercel-rest-openapi":                {"ae3e778a91928ef7753675d31ca52d89dbd5fb2247e4a6e404f4f986124ea027", "get", "/v9/projects/{idOrName}/domains"},
}

func catalogUpdateNativeOpenAPI(ctx context.Context, content []byte, url string, ref catalog.RefreshableSpecReference, pin string) catalogUpdateOpenAPIEvidence {
	evidence := catalogUpdateOpenAPIEvidence{RawStatus: CatalogRefreshInvalid}
	root := catalogUpdateNativeRoot(ctx, content)
	if root == nil || ctx.Err() != nil {
		return evidence
	}
	evidence.Version = strings.TrimSpace(stringValue(mapValue(root["info"])["version"]))
	if status := catalogUpdateStrictOpenAPI(root, url); status != "" {
		evidence.RawStatus = status
		return evidence
	}
	if !containsExternalRef(root, 0) {
		if _, ok := parseableInvalidOpenAPI3Metadata(root); ok {
			evidence.RawStatus = CatalogRefreshParseableOpenAPIInvalid
		} else if _, ok := parseableInvalidSwagger2Metadata(root); ok {
			evidence.RawStatus = CatalogRefreshParseableSwaggerInvalid
		}
	}
	if pin == "" {
		return evidence
	}
	digest := sha256.Sum256(content)
	if hex.EncodeToString(digest[:]) != pin {
		return evidence
	}
	reviewed := catalogUpdateReviewedOpenAPI[ref.ProviderID+"/"+ref.SpecRefID]
	if reviewed.digest == pin && catalogRefreshStatusNeedsStrictReview(evidence.RawStatus) {
		op := mapValue(mapValue(mapValue(root["paths"])[reviewed.path])[reviewed.method])
		if len(mapValue(op["responses"])) > 0 {
			evidence.ReviewedRoute = strings.ToUpper(reviewed.method) + " " + reviewed.path
			return evidence
		}
	}
	if notes := applyCatalogRefreshCorrection(ref, root); len(notes) != 0 {
		evidence.CorrectedStatus = catalogUpdateStrictOpenAPI(root, url)
		evidence.CorrectionNotes = notes
	}
	return evidence
}

func catalogUpdateStrictOpenAPI(root map[string]any, url string) string {
	metadata, ok := specMetadataRoot(root)
	if !ok && awsOpenAPIFallbackSource(url) {
		metadata, ok = awsOpenAPIMetadataRoot(root)
	}
	if !ok {
		return ""
	}
	if metadata.Swagger != "" {
		return CatalogRefreshValidSwagger
	}
	return CatalogRefreshValidOpenAPI
}

// Validation uses the larger structural budget. Shape projection retains its
// separate semantic work budget and alias-free decoder.
func catalogUpdateNativeRoot(ctx context.Context, content []byte) map[string]any {
	const kind = "catalog-update-openapi"
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return nil
	}
	if trimmed[0] == '{' {
		if err := sourceguard.CheckJSONWithLimits(ctx, kind, content, sourceguard.DefaultLimits()); err != nil {
			return nil
		}
		decoder := json.NewDecoder(bytes.NewReader(trimmed))
		decoder.UseNumber()
		root, err := shapeJSONValue(ctx, decoder)
		if err != nil {
			return nil
		}
		if _, err := decoder.Token(); err != io.EOF {
			return nil
		}
		return mapValue(root)
	}
	node, err := sourceguard.YAMLDocumentForMetadata(ctx, kind, content)
	if err != nil {
		return nil
	}
	mergeBudget := sourceguard.MaxStructuralItems
	root, err := shapeYAMLValue(node, nil, &mergeBudget)
	if err != nil {
		return nil
	}
	return mapValue(root)
}

func (s catalogUpdateSource) refreshReference() catalog.RefreshableSpecReference {
	return catalog.RefreshableSpecReference{ProviderID: s.ProviderID, SpecRefID: s.ReferenceID, Kind: s.Kind, URL: s.URL}
}
