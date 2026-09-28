package apitools

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
)

var negativeEffectContraction = regexp.MustCompile(`(?i)n['’]t\b`)

type effectSignal struct {
	effect    OperationEffect
	kind      string
	reference string
	word      string
}

func assessOperationEffect(operation OperationSummary, selector string) EffectAssessment {
	signals, conflict := semanticSignalsForOperation(operation, selector)
	return resolveEffectSignals(signals, conflict)
}

func assessOperationEffectWithNative(operation OperationSummary, selector string, native []effectSignal) EffectAssessment {
	textSignals, textConflict := semanticSignalsForOperation(operation, selector)
	if len(native) > 0 {
		filtered := textSignals[:0]
		for _, signal := range textSignals {
			if signal.kind == "operation.operation_id" && signal.word == "unrecognized-leading-action" {
				continue
			}
			filtered = append(filtered, signal)
		}
		textSignals = filtered
	}
	all := append(append([]effectSignal(nil), textSignals...), native...)
	return resolveEffectSignals(all, textConflict)
}

func semanticSignalsForOperation(operation OperationSummary, selector string) ([]effectSignal, bool) {
	ref := strings.TrimSpace(selector)
	if ref == "" {
		ref = firstNonEmpty(operation.OperationID, operation.ID, operation.Path)
	}
	var signals []effectSignal
	conflict := false
	for _, value := range []struct {
		text string
		kind string
	}{
		{operation.OperationID, "operation.operation_id"},
		{operation.Summary, "operation.summary"},
		{operation.Description, "operation.description"},
	} {
		found, ambiguous := semanticEffectSignals(value.text, value.kind, ref)
		signals = append(signals, found...)
		conflict = conflict || ambiguous
	}
	return signals, conflict
}

func resolveEffectSignals(signals []effectSignal, conflictingText bool) EffectAssessment {
	assessment := EffectAssessment{Class: OperationEffectUnknown}
	read, write := false, false
	for _, signal := range signals {
		if signal.effect == OperationEffectRead {
			read = true
		} else if signal.effect == OperationEffectWrite {
			write = true
		}
		if signal.kind != "" && signal.reference != "" {
			assessment.Evidence = append(assessment.Evidence, OperationEvidence{Kind: signal.kind, Reference: signal.reference})
		}
	}
	sort.SliceStable(assessment.Evidence, func(i, j int) bool {
		if assessment.Evidence[i].Kind != assessment.Evidence[j].Kind {
			return assessment.Evidence[i].Kind < assessment.Evidence[j].Kind
		}
		return assessment.Evidence[i].Reference < assessment.Evidence[j].Reference
	})
	if len(assessment.Evidence) > 1 {
		unique := assessment.Evidence[:1]
		for _, item := range assessment.Evidence[1:] {
			if item != unique[len(unique)-1] {
				unique = append(unique, item)
			}
		}
		assessment.Evidence = unique
	}
	if conflictingText || (read && write) {
		for _, signal := range signals {
			if signal.word == "text-over-budget" {
				assessment.Reasons = []string{"Operation meaning text exceeds the bounded effect-analysis budget; its effect remains unknown."}
				return assessment
			}
		}
		assessment.Reasons = []string{"Source meaning evidence is incomplete, negated, or conflicting; the operation effect is unknown."}
		return assessment
	}
	for _, signal := range signals {
		if signal.word == "unrecognized-leading-action" {
			assessment.Reasons = []string{"The leading documented action is not recognized; the operation effect remains unknown."}
			return assessment
		}
	}
	if read {
		assessment.Class = OperationEffectRead
		word := firstEffectWord(signals, OperationEffectRead)
		assessment.Reasons = []string{fmt.Sprintf("Operation metadata contains the read-style action %q; the HTTP method was not used as proof.", word)}
		return assessment
	}
	if write {
		assessment.Class = OperationEffectWrite
		word := firstEffectWord(signals, OperationEffectWrite)
		assessment.Reasons = []string{fmt.Sprintf("Operation metadata contains the mutating action %q; the HTTP method was not used as proof.", word)}
		return assessment
	}
	assessment.Reasons = []string{"The available operation meaning does not establish whether this operation reads or writes."}
	return assessment
}

func firstEffectWord(signals []effectSignal, effect OperationEffect) string {
	for _, signal := range signals {
		if signal.effect == effect && signal.word != "" {
			return signal.word
		}
	}
	return string(effect)
}

