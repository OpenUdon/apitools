package apitools

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

const DefaultMaxSchemaResolutionWork = 4_096

func requestBodySummary(ctx context.Context, root map[string]any, operation map[string]any, op *OperationSummary) (*RequestBodySummary, error) {
	if body := mapValue(operation["requestBody"]); len(body) > 0 {
		if ref := stringValue(body["$ref"]); strings.HasPrefix(ref, "#/") {
			if resolved := localObjectRef(root, ref); len(resolved) > 0 {
				body = mergeObjectRef(resolved, body)
			} else {
				addOperationIssue(op, "schema.ref_unresolved", "requestBody reference was not resolved", ref)
			}
		}
		summary := &RequestBodySummary{
			Description: stringValue(body["description"]),
			Required:    boolValue(body["required"]),
			Ref:         stringValue(body["$ref"]),
		}
		if summary.Ref != "" {
			addOperationIssue(op, "schema.ref_unresolved", "requestBody references a schema that was not resolved", summary.Ref)
		}
		content := mapValue(body["content"])
		summary.ContentTypes = sortedMapKeys(content)
		if len(summary.ContentTypes) > 1 {
			addOperationIssue(op, "schema.content_type_selection_required", "request body declares multiple media types; only one selected request schema is summarized", "")
		}
		if len(summary.ContentTypes) > 0 {
			media := mapValue(content[summary.ContentTypes[0]])
			rawSchema, err := resolveLocalSchemaRefs(ctx, root, mapValue(media["schema"]), op)
			if err != nil {
				return nil, err
			}
			schema := schemaSummary(rawSchema)
			summary.Schema = &schema
			summary.Fields = requestFieldSummaries(rawSchema, "", summary.Required, 0)
			if requestFieldSummaryWouldTruncate(rawSchema) {
				addOperationIssue(op, "schema.field_summary_truncated", "request-body fields exceeded the local summary count or depth bound", "")
			}
			if len(summary.Fields) == 0 && len(rawSchema) > 0 && !looksLikeCredentialName("body") {
				summary.Fields = []RequestFieldSummary{requestFieldSummary("body", summary.Required, rawSchema)}
			}
			summary.RequiredFieldPaths = requiredRequestFieldPaths(summary.Fields)
			if schema.Ref != "" {
				addOperationIssue(op, "schema.ref_unresolved", "request body schema reference was not resolved", schema.Ref)
			}
		}
		return summary, nil
	}
	for _, parameterValue := range sliceValue(operation["parameters"]) {
		parameter := mapValue(parameterValue)
		if ref := stringValue(parameter["$ref"]); strings.HasPrefix(ref, "#/") {
			if resolved := localObjectRef(root, ref); len(resolved) > 0 {
				parameter = mergeObjectRef(resolved, parameter)
			} else {
				addOperationIssue(op, "schema.ref_unresolved", "body parameter reference was not resolved", ref)
			}
		}
		if stringValue(parameter["in"]) != "body" {
			continue
		}
		rawSchema, err := resolveLocalSchemaRefs(ctx, root, mapValue(parameter["schema"]), op)
		if err != nil {
			return nil, err
		}
		schema := schemaSummary(rawSchema)
		if schema.Ref != "" {
			addOperationIssue(op, "schema.ref_unresolved", "body parameter schema reference was not resolved", schema.Ref)
		}
		fields := requestFieldSummaries(rawSchema, "", boolValue(parameter["required"]), 0)
		if requestFieldSummaryWouldTruncate(rawSchema) {
			addOperationIssue(op, "schema.field_summary_truncated", "request-body fields exceeded the local summary count or depth bound", "")
		}
		if len(fields) == 0 && len(rawSchema) > 0 && !looksLikeCredentialName("body") {
			fields = []RequestFieldSummary{requestFieldSummary("body", boolValue(parameter["required"]), rawSchema)}
		}
		return &RequestBodySummary{
			Description:        stringValue(parameter["description"]),
			Required:           boolValue(parameter["required"]),
			Schema:             &schema,
			Ref:                stringValue(parameter["$ref"]),
			Fields:             fields,
			RequiredFieldPaths: requiredRequestFieldPaths(fields),
		}, nil
	}
	return nil, nil
}

