package apitools

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/OpenUdon/apitools/internal/sourceguard"
	"github.com/OpenUdon/uws/binding"
)

// OperationEligibilityRequest selects metadata only. Truncated describes the
// supplying report; Candidate, when present, must identify the same native leaf.
// Source admission and independent source reproduction remain caller duties.
type OperationEligibilityRequest struct {
	Binding   binding.Binding     `json:"binding"`
	Candidate *OperationCandidate `json:"candidate,omitempty"`
	Truncated bool                `json:"truncated"`
}

// OperationEligibility preserves a private shape snapshot, including partial
// schemas and security alternatives. Compatible never grants product admission.
type OperationEligibility struct {
	Status  binding.Status          `json:"status"`
	Outcome binding.Outcome         `json:"outcome"`
	Codes   []string                `json:"codes"`
	Shape   *binding.OperationShape `json:"shape,omitempty"`
}

// AssessOperationEligibility resolves exact source/selector identity without
// choosing among operations or interpreting purpose, effects or credentials.
func AssessOperationEligibility(ctx context.Context, table binding.ShapeTable, request OperationEligibilityRequest) (OperationEligibility, error) {
	result := OperationEligibility{Outcome: binding.Indeterminate, Codes: []string{}}
	if ctx == nil {
		return result, ErrOperationShapeTable
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	resolver, err := binding.NewResolver(table)
	if err != nil {
		return result, ErrOperationShapeTable
	}
	resolved, err := resolver.Resolve(ctx, request.Binding)
	if err != nil {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}
		return result, ErrOperationShapeTable
	}
	result.Status, result.Shape = resolved.Status, resolved.Shape
	if resolved.Status != binding.Resolved || resolved.Shape == nil {
		result.Codes = append(result.Codes, "operation."+string(resolved.Status))
		if resolved.Status == binding.Missing {
			result.Outcome = binding.Incompatible
		}
		return result, nil
	}
	shape := resolved.Shape
	if candidate := request.Candidate; candidate != nil {
		selectorMatches := shape.Selector.Kind == "ref" && shape.Selector.Value == candidate.Source.Selector
		for _, alias := range shape.Aliases {
			selectorMatches = selectorMatches || alias.Kind == "ref" && alias.Value == candidate.Source.Selector
		}
		if string(candidate.Source.Kind) != shape.Source.Kind || candidate.Source.SHA256 != shape.Source.SHA256 || !selectorMatches {
			result.Outcome, result.Shape = binding.Incompatible, nil
			result.Codes = append(result.Codes, "operation.candidate_identity")
			return result, nil
		}
	}
	if !validOperationSourceKind(OperationSourceKind(shape.Source.Kind)) {
		result.Status = binding.Unsupported
		result.Codes = append(result.Codes, "operation.unsupported")
	}
	if request.Truncated {
		result.Codes = append(result.Codes, "operation.truncated")
	}
	if !shape.Complete {
		result.Codes = append(result.Codes, "operation.incomplete")
	}
	if !shape.Security.Known {
		result.Codes = append(result.Codes, "operation.security_unknown")
	}
	unknown := false
	for _, input := range shape.Inputs {
		unknown = unknown || !input.Schema.Known
	}
	for _, output := range shape.Outputs {
		unknown = unknown || !output.Schema.Known
	}
	if unknown {
		result.Codes = append(result.Codes, "operation.schema_unknown")
	}
	if len(result.Codes) == 0 {
		result.Outcome = binding.Compatible
	}
	return result, ctx.Err()
}

// OperationFieldReference names a whole native field. Direction is input or
// output; location and name are exact source metadata, not expression paths.
type OperationFieldReference struct {
	Direction string `json:"direction"`
	Location  string `json:"location"`
	Name      string `json:"name"`
}

// OperationFieldMapping contains schema evidence, never values. Required is
// the consumer input's presence guarantee (nil means unknown); outputs omit it.
// Placeholder preserves an unresolved pair without promoting it to a binding.
type OperationFieldMapping struct {
	OperationFieldReference
	Field       string         `json:"field"`
	Schema      binding.Schema `json:"schema"`
	Required    *bool          `json:"required,omitempty"`
	Placeholder bool           `json:"placeholder,omitempty"`
}

type OperationFieldMatch struct {
	Mapping           OperationFieldMapping `json:"mapping"`
	OperationRequired *bool                 `json:"operation_required,omitempty"`
	Outcome           binding.Outcome       `json:"outcome"`
	Code              string                `json:"code"`
}

type OperationMappingRequest struct {
	Selection OperationEligibilityRequest `json:"selection"`
	Mappings  []OperationFieldMapping     `json:"mappings"`
}

