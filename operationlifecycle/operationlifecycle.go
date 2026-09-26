package operationlifecycle

import (
	"net/url"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"unicode"

	"github.com/OpenUdon/apitools"
)

const (
	// MinimumSiblingScore is the minimum score required to select a lifecycle
	// sibling.
	MinimumSiblingScore = 60
	// HTTPVerbScore rewards an HTTP method that directly represents the role.
	HTTPVerbScore = 45
	// OperationNameScore rewards an operation name that names the role.
	OperationNameScore = 35
	// PatchMethodPreferenceScore prefers PATCH among otherwise valid updates.
	PatchMethodPreferenceScore = 8
	// PatchNamePreferenceScore rewards update operations explicitly named patch.
	PatchNamePreferenceScore = 6
	// FamilyMatchScore rewards operation-family agreement with the seed.
	FamilyMatchScore = 30
	// PathMatchScore rewards a matching collection/item path pair.
	PathMatchScore = 25
	// HighConfidenceScore is the minimum score reported as high confidence.
	HighConfidenceScore = 90
	// MediumConfidenceScore is the minimum score reported as medium confidence.
	MediumConfidenceScore = 70
)

// Options configures lifecycle expansion.
type Options struct {
	Goal         string
	DesiredState bool
}

// Expansion is one conservative lifecycle-role expansion.
type Expansion struct {
	SeedOperationID string          `json:"seed_operation_id,omitempty"`
	FamilyKey       string          `json:"family_key,omitempty"`
	Roles           []RoleCandidate `json:"roles,omitempty"`
	Diagnostics     []Diagnostic    `json:"diagnostics,omitempty"`
}

// RoleCandidate associates one API operation with a lifecycle role.
type RoleCandidate struct {
	Role       string                    `json:"role"`
	Operation  apitools.OperationSummary `json:"operation"`
	Confidence string                    `json:"confidence,omitempty"`
	Reason     string                    `json:"reason,omitempty"`
}

// Diagnostic explains why a possible lifecycle role was omitted.
type Diagnostic struct {
	Code     string `json:"code,omitempty"`
	Severity string `json:"severity,omitempty"`
	Message  string `json:"message,omitempty"`
}

type candidateScore struct {
	operation apitools.OperationSummary
	score     int
	reason    string
}

// Expand returns a conservative lifecycle expansion for seed. It only selects
// sibling operations when one same-source candidate is clearly stronger for a
// role; ambiguous role matches are diagnosed and omitted.
func Expand(operations []apitools.OperationSummary, seed apitools.OperationSummary, opts Options) Expansion {
	var seedDiagnostic Diagnostic
	seed, seedDiagnostic = normalizeSeed(operations, seed)
	seedID := operationID(seed)
	out := Expansion{SeedOperationID: seedID, FamilyKey: familyKey(seed)}
	if seedDiagnostic.Code != "" {
		out.Diagnostics = append(out.Diagnostics, seedDiagnostic)
		return out
	}
	if seedID == "" {
		out.Diagnostics = append(out.Diagnostics, Diagnostic{Code: "operation_lifecycle.seed_missing", Severity: "warning", Message: "seed operation is empty"})
		return out
	}
	seedPurpose := seedPrimaryPurpose(seed, opts.DesiredState)
	siblings := map[string]candidateScore{}
	for _, role := range []string{"read", "update", "delete"} {
		if role == seedPurpose {
			continue
		}
		if role == "update" && !goalWantsUpdate(opts.Goal) {
			continue
		}
		match, ok, diag := bestSibling(operations, seed, role)
		if diag.Code != "" {
			out.Diagnostics = append(out.Diagnostics, diag)
		}
		if ok {
			siblings[role] = match
		}
	}
	primaryRole := primaryRole(seed, opts, len(siblings) > 0)
	out.Roles = append(out.Roles, RoleCandidate{Role: primaryRole, Operation: seed, Confidence: "high", Reason: "selected seed operation"})
	if primaryRole == "read" {
		return out
	}
	for _, role := range []string{"read", "update", "delete"} {
		if role == primaryRole {
			continue
		}
		match, ok := siblings[role]
		if !ok {
			continue
		}
		out.Roles = append(out.Roles, RoleCandidate{Role: role, Operation: match.operation, Confidence: confidence(match.score), Reason: match.reason})
	}
	return out
}

