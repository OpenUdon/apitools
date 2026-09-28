package apitools

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

type loadedOperationSource struct {
	input       OperationSourceInput
	content     []byte
	digest      string
	readErr     error
	failureCode string
	inputIndex  int
	urlRedacted bool
}

type operationSourceAdapterResult struct {
	operationCount int
	title          string
	capabilities   []OperationCapability
	candidates     []OperationCandidate
	diagnostics    []Diagnostic
}

// BuildOperationCandidates parses explicitly supplied local API source
// documents and ranks their native operations against a step contract. URLs
// are provenance only: this function never fetches, resolves credentials, or
// executes an operation.
func BuildOperationCandidates(ctx context.Context, options OperationCandidateOptions) (OperationCandidateReport, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	report := OperationCandidateReport{SchemaVersion: OperationCandidatesSchemaVersion}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	budget := resolvedPromptBudget(options.PromptBudget)
	if options.MaxOperations < 0 || options.MaxCandidates < 0 || options.MaxBytes < 0 {
		return candidateRequestError(report, "candidate.limit_invalid", "source, operation, candidate, and byte limits cannot be negative")
	}
	if options.MaxBytes > DefaultMaxBytes {
		return candidateRequestError(report, "candidate.byte_limit_invalid", fmt.Sprintf("MaxBytes cannot exceed the parser limit of %d bytes", DefaultMaxBytes))
	}
	if options.MaxOperations > 0 {
		budget.MaxOperations = options.MaxOperations
	}
	if len(options.Sources) == 0 {
		return candidateRequestError(report, "source.required", "at least one explicit local API source is required")
	}
	if len(options.Sources) > budget.MaxCollectionItems {
		report.Truncated = true
		return candidateRequestError(report, "source.count_limit", fmt.Sprintf("%d sources exceed the configured source limit of %d", len(options.Sources), budget.MaxCollectionItems))
	}

	loaded := make([]loadedOperationSource, 0, len(options.Sources))
	maxBytes := options.MaxBytes
	if maxBytes <= 0 {
		maxBytes = DefaultMaxBytes
	}
	for i, input := range options.Sources {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		item := loadedOperationSource{input: input, inputIndex: i}
		cleanURL, urlRedacted, urlErr := sanitizeOperationSourceURL(input.URL)
		item.input.URL = cleanURL
		item.urlRedacted = urlRedacted
		if urlErr != nil {
			item.readErr = urlErr
			item.failureCode = "source.identity_invalid"
		} else if err := validateOperationSourceIdentity(item.input, budget); err != nil {
			item.readErr = err
			item.failureCode = "source.identity_invalid"
		} else {
			switch {
			case input.Content != nil:
				item.content = append([]byte(nil), input.Content...)
				item.readErr = validateInlineSpecContent(item.content, maxBytes, sourceLabel(input))
			case strings.TrimSpace(input.Path) != "":
				item.content, item.readErr = readLocalSpecFile(input.Path, maxBytes)
			case strings.TrimSpace(input.URL) != "":
				item.readErr = errors.New("a source URL is provenance only; provide local bytes or a local path")
			default:
				item.readErr = errors.New("source content or a local path is required")
			}
		}
		if item.readErr == nil {
			digest := sha256.Sum256(item.content)
			item.digest = hex.EncodeToString(digest[:])
		}
		loaded = append(loaded, item)
	}
	sort.SliceStable(loaded, func(i, j int) bool { return operationSourceLess(loaded[i], loaded[j]) })

	var candidates []OperationCandidate
	totalOperations := 0
	for _, source := range loaded {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		sourceReport := OperationSourceReport{
			Kind:         source.input.Kind,
			Name:         safeSourceValue(sourceName(source.input), budget.MaxIdentifierRunes),
			DocumentPath: safeSourceValue(source.input.Path, budget.MaxTextRunes),
			DocumentURL:  safeSourceValue(source.input.URL, budget.MaxTextRunes),
			SHA256:       source.digest,
		}
		if source.readErr != nil {
			report.Truncated = true
			code := source.failureCode
			message := "source bytes could not be read within the local source safety limits"
			if code == "" {
				code = "source.read"
			}
			if code == "source.identity_invalid" {
				message = "source name, path, or provenance URL is unsafe or exceeds the identifier budget"
			}
			diagnostic := sourceDiagnostic(source, code, message, source.readErr)
			sourceReport.Diagnostics = []Diagnostic{diagnostic}
			report.Diagnostics = append(report.Diagnostics, diagnostic)
			report.Sources = append(report.Sources, sourceReport)
			continue
		}
		if !validOperationSourceKind(source.input.Kind) {
			report.Truncated = true
			diagnostic := sourceDiagnostic(source, "source.kind_invalid", "source kind is not supported", nil)
			sourceReport.Diagnostics = []Diagnostic{diagnostic}
			report.Diagnostics = append(report.Diagnostics, diagnostic)
			report.Sources = append(report.Sources, sourceReport)
			continue
		}

		adapted, err := adaptOperationSource(ctx, source, budget)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return report, ctxErr
		}
		adapterErr := err
		var rationaleDiagnostics []Diagnostic
		var rationaleErr error
		adapted.candidates, rationaleDiagnostics, rationaleErr = sanitizeAdapterCandidates(adapted.candidates, budget)
		adapted.diagnostics = append(adapted.diagnostics, rationaleDiagnostics...)
		if rationaleErr != nil {
			err = rationaleErr
		} else {
			err = adapterErr
		}
		if err == nil {
			err = firstAdapterError(adapted.diagnostics)
		}
		sourceReport.Title = adapted.title
		if title, changed := sanitizePromptString(sourceReport.Title, budget.MaxTextRunes); changed {
			sourceReport.Title = title
			report.Diagnostics = append(report.Diagnostics, Diagnostic{Severity: "warning", Code: "source.title_sanitized", Message: "source title contained unsafe controls or exceeded the display text budget and was sanitized", Path: sourceReport.DocumentPath})
		}
		sourceReport.OperationCount = adapted.operationCount
		sourceReport.Capabilities = adapted.capabilities
		sourceReport.Diagnostics = adapted.diagnostics
		if source.urlRedacted {
			diagnostic := Diagnostic{Severity: "warning", Code: "source.url_redacted", Message: "source URL userinfo, query, or fragment was omitted from public provenance", Path: sourceReport.DocumentPath}
			sourceReport.Diagnostics = append(sourceReport.Diagnostics, diagnostic)
			report.Diagnostics = append(report.Diagnostics, diagnostic)
		}
		if err != nil {
			report.Truncated = true
			diagnostic := sourceDiagnostic(source, "source.parse", "source could not be parsed or summarized by its declared family", err)
			sourceReport.Diagnostics = append(sourceReport.Diagnostics, diagnostic)
			report.Diagnostics = append(report.Diagnostics, diagnostic)
		}
		sortDiagnostics(sourceReport.Diagnostics)
		report.Diagnostics = append(report.Diagnostics, adapted.diagnostics...)
		candidates = append(candidates, adapted.candidates...)
		totalOperations += adapted.operationCount
		report.Sources = append(report.Sources, sourceReport)
		if totalOperations > budget.MaxOperations || len(candidates) > budget.MaxOperations {
			diagnostic := Diagnostic{
				Severity: "error", Code: "candidate.operation_work_budget",
				Message:     fmt.Sprintf("sources contain more than the configured %d-operation ranking budget", budget.MaxOperations),
				Remediation: "Narrow the explicit source set or raise MaxOperations for a reviewed workload.",
			}
			report.Diagnostics = append(report.Diagnostics, diagnostic)
			report.Truncated = true
			sortDiagnostics(report.Diagnostics)
			return report, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
		}
	}

	ranked, rankErr := rankOperationCandidates(ctx, candidates, options.Contract, budget, options.MaxCandidates)
	if rankErr != nil {
		var diagnosticErr DiagnosticError
		if errors.As(rankErr, &diagnosticErr) {
			report.Diagnostics = append(report.Diagnostics, ranked.Diagnostics...)
			report.Candidates = ranked.Candidates
			report.Truncated = report.Truncated || ranked.Truncated
			return enforceCandidateContextBudget(report, budget)
		}
		return report, rankErr
	}
	report.Candidates = ranked.Candidates
	report.Diagnostics = append(report.Diagnostics, ranked.Diagnostics...)
	report.Truncated = report.Truncated || ranked.Truncated
	return enforceCandidateContextBudget(report, budget)
}