func responseBodySummary(ctx context.Context, root map[string]any, operation map[string]any, op *OperationSummary) (*ResponseBodySummary, error) {
	responses := mapValue(operation["responses"])
	if len(responses) == 0 {
		return nil, nil
	}
	statuses := successfulResponseStatuses(responses)
	if len(statuses) > 1 {
		addOperationIssue(op, "schema.response_selection_required", "operation declares multiple successful response statuses; only one response schema is summarized", "")
	}
	for _, status := range statuses {
		response := mapValue(responses[status])
		if len(response) == 0 {
			continue
		}
		if ref := stringValue(response["$ref"]); strings.HasPrefix(ref, "#/") {
			if resolved := localObjectRef(root, ref); len(resolved) > 0 {
				response = mergeObjectRef(resolved, response)
			} else {
				addOperationIssue(op, "schema.ref_unresolved", "response reference was not resolved", ref)
			}
		}
		summary := &ResponseBodySummary{
			StatusCode:  status,
			Description: stringValue(response["description"]),
			Ref:         stringValue(response["$ref"]),
		}
		if summary.Ref != "" {
			addOperationIssue(op, "schema.ref_unresolved", "response references a schema that was not resolved", summary.Ref)
		}
		if content := mapValue(response["content"]); len(content) > 0 {
			summary.ContentTypes = sortedMapKeys(content)
			if len(summary.ContentTypes) > 1 {
				addOperationIssue(op, "schema.response_content_type_selection_required", "successful response declares multiple media types; only one selected response schema is summarized", "")
			}
			for _, contentType := range preferredResponseContentTypes(summary.ContentTypes) {
				media := mapValue(content[contentType])
				rawSchema, err := resolveLocalSchemaRefs(ctx, root, mapValue(media["schema"]), op)
				if err != nil {
					return nil, err
				}
				if len(rawSchema) == 0 {
					continue
				}
				recordNullableResponseSchema(op, rawSchema)
				schema := schemaSummary(rawSchema)
				summary.Schema = &schema
				summary.Fields = responseFieldSummaries(rawSchema, true)
				if requestFieldSummaryWouldTruncate(rawSchema) {
					addOperationIssue(op, "schema.field_summary_truncated", "response fields exceeded the local summary count or depth bound", "")
				}
				if len(summary.Fields) == 0 && !looksLikeCredentialName("body") {
					summary.Fields = []RequestFieldSummary{requestFieldSummary("body", false, rawSchema)}
				}
				if schema.Ref != "" {
					addOperationIssue(op, "schema.ref_unresolved", "response schema reference was not resolved", schema.Ref)
				}
				return summary, nil
			}
			if len(summary.ContentTypes) > 0 {
				return summary, nil
			}
			continue
		}
		rawSchema, err := resolveLocalSchemaRefs(ctx, root, mapValue(response["schema"]), op)
		if err != nil {
			return nil, err
		}
		if len(rawSchema) == 0 {
			continue
		}
		recordNullableResponseSchema(op, rawSchema)
		schema := schemaSummary(rawSchema)
		summary.Schema = &schema
		summary.Fields = responseFieldSummaries(rawSchema, true)
		if requestFieldSummaryWouldTruncate(rawSchema) {
			addOperationIssue(op, "schema.field_summary_truncated", "response fields exceeded the local summary count or depth bound", "")
		}
		if len(summary.Fields) == 0 && !looksLikeCredentialName("body") {
			summary.Fields = []RequestFieldSummary{requestFieldSummary("body", false, rawSchema)}
		}
		if schema.Ref != "" {
			addOperationIssue(op, "schema.ref_unresolved", "response schema reference was not resolved", schema.Ref)
		}
		return summary, nil
	}
	return nil, nil
}

// Preserve an inventory warning while field summaries carry the exact
// nullable state needed for selected-output matching.
func recordNullableResponseSchema(op *OperationSummary, schema map[string]any) {
	visited := 0
	if !responseSchemaMayBeNullable(schema, 0, &visited) {
		return
	}
	op.ReadinessIssues = append(op.ReadinessIssues, ReadinessIssue{
		Severity:    "warning",
		Code:        "schema.response_nullable",
		Message:     "One or more response fields or schema ancestors permit null; inspect field-level nullable metadata before selecting outputs.",
		OperationID: op.OperationID,
		Path:        op.Provenance,
		Remediation: "Review the nullable response before requiring a non-null output.",
	})
}

