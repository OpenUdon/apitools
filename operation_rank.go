package apitools

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

const DefaultMaxOperationCandidates = 100

type flattenedStepContract struct {
	contract StepContract
	inputs   map[string]ContractValue
	outputs  map[string]ContractValue
}

type matchAccumulator struct {
	incompatible  bool
	indeterminate bool
}

type sourceValueIndex struct {
	values  []OperationValueSummary
	byAlias map[string][]int
}

func rankOperationCandidates(ctx context.Context, candidates []OperationCandidate, contract StepContract, promptBudget PromptBudget, maxCandidates int) (OperationCandidateReport, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	budget := resolvedPromptBudget(promptBudget)
	report := OperationCandidateReport{SchemaVersion: OperationCandidatesSchemaVersion}
	if err := ctx.Err(); err != nil {
		return report, err
	}
	if maxCandidates < 0 {
		diagnostic := Diagnostic{Severity: "error", Code: "candidate.limit_invalid", Message: "MaxCandidates cannot be negative"}
		report.Diagnostics = []Diagnostic{diagnostic}
		return report, DiagnosticError{Diagnostics: report.Diagnostics}
	}
	if maxCandidates == 0 {
		maxCandidates = DefaultMaxOperationCandidates
	}
	prepared, inputValues, outputValues, diagnostics := prepareStepContract(contract, budget)
	report.Diagnostics = append(report.Diagnostics, diagnostics...)
	if len(errorDiagnostics(diagnostics)) > 0 {
		return report, DiagnosticError{Diagnostics: errorDiagnostics(diagnostics)}
	}
	if len(candidates) > budget.MaxOperations {
		diagnostic := Diagnostic{
			Severity: "error", Code: "candidate.operation_work_budget",
			Message:     fmt.Sprintf("candidate set contains %d operations, exceeding the %d-operation ranking budget", len(candidates), budget.MaxOperations),
			Remediation: "Narrow source roots or increase MaxOperations for a reviewed ranking.",
		}
		report.Diagnostics = append(report.Diagnostics, diagnostic)
		report.Truncated = true
		return report, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	if len(candidates) == 0 {
		diagnostic := Diagnostic{Severity: "warning", Code: "candidate.no_matches", Message: "No operations were available to rank against the step contract."}
		report.Diagnostics = append(report.Diagnostics, diagnostic)
		return report, nil
	}

	ranked := append([]OperationCandidate(nil), candidates...)
	for i := range ranked {
		if i%32 == 0 {
			if err := ctx.Err(); err != nil {
				return report, err
			}
		}
		ranked[i].Match = compareCandidateToContract(ranked[i], prepared, inputValues, outputValues)
	}
	sort.SliceStable(ranked, func(i, j int) bool { return candidateLess(ranked[i], ranked[j]) })
	if len(ranked) > maxCandidates {
		diagnostic := Diagnostic{
			Severity: "error", Code: "candidate.result_limit",
			Message:     fmt.Sprintf("ranked operations exceed the %d-candidate result limit; the returned shortlist is incomplete", maxCandidates),
			Remediation: "Narrow the source roots or increase MaxCandidates before selecting an operation.",
		}
		ranked = ranked[:maxCandidates]
		report.Truncated = true
		report.Diagnostics = append(report.Diagnostics, diagnostic)
	}
	report.Candidates = ranked
	if !hasContractMatch(ranked, prepared.contract) {
		report.Diagnostics = append(report.Diagnostics, Diagnostic{
			Severity: "warning", Code: "candidate.no_compatible_match",
			Message: "No candidate has evidence compatible with every declared step constraint; review the dimension gaps and conflicts.",
		})
	}
	if len(ranked) > 1 && ranked[0].Match.Score > 0 && ranked[0].Match.Score == ranked[1].Match.Score && !hasIncompatibleDimension(ranked[0].Match) && !hasIncompatibleDimension(ranked[1].Match) {
		report.Diagnostics = append(report.Diagnostics, Diagnostic{
			Severity: "warning", Code: "candidate.top_tie",
			Message: "Multiple top-ranked operations have the same advisory score; selection remains ambiguous.",
		})
	}
	if len(report.Diagnostics) > 1 {
		sort.SliceStable(report.Diagnostics, func(i, j int) bool {
			if report.Diagnostics[i].Code != report.Diagnostics[j].Code {
				return report.Diagnostics[i].Code < report.Diagnostics[j].Code
			}
			return report.Diagnostics[i].Path < report.Diagnostics[j].Path
		})
	}
	for _, candidate := range report.Candidates {
		size, err := marshalSize(candidate)
		if err != nil {
			return report, fmt.Errorf("encode candidate operation for budget check: %w", err)
		}
		if size > budget.MaxOperationBytes {
			diagnostic := Diagnostic{
				Severity: "error", Code: "candidate.operation_byte_budget",
				Message:     fmt.Sprintf("ranked candidate exceeds the %d-byte per-operation result budget", budget.MaxOperationBytes),
				Path:        candidate.Source.Selector,
				Remediation: "Narrow source metadata or increase MaxOperationBytes for a reviewed candidate.",
			}
			report.Diagnostics = append(report.Diagnostics, diagnostic)
			report.Truncated = true
			return report, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
		}
	}
	contextSize, err := marshalSize(report)
	if err != nil {
		return report, fmt.Errorf("encode candidate report for budget check: %w", err)
	}
	if contextSize > budget.MaxContextBytes {
		diagnostic := Diagnostic{
			Severity: "error", Code: "candidate.context_byte_budget",
			Message:     fmt.Sprintf("ranked candidate report exceeds the %d-byte prompt-context budget", budget.MaxContextBytes),
			Remediation: "Narrow source roots or increase MaxContextBytes for a reviewed ranking.",
		}
		report.Diagnostics = append(report.Diagnostics, diagnostic)
		report.Truncated = true
		return report, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	if errors := errorDiagnostics(report.Diagnostics); len(errors) > 0 {
		return report, DiagnosticError{Diagnostics: errors}
	}
	return report, nil
}

func prepareStepContract(contract StepContract, budget PromptBudget) (flattenedStepContract, map[string]ContractValue, map[string]ContractValue, []Diagnostic) {
	prepared := flattenedStepContract{contract: contract}
	var diagnostics []Diagnostic
	purpose := strings.TrimSpace(contract.Purpose)
	if safe, changed := sanitizePromptString(purpose, budget.MaxTextRunes); changed {
		diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.purpose_unsafe", Message: "step purpose contains controls or exceeds the prompt text budget", Remediation: "Supply a bounded plain-text step purpose."})
	} else {
		prepared.contract.Purpose = safe
	}
	if contract.Effect != "" && contract.Effect != OperationEffectRead && contract.Effect != OperationEffectWrite && contract.Effect != OperationEffectUnknown {
		diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.effect_invalid", Message: "step effect must be read, write, unknown, or omitted"})
	}
	inputs, inputDiagnostics := flattenContractValues(contract.Inputs, "inputs", budget)
	outputs, outputDiagnostics := flattenContractValues(contract.Outputs, "outputs", budget)
	diagnostics = append(diagnostics, inputDiagnostics...)
	diagnostics = append(diagnostics, outputDiagnostics...)
	prepared.inputs, prepared.outputs = inputs, outputs
	return prepared, inputs, outputs, diagnostics
}

func flattenContractValues(values map[string]ContractValue, label string, budget PromptBudget) (map[string]ContractValue, []Diagnostic) {
	flat := map[string]ContractValue{}
	var diagnostics []Diagnostic
	work := 0
	var visit func(string, ContractValue, int)
	visit = func(path string, value ContractValue, depth int) {
		if len(diagnostics) > 0 {
			return
		}
		work++
		if work > 10_000 || depth > 32 || len(flat) >= budget.MaxFields {
			diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_budget", Message: fmt.Sprintf("step %s exceed the nested value count/depth budget", label), Remediation: "Reduce nested contract fields or increase MaxFields within reviewed bounds."})
			return
		}
		if path == "" || strings.TrimSpace(path) != path {
			diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_name_invalid", Message: fmt.Sprintf("step %s contains an empty or whitespace-padded value name", label)})
			return
		}
		if _, changed := sanitizePromptString(path, budget.MaxIdentifierRunes); changed {
			diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_name_unsafe", Message: fmt.Sprintf("step %s value names contain controls or exceed the identifier budget", label)})
			return
		}
		if existing, exists := flat[path]; exists {
			_ = existing
			diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_duplicate", Message: fmt.Sprintf("step %s contains duplicate flattened value %q", label, path)})
			return
		}
		for _, field := range []struct {
			name, text string
			limit      int
		}{
			{name: "type", text: value.Type, limit: budget.MaxIdentifierRunes},
			{name: "format", text: value.Format, limit: budget.MaxIdentifierRunes},
			{name: "description", text: value.Description, limit: budget.MaxTextRunes},
		} {
			if _, changed := sanitizePromptString(field.text, field.limit); changed {
				diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_unsafe", Message: fmt.Sprintf("step %s value %q %s is unsafe or exceeds its budget", label, path, field.name)})
				return
			}
		}
		flat[path] = ContractValue{Type: value.Type, Format: value.Format, Required: value.Required, Description: value.Description}
		keys := sortedContractValueKeys(value.Properties)
		for _, key := range keys {
			if strings.TrimSpace(key) == "" || strings.TrimSpace(key) != key {
				diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_name_invalid", Message: fmt.Sprintf("step %s contains an empty or whitespace-padded property name", label)})
				return
			}
			visit(path+"."+key, value.Properties[key], depth+1)
		}
		if value.Items != nil {
			visit(path+"[]", *value.Items, depth+1)
		}
	}
	for _, key := range sortedContractValueKeys(values) {
		if strings.TrimSpace(key) == "" || strings.TrimSpace(key) != key {
			diagnostics = append(diagnostics, Diagnostic{Severity: "error", Code: "contract.value_name_invalid", Message: fmt.Sprintf("step %s contains an empty or whitespace-padded value name", label)})
			break
		}
		visit(key, values[key], 0)
		if len(diagnostics) > 0 {
			break
		}
	}
	return flat, diagnostics
}

func sortedContractValueKeys(values map[string]ContractValue) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func compareCandidateToContract(candidate OperationCandidate, contract flattenedStepContract, stepInputs, stepOutputs map[string]ContractValue) StepContractMatch {
	match := StepContractMatch{}
	match.Purpose = matchPurpose(candidate, contract.contract.Purpose)
	match.Inputs = matchInputs(candidate, stepInputs)
	match.Outputs = matchOutputs(candidate, stepOutputs)
	match.Effect = matchEffect(candidate.Effect, contract.contract.Effect)
	match.Score = match.Purpose.Score + match.Inputs.Score + match.Outputs.Score + match.Effect.Score
	return match
}

func matchPurpose(candidate OperationCandidate, purpose string) ContractDimensionMatch {
	result := ContractDimensionMatch{Status: ContractMatchIndeterminate}
	queryTokens := purposeTokens(purpose)
	if len(queryTokens) == 0 {
		result.Missing = []string{"Step purpose is empty or contains no searchable terms."}
		return result
	}
	var text []string
	for _, evidence := range candidate.Summary.Evidence {
		if evidence.Kind == "operation.summary" || evidence.Kind == "operation.description" {
			text = append(text, candidate.Summary.Description)
			result.Evidence = append(result.Evidence, evidence)
			break
		}
	}
	if len(text) == 0 {
		text = append(text, candidate.Operation.OperationID, candidate.Operation.ID, candidate.Operation.Path)
		result.Evidence = append(result.Evidence, OperationEvidence{Kind: "operation.identity", Reference: candidate.Source.Selector})
		result.Gaps = append(result.Gaps, "Only operation identity is available; it does not establish the operation's purpose.")
	}
	available := make(map[string]bool)
	for _, value := range text {
		for _, token := range purposeTokens(value) {
			available[token] = true
		}
	}
	shared := 0
	for _, token := range queryTokens {
		if available[token] {
			shared++
		}
	}
	if shared == 0 {
		result.Missing = append(result.Missing, "No source purpose wording overlaps the requested step purpose.")
		result.Reasons = []string{"Purpose fit is indeterminate because the source wording has no overlapping content terms."}
		return result
	}
	result.Score = 20 * shared / len(queryTokens)
	if result.Score == 0 {
		result.Score = 1
	}
	if len(result.Gaps) > 0 {
		result.Missing = append(result.Missing, "The operation's documented purpose is unavailable; identity overlap is only a weak ranking hint.")
		result.Reasons = []string{"Identity-token overlap is a weak ranking hint, not evidence of the operation's purpose."}
		return result
	}
	result.Status = ContractMatchCompatible
	result.Reasons = []string{fmt.Sprintf("%d of %d step-purpose terms overlap documented operation wording; this is advisory lexical similarity.", shared, len(queryTokens))}
	return result
}

func purposeTokens(value string) []string {
	stop := map[string]bool{"a": true, "an": true, "and": true, "are": true, "by": true, "for": true, "from": true, "in": true, "into": true, "of": true, "on": true, "or": true, "the": true, "to": true, "with": true, "operation": true, "api": true, "step": true}
	set := map[string]bool{}
	for _, token := range effectWords(value) {
		if len(token) < 3 || stop[token] {
			continue
		}
		if strings.HasSuffix(token, "ies") && len(token) > 4 {
			token = strings.TrimSuffix(token, "ies") + "y"
		} else if strings.HasSuffix(token, "s") && !strings.HasSuffix(token, "ss") && len(token) > 3 {
			token = strings.TrimSuffix(token, "s")
		}
		set[token] = true
	}
	result := make([]string, 0, len(set))
	for token := range set {
		result = append(result, token)
	}
	sort.Strings(result)
	return result
}

func matchInputs(candidate OperationCandidate, stepValues map[string]ContractValue) ContractDimensionMatch {
	result := ContractDimensionMatch{Status: ContractMatchCompatible}
	capability := candidateCapability(candidate, "inputs")
	index := indexOperationValues(candidate.Summary.Inputs)
	accumulator := matchAccumulator{}
	result.Evidence = append(result.Evidence, capability.Evidence...)
	result.Gaps = append(result.Gaps, capability.Gaps...)
	matched := 0
	knownInputs := 0
	for i, operationValue := range index.values {
		if operationValue.Name == "body" && operationValue.Location == "body" && hasBodyProperties(index.values) {
			bodyRequired := operationValueRequired(operationValue, true)
			if value, _, ok, _, _ := stepValueForSource(stepValues, index, i); ok {
				knownInputs++
				if bodyRequired != nil && *bodyRequired && value.Required != nil && !*value.Required {
					accumulator.incompatible = true
					result.Conflicts = append(result.Conflicts, "Step body input is optional but the operation requires a body.")
				} else if bodyRequired != nil && *bodyRequired && value.Required == nil {
					accumulator.indeterminate = true
					result.Gaps = append(result.Gaps, "Step body input does not declare whether it is always available.")
				}
				status, reason := compareValueType(value, contractValueFromOperation(operationValue), true)
				applyTypeResult(&result, &accumulator, status, "body", reason)
				result.Evidence = append(result.Evidence, operationValue.Evidence...)
				if status == ContractMatchCompatible && !(bodyRequired != nil && *bodyRequired && (value.Required == nil || !*value.Required)) {
					matched++
				}
			} else if bodyRequired == nil || *bodyRequired {
				knownInputs++
				if bodyRequired != nil && *bodyRequired && hasRequiredBodyStepInput(stepValues, index) {
					matched++
				} else {
					result.Gaps = append(result.Gaps, "The required request-body container has no explicit body value; flattened fields may only partially establish it.")
					accumulator.indeterminate = true
				}
			}
			continue
		}
		knownInputs++
		value, key, found, ambiguous, parent := stepValueForSource(stepValues, index, i)
		if !found {
			if ambiguous {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Request input "+operationValue.Name+" has multiple possible step-input aliases.")
				continue
			}
			required := operationValueRequired(operationValue, false)
			if required == nil {
				accumulator.indeterminate = true
				result.Missing = append(result.Missing, "Request input "+operationValue.Name+" has unknown requiredness and no matching step input.")
			} else if parent {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Step input for "+operationValue.Name+" is covered only by a parent object without field detail.")
			} else if *required {
				if capability.Status == OperationCapabilitySupported {
					accumulator.incompatible = true
					result.Missing = append(result.Missing, "Required operation input "+operationValue.Name+" is not declared by the step contract.")
				} else {
					accumulator.indeterminate = true
					result.Missing = append(result.Missing, "Potentially required operation input "+operationValue.Name+" is not declared; the source adapter is incomplete.")
				}
			}
			continue
		}
		stepRequired := value.Required
		operationRequired := operationValueRequired(operationValue, false)
		if operationRequired != nil && *operationRequired && stepRequired != nil && !*stepRequired {
			accumulator.incompatible = true
			result.Conflicts = append(result.Conflicts, "Step input "+key+" is optional but the operation requires it.")
		} else if operationRequired == nil && (stepRequired == nil || !*stepRequired) {
			accumulator.indeterminate = true
			result.Gaps = append(result.Gaps, "Request input "+operationValue.Name+" requiredness is not fully established.")
		} else if operationRequired != nil && *operationRequired && stepRequired == nil {
			accumulator.indeterminate = true
			result.Gaps = append(result.Gaps, "Step input "+key+" does not declare whether it is always available.")
		}
		status, reason := compareValueType(value, contractValueFromOperation(operationValue), true)
		applyTypeResult(&result, &accumulator, status, key, reason)
		result.Evidence = append(result.Evidence, operationValue.Evidence...)
		if status == ContractMatchCompatible && !(operationRequired != nil && *operationRequired && (stepRequired == nil || !*stepRequired)) {
			matched++
		}
	}
	if len(index.values) == 0 {
		if capability.Status == OperationCapabilitySupported {
			result.Score = 30
			result.Reasons = []string{"The source establishes that this operation declares no request inputs."}
		} else {
			accumulator.indeterminate = true
			result.Missing = append(result.Missing, "The source adapter does not establish the operation's input contract.")
			result.Reasons = []string{"Input coverage cannot be scored because source capability is unavailable."}
		}
	} else if knownInputs > 0 {
		result.Score = 30 * matched / knownInputs
		result.Reasons = []string{fmt.Sprintf("%d of %d source-declared input values have compatible type and requiredness evidence.", matched, knownInputs)}
	}
	for name := range stepValues {
		if !index.hasAlias(name) {
			if hasSourceParent(index, name) {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Nested step input "+name+" is declared beneath a source value without matching field detail.")
			} else {
				result.Gaps = append(result.Gaps, "Step input "+name+" is not used by the operation's declared inputs.")
			}
		}
	}
	if capability.Status != OperationCapabilitySupported {
		accumulator.indeterminate = true
		result.Gaps = append(result.Gaps, "Input evidence is partial or unsupported for this source family.")
	}
	setDimensionStatus(&result, accumulator)
	return normalizeDimensionMatch(result)
}

func matchOutputs(candidate OperationCandidate, stepValues map[string]ContractValue) ContractDimensionMatch {
	result := ContractDimensionMatch{Status: ContractMatchCompatible}
	capability := candidateCapability(candidate, "outputs")
	index := indexOperationValues(candidate.Summary.Outputs)
	accumulator := matchAccumulator{}
	result.Evidence = append(result.Evidence, capability.Evidence...)
	result.Gaps = append(result.Gaps, capability.Gaps...)
	matched := 0
	outputAliases := map[int]string{}
	for _, name := range sortedContractValueKeys(stepValues) {
		expected := stepValues[name]
		_, source, found, ambiguous, parent := operationValueForStep(name, index)
		if !found {
			if ambiguous {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Expected output "+name+" matches multiple source values.")
			} else if expected.Required != nil && !*expected.Required {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Optional expected output "+name+" is not established by the source.")
			} else if parent {
				accumulator.indeterminate = true
				result.Missing = append(result.Missing, "Expected output "+name+" is covered only by a source parent without field detail.")
			} else if capability.Status == OperationCapabilitySupported {
				accumulator.incompatible = true
				result.Missing = append(result.Missing, "Expected output "+name+" is absent from the operation response contract.")
			} else {
				accumulator.indeterminate = true
				result.Missing = append(result.Missing, "Expected output "+name+" is not present in the partial source metadata.")
			}
			continue
		}
		sourceIndex := index.byAlias[name][0]
		duplicateAlias := false
		if previous, exists := outputAliases[sourceIndex]; exists && previous != name {
			duplicateAlias = true
			accumulator.indeterminate = true
			result.Gaps = append(result.Gaps, "Expected outputs "+previous+" and "+name+" resolve to the same source value.")
		} else {
			outputAliases[sourceIndex] = name
		}
		if expected.Required == nil {
			accumulator.indeterminate = true
			result.Gaps = append(result.Gaps, "Expected output "+name+" does not declare requiredness.")
		}
		if expected.Required != nil && *expected.Required {
			providedRequired := source.Required
			if providedRequired != nil && !*providedRequired {
				accumulator.incompatible = true
				result.Conflicts = append(result.Conflicts, "Expected output "+name+" is required but the operation marks it optional.")
			} else if providedRequired == nil {
				accumulator.indeterminate = true
				result.Gaps = append(result.Gaps, "Operation output "+name+" does not establish requiredness.")
			}
		}
		status, reason := compareValueType(expected, contractValueFromOperation(source), false)
		applyTypeResult(&result, &accumulator, status, name, reason)
		if source.Nullable {
			accumulator.indeterminate = true
			result.Gaps = append(result.Gaps, "Operation output "+name+" or one of its response schema ancestors permits null; the step contract does not establish that null is acceptable.")
		}
		result.Evidence = append(result.Evidence, source.Evidence...)
		if status == ContractMatchCompatible && !source.Nullable && !duplicateAlias && !hasConflictForName(result.Conflicts, name) {
			matched++
		}
	}
	if len(stepValues) == 0 {
		if capability.Status == OperationCapabilitySupported {
			result.Score = 0
		} else {
			accumulator.indeterminate = true
			result.Missing = append(result.Missing, "The source adapter does not establish response-output capabilities.")
		}
	} else {
		result.Score = 30 * matched / len(stepValues)
	}
	if len(stepValues) == 0 {
		result.Reasons = []string{"No expected output values were declared."}
	} else {
		result.Reasons = []string{fmt.Sprintf("%d of %d expected output values have compatible source type evidence.", matched, len(stepValues))}
	}
	if capability.Status != OperationCapabilitySupported {
		accumulator.indeterminate = true
		result.Gaps = append(result.Gaps, "Output evidence is partial or unsupported for this source family.")
	}
	setDimensionStatus(&result, accumulator)
	return normalizeDimensionMatch(result)
}

func matchEffect(assessment EffectAssessment, constraint OperationEffect) ContractDimensionMatch {
	result := ContractDimensionMatch{Status: ContractMatchIndeterminate, Evidence: append([]OperationEvidence(nil), assessment.Evidence...)}
	if constraint == "" {
		result.Gaps = []string{"The step contract does not constrain operation effect."}
		result.Reasons = []string{"Effect contributes no score because the step supplied no effect constraint."}
		return result
	}
	if constraint == OperationEffectUnknown {
		result.Missing = []string{"The step effect is unknown; no read/write compatibility claim can be made."}
		result.Reasons = []string{"An unknown step effect does not establish read/write compatibility."}
		return result
	}
	if assessment.Class == OperationEffectUnknown || assessment.Class == "" {
		result.Missing = []string{"The operation effect is unknown and cannot satisfy a declared read/write constraint."}
		result.Reasons = []string{"Unknown operation effect cannot satisfy a declared read/write constraint."}
		return result
	}
	if assessment.Class != constraint {
		result.Status = ContractMatchIncompatible
		result.Conflicts = []string{fmt.Sprintf("Step requires effect %q but the operation is classified as %q.", constraint, assessment.Class)}
		result.Reasons = []string{"The source-backed operation effect contradicts the requested effect."}
		return result
	}
	result.Status = ContractMatchCompatible
	result.Score = 20
	result.Reasons = append(result.Reasons, assessment.Reasons...)
	return normalizeDimensionMatch(result)
}

func candidateCapability(candidate OperationCandidate, dimension string) OperationCapability {
	for _, capability := range candidate.Capabilities {
		if capability.Dimension == dimension {
			return capability
		}
	}
	return OperationCapability{Dimension: dimension, Status: OperationCapabilityPartial, Gaps: []string{"Source capability was not declared."}}
}

func indexOperationValues(values []OperationValueSummary) sourceValueIndex {
	index := sourceValueIndex{values: append([]OperationValueSummary(nil), values...), byAlias: map[string][]int{}}
	for i, value := range index.values {
		for _, alias := range operationValueAliases(value) {
			index.byAlias[alias] = append(index.byAlias[alias], i)
		}
	}
	return index
}

func operationValueAliases(value OperationValueSummary) []string {
	name := strings.TrimSpace(value.Name)
	location := strings.TrimSpace(value.Location)
	if name == "" {
		return nil
	}
	var aliases []string
	if (location == "body" || location == "response") && name == "body" {
		aliases = append(aliases, "body")
	} else {
		switch location {
		case "body":
			aliases = append(aliases, "body."+name, name)
		case "response":
			aliases = append(aliases, "body."+name, name, "response."+name)
		case "":
			aliases = append(aliases, name)
		default:
			aliases = append(aliases, location+"."+name, name)
		}
	}
	return uniqueStringsStable(aliases)
}

func uniqueStringsStable(values []string) []string {
	seen := make(map[string]bool, len(values))
	out := make([]string, 0, len(values))
	for _, value := range values {
		if value != "" && !seen[value] {
			seen[value] = true
			out = append(out, value)
		}
	}
	return out
}

func (index sourceValueIndex) hasAlias(name string) bool { return len(index.byAlias[name]) > 0 }

func stepValueForSource(stepValues map[string]ContractValue, index sourceValueIndex, sourceIndex int) (ContractValue, string, bool, bool, bool) {
	if sourceIndex < 0 || sourceIndex >= len(index.values) {
		return ContractValue{}, "", false, false, false
	}
	value := index.values[sourceIndex]
	var matches []string
	for _, alias := range operationValueAliases(value) {
		if _, ok := stepValues[alias]; !ok {
			continue
		}
		if indices := index.byAlias[alias]; len(indices) > 1 {
			return ContractValue{}, alias, false, true, false
		} else if len(indices) == 1 && indices[0] == sourceIndex {
			matches = append(matches, alias)
		}
	}
	if len(matches) > 1 {
		return ContractValue{}, strings.Join(matches, " / "), false, true, false
	}
	if len(matches) == 1 {
		return stepValues[matches[0]], matches[0], true, false, false
	}
	return ContractValue{}, "", false, false, hasContractParent(stepValues, operationValueAliases(value))
}

func operationValueForStep(stepName string, index sourceValueIndex) (ContractValue, OperationValueSummary, bool, bool, bool) {
	indices := index.byAlias[stepName]
	if len(indices) == 1 {
		value := index.values[indices[0]]
		return contractValueFromOperation(value), value, true, false, false
	}
	if len(indices) > 1 {
		return ContractValue{}, OperationValueSummary{}, false, true, false
	}
	return ContractValue{}, OperationValueSummary{}, false, false, hasSourceParent(index, stepName)
}

func hasContractParent(values map[string]ContractValue, aliases []string) bool {
	for _, alias := range aliases {
		for _, parent := range parentValueNames(alias) {
			if _, ok := values[parent]; ok {
				return true
			}
		}
	}
	return false
}

func hasSourceParent(index sourceValueIndex, name string) bool {
	for _, parent := range parentValueNames(name) {
		if len(index.byAlias[parent]) > 0 {
			return true
		}
	}
	return false
}

func parentValueNames(value string) []string {
	var parents []string
	for current := value; current != ""; {
		dot := strings.LastIndex(current, ".")
		array := strings.LastIndex(current, "[]")
		if array > dot {
			current = current[:array]
		} else if dot >= 0 {
			current = current[:dot]
		} else {
			break
		}
		if current != "" {
			parents = append(parents, current)
		}
	}
	return parents
}

func contractValueFromOperation(value OperationValueSummary) ContractValue {
	return ContractValue{Type: value.Type, Format: value.Format, Required: operationValueRequired(value, false), Description: value.Description}
}

func operationValueRequired(value OperationValueSummary, container bool) *bool {
	if container || value.Name == "body" {
		return value.Required
	}
	if value.Location == "body" && value.ContainerRequired != nil {
		if !*value.ContainerRequired {
			return boolPointer(false)
		}
		if value.Required == nil {
			return nil
		}
	}
	return value.Required
}

func compareValueType(step, operation ContractValue, isInput bool) (ContractMatchStatus, string) {
	stepType, operationType := normalizeContractType(step.Type), normalizeContractType(operation.Type)
	if stepType == "" || operationType == "" {
		return ContractMatchIndeterminate, "Type evidence is missing on at least one side."
	}
	if !isKnownPrimitive(stepType) || !isKnownPrimitive(operationType) {
		return ContractMatchIndeterminate, fmt.Sprintf("Named or protocol types %q and %q cannot be compared without resolving source shapes.", stepType, operationType)
	}
	if stepType != operationType {
		if isInput && stepType == "integer" && operationType == "number" || !isInput && operationType == "integer" && stepType == "number" {
			// Integer values are assignable to a number contract in either data-flow direction.
		} else {
			return ContractMatchIncompatible, fmt.Sprintf("Type conflict: step has %q while operation declares %q.", stepType, operationType)
		}
	}
	stepFormat, operationFormat := strings.ToLower(strings.TrimSpace(step.Format)), strings.ToLower(strings.TrimSpace(operation.Format))
	if stepFormat != "" && operationFormat != "" && stepFormat != operationFormat {
		return ContractMatchIncompatible, fmt.Sprintf("Format conflict: step has %q while operation declares %q.", stepFormat, operationFormat)
	}
	if (stepFormat == "") != (operationFormat == "") {
		return ContractMatchIndeterminate, "Format evidence is present on only one side."
	}
	return ContractMatchCompatible, ""
}

func normalizeContractType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "int", "int32", "int64", "integer":
		return "integer"
	case "float", "double", "decimal", "number":
		return "number"
	case "bool":
		return "boolean"
	default:
		return value
	}
}

