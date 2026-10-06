package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

func validateAPIVersionHints(hints APIVersionHints) error {
	data, err := json.Marshal(hints)
	if err != nil || len(data) > 64<<10 || len(hints.Hints) > 8 || hints.SchemaVersion != APIVersionHintsSchemaVersion {
		return fmt.Errorf("invalid API version hints schema, size or count")
	}
	return nil
}

// VerifyHints verifies untrusted candidate URLs using a new bounded invocation.
// Producer notes and claimed versions never grant source authority.
func (c *Client) VerifyHints(ctx context.Context, request APIVersionDiscoveryRequest, hints APIVersionHints, options APIVersionDiscoveryOptions) (APIVersionDiscoveryReport, error) {
	if err := validateAPIVersionHints(hints); err != nil {
		return APIVersionDiscoveryReport{}, err
	}
	s, err := newAPIVersionSession(ctx, c, request, options)
	if err != nil {
		return APIVersionDiscoveryReport{}, err
	}
	defer s.cancel()
	s.verifyHints(hints)
	return s.finish(), nil
}

func (s *apiVersionSession) verifyHints(hints APIVersionHints) {
	if !s.opts.Network {
		s.gap("version.network_denied", "Hint verification was not enabled.")
		return
	}
	rows := append([]APIVersionHint(nil), hints.Hints...)
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Version != rows[j].Version {
			return apiVersionLess(rows[i].Version, rows[j].Version)
		}
		if rows[i].URL != rows[j].URL {
			return rows[i].URL < rows[j].URL
		}
		return rows[i].Producer < rows[j].Producer
	})
	seen := map[string]bool{}
	for _, hint := range rows {
		if seen[hint.URL] {
			continue
		}
		seen[hint.URL] = true
		if len(hint.URL) > 2048 {
			s.gap("version.hint_invalid", "Hint URL exceeds its budget.")
			continue
		}
		if err := s.client.versionScopeAllows(hint.URL, s.scopes); err != nil {
			s.gap("version.hint_invalid", err.Error())
			continue
		}
		l := APIVersionLocator{Kind: "direct", URL: hint.URL}
		if hint.Locator != nil {
			test := s.req
			test.Locators = []APIVersionLocator{*hint.Locator}
			if _, err := validateVersionRequest(s.client, test, s.opts); err != nil {
				s.gap("version.hint_recipe", err.Error())
				continue
			}
			l = *hint.Locator
		}
		s.fetchSource(hint.URL, l, versionTokenURL(hint.URL), "", "")
		producer, changed := sanitizePromptString(hint.Producer, 256)
		claim, claimChanged := sanitizePromptString(hint.Version, 256)
		_, noteChanged := sanitizePromptString(hint.Note, 2048)
		s.mu.Lock()
		for i := range s.report.Versions {
			if s.report.Versions[i].ID == versionRecordID(hint.URL) {
				s.report.Versions[i].HintProducer = producer
				s.report.Versions[i].VersionClaims.DocsVersion = claim
				if changed || claimChanged || noteChanged {
					s.report.Versions[i].Diagnostics = append(s.report.Versions[i].Diagnostics, versionDiagnostic("version.hint_shortened", "Untrusted hint text was sanitized or shortened."))
				}
			}
		}
		s.mu.Unlock()
	}
	s.tier("hints", "examined")
	if err := s.ctx.Err(); err != nil {
		s.gap("version.hint_deadline", err.Error())
	}
}

// DiscoverAPIVersionsWithAdapter is an optional pull integration for a
// context-cooperative adapter. The deterministic network budget is never reset.
func (c *Client) DiscoverAPIVersionsWithAdapter(ctx context.Context, request APIVersionDiscoveryRequest, options APIVersionDiscoveryOptions, adapter APIVersionHintAdapter) (APIVersionDiscoveryReport, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	timeout := options.AdapterTimeout
	if timeout == 0 {
		timeout = 15 * time.Second
	}
	if timeout < 0 || timeout > 60*time.Second {
		return APIVersionDiscoveryReport{}, fmt.Errorf("invalid adapter deadline")
	}
	adapterCtx, stop := context.WithTimeout(ctx, timeout)
	defer stop()
	s, err := newAPIVersionSession(adapterCtx, c, request, options)
	if err != nil {
		return APIVersionDiscoveryReport{}, err
	}
	defer s.cancel()
	s.discover()
	result := s.finish()
	if !options.Network || !options.AdapterEnabled || adapter == nil || (result.Status != "unexamined" && result.Status != "conflicting" && !s.incomplete) {
		return result, nil
	}
	// Independent copies keep an adapter from mutating the request or verified
	// result before it submits new, explicitly untrusted hints.
	data, _ := json.Marshal(APIVersionHintContext{Request: request, Report: result})
	var input APIVersionHintContext
	_ = json.Unmarshal(data, &input)
	hints, err := adapter.DiscoverHints(adapterCtx, input)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, versionDiagnostic("version.adapter", err.Error()))
		return result, nil
	}
	if err = validateAPIVersionHints(hints); err != nil {
		result.Diagnostics = append(result.Diagnostics, versionDiagnostic("version.adapter", err.Error()))
		return result, nil
	}
	s.reused = false
	s.verifyHints(hints)
	return s.finish(), nil
}
