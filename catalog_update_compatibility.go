package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
	"github.com/OpenUdon/uws/binding"
)

func (c *Client) proveCompatibleCatalogVersion(ctx context.Context, store *catalogupdate.Store,
	held, candidate catalogupdate.Version, kind OperationSourceKind, reviewedKeys []string,
) (catalogUpdateComparison, error) {
	if held.ProviderID != candidate.ProviderID || held.ReferenceID != candidate.ReferenceID ||
		held.Service == "" || held.Service != candidate.Service ||
		held.Locator != candidate.Locator || len(reviewedKeys) == 0 {
		return catalogUpdateUnexamined, nil
	}
	if _, err := c.versionURLSyntax(held.Locator); err != nil {
		return catalogUpdateUnexamined, err
	}
	scope := []APIVersionOfficialScope{scopeForVersionURL(held.Locator)}
	if held.Locator == "https://api.hubapi.com/public/api/spec/v1/specs" {
		scope[0].PathPrefix = "/public/api/spec/"
	}
	for _, raw := range []string{held.URL, candidate.URL} {
		if err := c.versionScopeAllows(raw, scope); err != nil {
			return catalogUpdateUnexamined, err
		}
	}
	if !sameCatalogUpdateResource(held.Locator, held.URL, candidate.URL) {
		return catalogUpdateUnexamined, nil
	}
	if held.Object.Bytes > DefaultMaxBytes || candidate.Object.Bytes > DefaultMaxBytes {
		return catalogUpdateUnexamined, nil
	}
	read := func(v catalogupdate.Version) (OperationSourceInput, error) {
		f, err := store.OpenObject(ctx, v.Object)
		if err != nil {
			return OperationSourceInput{}, err
		}
		defer f.Close()
		body, err := io.ReadAll(io.LimitReader(f, v.Object.Bytes+1))
		if err != nil {
			return OperationSourceInput{}, err
		}
		digest := sha256.Sum256(body)
		if int64(len(body)) != v.Object.Bytes || hex.EncodeToString(digest[:]) != v.Object.SHA256 {
			return OperationSourceInput{}, fmt.Errorf("%w: catalog version bytes changed during comparison", catalogupdate.ErrObjectIntegrity)
		}
		return OperationSourceInput{Kind: kind, URL: v.URL, Content: body}, nil
	}
	before, err := read(held)
	if err != nil {
		return catalogUpdateUnexamined, err
	}
	after, err := read(candidate)
	if err != nil {
		return catalogUpdateUnexamined, err
	}
	return compareCatalogUpdate(ctx, before, after, reviewedKeys), nil
}

func sameCatalogUpdateResource(locator, left, right string) bool {
	aURL, aErr := url.Parse(left)
	bURL, bErr := url.Parse(right)
	if aErr != nil || bErr != nil || versionOrigin(aURL) != versionOrigin(bURL) ||
		aURL.RawQuery != bURL.RawQuery || aURL.Fragment != bURL.Fragment {
		return false
	}
	if locator == "https://api.hubapi.com/public/api/spec/v1/specs" {
		a, aVersion, aOK := strings.Cut(aURL.EscapedPath(), "/version/")
		b, bVersion, bOK := strings.Cut(bURL.EscapedPath(), "/version/")
		return aOK && bOK && a == b && validCatalogUpdateIndexVersion(aVersion) && validCatalogUpdateIndexVersion(bVersion)
	}
	a, _ := catalogUpdateReferenceIdentity(catalog.SpecReference{URL: left})
	b, _ := catalogUpdateReferenceIdentity(catalog.SpecReference{URL: right})
	return a != "" && a == b
}

func validCatalogUpdateIndexVersion(version string) bool {
	if _, err := time.Parse("2006-01", version); err == nil {
		return true
	}
	_, ok := numericAPIVersion(version)
	return ok
}

type catalogUpdateComparison string

const (
	catalogUpdateCompatible catalogUpdateComparison = "compatible"
	catalogUpdateChanged    catalogUpdateComparison = "contract-changed"
	catalogUpdateUnexamined catalogUpdateComparison = "compatibility-unexamined"
)

