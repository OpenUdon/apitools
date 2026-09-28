package apitools

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func summarizeOperationForConsumer(operation OperationSummary, selector string, promptBudget PromptBudget) (ConsumerOperationSummary, []Diagnostic, error) {
	budget := resolvedPromptBudget(promptBudget)
	if count := operationInputCount(operation); count > budget.MaxFields {
		diagnostic := Diagnostic{
			Severity: "error", Code: "operation_summary.input_budget",
			Message:     fmt.Sprintf("operation declares %d input values, exceeding the %d-value consumer summary budget", count, budget.MaxFields),
			Path:        firstNonEmpty(selector, operation.ID),
			Remediation: "Narrow the selected operation or increase MaxFields for a reviewed summary.",
		}
		return ConsumerOperationSummary{}, []Diagnostic{diagnostic}, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	if len(operation.Parameters) > budget.MaxCollectionItems {
		diagnostic := Diagnostic{
			Severity: "error", Code: "operation_summary.parameter_budget",
			Message:     fmt.Sprintf("operation declares %d parameters, exceeding the %d-parameter consumer summary budget", len(operation.Parameters), budget.MaxCollectionItems),
			Path:        firstNonEmpty(selector, operation.ID),
			Remediation: "Narrow the selected operation or increase MaxCollectionItems for a reviewed summary.",
		}
		return ConsumerOperationSummary{}, []Diagnostic{diagnostic}, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	if count := operationOutputCount(operation); count > budget.MaxFields {
		diagnostic := Diagnostic{
			Severity: "error", Code: "operation_summary.output_budget",
			Message:     fmt.Sprintf("operation declares %d output values, exceeding the %d-value consumer summary budget", count, budget.MaxFields),
			Path:        firstNonEmpty(selector, operation.ID),
			Remediation: "Narrow the selected operation or increase MaxFields for a reviewed summary.",
		}
		return ConsumerOperationSummary{}, []Diagnostic{diagnostic}, DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}

	sanitized, err := SanitizeOperationSummaries([]OperationSummary{operation}, budget)
	if err != nil {
		return ConsumerOperationSummary{}, sanitized.Diagnostics, err
	}
	if len(sanitized.Operations) != 1 {
		return ConsumerOperationSummary{}, sanitized.Diagnostics, fmt.Errorf("operation summary sanitization returned %d operations", len(sanitized.Operations))
	}
	op := sanitized.Operations[0]
	ref := strings.TrimSpace(selector)
	if ref == "" {
		ref = firstNonEmpty(op.OperationID, op.ID, op.Path)
	}
	if safe, changed := sanitizePromptString(ref, budget.MaxTextRunes); changed || strings.TrimSpace(safe) == "" {
		diagnostic := Diagnostic{
			Severity: "error", Code: "operation_summary.selector_unsafe",
			Message:     "the operation selector is empty or exceeds the safe source-reference budget",
			Path:        op.ID,
			Remediation: "Use an exact bounded native selector for the operation; do not truncate source identity.",
		}
		return ConsumerOperationSummary{}, append(sanitized.Diagnostics, diagnostic), DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	} else {
		ref = safe
	}

	summary := ConsumerOperationSummary{}
	text := strings.TrimSpace(firstNonEmpty(op.Summary, op.Description))
	if text != "" {
		summary.Description = text
		kind := "operation.description"
		if strings.TrimSpace(op.Summary) != "" {
			kind = "operation.summary"
		}
		summary.Evidence = append(summary.Evidence, OperationEvidence{Kind: kind, Reference: ref})
	} else {
		identity := firstNonEmpty(op.OperationID, op.ID, "unnamed operation")
		summary.Description = fmt.Sprintf("Operation %q", identity)
		summary.Description += "; its documented purpose is unavailable."
		summary.Evidence = append(summary.Evidence, OperationEvidence{Kind: "operation.identity", Reference: ref})
		summary.Gaps = append(summary.Gaps, "The source does not provide an operation summary or description.")
	}

	for _, parameter := range op.Parameters {
		if strings.TrimSpace(parameter.Name) == "" {
			summary.Gaps = append(summary.Gaps, "A declared request parameter has no usable name.")
			continue
		}
		if looksLikeCredentialName(parameter.Name) {
			summary.Gaps = append(summary.Gaps, "Credential-shaped parameter "+parameter.Name+" is omitted from workflow data inputs; review the separate auth requirements.")
			continue
		}
		location := strings.TrimSpace(parameter.In)
		summary.Inputs = append(summary.Inputs, OperationValueSummary{
			Name: parameter.Name, Location: location, Type: parameter.Type, Format: parameter.Format,
			Required: boolPointer(parameter.Required), Description: parameter.Description,
			Evidence: []OperationEvidence{{Kind: "request.parameter", Reference: evidenceReference(ref, "parameters", location, parameter.Name)}},
		})
		if parameter.Ref != "" {
			summary.Gaps = append(summary.Gaps, "Request parameter "+parameter.Name+" contains an unresolved schema reference.")
		}
		if parameter.Type == "" {
			summary.Gaps = append(summary.Gaps, "Request parameter "+parameter.Name+" has no declared type.")
		}
	}

	if body := op.RequestBody; body != nil {
		if len(body.Fields) > 0 {
			bodyType := "object"
			bodyFormat := ""
			if body.Schema != nil {
				bodyType = firstNonEmpty(body.Schema.Type, bodyType)
				bodyFormat = body.Schema.Format
			}
			summary.Inputs = append(summary.Inputs, OperationValueSummary{
				Name: "body", Location: "body", Type: bodyType, Format: bodyFormat,
				Required: boolPointer(body.Required), Description: body.Description,
				Evidence: []OperationEvidence{{Kind: "request.body", Reference: evidenceReference(ref, "request")}},
			})
			if body.Schema != nil && body.Schema.Ref != "" {
				summary.Gaps = append(summary.Gaps, "The request body schema contains an unresolved reference.")
			}
			for _, field := range body.Fields {
				name := strings.TrimSpace(field.Path)
				if name == "" {
					summary.Gaps = append(summary.Gaps, "A request-body field has no usable path.")
					continue
				}
				summary.Inputs = append(summary.Inputs, OperationValueSummary{
					Name: name, Location: "body", Type: field.Type, Format: field.Format,
					Required: boolPointer(field.Required), ContainerRequired: boolPointer(body.Required), Description: field.Description,
					Evidence: []OperationEvidence{{Kind: "request.body_field", Reference: evidenceReference(ref, "request", "fields", name)}},
				})
				if field.Ref != "" {
					summary.Gaps = append(summary.Gaps, "Request field "+name+" contains an unresolved schema reference.")
				}
				if field.Type == "" {
					summary.Gaps = append(summary.Gaps, "Request field "+name+" has no declared type.")
				}
			}
		} else {
			value := OperationValueSummary{
				Name: "body", Location: "body", Required: boolPointer(body.Required),
				Description: body.Description,
				Evidence:    []OperationEvidence{{Kind: "request.body", Reference: evidenceReference(ref, "request")}},
			}
			if body.Schema != nil {
				value.Type, value.Format = body.Schema.Type, body.Schema.Format
				if body.Schema.Ref != "" {
					summary.Gaps = append(summary.Gaps, "The request body schema contains an unresolved reference.")
				}
			}
			summary.Inputs = append(summary.Inputs, value)
			if value.Type == "" {
				summary.Gaps = append(summary.Gaps, "The request body has no inspectable schema type or fields.")
			}
		}
	}
	if len(summary.Inputs) == 0 {
		summary.Gaps = append(summary.Gaps, "The source declares no request parameters or body.")
	}

	if body := op.ResponseBody; body != nil {
		if len(body.Fields) > 0 {
			bodyType := "object"
			bodyFormat := ""
			if body.Schema != nil {
				bodyType = firstNonEmpty(body.Schema.Type, bodyType)
				bodyFormat = body.Schema.Format
			}
			summary.Outputs = append(summary.Outputs, OperationValueSummary{
				Name: "body", Location: "response", Type: bodyType, Format: bodyFormat,
				Nullable:    body.Schema != nil && body.Schema.Nullable,
				Description: body.Description,
				Evidence:    []OperationEvidence{{Kind: "response.body", Reference: evidenceReference(ref, "response")}},
			})
			if body.Schema != nil && body.Schema.Ref != "" {
				summary.Gaps = append(summary.Gaps, "The response body schema contains an unresolved reference.")
			}
			for i, field := range body.Fields {
				name := strings.TrimSpace(field.Path)
				if name == "" {
					summary.Gaps = append(summary.Gaps, "A response field has no usable path.")
					continue
				}
				if syntheticResponseBodyField(body, i, field) {
					summary.Gaps = append(summary.Gaps, "The response schema does not establish whether a body value is always present.")
					continue
				}
				required := boolPointer(field.Required)
				summary.Outputs = append(summary.Outputs, OperationValueSummary{
					Name: name, Location: "response", Type: field.Type, Format: field.Format,
					Required: required, Nullable: field.Nullable, Description: field.Description,
					Evidence: []OperationEvidence{{Kind: "response.field", Reference: evidenceReference(ref, "response", "fields", name)}},
				})
				if field.Ref != "" {
					summary.Gaps = append(summary.Gaps, "Response field "+name+" contains an unresolved schema reference.")
				}
				if field.Type == "" {
					summary.Gaps = append(summary.Gaps, "Response field "+name+" has no declared type.")
				}
			}
		} else if body.Schema != nil {
			value := OperationValueSummary{
				Name: "body", Location: "response", Type: body.Schema.Type,
				Format: body.Schema.Format, Nullable: body.Schema.Nullable, Description: body.Description,
				Evidence: []OperationEvidence{{Kind: "response.body", Reference: evidenceReference(ref, "response")}},
			}
			if body.Schema.Ref != "" {
				summary.Gaps = append(summary.Gaps, "The response body schema contains an unresolved reference.")
			}
			summary.Outputs = append(summary.Outputs, value)
			if value.Type == "" {
				summary.Gaps = append(summary.Gaps, "The response body has no declared schema type or fields.")
			}
		} else {
			summary.Gaps = append(summary.Gaps, "A successful response is declared without inspectable body fields.")
		}
	} else {
		summary.Gaps = append(summary.Gaps, "The source does not declare an inspectable successful response schema.")
	}

	for _, issue := range op.ReadinessIssues {
		if issue.Severity == "error" && strings.TrimSpace(issue.Message) != "" {
			summary.Gaps = append(summary.Gaps, issue.Message)
		}
	}
	sort.SliceStable(summary.Inputs, func(i, j int) bool {
		if summary.Inputs[i].Location != summary.Inputs[j].Location {
			return summary.Inputs[i].Location < summary.Inputs[j].Location
		}
		return summary.Inputs[i].Name < summary.Inputs[j].Name
	})
	sort.SliceStable(summary.Outputs, func(i, j int) bool { return summary.Outputs[i].Name < summary.Outputs[j].Name })
	summary.Gaps = uniqueSortedStrings(summary.Gaps)
	encoded, err := json.Marshal(summary)
	if err != nil {
		return ConsumerOperationSummary{}, sanitized.Diagnostics, fmt.Errorf("encode consumer operation summary: %w", err)
	}
	if len(encoded) > budget.MaxOperationBytes {
		diagnostic := Diagnostic{
			Severity: "error", Code: "operation_summary.byte_budget",
			Message:     fmt.Sprintf("consumer operation summary is %d bytes, exceeding the %d-byte operation budget", len(encoded), budget.MaxOperationBytes),
			Path:        ref,
			Remediation: "Narrow the source schema or increase MaxOperationBytes for a reviewed summary.",
		}
		return ConsumerOperationSummary{}, append(sanitized.Diagnostics, diagnostic), DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
	}
	return summary, sanitized.Diagnostics, nil
}

func syntheticResponseBodyField(body *ResponseBodySummary, index int, field RequestFieldSummary) bool {
	return body != nil && index == 0 && len(body.Fields) == 1 &&
		strings.TrimSpace(field.Path) == "body" && body.Schema != nil && len(body.Schema.Properties) == 0
}

func operationInputCount(operation OperationSummary) int {
	count := len(operation.Parameters)
	if operation.RequestBody != nil {
		if len(operation.RequestBody.Fields) > 0 {
			count += 1 + len(operation.RequestBody.Fields)
		} else {
			count++
		}
	}
	return count
}

func operationOutputCount(operation OperationSummary) int {
	if operation.ResponseBody == nil {
		return 0
	}
	if len(operation.ResponseBody.Fields) > 0 {
		return 1 + len(operation.ResponseBody.Fields)
	}
	if operation.ResponseBody.Schema != nil {
		return 1
	}
	return 0
}

func boolPointer(value bool) *bool { return &value }

func evidenceSegment(values ...string) string {
	parts := make([]string, 0, len(values))
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			parts = append(parts, strings.ReplaceAll(strings.ReplaceAll(value, "~", "~0"), "/", "~1"))
		}
	}
	return strings.Join(parts, "/")
}

func evidenceReference(selector string, location ...string) string {
	suffix := evidenceSegment(location...)
	if strings.HasPrefix(selector, "#/") {
		return selector + "/" + suffix
	}
	return selector + "#/" + suffix
}

func uniqueSortedStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	sort.Strings(values)
	out := values[:0]
	for _, value := range values {
		if len(out) == 0 || out[len(out)-1] != value {
			out = append(out, value)
		}
	}
	return out
}