// OperationMappingReport describes declared schema relationships and ordinary
// input coverage only. Excluded inputs still need product-owned security policy.
type OperationMappingReport struct {
	Eligibility      OperationEligibility      `json:"eligibility"`
	Outcome          binding.Outcome           `json:"outcome"`
	Fields           []OperationFieldMatch     `json:"fields"`
	UnmappedRequired []OperationFieldReference `json:"unmapped_required"`
	ExcludedInputs   []OperationFieldReference `json:"excluded_inputs"`
}

// MatchOperationFields checks at most 64 explicit whole-field mappings. UWS
// checks directional schema containment; this helper never guesses field names,
// extracts nested schemas, emits expressions or satisfies authentication.
func MatchOperationFields(ctx context.Context, table binding.ShapeTable, request OperationMappingRequest) (OperationMappingReport, error) {
	report := OperationMappingReport{Outcome: binding.Indeterminate, Fields: []OperationFieldMatch{}, UnmappedRequired: []OperationFieldReference{}, ExcludedInputs: []OperationFieldReference{}}
	if len(request.Mappings) > 64 {
		return report, ErrOperationShapeTable
	}
	eligibility, err := AssessOperationEligibility(ctx, table, request.Selection)
	if err != nil {
		return report, err
	}
	report.Eligibility, report.Outcome = eligibility, eligibility.Outcome
	seen, supplied := map[string]bool{}, map[string]bool{}
	for _, mapping := range request.Mappings {
		if err := ctx.Err(); err != nil {
			return report, err
		}
		if (mapping.Direction != "input" && mapping.Direction != "output") || !shapeText(mapping.Field, 1024) || !shapeText(mapping.Location, 128) || !shapeText(mapping.Name, 1024) || mapping.Direction == "output" && mapping.Required != nil || !mappingSchemaValid(ctx, mapping.Schema) {
			if ctx.Err() != nil {
				return report, ctx.Err()
			}
			return OperationMappingReport{}, ErrOperationShapeTable
		}
		match := OperationFieldMatch{Mapping: mapping, Outcome: binding.Indeterminate, Code: "mapping.operation_unresolved"}
		key := mapping.Direction + "\x00" + mapping.Location + "\x00" + mapping.Name
		duplicateKey := key
		if mapping.Direction == "output" {
			duplicateKey = "output\x00" + mapping.Field
		}
		switch {
		case seen[duplicateKey]:
			match.Outcome, match.Code = binding.Incompatible, "mapping.duplicate"
		case eligibility.Shape != nil:
			match, err = matchOperationField(ctx, *eligibility.Shape, mapping)
			if err != nil {
				return report, err
			}
		}
		seen[duplicateKey] = true
		if mapping.Direction == "input" && match.OperationRequired != nil {
			supplied[key] = true
		}
		report.Fields = append(report.Fields, match)
		report.Outcome = mappingOutcome(report.Outcome, match.Outcome)
	}
	if shape := eligibility.Shape; shape != nil {
		for _, input := range shape.Inputs {
			ref := OperationFieldReference{Direction: "input", Location: input.Location, Name: input.Name}
			if excludedMappingField(*shape, ref, input.Schema) {
				report.ExcludedInputs = append(report.ExcludedInputs, ref)
				continue
			}
			if input.Required && !supplied["input\x00"+input.Location+"\x00"+input.Name] {
				report.UnmappedRequired = append(report.UnmappedRequired, ref)
				report.Outcome = mappingOutcome(report.Outcome, binding.Incompatible)
			}
		}
	}
	return report, ctx.Err()
}

func matchOperationField(ctx context.Context, shape binding.OperationShape, mapping OperationFieldMapping) (OperationFieldMatch, error) {
	match := OperationFieldMatch{Mapping: mapping, Outcome: binding.Indeterminate, Code: "mapping.field_missing"}
	var schema binding.Schema
	found := false
	if mapping.Direction == "input" {
		for _, input := range shape.Inputs {
			if input.Location == mapping.Location && input.Name == mapping.Name {
				schema, found = input.Schema, true
				required := input.Required
				match.OperationRequired = &required
			}
		}
	} else {
		for _, output := range shape.Outputs {
			if output.Location == mapping.Location && output.Name == mapping.Name {
				schema, found = output.Schema, true
			}
		}
	}
	switch {
	case !found:
		if shape.Complete {
			match.Outcome = binding.Incompatible
		}
	case excludedMappingField(shape, mapping.OperationFieldReference, schema) || looksLikeCredentialName(mapping.Field) || schemaHasCredentialFields(mapping.Schema.JSON):
		match.Code = "mapping.credential_excluded"
	case mapping.Placeholder:
		match.Code = "mapping.placeholder"
	case match.OperationRequired != nil && *match.OperationRequired && mapping.Required == nil:
		match.Code = "mapping.presence_unknown"
	case match.OperationRequired != nil && *match.OperationRequired && !*mapping.Required:
		match.Outcome, match.Code = binding.Incompatible, "mapping.required"
	case !schema.Known || !mapping.Schema.Known:
		match.Code = "mapping.schema_unknown"
	default:
		source, target := mapping.Schema, schema
		if mapping.Direction == "output" {
			source, target = schema, mapping.Schema
		}
		var err error
		match.Outcome, err = mappingSchemaContainment(ctx, shape, source, target)
		match.Code = "mapping.schema_" + string(match.Outcome)
		return match, err
	}
	return match, nil
}