func sanitizeAdapterCandidates(candidates []OperationCandidate, budget PromptBudget) ([]OperationCandidate, []Diagnostic, error) {
	sanitized := make([]OperationCandidate, 0, len(candidates))
	var diagnostics []Diagnostic
	var firstErr error
	for _, candidate := range candidates {
		if err := sanitizeCandidateRationale(&candidate, budget); err != nil {
			var diagnosticErr DiagnosticError
			if errors.As(err, &diagnosticErr) {
				diagnostics = append(diagnostics, diagnosticErr.Diagnostics...)
			}
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		sanitized = append(sanitized, candidate)
	}
	return sanitized, diagnostics, firstErr
}

func firstAdapterError(diagnostics []Diagnostic) error {
	if errorsFound := errorDiagnostics(diagnostics); len(errorsFound) > 0 {
		return DiagnosticError{Diagnostics: errorsFound}
	}
	return nil
}

func candidateRequestError(report OperationCandidateReport, code, message string) (OperationCandidateReport, error) {
	diagnostic := Diagnostic{Severity: "error", Code: code, Message: message}
	report.Diagnostics = []Diagnostic{diagnostic}
	return report, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
}

func enforceCandidateContextBudget(report OperationCandidateReport, budget PromptBudget) (OperationCandidateReport, error) {
	if len(report.Diagnostics) > 1 {
		sortDiagnostics(report.Diagnostics)
	}
	for _, diagnostic := range report.Diagnostics {
		if diagnostic.Code == "candidate.operation_byte_budget" {
			report.Candidates = nil
			break
		}
	}
	size, err := marshalSize(report)
	if err != nil {
		return report, fmt.Errorf("encode operation candidate report: %w", err)
	}
	if size > budget.MaxContextBytes {
		diagnostic := Diagnostic{
			Severity: "error", Code: "candidate.context_byte_budget",
			Message:     fmt.Sprintf("operation candidate report is %d bytes, exceeding the %d-byte context budget", size, budget.MaxContextBytes),
			Remediation: "Narrow sources or increase MaxContextBytes for a reviewed report.",
		}
		bounded := OperationCandidateReport{
			SchemaVersion: report.SchemaVersion,
			Diagnostics:   []Diagnostic{diagnostic},
			Truncated:     true,
		}
		boundedSize, sizeErr := marshalSize(bounded)
		if sizeErr != nil {
			return OperationCandidateReport{}, fmt.Errorf("encode bounded operation candidate error report: %w", sizeErr)
		}
		if boundedSize > budget.MaxContextBytes {
			return OperationCandidateReport{}, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
		}
		return bounded, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	if diagnostics := errorDiagnostics(report.Diagnostics); len(diagnostics) > 0 {
		return report, DiagnosticError{Diagnostics: diagnostics}
	}
	return report, nil
}

func sortDiagnostics(diagnostics []Diagnostic) {
	sort.SliceStable(diagnostics, func(i, j int) bool {
		left, right := diagnostics[i], diagnostics[j]
		if left.Code != right.Code {
			return left.Code < right.Code
		}
		if left.Path != right.Path {
			return left.Path < right.Path
		}
		if left.Severity != right.Severity {
			return left.Severity < right.Severity
		}
		if left.Message != right.Message {
			return left.Message < right.Message
		}
		return left.Remediation < right.Remediation
	})
}

func sourceDiagnostic(source loadedOperationSource, code, message string, cause error) Diagnostic {
	_ = cause // Do not echo untrusted parser or file error text into authoring output.
	message, _ = sanitizePromptString(message, DefaultPromptTextRunes)
	path, _ := sanitizePromptString(firstNonEmpty(source.input.Path, source.input.URL, sourceName(source.input)), DefaultPromptTextRunes)
	return Diagnostic{Severity: "error", Code: code, Message: message, Path: path}
}

func validateOperationSourceIdentity(input OperationSourceInput, budget PromptBudget) error {
	for _, value := range []struct {
		name  string
		text  string
		limit int
	}{
		{name: "source name", text: sourceName(input), limit: budget.MaxIdentifierRunes},
		{name: "source path", text: input.Path, limit: budget.MaxTextRunes},
	} {
		if _, changed := sanitizePromptString(value.text, value.limit); changed {
			return fmt.Errorf("%s is unsafe or exceeds its display budget", value.name)
		}
	}
	cleanURL, _, err := sanitizeOperationSourceURL(input.URL)
	if err != nil {
		return err
	}
	if _, textChanged := sanitizePromptString(cleanURL, budget.MaxTextRunes); textChanged {
		return errors.New("source URL exceeds its display budget")
	}
	return nil
}

func sanitizeOperationSourceURL(raw string) (string, bool, error) {
	if strings.TrimSpace(raw) == "" {
		return "", false, nil
	}
	parsed, err := url.Parse(raw)
	if err != nil || !parsed.IsAbs() || parsed.Host == "" || parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", false, errors.New("source URL must be an absolute HTTP(S) URL")
	}
	changed := false
	if parsed.User != nil {
		parsed.User = nil
		changed = true
	}
	if parsed.RawQuery != "" || parsed.ForceQuery {
		parsed.RawQuery = ""
		parsed.ForceQuery = false
		changed = true
	}
	if parsed.Fragment != "" || parsed.RawFragment != "" {
		parsed.Fragment = ""
		parsed.RawFragment = ""
		changed = true
	}
	return parsed.String(), changed, nil
}

func safeSourceValue(value string, limit int) string {
	value, _ = sanitizePromptString(value, limit)
	return value
}

func sourceLabel(input OperationSourceInput) string {
	return firstNonEmpty(input.Path, input.URL, input.Name, string(input.Kind))
}

func sourceName(input OperationSourceInput) string {
	if name := strings.TrimSpace(input.Name); name != "" {
		return name
	}
	if path := strings.TrimSpace(input.Path); path != "" {
		return filepath.Base(path)
	}
	return string(input.Kind)
}

func operationSourceLess(left, right loadedOperationSource) bool {
	leftKeys := []string{string(left.input.Kind), filepath.Clean(left.input.Path), strings.TrimSpace(left.input.URL), strings.TrimSpace(left.input.Name), left.digest, sourceLoadError(left)}
	rightKeys := []string{string(right.input.Kind), filepath.Clean(right.input.Path), strings.TrimSpace(right.input.URL), strings.TrimSpace(right.input.Name), right.digest, sourceLoadError(right)}
	for i := range leftKeys {
		if leftKeys[i] != rightKeys[i] {
			return leftKeys[i] < rightKeys[i]
		}
	}
	return left.inputIndex < right.inputIndex
}

func sourceLoadError(source loadedOperationSource) string {
	if source.readErr == nil {
		return ""
	}
	return source.readErr.Error()
}

func validOperationSourceKind(kind OperationSourceKind) bool {
	switch kind {
	case OperationSourceOpenAPI, OperationSourceGoogleDiscovery, OperationSourceAWSSmithy,
		OperationSourceAsyncAPI, OperationSourceGraphQL, OperationSourceOpenRPC,
		OperationSourceGRPCProtobuf, OperationSourceOData:
		return true
	default:
		return false
	}
}

func adaptOperationSource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	switch source.input.Kind {
	case OperationSourceOpenAPI:
		return adaptOpenAPISource(ctx, source, budget)
	case OperationSourceGoogleDiscovery:
		return adaptGoogleDiscoverySource(ctx, source, budget)
	case OperationSourceAWSSmithy:
		return adaptSmithySource(ctx, source, budget)
	case OperationSourceAsyncAPI:
		return adaptAsyncAPISource(ctx, source, budget)
	case OperationSourceGraphQL:
		return adaptGraphQLSource(ctx, source, budget)
	case OperationSourceOpenRPC:
		return adaptOpenRPCSource(ctx, source, budget)
	case OperationSourceGRPCProtobuf:
		return adaptGRPCSource(ctx, source, budget)
	case OperationSourceOData:
		return adaptODataSource(ctx, source, budget)
	default:
		return operationSourceAdapterResult{}, fmt.Errorf("unsupported source kind %q", source.input.Kind)
	}
}