func normalizeSeed(operations []apitools.OperationSummary, seed apitools.OperationSummary) (apitools.OperationSummary, Diagnostic) {
	seedID := operationID(seed)
	if seedID == "" {
		return seed, Diagnostic{}
	}
	var matches []apitools.OperationSummary
	for _, operation := range operations {
		if operationID(operation) == seedID {
			matches = append(matches, operation)
		}
	}
	if hasSourceIdentity(seed) {
		var sourceMatches []apitools.OperationSummary
		for _, operation := range matches {
			if sameSource(operation, seed) {
				sourceMatches = append(sourceMatches, operation)
			}
		}
		sourceMatches = narrowSeedMatchesByMethodAndPath(sourceMatches, seed)
		if len(sourceMatches) == 1 {
			return sourceMatches[0], Diagnostic{}
		}
		if len(sourceMatches) > 1 {
			return seed, Diagnostic{Code: "operation_lifecycle.seed_ambiguous", Severity: "warning", Message: "multiple operations with the seed identity are present in the same source"}
		}
		if len(matches) > 0 {
			return seed, Diagnostic{Code: "operation_lifecycle.seed_source_mismatch", Severity: "warning", Message: "the seed operation ID was found only in different document sources"}
		}
	}
	matches = narrowSeedMatchesByMethodAndPath(matches, seed)
	if len(matches) == 1 {
		return matches[0], Diagnostic{}
	}
	if len(matches) > 1 {
		return seed, Diagnostic{Code: "operation_lifecycle.seed_ambiguous", Severity: "warning", Message: "the seed operation ID occurs in multiple document sources"}
	}
	if strings.TrimSpace(seed.Method) == "" || strings.TrimSpace(seed.Path) == "" {
		return seed, Diagnostic{Code: "operation_lifecycle.seed_not_found", Severity: "warning", Message: "an operation-ID-only seed did not match a supplied operation"}
	}
	return seed, Diagnostic{}
}

// narrowSeedMatchesByMethodAndPath disambiguates same-ID, same-source
// candidates using the seed's own method and path when it carries them. A
// seed built from a real operation (rather than bare metadata) resolves to
// its exact match instead of being declared ambiguous merely because its
// operation ID recurs in the same document (for example, one route action
// among several sharing a base resource name).
func narrowSeedMatchesByMethodAndPath(candidates []apitools.OperationSummary, seed apitools.OperationSummary) []apitools.OperationSummary {
	method := strings.TrimSpace(seed.Method)
	seedPath := strings.TrimSpace(seed.Path)
	if method == "" || seedPath == "" || len(candidates) <= 1 {
		return candidates
	}
	var exact []apitools.OperationSummary
	for _, candidate := range candidates {
		if strings.EqualFold(strings.TrimSpace(candidate.Method), method) && strings.TrimSpace(candidate.Path) == seedPath {
			exact = append(exact, candidate)
		}
	}
	if len(exact) > 0 {
		return exact
	}
	return candidates
}

func primaryRole(seed apitools.OperationSummary, opts Options, expanded bool) string {
	if opts.DesiredState && expanded {
		if purpose := seedPrimaryPurpose(seed, opts.DesiredState); purpose != "" {
			return purpose
		}
	}
	return methodRole(seed)
}

// seedPrimaryPurpose determines the seed's own lifecycle role when the caller
// wants a desired-state interpretation. It is shared by Expand's sibling-role
// skip decision and primaryRole's final label, so both agree: a create-or-update
// PUT is "create" in both places, and Expand does not skip searching for an
// update sibling merely because ClassifyOperationPurpose's raw, method-based
// answer for PUT is "update".
func seedPrimaryPurpose(seed apitools.OperationSummary, desiredState bool) string {
	if desiredState && strings.EqualFold(seed.Method, "PUT") && operationHasAny(seed, "create", "createorupdate", "insert") {
		return "create"
	}
	return lifecyclePurpose(seed)
}

