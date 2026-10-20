package catalog

import (
	"sort"
	"strings"
)

// ProviderMatchRank records the most specific kind of provider text a search
// needed in order to match. A lower rank is a closer match.
type ProviderMatchRank int

const (
	// ProviderMatchName means every query word was found in the provider ID,
	// display name or an alias.
	ProviderMatchName ProviderMatchRank = 1
	// ProviderMatchCategory means every query word was found in the provider
	// ID, display name, aliases, category or workflow relevance, and at least one
	// needed the category or workflow relevance.
	ProviderMatchCategory ProviderMatchRank = 2
	// ProviderMatchOther means at least one query word was found only in other
	// catalog text: quirks, spec reference IDs, kinds, URLs and source notes, or
	// the provider's security overlays.
	ProviderMatchOther ProviderMatchRank = 3
)

// ProviderMatch is one provider found by a catalog search with its rank.
type ProviderMatch struct {
	Provider Provider          `json:"provider"`
	Rank     ProviderMatchRank `json:"rank"`
}

// searchPunctuation is trimmed from both ends of each query word.
const searchPunctuation = ",.:;\"'()[]{}"

// SearchBuiltInProviders searches the built-in catalog. See
// Catalog.SearchProviders.
func SearchBuiltInProviders(query string, limit int) []ProviderMatch {
	return BuiltInCatalog().SearchProviders(query, limit)
}

// SearchProviders returns the providers that match a free-text query, best
// match first. The query is lowercased and split on whitespace, and
// punctuation around each word is ignored. Every word must appear as a
// substring of the provider's text; a query that is empty after that returns
// nothing.
//
// A match's rank is the lowest-priority text any word needed: the provider ID,
// display name and aliases (ProviderMatchName), then the category and workflow
// relevance (ProviderMatchCategory), then everything else (ProviderMatchOther).
// A provider whose ID, display name or alias equals the whole query, as
// FindProvider normalizes it, comes first. The rest are ordered by rank and
// then provider ID. A limit of zero or less returns every match.
//
// The search reads only the catalog. It starts no network request and reads no
// file, and it returns independent copies of the providers.
func (c Catalog) SearchProviders(query string, limit int) []ProviderMatch {
	words := searchWords(query)
	if len(words) == 0 {
		return nil
	}
	exactKey := normalizeKey(query)

	overlayText := map[string][]string{}
	for _, overlay := range c.SecurityOverlays {
		overlayText[overlay.ProviderID] = append(overlayText[overlay.ProviderID], overlaySearchText(overlay))
	}

	type candidate struct {
		match ProviderMatch
		exact bool
	}
	var found []candidate
	for _, provider := range c.Providers {
		names := strings.ToLower(strings.Join(append([]string{provider.ID, provider.DisplayName}, provider.Aliases...), " "))
		categories := strings.ToLower(provider.Category + " " + provider.WorkflowRelevance)
		others := strings.ToLower(strings.Join(append(providerOtherText(provider), overlayText[provider.ID]...), " "))

		exact := false
		for _, key := range provider.lookupKeys() {
			if key == exactKey {
				exact = true
				break
			}
		}

		rank := ProviderMatchName
		matched := true
		for _, word := range words {
			switch {
			case strings.Contains(names, word):
			case strings.Contains(categories, word):
				if rank < ProviderMatchCategory {
					rank = ProviderMatchCategory
				}
			case strings.Contains(others, word):
				rank = ProviderMatchOther
			default:
				matched = false
			}
			if !matched {
				break
			}
		}
		if !matched && !exact {
			continue
		}
		if !matched {
			rank = ProviderMatchName
		}
		found = append(found, candidate{match: ProviderMatch{Provider: cloneProvider(provider), Rank: rank}, exact: exact})
	}

	sort.SliceStable(found, func(i, j int) bool {
		a, b := found[i], found[j]
		if a.exact != b.exact {
			return a.exact
		}
		if a.match.Rank != b.match.Rank {
			return a.match.Rank < b.match.Rank
		}
		return a.match.Provider.ID < b.match.Provider.ID
	})
	if limit > 0 && len(found) > limit {
		found = found[:limit]
	}
	out := make([]ProviderMatch, len(found))
	for i, item := range found {
		out[i] = item.match
	}
	return out
}

func searchWords(query string) []string {
	var words []string
	for _, field := range strings.Fields(strings.ToLower(query)) {
		if word := strings.Trim(field, searchPunctuation); word != "" {
			words = append(words, word)
		}
	}
	return words
}

// providerOtherText lists the provider text that ranks below names and
// categories.
func providerOtherText(provider Provider) []string {
	out := append([]string(nil), provider.Quirks...)
	for _, ref := range provider.SpecReferences {
		out = append(out, ref.ID, string(ref.Kind), ref.URL, ref.SourceNote)
	}
	return out
}

func overlaySearchText(overlay SecurityOverlay) string {
	parts := []string{overlay.ID, string(overlay.Status), overlay.SourceNote}
	for _, scheme := range overlay.SecuritySchemes {
		parts = append(parts, scheme.Name, string(scheme.Type), scheme.Description, scheme.ParameterName)
	}
	return strings.Join(parts, " ")
}