func adaptOpenAPISource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	doc := InventoryDocument{
		Name: sourceName(source.input), Path: source.input.Path, URL: source.input.URL,
		Content: source.content,
	}
	inventory, err := buildOperationInventory(ctx, InventoryOptions{
		Documents: []InventoryDocument{doc}, MaxBytes: int64(len(source.content)), MaxOperations: budget.MaxOperations,
	}, budget)
	result := operationSourceAdapterResult{}
	if len(inventory.Documents) > 0 {
		result.title = firstNonEmpty(inventory.Documents[0].Title, inventory.Documents[0].Name)
	}
	result.diagnostics = append(result.diagnostics, inventory.Diagnostics...)
	result.operationCount = inventory.VisitedOperations
	for _, operation := range inventory.Operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		selector := openAPIOperationSelector(operation.Method, operation.Path)
		truncatedSchema := false
		for _, issue := range operation.ReadinessIssues {
			if issue.Code == "schema.field_summary_truncated" {
				truncatedSchema = true
			}
		}
		if truncatedSchema {
			diagnostic := Diagnostic{Severity: "error", Code: "candidate.schema_field_truncated", Message: "OpenAPI request or response fields exceed the inventory field-depth/count bounds; no incomplete candidate was ranked", Path: selector, Remediation: "Narrow the operation or use a reviewed source with a bounded complete schema summary."}
			result.diagnostics = append(result.diagnostics, diagnostic)
			continue
		}
		candidate, diagnostics, summaryErr := candidateFromOperation(source, operation, selector, assessOperationEffect(operation, selector), openAPICapabilities(operation, selector), budget)
		if summaryErr == nil {
			for _, capability := range candidate.Capabilities {
				if capability.Dimension == "inputs" || capability.Dimension == "outputs" {
					candidate.Summary.Gaps = append(candidate.Summary.Gaps, capability.Gaps...)
				}
			}
			candidate.Summary.Gaps = uniqueSortedStrings(candidate.Summary.Gaps)
		}
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if summaryErr != nil {
			continue
		}
		result.candidates = append(result.candidates, candidate)
	}
	result.capabilities = openAPISourceCapabilities()
	if err != nil {
		return result, err
	}
	if diagnostics := errorDiagnostics(inventory.Diagnostics); len(diagnostics) > 0 {
		return result, DiagnosticError{Diagnostics: diagnostics}
	}
	if diagnostics := errorDiagnostics(result.diagnostics); len(diagnostics) > 0 {
		return result, DiagnosticError{Diagnostics: diagnostics}
	}
	return result, nil
}

