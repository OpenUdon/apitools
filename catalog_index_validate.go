package apitools

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/OpenUdon/apitools/catalog"
)

// Bound recursive typed decoding before it allocates candidate structures.
// This is a byte/depth/work budget, not a process memory ceiling.
func checkCatalogIndexJSON(ctx context.Context, data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	depth := 0
	for count := 0; ; count++ {
		if err := ctx.Err(); err != nil {
			return err
		}
		if count > 32_000_000 {
			return fmt.Errorf("catalog index exceeds structural budget")
		}
		token, err := decoder.Token()
		if err == io.EOF {
			if depth != 0 {
				return fmt.Errorf("incomplete catalog index")
			}
			return nil
		}
		if err != nil {
			return fmt.Errorf("invalid catalog index JSON")
		}
		if delimiter, ok := token.(json.Delim); ok {
			if delimiter == '{' || delimiter == '[' {
				depth++
				if depth > 100 {
					return fmt.Errorf("catalog index exceeds nesting budget")
				}
			} else {
				depth--
			}
		}
	}
}

func catalogIndexDigest(value string) bool {
	_, err := hex.DecodeString(value)
	return len(value) == 64 && err == nil && value == strings.ToLower(value)
}

func catalogIndexRelativePath(value string) bool {
	if !filepath.IsLocal(filepath.FromSlash(value)) {
		return false
	}
	for _, part := range strings.Split(filepath.FromSlash(value), string(filepath.Separator)) {
		if part == "" || part == "." || part == ".." {
			return false
		}
	}
	return true
}

// Validate derived identities and complete coverage even when a saved snapshot
// is stale. Current registration bindings apply only to a matching snapshot.
// The index is trusted installation-owned derived metadata, not a signature.
func validateCatalogIndexBody(index CatalogOperationIndex, cat catalog.Catalog, rows []catalog.CatalogSpecArtifact, bind bool) error {
	bad := func() error { return fmt.Errorf("invalid catalog operation index metadata or coverage") }
	if len(index.Artifacts) > MaxCatalogIndexArtifacts || len(index.Coverage) > MaxCatalogIndexArtifacts*4 {
		return bad()
	}
	registered := map[string]catalog.CatalogSpecArtifact{}
	for _, row := range rows {
		registered[row.ProviderID+"\x00"+row.ArtifactID] = row
	}
	used := map[string]bool{}
	expected := map[string]string{}
	coveredRefs := map[string]bool{}
	groups := map[string]bool{}
	total := 0
	for _, artifact := range index.Artifacts {
		if artifact.ArtifactID == "" || !catalogIndexDigest(artifact.SHA256) || artifact.Bytes <= 0 || !catalogIndexRelativePath(artifact.Path) || len(artifact.Links) == 0 || len(artifact.Links) > MaxCatalogIndexArtifacts || artifact.OperationCount < len(artifact.Operations) || artifact.OperationCount > maxCatalogArtifactOperations {
			return bad()
		}
		switch artifact.State {
		case "indexed":
			if artifact.OperationCount != len(artifact.Operations) || !validOperationSourceKind(artifact.Kind) {
				return bad()
			}
		case "missing", "digest_mismatch", "oversize", "parse_failure", "unsupported":
		default:
			return bad()
		}
		if artifact.State != "indexed" && artifact.Reason == "" {
			return bad()
		}
		group := artifact.ArtifactID + "\x00" + string(artifact.Kind) + "\x00" + artifact.SHA256 + fmt.Sprintf("/%d", artifact.Bytes)
		if groups[group] {
			return bad()
		}
		groups[group] = true
		pathBound := false
		for _, link := range artifact.Links {
			key := link.ProviderID + "\x00" + link.ArtifactID
			provider, ok := cat.FindProvider(link.ProviderID)
			if !ok || provider.ID != link.ProviderID || link.ArtifactID != artifact.ArtifactID || link.SpecRefID == "" || used[key] || catalogRegisteredSourceKind(link.RegisteredKind) != artifact.Kind || link.Advisory != (link.RegisteredKind == "advisory-overlay") {
				return bad()
			}
			refFound := false
			for _, ref := range provider.SpecReferences {
				if ref.ID == link.SpecRefID {
					refFound = true
					if !link.Advisory && catalogRegisteredSourceKind(string(ref.Kind)) != artifact.Kind {
						return bad()
					}
				}
			}
			if !refFound && !link.Advisory {
				return bad()
			}
			if bind {
				row, ok := registered[key]
				if !ok || row.SpecRefID != link.SpecRefID || row.Kind != link.RegisteredKind || row.SHA256 != artifact.SHA256 || row.Bytes != artifact.Bytes {
					return bad()
				}
				pathBound = pathBound || row.Path == artifact.Path
			}
			used[key] = true
			coveredRefs[link.ProviderID+"\x00"+link.SpecRefID] = true
			expected[link.ProviderID+"\x00"+link.SpecRefID+"\x00"+link.ArtifactID] = artifact.State
		}
		if bind && !pathBound {
			return bad()
		}
		selectors := map[string]bool{}
		for _, candidate := range artifact.Operations {
			if candidate.Source.Kind != artifact.Kind || candidate.Source.SHA256 != artifact.SHA256 || candidate.Source.Selector == "" || selectors[candidate.Source.Selector] || candidate.Source.DocumentPath != artifact.Path || candidate.Operation.DocumentPath != artifact.Path {
				return bad()
			}
			cleanURL, changed, err := sanitizeOperationSourceURL(candidate.Source.DocumentURL)
			if err != nil || changed || cleanURL != candidate.Source.DocumentURL {
				return bad()
			}
			selectors[candidate.Source.Selector] = true
		}
		total += len(artifact.Operations)
		if total > MaxCatalogIndexOperations {
			return bad()
		}
	}
	if bind && len(used) != len(registered) {
		return bad()
	}
	for _, provider := range cat.ListProviders() {
		for _, ref := range provider.SpecReferences {
			if !coveredRefs[provider.ID+"\x00"+ref.ID] {
				state := "missing"
				if ref.Kind == catalog.SpecKindHumanDocs {
					state = "unsupported"
				}
				expected[provider.ID+"\x00"+ref.ID+"\x00"] = state
			}
		}
	}
	if len(expected) != len(index.Coverage) {
		return bad()
	}
	for _, coverage := range index.Coverage {
		key := coverage.ProviderID + "\x00" + coverage.SpecRefID + "\x00" + coverage.ArtifactID
		state, ok := expected[key]
		if !ok || state != coverage.State || (state != "indexed" && coverage.Reason == "") {
			return bad()
		}
		delete(expected, key)
	}
	return nil
}