func responseSchemaMayBeNullable(schema map[string]any, depth int, visited *int) bool {
	if len(schema) == 0 || depth > maxRequestFieldDepth || *visited >= maxRequestFields+1 {
		return false
	}
	(*visited)++
	if boolValue(schema["nullable"]) {
		return true
	}
	for _, name := range sortedMapKeys(mapValue(schema["properties"])) {
		if responseSchemaMayBeNullable(mapValue(mapValue(schema["properties"])[name]), depth+1, visited) {
			return true
		}
	}
	return responseSchemaMayBeNullable(mapValue(schema["items"]), depth+1, visited)
}

func successfulResponseStatuses(responses map[string]any) []string {
	var statuses []string
	for _, status := range sortedMapKeys(responses) {
		trimmed := strings.TrimSpace(status)
		if len(trimmed) == 3 && trimmed[0] == '2' {
			statuses = append(statuses, trimmed)
		}
	}
	return statuses
}

func preferredResponseContentTypes(contentTypes []string) []string {
	out := append([]string(nil), contentTypes...)
	sort.SliceStable(out, func(i, j int) bool {
		return responseContentTypeScore(out[i]) > responseContentTypeScore(out[j])
	})
	return out
}

func responseContentTypeScore(contentType string) int {
	contentType = strings.ToLower(strings.TrimSpace(contentType))
	switch {
	case strings.Contains(contentType, "json"):
		return 3
	case strings.Contains(contentType, "yaml") || strings.Contains(contentType, "xml"):
		return 2
	default:
		return 1
	}
}

type schemaResolutionState struct {
	ctx            context.Context
	root           map[string]any
	op             *OperationSummary
	activeRefs     map[string]bool
	work           int
	limitReported  bool
	branchReported bool
	depthReported  bool
}

func resolveLocalSchemaRefs(ctx context.Context, root map[string]any, schema map[string]any, op *OperationSummary) (map[string]any, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	state := &schemaResolutionState{ctx: ctx, root: root, op: op, activeRefs: map[string]bool{}}
	return state.resolve(schema, 0)
}

func (state *schemaResolutionState) resolve(schema map[string]any, depth int) (map[string]any, error) {
	if err := state.ctx.Err(); err != nil {
		return nil, err
	}
	if len(schema) == 0 {
		return schema, nil
	}
	if depth > maxRequestFieldDepth {
		if !state.depthReported {
			state.depthReported = true
			addOperationIssue(state.op, "schema.resolution_depth", fmt.Sprintf("schema resolution exceeded depth %d", maxRequestFieldDepth), "")
		}
		return copyStringAnyMap(schema), nil
	}
	state.work++
	if state.work > DefaultMaxSchemaResolutionWork {
		if !state.limitReported {
			state.limitReported = true
			addOperationIssue(state.op, "schema.resolution_limit", fmt.Sprintf("schema resolution exceeded %d visited nodes", DefaultMaxSchemaResolutionWork), "")
		}
		return copyStringAnyMap(schema), nil
	}

	out := copyStringAnyMap(schema)
	if ref := stringValue(out["$ref"]); strings.HasPrefix(ref, "#/") {
		if state.activeRefs[ref] {
			addOperationIssue(state.op, "schema.ref_cycle", "schema reference cycle was not expanded", ref)
			return out, nil
		}
		if resolved := localObjectRef(state.root, ref); len(resolved) > 0 {
			state.activeRefs[ref] = true
			base, err := state.resolve(resolved, depth+1)
			delete(state.activeRefs, ref)
			if err != nil {
				return nil, err
			}
			out = mergeObjectRef(base, out)
		} else {
			addOperationIssue(state.op, "schema.ref_unresolved", "schema reference was not resolved", ref)
		}
	}

	if branches := sliceValue(out["allOf"]); len(branches) > 0 {
		delete(out, "allOf")
		for _, branchValue := range branches {
			branch, err := state.resolve(mapValue(branchValue), depth+1)
			if err != nil {
				return nil, err
			}
			out = mergeSchema(out, branch, true)
		}
	}
	for _, keyword := range []string{"oneOf", "anyOf"} {
		branches := sliceValue(out[keyword])
		if len(branches) == 0 {
			continue
		}
		if !state.branchReported {
			state.branchReported = true
			addOperationIssue(state.op, "schema.branch_selection_required", keyword+" requires an explicit schema branch selection", "")
		}
		delete(out, keyword)
		for _, branchValue := range branches {
			branch, err := state.resolve(mapValue(branchValue), depth+1)
			if err != nil {
				return nil, err
			}
			out = mergeSchema(out, branch, false)
		}
	}

	if properties := mapValue(out["properties"]); len(properties) > 0 {
		resolvedProperties := make(map[string]any, len(properties))
		for _, name := range sortedMapKeys(properties) {
			resolved, err := state.resolve(mapValue(properties[name]), depth+1)
			if err != nil {
				return nil, err
			}
			resolvedProperties[name] = resolved
		}
		out["properties"] = resolvedProperties
	}
	if items := mapValue(out["items"]); len(items) > 0 {
		resolved, err := state.resolve(items, depth+1)
		if err != nil {
			return nil, err
		}
		out["items"] = resolved
	}
	return out, nil
}

