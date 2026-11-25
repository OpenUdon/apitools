package apitools

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestReviewedProviderSearchExactAndCoveredOperations(t *testing.T) {
	client := &Client{}
	alias, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: " SLACK web api "})
	if err != nil || len(alias.Results) == 0 || alias.Results[0].Profile.ProviderID != "slack" ||
		alias.Results[0].Match != ReviewedProviderExactMatch {
		t.Fatalf("exact reviewed alias: %+v, %v", alias.Results, err)
	}

	message, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "send slack message", Limit: 1})
	if err != nil || len(message.Results) != 1 || message.Results[0].Profile.ProviderID != "slack" ||
		message.Results[0].Match != ReviewedProviderOperationMatch || len(message.Results[0].MatchedCapabilities) != 1 {
		t.Fatalf("Slack operation lookup: %+v, %v", message.Results, err)
	}
	capability := message.Results[0].MatchedCapabilities[0]
	if capability.Selector != "#/paths/~1chat.postMessage/post" || capability.ArtifactID != "slack-web-message-history-overlay" ||
		len(capability.SHA256) != 64 || message.Results[0].Profile.Coverage.Sources[0].VerifiedAt == "" {
		t.Fatalf("Slack operation provenance: %+v", capability)
	}
	if len(message.MetadataSHA256) != 64 || message.MetadataSHA256 != alias.MetadataSHA256 {
		t.Fatalf("inconsistent metadata digest: %q vs %q", message.MetadataSHA256, alias.MetadataSHA256)
	}

	for _, query := range []string{"refund a payment", "refund"} {
		refunds, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: query, Limit: 25})
		if err != nil {
			t.Fatal(err)
		}
		selected := map[string]bool{}
		for _, result := range refunds.Results {
			if result.Match == ReviewedProviderOperationMatch {
				selected[result.Profile.ProviderID] = true
			}
		}
		if !selected["stripe"] || !selected["square"] || !selected["mercado-pago"] || selected["razorpay"] {
			t.Fatalf("unexpected capability results for %q: %v", query, selected)
		}
	}

	identity, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "intuit quickbook"})
	if err != nil || len(identity.Results) == 0 || identity.Results[0].Profile.ProviderID != "quickbooks" ||
		identity.Results[0].Match != ReviewedProviderMetadataMatch || len(identity.Results[0].MatchedCapabilities) != 0 {
		t.Fatalf("stemmed identity invented operation evidence: %+v, %v", identity, err)
	}

	for _, query := range []string{"galactic zeppelin API", "(-*~) <>", "widgets \" OR x"} {
		got, err := client.SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: query})
		if err != nil || len(got.Results) != 0 {
			t.Fatalf("unexpected matches or syntax error on %q: %v, %v", query, len(got.Results), err)
		}
	}
}

func TestReviewedProviderPublisherObservationIdentity(t *testing.T) {
	source := catalog.ProviderSearchSource{SpecRefID: "example-spec", SourceURL: "https://example.test/spec", SHA256: "digest-a"}
	observed := catalogupdate.Version{
		ReferenceID: source.SpecRefID, Locator: source.SourceURL,
		Object: catalogupdate.Object{SHA256: source.SHA256},
	}
	if !reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("an exact legacy object did not match the reviewed source")
	}
	observed.Object.SHA256 = "digest-new"
	if reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("unrelated object matched")
	}
	observed.CatalogPin = source.SHA256
	if !reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("accepted successor did not match its anchor")
	}
	observed.Locator += "/other"
	if reviewedPublisherObservationMatches(source, observed) {
		t.Fatal("unrelated locator matched")
	}
}