func candidateFromOperation(source loadedOperationSource, operation OperationSummary, selector string, effect EffectAssessment, capabilities []OperationCapability, budget PromptBudget) (OperationCandidate, []Diagnostic, error) {
	operation, omissionGaps, operationDiagnostics, err := sanitizeCandidateOperation(operation, budget)
	if err != nil {
		return OperationCandidate{}, operationDiagnostics, err
	}
	summary, diagnostics, err := summarizeOperationForConsumer(operation, selector, budget)
	diagnostics = append(operationDiagnostics, diagnostics...)
	if err != nil {
		return OperationCandidate{}, diagnostics, err
	}
	summary.Gaps = uniqueSortedStrings(append(summary.Gaps, omissionGaps...))
	if effect.Class == "" {
		effect.Class = OperationEffectUnknown
	}
	candidate := OperationCandidate{
		Source: OperationSourceIdentity{
			Kind: source.input.Kind, DocumentPath: source.input.Path, DocumentURL: source.input.URL,
			SHA256: source.digest, Selector: selector,
		},
		Operation: operation, Summary: summary, Capabilities: capabilities, Effect: effect,
	}
	if err := sanitizeCandidateRationale(&candidate, budget); err != nil {
		var diagnosticErr DiagnosticError
		if errors.As(err, &diagnosticErr) {
			diagnostics = append(diagnostics, diagnosticErr.Diagnostics...)
		}
		return OperationCandidate{}, diagnostics, err
	}
	return candidate, diagnostics, nil
}

