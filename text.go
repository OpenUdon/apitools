package apitools

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

func sortResults(results []Result) {
	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Score != results[j].Score {
			return results[i].Score > results[j].Score
		}
		if results[i].Title != results[j].Title {
			return results[i].Title < results[j].Title
		}
		return results[i].SpecURL < results[j].SpecURL
	})
}

// ScoreText returns a deterministic relevance score for query terms in text.
// Higher scores indicate stronger textual overlap.
func ScoreText(query, text string) int {
	return scoreText(query, text)
}

func scoreText(query, text string) int {
	query = strings.ToLower(strings.TrimSpace(query))
	text = strings.ToLower(text)
	if query == "" || text == "" {
		return 0
	}
	score := 0
	if strings.Contains(text, query) {
		score += 10
	}
	for _, token := range tokenPattern.FindAllString(query, -1) {
		if len(token) < 2 {
			continue
		}
		if strings.Contains(text, token) {
			score += len(token)
		}
	}
	return score
}

var tokenPattern = regexp.MustCompile(`[a-z0-9]+`)

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func nonEmptyStrings(values ...string) []string {
	var out []string
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			out = append(out, strings.TrimSpace(value))
		}
	}
	return out
}

// looksLikeCredentialName reports whether name looks like a field or path that
// holds credential material. Inventory rendering uses this denylist to omit
// credential-shaped request fields from prompt-safe output before summaries are
// sent to LLM-backed authoring flows. It recognizes credential words across
// separator and camel-case boundaries while treating ambiguous business-key
// names such as "projectKey" as data, not credentials.
func looksLikeCredentialName(name string) bool {
	lower := strings.ToLower(name)
	if strings.Contains(lower, "api_key") || strings.Contains(lower, "apikey") || strings.Contains(lower, "api-key") {
		return true
	}
	type namePart struct {
		value               string
		camelBoundaryBefore bool
	}
	var parts []namePart
	runes := []rune(name)
	start := 0
	partCamelBoundary := false
	flush := func(end int, nextPartCamelBoundary bool) {
		if end > start {
			parts = append(parts, namePart{
				value:               strings.ToLower(string(runes[start:end])),
				camelBoundaryBefore: partCamelBoundary,
			})
		}
		start = end
		partCamelBoundary = nextPartCamelBoundary
	}
	for i, current := range runes {
		if current == '.' || current == '_' || current == '-' || current == '[' || current == ']' {
			flush(i, false)
			start = i + 1
			partCamelBoundary = false
			continue
		}
		if i == start {
			continue
		}
		previous := runes[i-1]
		nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
		camelBoundary := unicode.IsUpper(current) && (unicode.IsLower(previous) || unicode.IsDigit(previous) || unicode.IsUpper(previous) && nextLower)
		digitBoundary := unicode.IsDigit(current) != unicode.IsDigit(previous)
		if camelBoundary || digitBoundary {
			flush(i, camelBoundary)
		}
	}
	flush(len(runes), false)
	for i, part := range parts {
		switch part.value {
		case "secret", "password", "passwd", "pwd", "token", "authorization", "auth", "credential", "credentials":
			return true
		case "key":
			if !part.camelBoundaryBefore || (i > 0 && isCredentialKeyContext(parts[i-1].value)) || (i+1 < len(parts) && isCredentialKeyContext(parts[i+1].value)) {
				return true
			}
		}
	}
	return false
}

func isCredentialKeyContext(part string) bool {
	switch part {
	case "access", "api", "auth", "bearer", "client", "consumer", "credential", "credentials", "encryption", "master", "oauth", "private", "public", "refresh", "secret", "session", "signing", "subscription", "webhook":
		return true
	default:
		return false
	}
}
