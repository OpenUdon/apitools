package apitools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/OpenUdon/apitools/catalog"
	"github.com/OpenUdon/apitools/internal/sourceguard"
)

const (
	APIVersionDiscoverySchemaVersion = "apitools.api-version-discovery/v1"
	APIVersionHintsSchemaVersion     = "apitools.api-version-hints/v1"
	APIVersionStateSchemaVersion     = "apitools.api-version-state/v1"
	DefaultAPIVersionTimeout         = 5 * time.Second
	apiVersionMaxSourceBytes         = 8 << 20
	apiVersionMaxReportBytes         = 512 << 10
)

type APIVersionBaseline struct {
	SourceKind     OperationSourceKind `json:"source_kind"`
	ProviderKey    string              `json:"provider_key,omitempty"`
	SourceURL      string              `json:"source_url,omitempty"`
	Version        string              `json:"version"`
	SHA256         string              `json:"sha256"`
	CatalogUpdated *time.Time          `json:"catalog_updated,omitempty"`
	ETag           string              `json:"etag,omitempty"`
	LastModified   string              `json:"last_modified,omitempty"`
}

type APIVersionOfficialScope struct {
	Origin     string `json:"origin"`
	PathPrefix string `json:"path_prefix"`
	Repository string `json:"repository,omitempty"`
}

type APIVersionLocator struct {
	Kind          string `json:"kind"`
	URL           string `json:"url,omitempty"`
	URLTemplate   string `json:"url_template,omitempty"`
	BaselineToken string `json:"baseline_token,omitempty"`
	Repository    string `json:"repository,omitempty"`
	Ref           string `json:"ref,omitempty"`
	PathPattern   string `json:"path_pattern,omitempty"`
}

type APIVersionDiscoveryRequest struct {
	SchemaVersion  string                    `json:"schema_version"`
	Known          APIVersionBaseline        `json:"known"`
	OfficialScopes []APIVersionOfficialScope `json:"official_scopes,omitempty"`
	DocsURL        string                    `json:"docs_url,omitempty"`
	Locators       []APIVersionLocator       `json:"locators,omitempty"`
	Contract       *StepContract             `json:"contract,omitempty"`
}

type APIVersionDiscoveryOptions struct {
	Network                 bool
	Timeout                 time.Duration
	StatePath               string
	ListCachePath           string
	SaveDir                 string
	BaselineInventory       *OperationInventory
	BaselineInventorySHA256 string
	Previous                *APIVersionDiscoveryReport
	Catalog                 *catalog.Catalog
	AdapterEnabled          bool
	AdapterTimeout          time.Duration
}

type APIVersionClaims struct {
	URLToken    string `json:"url_token,omitempty"`
	InfoVersion string `json:"info_version,omitempty"`
	DocsVersion string `json:"docs_version,omitempty"`
}

type APIVersionOperationDiff struct {
	ChangedOperationIDs []string `json:"changed_operation_ids,omitempty"`
	Added               []string `json:"added,omitempty"`
	Removed             []string `json:"removed,omitempty"`
	Retained            []string `json:"retained,omitempty"`
}

type APIVersionRecord struct {
	SourceKind     OperationSourceKind      `json:"source_kind,omitempty"`
	CatalogUpdated *time.Time               `json:"catalog_updated,omitempty"`
	ID             string                   `json:"id"`
	Locator        APIVersionLocator        `json:"locator"`
	Recipe         APIVersionLocator        `json:"recipe"`
	Evidence       string                   `json:"evidence"`
	VersionClaims  APIVersionClaims         `json:"version_claims"`
	SourceURL      string                   `json:"source_url,omitempty"`
	FinalURL       string                   `json:"final_url,omitempty"`
	SHA256         string                   `json:"sha256,omitempty"`
	BytesObserved  int64                    `json:"bytes_observed"`
	BytesDeclared  *int64                   `json:"bytes_declared,omitempty"`
	FetchedAt      *time.Time               `json:"fetched_at,omitempty"`
	CheckedAt      time.Time                `json:"checked_at"`
	Validation     string                   `json:"validation"`
	Comparison     string                   `json:"comparison"`
	Diff           *APIVersionOperationDiff `json:"diff,omitempty"`
	Candidates     []OperationCandidate     `json:"candidates,omitempty"`
	Saved          *ImportedSpec            `json:"saved,omitempty"`
	LocatorStatus  string                   `json:"locator_status"`
	ETag           string                   `json:"etag,omitempty"`
	LastModified   string                   `json:"last_modified,omitempty"`
	HintProducer   string                   `json:"hint_producer,omitempty"`
	Diagnostics    []Diagnostic             `json:"diagnostics,omitempty"`
}