func sanitizeCandidateOperation(operation OperationSummary, budget PromptBudget) (OperationSummary, []string, []Diagnostic, error) {
	report, err := SanitizeOperationSummaries([]OperationSummary{operation}, budget)
	if err != nil {
		return OperationSummary{}, nil, report.Diagnostics, err
	}
	if len(report.Operations) != 1 {
		return OperationSummary{}, nil, report.Diagnostics, fmt.Errorf("candidate operation sanitization returned %d operations", len(report.Operations))
	}
	operation = report.Operations[0]
	var gaps []string
	removed := false
	parameters := operation.Parameters[:0]
	for _, parameter := range operation.Parameters {
		if looksLikeCredentialName(parameter.Name) {
			removed = true
			continue
		}
		parameters = append(parameters, parameter)
	}
	operation.Parameters = parameters
	if operation.RequestBody != nil {
		operation.RequestBody.Fields, removed = filterCredentialFields(operation.RequestBody.Fields, removed)
		operation.RequestBody.RequiredFieldPaths = filterCredentialNames(operation.RequestBody.RequiredFieldPaths)
		if operation.RequestBody.Schema != nil {
			operation.RequestBody.Schema.Properties, removed = filterCredentialProperties(operation.RequestBody.Schema.Properties, removed)
			operation.RequestBody.Schema.Required = filterCredentialNames(operation.RequestBody.Schema.Required)
		}
	}
	if operation.ResponseBody != nil {
		operation.ResponseBody.Fields, removed = filterCredentialFields(operation.ResponseBody.Fields, removed)
		if operation.ResponseBody.Schema != nil {
			operation.ResponseBody.Schema.Properties, removed = filterCredentialProperties(operation.ResponseBody.Schema.Properties, removed)
			operation.ResponseBody.Schema.Required = filterCredentialNames(operation.ResponseBody.Schema.Required)
		}
	}
	if removed {
		gaps = append(gaps, "Credential-shaped request or response fields were omitted from workflow data metadata and need separate authentication review.")
	}
	return operation, gaps, report.Diagnostics, nil
}