func localObjectRef(root map[string]any, ref string) map[string]any {
	if !strings.HasPrefix(ref, "#/") {
		return nil
	}
	var current any = root
	for _, encoded := range strings.Split(strings.TrimPrefix(ref, "#/"), "/") {
		segment, err := url.PathUnescape(encoded)
		if err != nil {
			return nil
		}
		key := strings.ReplaceAll(strings.ReplaceAll(segment, "~1", "/"), "~0", "~")
		object, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current, ok = object[key]
		if !ok {
			return nil
		}
	}
	return mapValue(current)
}

func mergeObjectRef(base, overlay map[string]any) map[string]any {
	out := copyStringAnyMap(base)
	for key, value := range overlay {
		if key != "$ref" {
			out[key] = value
		}
	}
	return out
}

func mergeSchema(base, addition map[string]any, includeRequired bool) map[string]any {
	out := copyStringAnyMap(base)
	for key, value := range addition {
		switch key {
		case "properties":
			properties := copyStringAnyMap(mapValue(out[key]))
			for name, property := range mapValue(value) {
				properties[name] = property
			}
			out[key] = properties
		case "required":
			if includeRequired {
				out[key] = sortedUniqueStrings(append(stringSlice(out[key]), stringSlice(value)...))
			}
		default:
			if _, exists := out[key]; !exists {
				out[key] = value
			}
		}
	}
	return out
}

func copyStringAnyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for key, value := range in {
		out[key] = value
	}
	return out
}

func schemaSummary(schema map[string]any) SchemaSummary {
	if len(schema) == 0 {
		return SchemaSummary{}
	}
	required := stringSlice(schema["required"])
	requiredSet := make(map[string]bool, len(required))
	for _, name := range required {
		requiredSet[name] = true
	}
	summary := SchemaSummary{
		Type:        schemaType(schema["type"]),
		Nullable:    boolValue(schema["nullable"]),
		Format:      stringValue(schema["format"]),
		Ref:         stringValue(schema["$ref"]),
		Description: stringValue(schema["description"]),
		Required:    required,
	}
	properties := mapValue(schema["properties"])
	for _, name := range sortedMapKeys(properties) {
		property := schemaSummary(mapValue(properties[name]))
		summary.Properties = append(summary.Properties, PropertySummary{
			Name:        name,
			Type:        property.Type,
			Format:      property.Format,
			Ref:         property.Ref,
			Description: property.Description,
			Required:    requiredSet[name],
		})
	}
	return summary
}

func schemaType(value any) string {
	if text := stringValue(value); text != "" {
		return text
	}
	values := stringSlice(value)
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, "|")
}

const (
	maxRequestFieldDepth = 6
	maxRequestFields     = 60
)

// maxResponseFieldWork bounds the nested response fields collected before the
// final selection. Top-level fields are always collected.
const maxResponseFieldWork = 1000

type collectedField struct {
	field RequestFieldSummary
	depth int
}

// fieldCollector walks a schema depth-first in property-name order. With
// exemptTop set, top-level fields do not count toward limit, so a large nested
// field cannot hide later top-level fields.
type fieldCollector struct {
	limit     int
	exemptTop bool
	counted   int
	out       []collectedField
}

func (c *fieldCollector) add(field RequestFieldSummary, depth int) {
	c.out = append(c.out, collectedField{field: field, depth: depth})
	if !(c.exemptTop && depth <= 1) {
		c.counted++
	}
}

func (c *fieldCollector) full(depth int) bool {
	return c.counted >= c.limit && !(c.exemptTop && depth <= 1)
}

