package apitools

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
	_ "modernc.org/sqlite"
)

const catalogProviderSearchSchema = "apitools.reviewed-provider-search.v1"

// ReviewedProviderSearchOptions selects a lookup over the reviewed
// catalog. An empty IndexDB uses disposable SQLite. StateDB optionally supplies
// accepted publisher checks from an existing maintenance database.
type ReviewedProviderSearchOptions struct {
	Query      string
	Limit      int
	IndexDB    string
	StateDB    string
	Offline    bool
	Embeddings *ReviewedEmbeddingConfig
}

// ReviewedProviderMatchKind describes the evidence that contributed to a result.
type ReviewedProviderMatchKind string

const (
	ReviewedProviderExactMatch     ReviewedProviderMatchKind = "exact-provider"
	ReviewedProviderOperationMatch ReviewedProviderMatchKind = "covered-operation"
	ReviewedProviderSemanticMatch  ReviewedProviderMatchKind = "semantic-provider"
	ReviewedProviderMetadataMatch  ReviewedProviderMatchKind = "provider-metadata"
)

// ReviewedCapabilityMatch records the exact reviewed operation behind a match.
type ReviewedCapabilityMatch struct {
	SpecRefID  string `json:"spec_ref_id"`
	ArtifactID string `json:"artifact_id,omitempty"`
	SHA256     string `json:"sha256"`
	catalog.ProviderSearchOperation
}

// ReviewedPublisherObservation reports an accepted publisher artifact check
// for the linked reference. CheckedAt is an observation time.
type ReviewedPublisherObservation struct {
	SpecRefID string    `json:"spec_ref_id"`
	Service   string    `json:"service"`
	Version   string    `json:"version,omitempty"`
	Revision  string    `json:"revision,omitempty"`
	SHA256    string    `json:"sha256"`
	CheckedAt time.Time `json:"checked_at"`
}

// ReviewedProviderSearchResult separates provider text from operation evidence.
type ReviewedProviderSearchResult struct {
	Profile             catalog.ProviderSearchProfile  `json:"profile"`
	Match               ReviewedProviderMatchKind      `json:"match"`
	MatchedCapabilities []ReviewedCapabilityMatch      `json:"matched_capabilities,omitempty"`
	PublisherChecks     []ReviewedPublisherObservation `json:"publisher_checks,omitempty"`
}

// ReviewedProviderSearchReport is distinct from the public directory response.
type ReviewedProviderSearchReport struct {
	Schema         string                         `json:"schema"`
	Query          string                         `json:"query"`
	Source         string                         `json:"source"`
	MetadataSHA256 string                         `json:"metadata_sha256"`
	PublisherState string                         `json:"publisher_state"`
	Embeddings     *ReviewedEmbeddingUsage        `json:"embeddings,omitempty"`
	Results        []ReviewedProviderSearchResult `json:"results"`
}

// SearchReviewedProviders queries the reviewed provider catalog. Source
// references and publisher observations are read as metadata.
func (c *Client) SearchReviewedProviders(ctx context.Context, options ReviewedProviderSearchOptions) (ReviewedProviderSearchReport, error) {
	query := strings.TrimSpace(options.Query)
	report := ReviewedProviderSearchReport{
		Schema: catalogProviderSearchSchema, Source: "reviewed-catalog", Query: query,
		PublisherState: "unconfigured", Results: []ReviewedProviderSearchResult{},
	}
	if len(query) < 2 {
		return report, fmt.Errorf("query must be at least two characters")
	}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	profiles := catalog.BuiltInProviderSearchProfiles()
	db, err := openReviewedProviderIndex(options.IndexDB)
	if err != nil {
		return report, err
	}
	defer db.Close()
	ranked, digest, err := rankReviewedProviders(ctx, db, profiles, query)
	if err != nil {
		return report, err
	}
	report.MetadataSHA256 = digest
	if options.Embeddings != nil {
		ranked, report.Embeddings = c.rankReviewedEmbeddings(ctx, db, options, profiles, ranked, query)
		if err := ctx.Err(); err != nil {
			return report, err
		}
	}
	limit := options.Limit
	if limit <= 0 {
		limit = 10
	}
	observations, state := readReviewedPublisherObservations(ctx, options.StateDB)
	if err := ctx.Err(); err != nil {
		return report, err
	}
	report.PublisherState = state
	for _, item := range ranked[:min(limit, len(ranked))] {
		result := item.result
		for _, source := range result.Profile.Coverage.Sources {
			if source.Kind == catalog.ProviderSearchAdvisoryOverlay {
				continue
			}
			for _, observed := range observations[result.Profile.ProviderID] {
				if reviewedPublisherObservationMatches(source, observed) {
					result.PublisherChecks = append(result.PublisherChecks, ReviewedPublisherObservation{
						SpecRefID: source.SpecRefID, Service: observed.Service, Version: observed.Version,
						Revision: observed.Revision, SHA256: observed.Object.SHA256, CheckedAt: observed.CheckedAt,
					})
				}
			}
		}
		report.Results = append(report.Results, result)
	}
	return report, nil
}