func filterCredentialFields(fields []RequestFieldSummary, removed bool) ([]RequestFieldSummary, bool) {
	filtered := fields[:0]
	for _, field := range fields {
		if looksLikeCredentialName(field.Path) {
			removed = true
			continue
		}
		filtered = append(filtered, field)
	}
	return filtered, removed
}

func filterCredentialProperties(properties []PropertySummary, removed bool) ([]PropertySummary, bool) {
	filtered := properties[:0]
	for _, property := range properties {
		if looksLikeCredentialName(property.Name) {
			removed = true
			continue
		}
		filtered = append(filtered, property)
	}
	return filtered, removed
}

func filterCredentialNames(names []string) []string {
	filtered := names[:0]
	for _, name := range names {
		if !looksLikeCredentialName(name) {
			filtered = append(filtered, name)
		}
	}
	return filtered
}

func sanitizeCandidateRationale(candidate *OperationCandidate, budget PromptBudget) error {
	if candidate == nil {
		return nil
	}
	for i, gap := range candidate.Summary.Gaps {
		clean, changed := sanitizePromptString(gap, budget.MaxTextRunes)
		if changed {
			diagnostic := Diagnostic{Severity: "error", Code: "candidate.rationale_unsafe", Message: "candidate explanation text contains unsafe controls or exceeds the display text budget", Path: candidate.Source.Selector, Remediation: "Review or narrow the source wording; the candidate was not ranked."}
			return DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
		}
		candidate.Summary.Gaps[i] = clean
	}
	candidate.Summary.Gaps = uniqueSortedStrings(candidate.Summary.Gaps)
	for i := range candidate.Capabilities {
		for j, gap := range candidate.Capabilities[i].Gaps {
			clean, changed := sanitizePromptString(gap, budget.MaxTextRunes)
			if changed {
				diagnostic := Diagnostic{Severity: "error", Code: "candidate.rationale_unsafe", Message: "candidate capability explanation contains unsafe controls or exceeds the display text budget", Path: candidate.Source.Selector, Remediation: "Review or narrow the source wording; the candidate was not ranked."}
				return DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
			}
			candidate.Capabilities[i].Gaps[j] = clean
		}
		candidate.Capabilities[i].Gaps = uniqueSortedStrings(candidate.Capabilities[i].Gaps)
	}
	for i, reason := range candidate.Effect.Reasons {
		clean, changed := sanitizePromptString(reason, budget.MaxTextRunes)
		if changed {
			diagnostic := Diagnostic{Severity: "error", Code: "candidate.rationale_unsafe", Message: "candidate effect explanation contains unsafe controls or exceeds the display text budget", Path: candidate.Source.Selector}
			return DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
		}
		candidate.Effect.Reasons[i] = clean
	}
	candidate.Effect.Reasons = uniqueSortedStrings(candidate.Effect.Reasons)
	return nil
}