type APIVersionTier struct {
	Tier            string       `json:"tier"`
	Status          string       `json:"status"`
	RequestCount    int          `json:"request_count"`
	SourceBodyCount int          `json:"source_body_count"`
	CheckedAt       time.Time    `json:"checked_at"`
	Diagnostics     []Diagnostic `json:"diagnostics,omitempty"`
}

type APIVersionDiscoveryReport struct {
	QueryIdentity    string                    `json:"x-apitools-query-identity,omitempty"`
	SchemaVersion    string                    `json:"schema_version"`
	Known            APIVersionBaseline        `json:"known"`
	Status           string                    `json:"status"`
	CheckedAt        time.Time                 `json:"checked_at"`
	CheckedScope     []APIVersionOfficialScope `json:"checked_scope,omitempty"`
	Versions         []APIVersionRecord        `json:"versions,omitempty"`
	Tiers            []APIVersionTier          `json:"tiers,omitempty"`
	Diagnostics      []Diagnostic              `json:"diagnostics,omitempty"`
	Truncated        bool                      `json:"truncated"`
	Preferred        string                    `json:"preferred,omitempty"`
	DirectoryRecords map[string]any            `json:"directory_records,omitempty"`
}

type APIVersionHint struct {
	URL      string             `json:"url"`
	Version  string             `json:"version,omitempty"`
	Note     string             `json:"note,omitempty"`
	Producer string             `json:"producer,omitempty"`
	Locator  *APIVersionLocator `json:"locator,omitempty"`
}

type APIVersionHints struct {
	SchemaVersion string           `json:"schema_version"`
	Hints         []APIVersionHint `json:"hints"`
}

type APIVersionHintContext struct {
	Request APIVersionDiscoveryRequest `json:"request"`
	Report  APIVersionDiscoveryReport  `json:"report"`
}

type APIVersionHintAdapter interface {
	DiscoverHints(context.Context, APIVersionHintContext) (APIVersionHints, error)
}

// DecodeAPIVersionDiscoveryRequest decodes a bounded strict request. Runtime
// paths, clients and network permission are supplied separately as Go options.
func DecodeAPIVersionDiscoveryRequest(data []byte) (APIVersionDiscoveryRequest, error) {
	var request APIVersionDiscoveryRequest
	if err := decodeVersionJSON(data, &request, 64<<10); err != nil {
		return request, err
	}
	if request.SchemaVersion != APIVersionDiscoverySchemaVersion {
		return request, fmt.Errorf("unsupported API version request schema %q", request.SchemaVersion)
	}
	return request, nil
}

func DecodeAPIVersionHints(data []byte) (APIVersionHints, error) {
	var hints APIVersionHints
	if err := decodeVersionJSON(data, &hints, 64<<10); err != nil {
		return hints, err
	}
	if hints.SchemaVersion != APIVersionHintsSchemaVersion || len(hints.Hints) > 8 {
		return hints, fmt.Errorf("invalid API version hints schema or count")
	}
	return hints, nil
}

func decodeVersionJSON(data []byte, target any, max int) error {
	if len(data) > max {
		return fmt.Errorf("API version metadata exceeds %d bytes", max)
	}
	if err := sourceguard.CheckJSON("api-version-metadata", data); err != nil {
		return err
	}
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return err
	}
	if err := d.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("API version metadata contains trailing data")
	}
	return nil
}