func lifecyclePurpose(operation apitools.OperationSummary) string {
	purpose := apitools.ClassifyOperationPurpose(operation, apitools.OperationSelectionHints{})
	if !strings.EqualFold(operation.Method, "POST") || purpose != "create" {
		return purpose
	}
	// Explicit create verbs in the operation ID itself outrank a summary or
	// tag that happens to mention update/patch wording; only the operation ID
	// is checked here; operationNameMatchesRole (below) also considers
	// summary/tag text for the update fallback.
	if operationIDHasAny(operation, "createorupdate") {
		return "create"
	}
	if operationIDHasAny(operation, "create", "insert", "add") {
		return "create"
	}
	if operationNameMatchesRole(operation, "update") {
		return "update"
	}
	if postOperationIsAction(operation) {
		return methodRole(operation)
	}
	return purpose
}

func postOperationIsAction(operation apitools.OperationSummary) bool {
	if operationHasAny(operation,
		"activate", "acknowledge", "approve", "archive", "cancel", "close",
		"complete", "confirm", "deactivate", "disable", "enable", "execute",
		"invoke", "lock", "pause", "publish", "refresh", "reject", "release",
		"reopen", "reset", "restore", "resume", "restart", "retry", "rotate",
		"run", "send", "start", "stop", "submit", "trigger", "unlock",
		"validate",
	) {
		return true
	}
	parts := pathSegments(normalizePath(operation.Path, isGoogleDiscovery(operation)))
	if len(parts) > 0 && strings.Contains(parts[len(parts)-1], ":") {
		return true
	}
	// A POST directly against a parameterized item (the route ends in a path
	// parameter, such as POST /widgets/{id}) is not sufficient evidence of
	// resource creation, so it stays a generic POST rather than an invented
	// create disposition. A parameter earlier in the path only scopes a
	// parent resource (for example the {projectId} in POST
	// /projects/{projectId}/children, a nested-collection create) and does
	// not by itself make the operation an item-level action.
	if len(parts) > 0 && isPathParameter(parts[len(parts)-1]) {
		return true
	}
	return false
}

func methodRole(op apitools.OperationSummary) string {
	switch strings.ToUpper(strings.TrimSpace(op.Method)) {
	case "GET", "HEAD":
		return "read"
	case "DELETE":
		return "delete"
	case "POST":
		return "post"
	case "PUT":
		return "put"
	case "PATCH":
		return "patch"
	default:
		return "post"
	}
}

func bestSibling(operations []apitools.OperationSummary, seed apitools.OperationSummary, role string) (candidateScore, bool, Diagnostic) {
	var matches []candidateScore
	for _, operation := range operations {
		if sameOperation(operation, seed) || !sameSource(operation, seed) {
			continue
		}
		if lifecyclePurpose(operation) != role {
			continue
		}
		score, reason := siblingScore(operation, seed, role)
		if score < MinimumSiblingScore {
			continue
		}
		matches = append(matches, candidateScore{operation: operation, score: score, reason: reason})
	}
	if len(matches) == 0 {
		return candidateScore{}, false, Diagnostic{}
	}
	slices.SortStableFunc(matches, func(a, b candidateScore) int {
		if a.score != b.score {
			return b.score - a.score
		}
		return strings.Compare(operationID(a.operation), operationID(b.operation))
	})
	if len(matches) > 1 && matches[0].score == matches[1].score {
		return candidateScore{}, false, Diagnostic{Code: "operation_lifecycle.ambiguous_" + role, Severity: "warning", Message: "multiple same-source operations match lifecycle role " + role}
	}
	return matches[0], true, Diagnostic{}
}