func openAPIOperationSelector(method, path string) string {
	return "#/paths/" + escapeJSONPointer(strings.TrimSpace(path)) + "/" + strings.ToLower(strings.TrimSpace(method))
}

func escapeJSONPointer(value string) string {
	value = strings.ReplaceAll(value, "~", "~0")
	return strings.ReplaceAll(value, "/", "~1")
}

func openAPISourceCapabilities() []OperationCapability {
	return []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilitySupported},
		{Dimension: "outputs", Status: OperationCapabilitySupported},
		{Dimension: "auth", Status: OperationCapabilitySupported},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"Read/write effect is inferred only from operation meaning; the HTTP method is not evidence."}},
	}
}

func openAPICapabilities(operation OperationSummary, selector string) []OperationCapability {
	input, output := OperationCapability{
		Dimension: "inputs", Status: OperationCapabilitySupported,
		Evidence: []OperationEvidence{{Kind: "operation.inputs", Reference: selector}},
	}, OperationCapability{
		Dimension: "outputs", Status: OperationCapabilitySupported,
		Evidence: []OperationEvidence{{Kind: "operation.outputs", Reference: selector}},
	}
	if operation.ResponseBody == nil || operation.ResponseBody.Schema == nil && len(operation.ResponseBody.Fields) == 0 {
		output.Status = OperationCapabilityPartial
		output.Gaps = append(output.Gaps, "The operation has no inspectable successful response schema.")
	}
	for _, issue := range operation.ReadinessIssues {
		if issue.Code == "schema.response_nullable" {
			// Nullable values are represented on the selected response fields.
			// An unrelated nullable sibling must not make this capability partial.
			continue
		}
		if strings.HasPrefix(issue.Code, "schema.") {
			input.Status = OperationCapabilityPartial
			output.Status = OperationCapabilityPartial
			input.Gaps = append(input.Gaps, issue.Message)
			output.Gaps = append(output.Gaps, issue.Message)
		}
	}
	auth := OperationCapability{Dimension: "auth", Status: OperationCapabilitySupported, Evidence: []OperationEvidence{{Kind: "openapi.security", Reference: selector}}}
	for _, alternative := range operation.SecurityRequirementSets {
		for _, requirement := range alternative.Requirements {
			if strings.TrimSpace(requirement.Name) == "" || strings.TrimSpace(requirement.Type) == "" {
				auth.Status = OperationCapabilityPartial
				auth.Gaps = append(auth.Gaps, "An OpenAPI security requirement references missing or incomplete security scheme metadata.")
			}
		}
	}
	effect := OperationCapability{Dimension: "effect", Status: OperationCapabilityPartial, Evidence: []OperationEvidence{{Kind: "operation.meaning", Reference: selector}}, Gaps: []string{"Textual operation meaning does not prove runtime effects or authorization."}}
	return []OperationCapability{input, output, auth, effect}
}
