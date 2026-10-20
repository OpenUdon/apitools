package catalog

import (
	"reflect"
	"testing"
)

func matchIDs(matches []ProviderMatch) []string {
	ids := make([]string, 0, len(matches))
	for _, match := range matches {
		ids = append(ids, match.Provider.ID)
	}
	return ids
}

func findMatch(matches []ProviderMatch, id string) (ProviderMatch, bool) {
	for _, match := range matches {
		if match.Provider.ID == id {
			return match, true
		}
	}
	return ProviderMatch{}, false
}

func TestSearchBuiltInProvidersRanksTheProviderItselfFirst(t *testing.T) {
	for _, name := range []string{"github", "stripe", "openweathermap", "cloudflare", "slack"} {
		matches := SearchBuiltInProviders(name, 0)
		if len(matches) == 0 {
			t.Fatalf("SearchBuiltInProviders(%q) returned nothing", name)
		}
		if matches[0].Provider.ID != name || matches[0].Rank != ProviderMatchName {
			t.Fatalf("SearchBuiltInProviders(%q)[0] = %s rank %d, want the provider itself at rank %d",
				name, matches[0].Provider.ID, matches[0].Rank, ProviderMatchName)
		}
	}
}

func TestSearchBuiltInProvidersPutsURLOnlyMatchesLast(t *testing.T) {
	matches := SearchBuiltInProviders("github", 0)
	// Many providers only mention github.com in a spec URL or note.
	if len(matches) < 20 {
		t.Fatalf("expected many providers to mention github, got %d", len(matches))
	}
	seenOther := false
	for i, match := range matches {
		if i > 0 && match.Rank < matches[i-1].Rank {
			t.Fatalf("rank order broken at %d: %s rank %d after rank %d", i, match.Provider.ID, match.Rank, matches[i-1].Rank)
		}
		if match.Rank == ProviderMatchOther {
			seenOther = true
		}
	}
	if !seenOther {
		t.Fatal("expected at least one provider matched only through spec or overlay text")
	}
	if got := matches[0].Provider.ID; got != "github" {
		t.Fatalf("first match = %q, want github", got)
	}
}

func TestSearchBuiltInProvidersRequiresEveryWord(t *testing.T) {
	both := SearchBuiltInProviders("cloudflare dns", 0)
	if len(both) == 0 {
		t.Fatal("expected cloudflare dns to match cloudflare")
	}
	only := SearchBuiltInProviders("cloudflare", 0)
	if len(both) > len(only) {
		t.Fatalf("two words matched %d providers, one word matched %d", len(both), len(only))
	}
	for _, match := range both {
		single, ok := findMatch(only, match.Provider.ID)
		if !ok {
			t.Fatalf("%s matched \"cloudflare dns\" but not \"cloudflare\"", match.Provider.ID)
		}
		if match.Rank < single.Rank {
			t.Fatalf("%s: two-word rank %d is better than one-word rank %d", match.Provider.ID, match.Rank, single.Rank)
		}
	}
	if got := SearchBuiltInProviders("cloudflare zzzznotaword", 0); len(got) != 0 {
		t.Fatalf("a word nothing contains still matched: %v", matchIDs(got))
	}
}

func TestSearchProvidersIgnoresCaseAndSurroundingPunctuation(t *testing.T) {
	want := matchIDs(SearchBuiltInProviders("cloudflare dns", 0))
	for _, query := range []string{"Cloudflare DNS", "  cloudflare   dns  ", "(cloudflare), \"dns\".", "CLOUDFLARE DNS."} {
		got := matchIDs(SearchBuiltInProviders(query, 0))
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("SearchBuiltInProviders(%q) = %v, want %v", query, got, want)
		}
	}
}

func TestSearchProvidersMatchesDisplayNameAndAliases(t *testing.T) {
	provider, ok := FindBuiltInProvider("agile-crm")
	if !ok {
		t.Fatal("agile-crm missing from the built-in catalog")
	}
	for _, query := range []string{provider.DisplayName, "agile crm", "agile-crm"} {
		matches := SearchBuiltInProviders(query, 0)
		if len(matches) == 0 || matches[0].Provider.ID != "agile-crm" || matches[0].Rank != ProviderMatchName {
			t.Fatalf("SearchBuiltInProviders(%q) did not put agile-crm first at rank 1: %v", query, matchIDs(matches))
		}
	}
}

func TestSearchProvidersLimit(t *testing.T) {
	all := SearchBuiltInProviders("github", 0)
	if len(all) < 5 {
		t.Fatalf("expected more than 5 matches, got %d", len(all))
	}
	for _, limit := range []int{1, 5} {
		got := SearchBuiltInProviders("github", limit)
		if !reflect.DeepEqual(matchIDs(got), matchIDs(all[:limit])) {
			t.Fatalf("limit %d = %v, want the first %d of %v", limit, matchIDs(got), limit, matchIDs(all))
		}
	}
	if got := SearchBuiltInProviders("github", -1); len(got) != len(all) {
		t.Fatalf("negative limit returned %d, want all %d", len(got), len(all))
	}
}

func TestSearchProvidersEmptyQueryReturnsNothing(t *testing.T) {
	for _, query := range []string{"", "   ", "()", ", ."} {
		if got := SearchBuiltInProviders(query, 0); len(got) != 0 {
			t.Fatalf("SearchBuiltInProviders(%q) = %v, want none", query, matchIDs(got))
		}
	}
}