func semanticEffectSignals(value, evidenceKind, reference string) ([]effectSignal, bool) {
	if exceedsEffectTextBudget(value) {
		return []effectSignal{{kind: evidenceKind, reference: reference, word: "text-over-budget"}}, true
	}
	if evidenceKind == "operation.operation_id" && strings.HasPrefix(reference, "#/methods/") {
		// Google Discovery method IDs are service/resource/method names. Only
		// the final component names the operation action.
		if separator := strings.LastIndex(value, "."); separator >= 0 {
			value = value[separator+1:]
		}
	}
	tokens := effectWords(value)
	if len(tokens) == 0 {
		return nil, false
	}
	if evidenceKind == "operation.operation_id" {
		if _, ok := effectForVerb(tokens[0]); !ok {
			return []effectSignal{{kind: evidenceKind, reference: reference, word: "unrecognized-leading-action"}}, false
		}
		var signals []effectSignal
		read, write := false, false
		first := -1
		indices := make([]bool, len(tokens))
		for i, token := range tokens {
			_, ok := effectForVerb(token)
			if !ok {
				continue
			}
			if effectNegated(tokens, i) {
				return []effectSignal{{kind: evidenceKind, reference: reference}}, true
			}
			if first < 0 {
				first = i
			}
			if i > 0 && isEffectConnector(tokens[i-1]) {
				indices[i] = true
			}
		}
		if first >= 0 {
			indices[first] = true
		}
		for index, include := range indices {
			if !include {
				continue
			}
			token := tokens[index]
			effect, ok := effectForVerb(token)
			if !ok {
				continue
			}
			signals = append(signals, effectSignal{effect: effect, kind: evidenceKind, reference: reference, word: token})
			read = read || effect == OperationEffectRead
			write = write || effect == OperationEffectWrite
		}
		return signals, read && write
	}

	start := 0
	for start < len(tokens) && isEffectLeadIn(tokens[start]) {
		start++
	}
	if start >= len(tokens) {
		return nil, false
	}
	if _, ok := effectForVerb(tokens[start]); !ok {
		return []effectSignal{{kind: evidenceKind, reference: reference, word: "unrecognized-leading-action"}}, false
	}
	var signals []effectSignal
	for i := start; i < len(tokens) && i < start+5; i++ {
		effect, ok := effectForVerb(tokens[i])
		if !ok {
			continue
		}
		if effectNegated(tokens, i) {
			return []effectSignal{{kind: evidenceKind, reference: reference}}, true
		}
		signals = append(signals, effectSignal{effect: effect, kind: evidenceKind, reference: reference, word: tokens[i]})
		// Once the leading action is known, inspect the rest of the bounded
		// wording for conflicting actions. A connector can be separated from
		// its verb by objects or modifiers, and prose can use two sentences.
		for j := i + 1; j < len(tokens); j++ {
			nextEffect, nextOK := effectForVerb(tokens[j])
			if !nextOK {
				continue
			}
			if effectNegated(tokens, j) {
				return []effectSignal{{kind: evidenceKind, reference: reference}}, true
			}
			signals = append(signals, effectSignal{effect: nextEffect, kind: evidenceKind, reference: reference, word: tokens[j]})
		}
		break
	}
	read, write := false, false
	for _, signal := range signals {
		read = read || signal.effect == OperationEffectRead
		write = write || signal.effect == OperationEffectWrite
	}
	return signals, read && write
}

func exceedsEffectTextBudget(value string) bool {
	count := 0
	for range value {
		count++
		if count > DefaultPromptTextRunes {
			return true
		}
	}
	return false
}

func effectWords(value string) []string {
	if exceedsEffectTextBudget(value) {
		return nil
	}
	value = negativeEffectContraction.ReplaceAllString(value, " not")
	runes := []rune(value)
	var tokens []string
	var token []rune
	flush := func() {
		if len(token) > 0 {
			tokens = append(tokens, strings.ToLower(string(token)))
			token = token[:0]
		}
	}
	for i, current := range runes {
		if !unicode.IsLetter(current) && !unicode.IsDigit(current) {
			flush()
			continue
		}
		if len(token) > 0 {
			previous := runes[i-1]
			nextLower := i+1 < len(runes) && unicode.IsLower(runes[i+1])
			camelBoundary := unicode.IsLower(previous) && unicode.IsUpper(current)
			acronymBoundary := unicode.IsUpper(previous) && unicode.IsUpper(current) && nextLower
			digitBoundary := unicode.IsDigit(previous) != unicode.IsDigit(current)
			if camelBoundary || acronymBoundary || digitBoundary {
				flush()
			}
		}
		token = append(token, current)
	}
	flush()
	return tokens
}

func effectNegated(tokens []string, action int) bool {
	start := action - 3
	if start < 0 {
		start = 0
	}
	for i := start; i < action; i++ {
		switch tokens[i] {
		case "not", "never", "no", "without", "cannot":
			return true
		}
	}
	return false
}

func effectForVerb(value string) (OperationEffect, bool) {
	switch strings.ToLower(value) {
	case "get", "gets", "retrieve", "retrieves", "retrieved", "list", "lists", "listed", "search", "searches", "searched", "find", "finds", "found", "fetch", "fetches", "fetched", "read", "reads", "query", "queries", "queried", "describe", "describes", "described", "check", "checks", "checked", "verify", "verifies", "verified", "view", "views", "show", "shows", "lookup":
		return OperationEffectRead, true
	case "create", "creates", "created", "add", "adds", "added", "insert", "inserts", "inserted", "update", "updates", "updated", "modify", "modifies", "modified", "patch", "patches", "patched", "put", "puts", "post", "posts", "posted", "write", "writes", "wrote", "delete", "deletes", "deleted", "remove", "removes", "removed", "destroy", "destroys", "destroyed", "send", "sends", "sent", "submit", "submits", "submitted", "execute", "executes", "executed", "run", "runs", "trigger", "triggers", "triggered", "cancel", "cancels", "cancelled", "canceled", "publish", "publishes", "published", "revoke", "revokes", "revoked", "rotate", "rotates", "rotated", "enable", "enables", "enabled", "disable", "disables", "disabled", "approve", "approves", "approved", "reject", "rejects", "rejected", "subscribe", "subscribes", "subscribed":
		return OperationEffectWrite, true
	default:
		return OperationEffectUnknown, false
	}
}

func isEffectLeadIn(word string) bool {
	switch word {
	case "a", "an", "the", "this", "that", "operation", "method", "endpoint", "api", "it", "will", "can", "to", "for", "used", "use", "allows", "allow":
		return true
	default:
		return false
	}
}

func isEffectConnector(word string) bool {
	switch word {
	case "and", "or", "then", "but":
		return true
	default:
		return false
	}
}