func isKnownPrimitive(value string) bool {
	switch value {
	case "string", "integer", "number", "boolean", "object", "array", "null":
		return true
	default:
		return false
	}
}

func applyTypeResult(result *ContractDimensionMatch, accumulator *matchAccumulator, status ContractMatchStatus, name, reason string) {
	switch status {
	case ContractMatchIncompatible:
		accumulator.incompatible = true
		result.Conflicts = append(result.Conflicts, "Value "+name+": "+reason)
	case ContractMatchIndeterminate:
		accumulator.indeterminate = true
		result.Gaps = append(result.Gaps, "Value "+name+": "+reason)
	}
}

func setDimensionStatus(result *ContractDimensionMatch, accumulator matchAccumulator) {
	if accumulator.incompatible {
		result.Status = ContractMatchIncompatible
	} else if accumulator.indeterminate {
		result.Status = ContractMatchIndeterminate
	} else {
		result.Status = ContractMatchCompatible
	}
}

func normalizeDimensionMatch(result ContractDimensionMatch) ContractDimensionMatch {
	result.Reasons = uniqueSortedStrings(result.Reasons)
	result.Missing = uniqueSortedStrings(result.Missing)
	result.Conflicts = uniqueSortedStrings(result.Conflicts)
	result.Gaps = uniqueSortedStrings(result.Gaps)
	result.Evidence = uniqueOperationEvidence(result.Evidence)
	return result
}