func siblingScore(op, seed apitools.OperationSummary, role string) (int, string) {
	if !lifecyclePathsMatch(seed, op, role) {
		return 0, ""
	}
	score := 0
	var reasons []string
	if verbMatchesRole(op.Method, role) {
		score += HTTPVerbScore
		reasons = append(reasons, "HTTP verb matches "+role)
	}
	if operationNameMatchesRole(op, role) {
		score += OperationNameScore
		reasons = append(reasons, "operation id/name matches "+role)
	}
	if role == "update" && strings.EqualFold(op.Method, "PATCH") {
		score += PatchMethodPreferenceScore
		reasons = append(reasons, "PATCH update preferred")
	}
	if role == "update" && operationHasAny(op, "patch") {
		score += PatchNamePreferenceScore
		reasons = append(reasons, "patch operation preferred")
	}
	if sameFamily(seed, op) {
		score += FamilyMatchScore
		reasons = append(reasons, "operation family matches seed")
	}
	score += PathMatchScore
	reasons = append(reasons, "collection/item path matches seed")
	return score, strings.Join(reasons, "; ")
}

func verbMatchesRole(method, role string) bool {
	switch role {
	case "read":
		return strings.EqualFold(method, "GET") || strings.EqualFold(method, "HEAD")
	case "update":
		return strings.EqualFold(method, "PUT") || strings.EqualFold(method, "PATCH")
	case "delete":
		return strings.EqualFold(method, "DELETE")
	default:
		return false
	}
}

func operationNameMatchesRole(op apitools.OperationSummary, role string) bool {
	switch role {
	case "read":
		return operationHasAny(op, "get", "read", "show", "describe")
	case "update":
		return operationHasAny(op, "update", "patch", "replace", "createorupdate", "put")
	case "delete":
		return operationHasAny(op, "delete", "remove")
	default:
		return false
	}
}

func operationHasAny(op apitools.OperationSummary, terms ...string) bool {
	tokens := operationTokens(op)
	for _, term := range terms {
		if tokens[strings.ToLower(term)] {
			return true
		}
	}
	return false
}

// operationIDHasAny is like operationHasAny but tokenizes only the operation
// ID, not its summary or tags. Use it where an explicit verb named directly on
// the operation must outrank free-text wording elsewhere on the operation.
func operationIDHasAny(op apitools.OperationSummary, terms ...string) bool {
	tokens := operationIDTokens(op)
	for _, term := range terms {
		if tokens[strings.ToLower(term)] {
			return true
		}
	}
	return false
}

// sourceIdentity holds every document-identity signal an operation summary
// may carry, in priority order (absolute path, then URL, then relative path
// or name-like fallback). A field is empty when the operation does not carry
// that signal at all, which is distinct from two operations disagreeing on it.
type sourceIdentity struct {
	absPath string
	url     string
	relPath string
}

func sourceIdentityOf(op apitools.OperationSummary) sourceIdentity {
	var id sourceIdentity
	if documentPath := strings.TrimSpace(op.DocumentPath); filepath.IsAbs(documentPath) {
		id.absPath = filepath.Clean(documentPath)
	}
	if documentURL := absoluteURLIdentity(op.DocumentURL); documentURL != "" {
		id.url = documentURL
	} else if documentURL := absoluteURLIdentity(op.DocumentPath); documentURL != "" {
		id.url = documentURL
	}
	nonAbsPath := ""
	if documentPath := strings.TrimSpace(op.DocumentPath); documentPath != "" && !filepath.IsAbs(documentPath) {
		nonAbsPath = documentPath
	}
	id.relPath = firstNonEmpty(op.DocumentRelativePath, nonAbsPath, op.DocumentName)
	return id
}

func (id sourceIdentity) isEmpty() bool {
	return id.absPath == "" && id.url == "" && id.relPath == ""
}