func requestFieldSummaries(schema map[string]any, path string, required bool, depth int) []RequestFieldSummary {
	c := fieldCollector{limit: maxRequestFields}
	collectRequestFields(schema, path, required, false, depth, &c)
	out := make([]RequestFieldSummary, 0, len(c.out))
	for _, item := range c.out {
		out = append(out, item.field)
	}
	sort.SliceStable(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	if len(out) > maxRequestFields {
		out = out[:maxRequestFields]
	}
	return out
}

// responseFieldSummaries summarizes response fields within the existing depth
// bound. When more than maxRequestFields remain, it keeps every top-level field
// first (by path, at most maxRequestFields), then nested fields by path, and
// returns the kept fields sorted by path.
func responseFieldSummaries(schema map[string]any, required bool) []RequestFieldSummary {
	c := fieldCollector{limit: maxResponseFieldWork, exemptTop: true}
	collectRequestFields(schema, "", required, false, 0, &c)
	sort.SliceStable(c.out, func(i, j int) bool {
		return c.out[i].field.Path < c.out[j].field.Path
	})
	if len(c.out) <= maxRequestFields {
		out := make([]RequestFieldSummary, 0, len(c.out))
		for _, item := range c.out {
			out = append(out, item.field)
		}
		return out
	}
	keep := make([]bool, len(c.out))
	kept := 0
	for i, item := range c.out {
		if item.depth <= 1 && kept < maxRequestFields {
			keep[i] = true
			kept++
		}
	}
	for i, item := range c.out {
		if kept >= maxRequestFields {
			break
		}
		if !keep[i] && item.depth > 1 {
			keep[i] = true
			kept++
		}
	}
	out := make([]RequestFieldSummary, 0, kept)
	for i, item := range c.out {
		if keep[i] {
			out = append(out, item.field)
		}
	}
	return out
}

func collectRequestFields(schema map[string]any, path string, required, nullable bool, depth int, c *fieldCollector) {
	if c.full(depth) || depth > maxRequestFieldDepth {
		return
	}
	nullable = nullable || boolValue(schema["nullable"])
	if len(schema) == 0 {
		if path != "" && !looksLikeCredentialName(path) {
			c.add(RequestFieldSummary{Path: path, Required: required, Nullable: nullable}, depth)
		}
		return
	}
	if path != "" && !looksLikeCredentialName(path) {
		field := requestFieldSummary(path, required, schema)
		field.Nullable = nullable
		c.add(field, depth)
		if c.counted >= c.limit || depth == maxRequestFieldDepth {
			return
		}
	}
	properties := mapValue(schema["properties"])
	if len(properties) > 0 {
		requiredSet := make(map[string]bool)
		for _, name := range stringSlice(schema["required"]) {
			requiredSet[name] = true
		}
		for _, name := range sortedMapKeys(properties) {
			childPath := name
			if path != "" {
				childPath = path + "." + name
			}
			collectRequestFields(mapValue(properties[name]), childPath, required && requiredSet[name], nullable, depth+1, c)
			if c.counted >= c.limit && !(c.exemptTop && depth == 0) {
				return
			}
		}
		return
	}
	if items := mapValue(schema["items"]); len(items) > 0 {
		itemPath := "body[]"
		if path != "" {
			itemPath = path + "[]"
		}
		collectRequestFields(items, itemPath, required, nullable, depth+1, c)
	}
}

func requestFieldSummary(path string, required bool, schema map[string]any) RequestFieldSummary {
	return RequestFieldSummary{
		Path:        path,
		Required:    required,
		Nullable:    boolValue(schema["nullable"]),
		Type:        schemaType(schema["type"]),
		Format:      stringValue(schema["format"]),
		Ref:         stringValue(schema["$ref"]),
		Description: stringValue(schema["description"]),
	}
}

func requestFieldSummaryWouldTruncate(schema map[string]any) bool {
	work := 0
	depthExceeded := false
	var visit func(map[string]any, int)
	visit = func(current map[string]any, depth int) {
		if depth > maxRequestFieldDepth {
			depthExceeded = true
			return
		}
		if work > maxRequestFields {
			return
		}
		properties := mapValue(current["properties"])
		for _, name := range sortedMapKeys(properties) {
			work++
			if work > maxRequestFields {
				return
			}
			visit(mapValue(properties[name]), depth+1)
		}
		if items := mapValue(current["items"]); len(items) > 0 {
			work++
			visit(items, depth+1)
		}
	}
	visit(schema, 0)
	return work > maxRequestFields || depthExceeded
}

func requiredRequestFieldPaths(fields []RequestFieldSummary) []string {
	var out []string
	for _, field := range fields {
		if field.Required {
			out = append(out, field.Path)
		}
	}
	sort.Strings(out)
	return out
}