func uniqueOperationEvidence(values []OperationEvidence) []OperationEvidence {
	if len(values) == 0 {
		return nil
	}
	sort.SliceStable(values, func(i, j int) bool {
		if values[i].Kind != values[j].Kind {
			return values[i].Kind < values[j].Kind
		}
		return values[i].Reference < values[j].Reference
	})
	out := values[:1]
	for _, value := range values[1:] {
		if value != out[len(out)-1] {
			out = append(out, value)
		}
	}
	return out
}

func hasContractMatch(candidates []OperationCandidate, contract StepContract) bool {
	if strings.TrimSpace(contract.Purpose) == "" && len(contract.Inputs) == 0 && len(contract.Outputs) == 0 && contract.Effect == "" {
		return false
	}
	for _, candidate := range candidates {
		match := candidate.Match
		if contract.Purpose != "" && match.Purpose.Status != ContractMatchCompatible {
			continue
		}
		if len(contract.Inputs) > 0 && match.Inputs.Status != ContractMatchCompatible {
			continue
		}
		if len(contract.Outputs) > 0 && match.Outputs.Status != ContractMatchCompatible {
			continue
		}
		if contract.Effect != "" && match.Effect.Status != ContractMatchCompatible {
			continue
		}
		if !hasIncompatibleDimension(match) {
			return true
		}
	}
	return false
}