// sameSourceIdentity reports whether a and b identify the same document
// source. It compares the highest-priority identity kind that both sides have
// populated; a disagreement at that kind is decisive even when a
// lower-priority kind happens to agree, which is what keeps two different
// documents that share a relative path or name from being treated as one
// source. When one side lacks a higher-priority kind entirely, comparison
// falls through to the next kind, so an under-specified seed (for example,
// one supplied with only a relative path) can still resolve against a fully
// described operation from the same document.
func sameSourceIdentity(a, b sourceIdentity) bool {
	if a.isEmpty() || b.isEmpty() {
		return false
	}
	if a.absPath != "" && b.absPath != "" {
		return a.absPath == b.absPath
	}
	if a.url != "" && b.url != "" {
		return a.url == b.url
	}
	if a.relPath != "" && b.relPath != "" {
		return a.relPath == b.relPath
	}
	return false
}

func hasSourceIdentity(op apitools.OperationSummary) bool {
	return !sourceIdentityOf(op).isEmpty()
}

func sameSource(a, b apitools.OperationSummary) bool {
	return sameSourceIdentity(sourceIdentityOf(a), sourceIdentityOf(b))
}

func absoluteURLIdentity(value string) string {
	parsed, err := url.Parse(strings.TrimSpace(value))
	if err != nil || !parsed.IsAbs() || parsed.Host == "" {
		return ""
	}
	return parsed.String()
}

func sameOperation(a, b apitools.OperationSummary) bool {
	aID, bID := operationID(a), operationID(b)
	return aID != "" && aID == bID && sameSource(a, b)
}

func operationID(op apitools.OperationSummary) string {
	return strings.TrimSpace(firstNonEmpty(op.OperationID, op.ID))
}

func familyKey(op apitools.OperationSummary) string { return strings.Join(familyTokens(op), ".") }

func sameFamily(a, b apitools.OperationSummary) bool {
	aTokens, bTokens := familyTokens(a), familyTokens(b)
	if len(aTokens) == 0 || len(bTokens) == 0 {
		return false
	}
	bSet := map[string]bool{}
	for _, token := range bTokens {
		bSet[token] = true
	}
	common := 0
	for _, token := range aTokens {
		if bSet[token] {
			common++
		}
	}
	return common >= min(2, min(len(aTokens), len(bTokens)))
}

func familyTokens(op apitools.OperationSummary) []string {
	raw := wordTokenSet(operationID(op) + " " + op.Path)
	var out []string
	for token := range raw {
		if lifecycleWord(token) || familyStopWord(token) {
			continue
		}
		out = append(out, token)
	}
	slices.Sort(out)
	return out
}

func operationTokens(op apitools.OperationSummary) map[string]bool {
	text := operationID(op) + " " + op.Summary + " " + strings.Join(op.Tags, " ")
	out := wordTokenSet(text)
	joined := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(operationID(op)))
	if strings.Contains(joined, "createorupdate") {
		out["createorupdate"], out["create"], out["update"] = true, true, true
	}
	return out
}

// operationIDTokens is operationTokens restricted to the operation ID, with no
// summary or tag text mixed in.
func operationIDTokens(op apitools.OperationSummary) map[string]bool {
	out := wordTokenSet(operationID(op))
	joined := strings.ToLower(strings.NewReplacer("_", "", "-", "", ".", "").Replace(operationID(op)))
	if strings.Contains(joined, "createorupdate") {
		out["createorupdate"], out["create"], out["update"] = true, true, true
	}
	return out
}

func wordTokenSet(text string) map[string]bool {
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	out := map[string]bool{}
	for _, part := range parts {
		for _, token := range splitCamel(part) {
			token = strings.ToLower(strings.TrimSpace(token))
			if token != "" {
				out[token] = true
			}
		}
	}
	return out
}

func familyStopWord(token string) bool {
	switch token {
	case "api", "apis", "endpoint", "endpoints", "entity", "entities",
		"id", "ids", "item", "items", "name", "object", "objects",
		"operation", "operations", "request", "requests", "resource",
		"resources", "response", "responses", "service", "services",
		"v1", "v2", "v3", "v4":
		return true
	default:
		return false
	}
}

func splitCamel(value string) []string {
	var out []string
	start := 0
	runes := []rune(value)
	for i := 1; i < len(runes); i++ {
		if unicode.IsUpper(runes[i]) && (unicode.IsLower(runes[i-1]) || (i+1 < len(runes) && unicode.IsLower(runes[i+1]))) {
			out = append(out, string(runes[start:i]))
			start = i
		}
	}
	return append(out, string(runes[start:]))
}