func reviewedPublisherObservationMatches(source catalog.ProviderSearchSource, observed catalogupdate.Version) bool {
	return observed.ReferenceID == source.SpecRefID && observed.Locator == source.SourceURL &&
		(source.SHA256 == "" || observed.CatalogPin == source.SHA256 || observed.Object.SHA256 == source.SHA256)
}

func readReviewedPublisherObservations(ctx context.Context, stateDB string) (map[string][]catalogupdate.Version, string) {
	if stateDB == "" {
		return nil, "unconfigured"
	}
	raw, err := catalogupdate.ReadVersionExport(ctx, stateDB)
	if err != nil {
		return nil, "unavailable"
	}
	var export struct {
		Versions []catalogupdate.Version `json:"versions"`
	}
	if err := json.Unmarshal(raw, &export); err != nil {
		return nil, "unavailable"
	}
	versions := make(map[string][]catalogupdate.Version)
	for _, version := range export.Versions {
		versions[version.ProviderID] = append(versions[version.ProviderID], version)
	}
	return versions, "available"
}

const reviewedProviderIndexApplicationID = 0x41505343

func openReviewedProviderIndex(path string) (*sql.DB, error) {
	dsn := ":memory:"
	if path != "" {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		file, err := os.OpenFile(absolute, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o600)
		if err == nil {
			if err := file.Close(); err != nil {
				return nil, err
			}
		} else if !os.IsExist(err) {
			return nil, err
		}
		info, err := os.Lstat(absolute)
		if err != nil {
			return nil, err
		}
		if !info.Mode().IsRegular() {
			return nil, fmt.Errorf("reviewed provider index must be a regular file")
		}
		location := url.URL{Scheme: "file", Path: filepath.ToSlash(absolute), RawQuery: "mode=rw&_txlock=immediate&_pragma=busy_timeout(5000)"}
		dsn = location.String()
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

const reviewedProviderIndexSchema = `
CREATE TABLE search_metadata (id INTEGER PRIMARY KEY CHECK (id = 1), digest TEXT NOT NULL);
CREATE TABLE search_profile (id TEXT PRIMARY KEY, metadata BLOB NOT NULL);
CREATE VIRTUAL TABLE search_provider USING fts5(provider UNINDEXED, identity, category, keywords, description, tokenize='porter unicode61');
CREATE VIRTUAL TABLE search_operation USING fts5(provider UNINDEXED, source_position UNINDEXED, operation_position UNINDEXED, identity, label, tokenize='porter unicode61');
PRAGMA user_version=1;`

func prepareReviewedProviderIndex(ctx context.Context, tx *sql.Tx, profiles []catalog.ProviderSearchProfile) (string, error) {
	var app, schema int
	if err := tx.QueryRowContext(ctx, "PRAGMA application_id").Scan(&app); err != nil {
		return "", err
	}
	if err := tx.QueryRowContext(ctx, "PRAGMA user_version").Scan(&schema); err != nil {
		return "", err
	}
	if app == 0 && schema == 0 {
		var count int
		if err := tx.QueryRowContext(ctx, "SELECT count(*) FROM sqlite_schema WHERE name NOT GLOB 'sqlite_*'").Scan(&count); err != nil {
			return "", err
		}
		if count != 0 {
			return "", fmt.Errorf("SQLite file belongs to another application")
		}
		if _, err := tx.ExecContext(ctx, reviewedProviderIndexSchema+fmt.Sprintf("PRAGMA application_id=%d;", reviewedProviderIndexApplicationID)); err != nil {
			return "", err
		}
	} else if app != reviewedProviderIndexApplicationID || schema != 1 {
		return "", fmt.Errorf("unrecognized reviewed provider index")
	}
	canonical, err := json.Marshal(profiles)
	if err != nil {
		return "", err
	}
	digest := fmt.Sprintf("%x", sha256.Sum256(canonical))
	var prior string
	err = tx.QueryRowContext(ctx, "SELECT digest FROM search_metadata WHERE id=1").Scan(&prior)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}
	if prior == digest {
		return digest, nil
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM search_profile; DELETE FROM search_provider; DELETE FROM search_operation;"); err != nil {
		return "", err
	}
	profileInsert, err := tx.PrepareContext(ctx, "INSERT INTO search_profile VALUES (?, ?)")
	if err != nil {
		return "", err
	}
	defer profileInsert.Close()
	providerInsert, err := tx.PrepareContext(ctx, "INSERT INTO search_provider VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer providerInsert.Close()
	operationInsert, err := tx.PrepareContext(ctx, "INSERT INTO search_operation VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return "", err
	}
	defer operationInsert.Close()
	for _, profile := range profiles {
		identity := strings.Join(reviewedProviderIdentity(profile), " ")
		raw, err := json.Marshal(profile)
		if err != nil {
			return "", err
		}
		if _, err := profileInsert.ExecContext(ctx, profile.ProviderID, raw); err != nil {
			return "", err
		}
		if _, err := providerInsert.ExecContext(ctx, profile.ProviderID, identity, profile.Category,
			strings.Join(profile.Keywords, " "), profile.Description); err != nil {
			return "", err
		}
		for i, source := range profile.Coverage.Sources {
			for j, operation := range source.Operations {
				if _, err := operationInsert.ExecContext(ctx, profile.ProviderID, i, j, identity, operation.Label); err != nil {
					return "", err
				}
			}
		}
	}
	_, err = tx.ExecContext(ctx, "INSERT OR REPLACE INTO search_metadata VALUES (1, ?)", digest)
	return digest, err
}

func reviewedProviderIdentity(profile catalog.ProviderSearchProfile) []string {
	return append([]string{profile.ProviderID, profile.DisplayName}, profile.Aliases...)
}

type rankedReviewedProvider struct {
	result ReviewedProviderSearchResult
	rank   int
	score  float64
}

func rankReviewedProviders(ctx context.Context, db *sql.DB, profiles []catalog.ProviderSearchProfile, query string) ([]rankedReviewedProvider, string, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, "", err
	}
	defer tx.Rollback()
	digest, err := prepareReviewedProviderIndex(ctx, tx, profiles)
	if err != nil {
		return nil, "", err
	}
	byID := make(map[string]catalog.ProviderSearchProfile, len(profiles))
	matches := make(map[string]rankedReviewedProvider)
	for _, profile := range profiles {
		byID[profile.ProviderID] = profile
		if profile.MatchesIdentity(query) {
			matches[profile.ProviderID] = rankedReviewedProvider{
				result: ReviewedProviderSearchResult{Profile: profile, Match: ReviewedProviderExactMatch},
			}
		}
	}
	var terms []string
	if len(query) <= 1024 {
		terms = reviewedSearchWords(query)
	}
	if len(terms) > 0 && len(terms) <= 24 {
		quoted := make([]string, 0, len(terms))
		for _, term := range terms {
			quoted = append(quoted, fmt.Sprintf("%q", term))
		}
		expression := strings.Join(quoted, " AND ")
		if err := selectReviewedProviderMatches(ctx, tx, matches, byID, expression); err != nil {
			return nil, "", err
		}
		if err := selectReviewedOperationMatches(ctx, tx, matches, byID, expression, terms); err != nil {
			return nil, "", err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, "", err
	}
	ranked := make([]rankedReviewedProvider, 0, len(matches))
	for _, match := range matches {
		ranked = append(ranked, match)
	}
	sortRankedReviewedProviders(ranked)
	return ranked, digest, nil
}

func sortRankedReviewedProviders(ranked []rankedReviewedProvider) {
	slices.SortFunc(ranked, func(a, b rankedReviewedProvider) int {
		if a.rank != b.rank {
			return a.rank - b.rank
		}
		if a.score < b.score {
			return -1
		}
		if a.score > b.score {
			return 1
		}
		return strings.Compare(a.result.Profile.ProviderID, b.result.Profile.ProviderID)
	})
}

func selectReviewedProviderMatches(ctx context.Context, tx *sql.Tx, matches map[string]rankedReviewedProvider, profiles map[string]catalog.ProviderSearchProfile, expression string) error {
	rows, err := tx.QueryContext(ctx, "SELECT provider,bm25(search_provider,0,8,2,3,1) FROM search_provider WHERE search_provider MATCH ?", expression)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var weight float64
		if err := rows.Scan(&id, &weight); err != nil {
			return err
		}
		if profile, ok := profiles[id]; ok {
			if _, matched := matches[id]; !matched {
				matches[id] = rankedReviewedProvider{
					result: ReviewedProviderSearchResult{Profile: profile, Match: ReviewedProviderMetadataMatch}, rank: 3, score: weight,
				}
			}
		}
	}
	return rows.Err()
}

func selectReviewedOperationMatches(ctx context.Context, tx *sql.Tx, matches map[string]rankedReviewedProvider, profiles map[string]catalog.ProviderSearchProfile, expression string, terms []string) error {
	rows, err := tx.QueryContext(ctx, `
SELECT provider,source_position,operation_position,bm25(search_operation,0,0,0,0.5,4)
FROM search_operation WHERE search_operation MATCH ? ORDER BY 4,1,2,3`, expression)
	if err != nil {
		return err
	}
	defer rows.Close()
	identityWords := make(map[string][]string)
	for rows.Next() {
		var id string
		var i, j int
		var weight float64
		if err := rows.Scan(&id, &i, &j, &weight); err != nil {
			return err
		}
		profile, exists := profiles[id]
		if !exists || i < 0 || i >= len(profile.Coverage.Sources) {
			return fmt.Errorf("reviewed provider index contains an unknown source")
		}
		source := profile.Coverage.Sources[i]
		if j < 0 || j >= len(source.Operations) {
			return fmt.Errorf("reviewed provider index contains an unknown operation")
		}
		match, found := matches[id]
		if found && match.rank == 0 {
			continue
		}
		identity, known := identityWords[id]
		if !known {
			identity = reviewedSearchWords(strings.Join(reviewedProviderIdentity(profile), " "))
			identityWords[id] = identity
		}
		if !reviewedOperationEvidence(terms, identity, reviewedSearchWords(source.Operations[j].Label)) {
			continue
		}
		if !found || match.rank != 1 {
			match = rankedReviewedProvider{
				result: ReviewedProviderSearchResult{Profile: profile, Match: ReviewedProviderOperationMatch}, rank: 1, score: weight,
			}
		}
		if len(match.result.MatchedCapabilities) < 10 {
			match.result.MatchedCapabilities = append(match.result.MatchedCapabilities, ReviewedCapabilityMatch{
				SpecRefID: source.SpecRefID, ArtifactID: source.ArtifactID, SHA256: source.SHA256,
				ProviderSearchOperation: source.Operations[j],
			})
		}
		matches[id] = match
	}
	return rows.Err()
}

// FTS stemming proposes candidates. Exact words and simple plurals establish
// operation evidence from the reviewed label.
func reviewedOperationEvidence(query, identity, label []string) bool {
	evidence := false
	for _, term := range query {
		if slices.Contains(identity, term) {
			continue
		}
		matched := false
		for i, word := range label {
			// A label can repeat the provider's name. It is still identity,
			// so plural equivalence cannot turn it into operation evidence.
			if term != word && slices.Contains(identity, word) {
				continue
			}
			if !reviewedOperationWordMatches(term, word) {
				continue
			}
			matched = true
			if isGenericReviewedOperationWord(term) || isGenericReviewedOperationWord(word) {
				continue
			}
			switch word {
			case "user", "users", "product", "products":
				evidence = evidence || reviewedOperationTarget(label, identity, i)
			default:
				evidence = true
			}
		}
		if !matched {
			return false
		}
	}
	return evidence
}

func isGenericReviewedOperationWord(word string) bool {
	switch word {
	case "account", "accounts", "api", "apis", "add", "authenticated", "create", "current", "delete",
		"fetch", "find", "get", "list", "lists", "listing", "listings", "make", "map", "maps",
		"new", "one", "read", "remove", "send", "single", "status", "time", "times", "update":
		return true
	default:
		return false
	}
}

// A user or product can qualify another resource or occur in a description
// after that resource. Only a direct target supplies evidence on its own.
func reviewedOperationTarget(label, identity []string, target int) bool {
	if target+1 < len(label) {
		switch label[target+1] {
		case "by", "or":
		default:
			return false
		}
	}
	for i := 1; i < target; i++ {
		switch label[i] {
		case "authenticated", "current", "one", "single":
		default:
			if !slices.Contains(identity, label[i]) {
				return false
			}
		}
	}
	return true
}

func reviewedOperationWordMatches(a, b string) bool {
	if a == b {
		return true
	}
	if len(a) > len(b) {
		a, b = b, a
	}
	if len(a) < 3 || (a == "new" && b == "news") {
		return false
	}
	switch {
	case strings.HasSuffix(a, "y") && !strings.ContainsRune("aeiou", rune(a[len(a)-2])):
		return b == strings.TrimSuffix(a, "y")+"ies"
	case strings.HasSuffix(a, "s"), strings.HasSuffix(a, "x"), strings.HasSuffix(a, "z"),
		strings.HasSuffix(a, "ch"), strings.HasSuffix(a, "sh"):
		return b == a+"es"
	default:
		return b == a+"s"
	}
}

func reviewedSearchWords(text string) []string {
	var words []string
	for _, word := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	}) {
		switch word {
		case "a", "an", "the", "and", "to", "in", "on", "for", "of", "with", "using", "via", "my":
			continue
		}
		if !slices.Contains(words, word) {
			words = append(words, word)
		}
	}
	return words
}