func hasIncompatibleDimension(match StepContractMatch) bool {
	return match.Purpose.Status == ContractMatchIncompatible || match.Inputs.Status == ContractMatchIncompatible || match.Outputs.Status == ContractMatchIncompatible || match.Effect.Status == ContractMatchIncompatible
}

func hasConflictForName(conflicts []string, name string) bool {
	for _, conflict := range conflicts {
		if strings.Contains(conflict, "Value "+name+":") || strings.Contains(conflict, "Expected output "+name+" ") {
			return true
		}
	}
	return false
}

func candidateLess(left, right OperationCandidate) bool {
	if left.Match.Score != right.Match.Score {
		return left.Match.Score > right.Match.Score
	}
	leftKeys := []string{string(left.Source.Kind), left.Source.DocumentPath, left.Source.DocumentURL, left.Source.SHA256, left.Source.Selector, left.Operation.Method, left.Operation.Path, left.Operation.ID}
	rightKeys := []string{string(right.Source.Kind), right.Source.DocumentPath, right.Source.DocumentURL, right.Source.SHA256, right.Source.Selector, right.Operation.Method, right.Operation.Path, right.Operation.ID}
	for i := range leftKeys {
		if leftKeys[i] != rightKeys[i] {
			return leftKeys[i] < rightKeys[i]
		}
	}
	return false
}

func hasBodyProperties(values []OperationValueSummary) bool {
	for _, value := range values {
		if value.Location == "body" && value.Name != "body" {
			return true
		}
	}
	return false
}

func marshalSize(value any) (int, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func hasRequiredBodyStepInput(stepValues map[string]ContractValue, index sourceValueIndex) bool {
	for i, value := range index.values {
		if value.Location != "body" || value.Name == "body" {
			continue
		}
		required := operationValueRequired(value, false)
		if required == nil || !*required {
			continue
		}
		step, _, found, _, _ := stepValueForSource(stepValues, index, i)
		if found && step.Required != nil && *step.Required {
			return true
		}
	}
	return false
}