func lifecycleWord(token string) bool {
	switch token {
	case "create", "insert", "post", "put", "get", "read", "show", "describe", "list", "update", "patch", "replace", "delete", "remove", "createorupdate":
		return true
	default:
		return false
	}
}

func lifecyclePathsMatch(seed, operation apitools.OperationSummary, role string) bool {
	if role != "read" && role != "update" && role != "delete" {
		return false
	}
	seedPath := normalizePath(seed.Path, isGoogleDiscovery(seed))
	opPath := normalizePath(operation.Path, isGoogleDiscovery(operation))
	if seedPath == "" || opPath == "" {
		return false
	}
	seedCollection, seedIsItem := collectionPath(seedPath, isGoogleDiscovery(seed))
	opCollection, opIsItem := collectionPath(opPath, isGoogleDiscovery(operation))
	if !opIsItem {
		return false
	}
	if seedIsItem && seedCollection == opCollection {
		return true
	}
	if !seedIsItem && seedPath == opCollection {
		return true
	}
	return googleDiscoveryResourcePathMatches(seed, operation, opPath)
}

func isGoogleDiscovery(op apitools.OperationSummary) bool {
	return strings.EqualFold(strings.TrimSpace(op.Extensions["x-uws-source-kind"]), apitools.APISourceKindGoogleDiscovery)
}

func normalizePath(value string, stripDiscoveryUpload bool) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = "/" + strings.Trim(path.Clean("/"+value), "/")
	if stripDiscoveryUpload {
		if value == "/upload" {
			return "/"
		}
		value = strings.TrimPrefix(value, "/upload/")
		if !strings.HasPrefix(value, "/") {
			value = "/" + value
		}
	}
	return value
}

func collectionPath(value string, googleDiscovery bool) (string, bool) {
	parts := pathSegments(value)
	if len(parts) == 0 || !isPathParameter(parts[len(parts)-1]) ||
		(googleDiscovery && parts[len(parts)-1] == "{+parent}") {
		return value, false
	}
	parts = parts[:len(parts)-1]
	if len(parts) == 0 {
		return "/", true
	}
	return "/" + strings.Join(parts, "/"), true
}

func pathSegments(value string) []string {
	value = strings.Trim(value, "/")
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func isPathParameter(segment string) bool {
	return strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") && len(segment) > 2
}

func googleDiscoveryResourcePathMatches(seed, operation apitools.OperationSummary, operationPath string) bool {
	if !isGoogleDiscovery(seed) || !isGoogleDiscovery(operation) ||
		!hasTrailingGoogleNameParameter(operationPath) {
		return false
	}
	seedResource, seedOK := googleDiscoveryMethodResourceIdentity(seed)
	operationResource, operationOK := googleDiscoveryMethodResourceIdentity(operation)
	return seedOK && operationOK && seedResource == operationResource
}

func hasTrailingGoogleNameParameter(value string) bool {
	parts := pathSegments(value)
	return len(parts) > 0 && parts[len(parts)-1] == "{+name}"
}

func googleDiscoveryMethodResourceIdentity(operation apitools.OperationSummary) (string, bool) {
	if !isGoogleDiscovery(operation) {
		return "", false
	}
	id := operationID(operation)
	methodSeparator := strings.LastIndexByte(id, '.')
	if methodSeparator <= 0 || methodSeparator == len(id)-1 {
		return "", false
	}
	return id[:methodSeparator], true
}

func confidence(score int) string {
	switch {
	case score >= HighConfidenceScore:
		return "high"
	case score >= MediumConfidenceScore:
		return "medium"
	default:
		return "low"
	}
}

func goalWantsUpdate(goal string) bool {
	tokens := wordTokenSet(goal)
	for _, word := range []string{"update", "updates", "updated", "patch", "patches", "modify", "modifies", "modifying", "replace", "replaces", "replacing"} {
		if tokens[word] {
			return true
		}
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}