func TestReviewedProviderIndexRefreshAndAtomicRollback(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reviewed.sqlite")
	db, err := openReviewedProviderIndex(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	profiles := []catalog.ProviderSearchProfile{{
		ProviderID: "example", DisplayName: "Example", Description: "opal warehousing", Aliases: []string{"example api"},
	}}
	first, digest, err := rankReviewedProviders(t.Context(), db, profiles, "opal warehousing")
	if err != nil || len(first) != 1 {
		t.Fatalf("first snapshot: %d, %q, %v", len(first), digest, err)
	}
	profiles[0].Description = "birch carpentry"
	updated, current, err := rankReviewedProviders(t.Context(), db, profiles, "birch carpentry")
	if err != nil || len(updated) != 1 || current == digest {
		t.Fatalf("updated snapshot: %d, %q, %v", len(updated), current, err)
	}
	if stale, _, err := rankReviewedProviders(t.Context(), db, profiles, "opal warehousing"); err != nil || len(stale) != 0 {
		t.Fatalf("stale FTS records survived: %d, %v", len(stale), err)
	}
	rejected := append(append([]catalog.ProviderSearchProfile{}, profiles...), profiles[0])
	if _, _, err := rankReviewedProviders(t.Context(), db, rejected, "birch carpentry"); err == nil {
		t.Fatal("duplicate IDs unexpectedly replaced the good snapshot")
	}
	var digestAfterFailure string
	if err := db.QueryRow("SELECT digest FROM search_metadata").Scan(&digestAfterFailure); err != nil || digestAfterFailure != current {
		t.Fatalf("failed rebuild altered the snapshot: %q, %v", digestAfterFailure, err)
	}
	if got, _, err := rankReviewedProviders(t.Context(), db, profiles, "birch carpentry"); err != nil || len(got) != 1 {
		t.Fatalf("good snapshot lost after failure: %d, %v", len(got), err)
	}
}

func TestReviewedProviderIndexRejectsForeignDatabase(t *testing.T) {
	path := filepath.Join(t.TempDir(), "foreign.sqlite")
	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.Exec("CREATE TABLE example (note TEXT); INSERT INTO example VALUES ('unchanged');"); err != nil {
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	before, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	_, err = (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{Query: "slack", IndexDB: path})
	if err == nil || !strings.Contains(err.Error(), "another application") {
		t.Fatalf("foreign database check: %v", err)
	}
	after, err := os.ReadFile(path)
	if err != nil || string(before) != string(after) {
		t.Fatalf("foreign database changed: %v", err)
	}
}

func TestReviewedProviderSearchConcurrentInitialization(t *testing.T) {
	path := filepath.Join(t.TempDir(), "reviewed.sqlite")
	missingState := filepath.Join(t.TempDir(), "missing.sqlite")
	var wg sync.WaitGroup
	for _, query := range []string{"slack web api", "refund payment"} {
		wg.Go(func() {
			report, err := (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
				Query: query, IndexDB: path, StateDB: missingState, Limit: 1,
			})
			if err != nil || len(report.Results) != 1 || report.PublisherState != "unavailable" {
				t.Errorf("concurrent lookup %q: %d, %q, %v", query, len(report.Results), report.PublisherState, err)
			}
		})
	}
	wg.Wait()
	if _, err := os.Stat(missingState); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("publisher state was unexpectedly created: %v", err)
	}
}

func TestReviewedProviderSearchReadsAcceptedStateWithoutChangingIt(t *testing.T) {
	config := newCatalogRepairFixture(t)
	provider, ok := catalog.FindBuiltInProvider("adyen")
	if !ok {
		t.Fatal("missing Adyen fixture provider")
	}
	source := catalogUpdateReferenceSources(provider)[0]
	seedCatalogUpdateFixture(t, config, source, catalogRepairFixture)
	before, err := os.ReadFile(config.StateDB)
	if err != nil {
		t.Fatal(err)
	}
	report, err := (&Client{}).SearchReviewedProviders(t.Context(), ReviewedProviderSearchOptions{
		Query: "adyen", StateDB: config.StateDB, Limit: 1,
	})
	if err != nil || report.PublisherState != "available" || len(report.Results) != 1 ||
		len(report.Results[0].PublisherChecks) != 1 {
		t.Fatalf("accepted-state lookup did not match the exact reference: %+v, %v", report, err)
	}
	check := report.Results[0].PublisherChecks[0]
	if check.SpecRefID != source.ReferenceID || check.CheckedAt.IsZero() || len(check.SHA256) != 64 {
		t.Fatalf("incorrect accepted-state evidence: %+v", check)
	}
	after, err := os.ReadFile(config.StateDB)
	if err != nil || string(before) != string(after) {
		t.Fatalf("state file changed during search: %v", err)
	}
}