func TestSearchProvidersIsDeterministic(t *testing.T) {
	first := SearchBuiltInProviders("api", 0)
	for i := 0; i < 3; i++ {
		if !reflect.DeepEqual(first, SearchBuiltInProviders("api", 0)) {
			t.Fatal("two searches returned different results")
		}
	}
}

func TestSearchProvidersRanksByTheLowestPriorityFieldNeeded(t *testing.T) {
	cat := Catalog{
		Providers: []Provider{
			{ID: "alpha", DisplayName: "Alpha", Aliases: []string{"first-letter"}, Category: "weather", WorkflowRelevance: "forecast lookups"},
			{ID: "beta", DisplayName: "Beta", Category: "payments", Quirks: []string{"mentions weather in a note"}},
			{ID: "weather-hub", DisplayName: "Weather Hub"},
			{ID: "gamma", DisplayName: "Gamma", SpecReferences: []SpecReference{{ID: "gamma-docs", Kind: SpecKindHumanDocs, URL: "https://example.test/weather/docs", SourceNote: "n"}}},
			{ID: "delta", DisplayName: "Delta"},
		},
		SecurityOverlays: []SecurityOverlay{
			{ID: "delta-auth", ProviderID: "delta", SourceNote: "uses a weather api key"},
		},
	}
	matches := cat.SearchProviders("weather", 0)
	want := []struct {
		id   string
		rank ProviderMatchRank
	}{
		{"weather-hub", ProviderMatchName},
		{"alpha", ProviderMatchCategory},
		{"beta", ProviderMatchOther},
		{"delta", ProviderMatchOther},
		{"gamma", ProviderMatchOther},
	}
	if len(matches) != len(want) {
		t.Fatalf("got %v, want %d matches", matchIDs(matches), len(want))
	}
	for i, w := range want {
		if matches[i].Provider.ID != w.id || matches[i].Rank != w.rank {
			t.Fatalf("match %d = %s rank %d, want %s rank %d (all: %v)", i, matches[i].Provider.ID, matches[i].Rank, w.id, w.rank, matchIDs(matches))
		}
	}
}

func TestSearchProvidersMixesRanksAcrossWords(t *testing.T) {
	cat := Catalog{Providers: []Provider{
		{ID: "acme", DisplayName: "Acme", Category: "billing", WorkflowRelevance: "invoices"},
		{ID: "acme-two", DisplayName: "Acme Two", Quirks: []string{"invoices note"}},
	}}
	matches := cat.SearchProviders("acme invoices", 0)
	if len(matches) != 2 {
		t.Fatalf("got %v, want both providers", matchIDs(matches))
	}
	if matches[0].Provider.ID != "acme" || matches[0].Rank != ProviderMatchCategory {
		t.Fatalf("first = %s rank %d, want acme at rank %d", matches[0].Provider.ID, matches[0].Rank, ProviderMatchCategory)
	}
	if matches[1].Provider.ID != "acme-two" || matches[1].Rank != ProviderMatchOther {
		t.Fatalf("second = %s rank %d, want acme-two at rank %d", matches[1].Provider.ID, matches[1].Rank, ProviderMatchOther)
	}
}

func TestSearchProvidersPutsAnExactProviderKeyFirst(t *testing.T) {
	cat := Catalog{Providers: []Provider{
		{ID: "aaa-box-tools", DisplayName: "Box Tools"},
		{ID: "box", DisplayName: "Box"},
		{ID: "aaa-sandbox", DisplayName: "Sandbox"},
	}}
	got := matchIDs(cat.SearchProviders("box", 0))
	want := []string{"box", "aaa-box-tools", "aaa-sandbox"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	// A multi-word query equals the normalized key of the display name.
	cat = Catalog{Providers: []Provider{{ID: "ab-cd", DisplayName: "AB CD"}}}
	if got := cat.SearchProviders("ab cd", 0); len(got) != 1 || got[0].Rank != ProviderMatchName {
		t.Fatalf("got %v", got)
	}
}

func TestSearchProvidersReturnsCopies(t *testing.T) {
	matches := SearchBuiltInProviders("stripe", 1)
	if len(matches) != 1 {
		t.Fatalf("got %d matches", len(matches))
	}
	original, _ := FindBuiltInProvider("stripe")
	if len(matches[0].Provider.SpecReferences) == 0 {
		t.Fatal("stripe has no spec references to mutate")
	}
	matches[0].Provider.DisplayName = "changed"
	matches[0].Provider.SpecReferences[0].URL = "https://changed.test/"
	matches[0].Provider.Quirks = append(matches[0].Provider.Quirks, "added")
	after, _ := FindBuiltInProvider("stripe")
	if !reflect.DeepEqual(original, after) {
		t.Fatal("mutating a search result changed the built-in catalog")
	}
	again := SearchBuiltInProviders("stripe", 1)
	if again[0].Provider.DisplayName == "changed" {
		t.Fatal("mutating a search result changed a later search")
	}

	cat := Catalog{Providers: []Provider{{ID: "p", DisplayName: "P", Aliases: []string{"pp"}}}}
	got := cat.SearchProviders("p", 0)
	got[0].Provider.Aliases[0] = "changed"
	if cat.Providers[0].Aliases[0] != "pp" {
		t.Fatal("mutating a result alias changed the caller's catalog")
	}
}

func TestSearchProvidersOnAnEmptyCatalog(t *testing.T) {
	if got := (Catalog{}).SearchProviders("github", 0); len(got) != 0 {
		t.Fatalf("empty catalog returned %v", matchIDs(got))
	}
}