func mappingSchemaValid(ctx context.Context, schema binding.Schema) bool {
	if len(schema.JSON) > binding.MaxSchemaBytes || schema.Known && len(schema.JSON) == 0 {
		return false
	}
	if len(schema.JSON) == 0 {
		return true
	}
	if sourceguard.CheckJSONWithLimits(ctx, "mapping-schema", schema.JSON, sourceguard.DefaultLimits()) != nil {
		return false
	}
	var value any
	if json.Unmarshal(schema.JSON, &value) != nil {
		return false
	}
	switch value.(type) {
	case map[string]any, bool:
		return true
	}
	return false
}

func excludedMappingField(shape binding.OperationShape, ref OperationFieldReference, schema binding.Schema) bool {
	if looksLikeCredentialName(ref.Name) || schemaHasCredentialFields(schema.JSON) {
		return true
	}
	if ref.Location == "header" {
		lower := strings.ToLower(ref.Name)
		switch lower {
		case "authorization", "proxy-authorization", "cookie", "set-cookie", "host", "content-type", "content-length", "accept", "user-agent":
			return true
		}
		// Match the conservative header policy already used by resource
		// authoring. Session/signature and opaque key headers may carry secrets
		// even when a source omits their security scheme.
		for _, word := range []string{"token", "secret", "password", "passwd", "key", "auth", "credential", "session", "signature"} {
			if strings.Contains(lower, word) {
				return true
			}
		}
	}
	if ref.Direction == "input" {
		for _, alt := range shape.Security.Alternatives {
			for _, requirement := range alt.Requirements {
				if ref.Location == requirement.Location && (ref.Name == requirement.Name || ref.Location == "header" && strings.EqualFold(ref.Name, requirement.Name)) {
					return true
				}
			}
		}
	}
	return false
}

func schemaHasCredentialFields(raw json.RawMessage) bool {
	var value any
	if len(raw) == 0 || json.Unmarshal(raw, &value) != nil {
		return false
	}
	stack := []any{value}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		switch node := current.(type) {
		case map[string]any:
			if properties, ok := node["properties"].(map[string]any); ok {
				for name := range properties {
					if looksLikeCredentialName(name) {
						return true
					}
				}
			}
			for _, child := range node {
				stack = append(stack, child)
			}
		case []any:
			stack = append(stack, node...)
		}
	}
	return false
}

// This private one-field projection asks the existing UWS validator only about
// schema containment. It is never returned as a native shape or an auth claim.
func mappingSchemaContainment(ctx context.Context, native binding.OperationShape, source, target binding.Schema) (binding.Outcome, error) {
	shape := binding.OperationShape{Source: native.Source, Selector: native.Selector, Protocol: "schema-comparison", Complete: true,
		Inputs:   []binding.Input{{Location: "mapping", Name: "value", Required: true, Schema: target}},
		Security: binding.Security{Known: true, Alternatives: []binding.SecurityAlternative{{}}}}
	resolver, err := binding.NewResolver(binding.ShapeTable{Version: binding.TableVersion, Sources: []binding.Source{shape.Source}, Operations: []binding.OperationShape{shape}})
	if err != nil {
		return binding.Indeterminate, ErrOperationShapeTable
	}
	const expression = "$inputs.value"
	report, err := binding.ValidateBinding(ctx, resolver, binding.Request{
		Binding: binding.Binding{Source: shape.Source, SelectorKind: shape.Selector.Kind, SelectorValue: shape.Selector.Value},
		Inputs:  []binding.BoundInput{{Location: "mapping", Name: "value", Value: expression}}, ExpressionTypes: map[string]binding.Schema{expression: source},
	})
	return report.Outcome, err
}

func mappingOutcome(a, b binding.Outcome) binding.Outcome {
	if a == binding.Incompatible || b == binding.Incompatible {
		return binding.Incompatible
	}
	if a == binding.Indeterminate || b == binding.Indeterminate {
		return binding.Indeterminate
	}
	return binding.Compatible
}