// compareCatalogUpdate validates the chosen local sources and compares the
// previously reviewed operation shapes. An unproven projection is kept for
// review. Exact contract equality is a sufficient compatibility proof.
func compareCatalogUpdate(ctx context.Context, baseline, candidate OperationSourceInput, reviewedKeys []string) catalogUpdateComparison {
	if baseline.Kind != OperationSourceOpenAPI || baseline.Kind != candidate.Kind ||
		len(reviewedKeys) == 0 || baseline.Content == nil || candidate.Content == nil {
		return catalogUpdateUnexamined
	}
	before, err := BuildOperationShapeTable(ctx, OperationShapeOptions{
		Sources: []ShapeSourceInput{{ID: "catalog-update", OperationSourceInput: baseline}},
	})
	if err != nil {
		return catalogUpdateUnexamined
	}
	after, err := BuildOperationShapeTable(ctx, OperationShapeOptions{
		Sources: []ShapeSourceInput{{ID: "catalog-update", OperationSourceInput: candidate}},
	})
	if err != nil {
		return catalogUpdateUnexamined
	}
	if comparison := compareReviewedCatalogShapes(before.Operations, after.Operations, reviewedKeys); comparison != catalogUpdateCompatible {
		return comparison
	}
	// Native contract comparison catches information intentionally omitted
	// by shape projection, including the HTTP scheme, OAuth endpoints,
	// response status codes and media types.
	for _, source := range []OperationSourceInput{baseline, candidate} {
		if _, ok := downloadedSpecMetadata(ctx, source.Content, source.URL); !ok {
			return catalogUpdateUnexamined
		}
	}
	left, err := decodeShapeDocument(ctx, baseline.Content)
	if err != nil {
		return catalogUpdateUnexamined
	}
	right, err := decodeShapeDocument(ctx, candidate.Content)
	if err != nil {
		return catalogUpdateUnexamined
	}
	for _, root := range []map[string]any{left, right} {
		info := mapValue(root["info"])
		if !catalogUpdateReferencesInfo(root) {
			for _, key := range []string{"title", "version", "summary", "description"} {
				delete(info, key)
			}
		}
	}
	if !removeUnrelatedCatalogUpdatePaths(ctx, mapValue(left["paths"]), mapValue(right["paths"])) {
		return catalogUpdateUnexamined
	}
	if !reflect.DeepEqual(left, right) {
		return catalogUpdateChanged
	}
	return catalogUpdateCompatible
}

func catalogUpdateReferencesInfo(value any) bool {
	switch node := value.(type) {
	case map[string]any:
		if ref, ok := node["$ref"].(string); ok {
			decoded, _ := url.PathUnescape(ref)
			if decoded == "#/info" || strings.HasPrefix(decoded, "#/info/") {
				return true
			}
		}
		for _, child := range node {
			if catalogUpdateReferencesInfo(child) {
				return true
			}
		}
	case []any:
		for _, child := range node {
			if catalogUpdateReferencesInfo(child) {
				return true
			}
		}
	}
	return false
}

// Ignore only a new literal path whose routing is distinct from every old
// path. Retain changes to native definitions and all existing paths as
// review items, including potentially overlapping parameterized additions.
func removeUnrelatedCatalogUpdatePaths(ctx context.Context, before, after map[string]any) bool {
	const workLimit = 100_000
	work := 0
	previousPaths := make([][]string, 0, len(before))
	for path := range before {
		work += len(path) + 1
		if work > workLimit || ctx.Err() != nil {
			return false
		}
		if strings.HasPrefix(path, "/") {
			previousPaths = append(previousPaths, strings.Split(path, "/"))
		}
	}
	for path := range after {
		work += len(path) + 1
		if work > workLimit || ctx.Err() != nil {
			return false
		}
		if _, held := before[path]; held || !strings.HasPrefix(path, "/") || strings.Contains(path, "{") {
			continue
		}
		pathParts := strings.Split(path, "/")
		separate := true
		for _, previousParts := range previousPaths {
			work += len(previousParts) + 1
			if work > workLimit || ctx.Err() != nil {
				return false
			}
			if len(previousParts) != len(pathParts) {
				continue
			}
			distinct := false
			for i, part := range previousParts {
				if !strings.Contains(part, "{") && part != pathParts[i] {
					distinct = true
					break
				}
			}
			if !distinct {
				separate = false
				break
			}
		}
		if separate {
			delete(after, path)
		}
	}
	return true
}

func compareReviewedCatalogShapes(before, after []binding.OperationShape, reviewedKeys []string) catalogUpdateComparison {
	left := indexUpdateShapes(before)
	right := indexUpdateShapes(after)
	for _, key := range reviewedKeys {
		original, ok := left[key]
		if !ok || !knownCatalogUpdateShape(original) {
			return catalogUpdateUnexamined
		}
		changed, ok := right[key]
		if !ok {
			return catalogUpdateChanged
		}
		if !knownCatalogUpdateShape(changed) {
			return catalogUpdateUnexamined
		}
		original.Source = binding.Source{}
		changed.Source = binding.Source{}
		if !reflect.DeepEqual(original, changed) {
			return catalogUpdateChanged
		}
	}
	return catalogUpdateCompatible
}

func indexUpdateShapes(shapes []binding.OperationShape) map[string]binding.OperationShape {
	result := make(map[string]binding.OperationShape, len(shapes))
	aliases := make(map[binding.Selector][]string)
	for _, shape := range shapes {
		if _, duplicate := result[shape.Selector.Key]; duplicate {
			result[shape.Selector.Key] = binding.OperationShape{}
			continue
		}
		result[shape.Selector.Key] = shape
		for _, alias := range shape.Aliases {
			selector := binding.Selector{Kind: alias.Kind, Value: alias.Value}
			aliases[selector] = append(aliases[selector], shape.Selector.Key)
		}
	}
	for _, keys := range aliases {
		if len(keys) > 1 {
			for _, key := range keys {
				result[key] = binding.OperationShape{}
			}
		}
	}
	return result
}

func knownCatalogUpdateShape(shape binding.OperationShape) bool {
	if !shape.Complete || !shape.Security.Known {
		return false
	}
	for _, input := range shape.Inputs {
		if !input.Schema.Known {
			return false
		}
	}
	for _, output := range shape.Outputs {
		if !output.Schema.Known {
			return false
		}
	}
	return true
}
