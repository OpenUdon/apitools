package apitools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/OpenUdon/apitools/graphql"
	"github.com/OpenUdon/apitools/grpcproto"
	"github.com/OpenUdon/apitools/internal/sourceguard"
	"github.com/OpenUdon/apitools/odata"
	"github.com/OpenUdon/apitools/openrpc"
	asyncapiparser "github.com/OpenUdon/asyncapi"
	smithyparser "github.com/OpenUdon/awssmithy"
	googleparser "github.com/OpenUdon/googlediscovery"
)

type nativeSchemaResolver func(string) (map[string]any, bool)

type nativeSchemaWalk struct {
	fields    []RequestFieldSummary
	gaps      []string
	active    map[string]bool
	work      int
	maxFields int
	resolver  nativeSchemaResolver
}

func adaptGoogleDiscoverySource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := googleparser.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	result := operationSourceAdapterResult{title: model.Title, operationCount: len(model.Operations)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{"Only declared Discovery parameters and locally available referenced request-schema fields are summarized; external or unsupported shapes remain gaps."}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{"Only locally available referenced Discovery response-schema fields are summarized; absent or external references remain gaps."}},
		{Dimension: "auth", Status: OperationCapabilityPartial, Gaps: []string{"Google Discovery OAuth scopes identify scope hints, not a complete auth or account-selection policy."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"Read/write effect is inferred only from operation meaning; the HTTP method is not evidence."}},
	}
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native == nil || strings.TrimSpace(native.ID) == "" {
			continue
		}
		selector := "#/methods/" + escapeJSONPointer(native.ID)
		operation := nativeOperationSummary(source, native.ID, selector, native.HTTPMethod, native.Path, native.Summary, native.Description)
		operation.OperationID = native.OperationID
		operation.Tags = append([]string(nil), native.Tags...)
		for _, parameter := range native.Parameters {
			if parameter == nil {
				continue
			}
			typeName, format := schemaTypeAndFormat(parameter.Schema)
			operation.Parameters = append(operation.Parameters, ParameterSummary{
				Name: firstNonEmpty(parameter.Name, parameter.OriginalName), In: parameter.Location,
				Required: parameter.Required, Description: parameter.Description, Type: typeName, Format: format,
			})
		}
		var inputGaps, outputGaps []string
		inputStatus, outputStatus := OperationCapabilitySupported, OperationCapabilitySupported
		if native.RequestRef != "" {
			schema, ok := model.Schemas[native.RequestRef]
			if ok {
				fields, gaps, fieldErr := nativeSchemaFields(schema, budget, func(ref string) (map[string]any, bool) {
					resolved, found := model.Schemas[ref]
					return resolved, found
				})
				if fieldErr != nil {
					return result, fieldErr
				}
				inputGaps = append(inputGaps, gaps...)
				body := &RequestBodySummary{Required: true, Schema: schemaPointer(schema), Fields: fields, Ref: native.RequestRef}
				operation.RequestBody = body
			} else {
				inputStatus = OperationCapabilityPartial
				operation.ReadinessIssues = append(operation.ReadinessIssues, ReadinessIssue{Severity: "error", Code: "schema.ref_unresolved", Message: "Google Discovery request schema reference is not present in this document", Path: native.RequestRef})
			}
			if len(inputGaps) > 0 {
				inputStatus = OperationCapabilityPartial
			}
		}
		if native.ResponseRef != "" {
			schema, ok := model.Schemas[native.ResponseRef]
			if ok {
				fields, gaps, fieldErr := nativeSchemaFields(schema, budget, func(ref string) (map[string]any, bool) {
					resolved, found := model.Schemas[ref]
					return resolved, found
				})
				if fieldErr != nil {
					return result, fieldErr
				}
				outputGaps = append(outputGaps, gaps...)
				operation.ResponseBody = &ResponseBodySummary{Schema: schemaPointer(schema), Fields: fields, Ref: native.ResponseRef}
			} else {
				outputStatus = OperationCapabilityPartial
				operation.ReadinessIssues = append(operation.ReadinessIssues, ReadinessIssue{Severity: "error", Code: "schema.ref_unresolved", Message: "Google Discovery response schema reference is not present in this document", Path: native.ResponseRef})
			}
			if len(outputGaps) > 0 {
				outputStatus = OperationCapabilityPartial
			}
		} else {
			outputStatus = OperationCapabilityPartial
			outputGaps = append(outputGaps, "The Discovery method has no response schema reference.")
		}
		operation.SecurityRequirementSets = googleScopeRequirementSets(native.Scopes, model.OAuth2Scopes)
		capabilities := nativeCapabilities(
			inputStatus, outputStatus,
			OperationCapabilityPartial, OperationCapabilityPartial, selector,
		)
		capabilities[0].Gaps = append(capabilities[0].Gaps, inputGaps...)
		capabilities[1].Gaps = append(capabilities[1].Gaps, outputGaps...)
		if len(native.Scopes) > 0 {
			capabilities[2].Evidence = []OperationEvidence{{Kind: "google.oauth_scopes", Reference: selector}}
		}
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffect(operation, selector), capabilities, budget)
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, append(inputGaps, outputGaps...)...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func adaptSmithySource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := smithyparser.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	result := operationSourceAdapterResult{title: model.Title, operationCount: len(model.Operations)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilitySupported},
		{Dimension: "outputs", Status: OperationCapabilitySupported},
		{Dimension: "auth", Status: OperationCapabilityPartial, Gaps: []string{"Smithy protocol and signing traits do not establish a complete credential, account, or authorization policy."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"Only an explicit smithy.api#readonly trait or operation meaning can establish a read hint."}},
	}
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native == nil || strings.TrimSpace(native.ID) == "" {
			continue
		}
		selector := "#/shapes/" + escapeJSONPointer(native.ID)
		operation := nativeOperationSummary(source, native.ID, selector, model.Protocol, native.Path, "", "")
		var inputGaps, outputGaps []string
		inputStatus, outputStatus := OperationCapabilitySupported, OperationCapabilitySupported
		if input := model.Shapes[native.Input]; input != nil {
			fields, shapeGaps, fieldErr := smithyShapeFields(model, input, budget, "body")
			if fieldErr != nil {
				return result, fieldErr
			}
			inputGaps = append(inputGaps, shapeGaps...)
			operation.RequestBody = &RequestBodySummary{Required: true, Schema: &SchemaSummary{Type: "object", Ref: native.Input}, Fields: fields, Ref: native.Input}
		} else if native.Input != "" {
			inputStatus = OperationCapabilityPartial
			inputGaps = append(inputGaps, "The Smithy input shape is not available in the parsed document.")
		}
		if output := model.Shapes[native.Output]; output != nil {
			fields, shapeGaps, fieldErr := smithyShapeFields(model, output, budget, "body")
			if fieldErr != nil {
				return result, fieldErr
			}
			outputGaps = append(outputGaps, shapeGaps...)
			operation.ResponseBody = &ResponseBodySummary{Schema: &SchemaSummary{Type: "object", Ref: native.Output}, Fields: fields, Ref: native.Output}
		} else if native.Output != "" {
			outputStatus = OperationCapabilityPartial
			outputGaps = append(outputGaps, "The Smithy output shape is not available in the parsed document.")
		}
		if len(inputGaps) > 0 {
			inputStatus = OperationCapabilityPartial
		}
		if len(outputGaps) > 0 {
			outputStatus = OperationCapabilityPartial
		}
		var nativeSignals []effectSignal
		if operationShape, ok := model.Shape(native.ID); ok && operationShape != nil {
			if _, readonly := operationShape.Traits["smithy.api#readonly"]; readonly {
				nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectRead, kind: "smithy.readonly", reference: selector, word: "readonly"})
			}
		}
		capabilities := nativeCapabilities(inputStatus, outputStatus, OperationCapabilityPartial, OperationCapabilityPartial, selector)
		capabilities[0].Gaps = append(capabilities[0].Gaps, inputGaps...)
		capabilities[1].Gaps = append(capabilities[1].Gaps, outputGaps...)
		if len(nativeSignals) > 0 {
			capabilities[3].Evidence = []OperationEvidence{{Kind: nativeSignals[0].kind, Reference: selector}}
		}
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffectWithNative(operation, selector, nativeSignals), capabilities, budget)
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, append(inputGaps, outputGaps...)...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func nativeOperationSummary(source loadedOperationSource, id, selector, method, path, summary, description string) OperationSummary {
	return OperationSummary{
		ID: id, DocumentName: sourceName(source.input), DocumentPath: source.input.Path,
		DocumentURL: source.input.URL, OperationID: id, Method: strings.TrimSpace(method),
		Path: firstNonEmpty(path, selector), Summary: summary, Description: description,
		Provenance: selector,
	}
}

func nativeCapabilities(inputs, outputs, auth, effect OperationCapabilityStatus, selector string) []OperationCapability {
	return []OperationCapability{
		{Dimension: "inputs", Status: inputs, Evidence: []OperationEvidence{{Kind: "operation.inputs", Reference: selector}}},
		{Dimension: "outputs", Status: outputs, Evidence: []OperationEvidence{{Kind: "operation.outputs", Reference: selector}}},
		{Dimension: "auth", Status: auth, Evidence: []OperationEvidence{{Kind: "operation.auth", Reference: selector}}},
		{Dimension: "effect", Status: effect, Evidence: []OperationEvidence{{Kind: "operation.meaning", Reference: selector}}},
	}
}

func enforceNativeOperationBudget(result *operationSourceAdapterResult, budget PromptBudget) error {
	if result == nil || result.operationCount <= budget.MaxOperations {
		return nil
	}
	diagnostic := Diagnostic{
		Severity: "error", Code: "candidate.operation_work_budget",
		Message:     fmt.Sprintf("source contains %d operations, exceeding the configured %d-operation ranking budget", result.operationCount, budget.MaxOperations),
		Remediation: "Narrow the explicit source set or raise MaxOperations for a reviewed workload.",
	}
	result.diagnostics = append(result.diagnostics, diagnostic)
	return DiagnosticError{Diagnostics: []Diagnostic{diagnostic}}
}

func schemaPointer(schema map[string]any) *SchemaSummary {
	if len(schema) == 0 {
		return nil
	}
	summary := schemaSummary(schema)
	return &summary
}

func googleScopeRequirementSets(scopes []string, descriptions map[string]string) []SecurityRequirementSetSummary {
	sets := make([]SecurityRequirementSetSummary, 0, len(scopes))
	for _, scope := range sortedUniqueStrings(append([]string(nil), scopes...)) {
		if strings.TrimSpace(scope) == "" {
			continue
		}
		sets = append(sets, SecurityRequirementSetSummary{Requirements: []SecuritySummary{{
			Name: "google-oauth2", Type: "oauth2", Scopes: []string{scope}, Description: descriptions[scope],
		}}})
	}
	return sets
}

func schemaTypeAndFormat(schema map[string]any) (string, string) {
	if len(schema) == 0 {
		return "", ""
	}
	typeName := schemaType(schema["type"])
	if typeName == "" && stringValue(schema["$ref"]) != "" {
		typeName = stringValue(schema["$ref"])
	}
	return typeName, stringValue(schema["format"])
}

func nativeSchemaFields(schema map[string]any, budget PromptBudget, resolver nativeSchemaResolver) ([]RequestFieldSummary, []string, error) {
	walk := &nativeSchemaWalk{active: map[string]bool{}, maxFields: budget.MaxFields, resolver: resolver}
	if err := walkSchemaFields(schema, "", true, 0, walk); err != nil {
		return nil, walk.gaps, err
	}
	sort.SliceStable(walk.fields, func(i, j int) bool { return walk.fields[i].Path < walk.fields[j].Path })
	return walk.fields, uniqueSortedStrings(walk.gaps), nil
}

func walkSchemaFields(schema map[string]any, prefix string, parentRequired bool, depth int, walk *nativeSchemaWalk) error {
	walk.work++
	if walk.work > 10_000 || depth > 32 {
		return fmt.Errorf("native schema exceeds the 10,000-node or 32-level summary bound")
	}
	if ref := firstNonEmpty(stringValue(schema["$ref"]), stringValue(schema["ref"])); ref != "" {
		if walk.resolver == nil || walk.active[ref] {
			walk.gaps = append(walk.gaps, "Schema reference "+ref+" is external, unresolved, or recursive and was not fetched.")
			return nil
		}
		resolved, ok := walk.resolver(ref)
		if !ok {
			walk.gaps = append(walk.gaps, "Schema reference "+ref+" is not available locally and was not fetched.")
			return nil
		}
		walk.active[ref] = true
		defer delete(walk.active, ref)
		merged := mergeObjectRef(resolved, schema)
		delete(merged, "$ref")
		delete(merged, "ref")
		return walkSchemaFields(merged, prefix, parentRequired, depth+1, walk)
	}
	for _, keyword := range []string{"allOf", "oneOf", "anyOf"} {
		if len(sliceValue(schema[keyword])) > 0 {
			walk.gaps = append(walk.gaps, "Schema "+keyword+" composition is not expanded into operation fields.")
		}
	}
	if additional, exists := schema["additionalProperties"]; exists {
		if allowed, isBoolean := additional.(bool); !isBoolean || allowed {
			walk.gaps = append(walk.gaps, "Schema additionalProperties dictionary values are not summarized as named operation fields.")
		}
	}
	properties := mapValue(schema["properties"])
	if len(properties) > 0 {
		requiredSet := make(map[string]bool)
		for _, name := range stringSlice(schema["required"]) {
			requiredSet[name] = true
		}
		for _, name := range sortedMapKeys(properties) {
			child := mapValue(properties[name])
			path := name
			if prefix != "" {
				path = prefix + "." + name
			}
			if looksLikeCredentialName(path) {
				walk.gaps = append(walk.gaps, "Credential-shaped schema field "+path+" is omitted from workflow data inputs and requires separate auth review.")
				continue
			}
			required := parentRequired && requiredSet[name]
			if explicit, exists := child["required"].(bool); exists {
				required = parentRequired && explicit
			}
			typeName, format := schemaTypeAndFormat(child)
			if err := addNativeField(walk, RequestFieldSummary{
				Path: path, Required: required, Type: typeName, Format: format,
				Ref: firstNonEmpty(stringValue(child["$ref"]), stringValue(child["ref"])), Description: stringValue(child["description"]),
			}); err != nil {
				return err
			}
			if err := walkSchemaFields(child, path, required, depth+1, walk); err != nil {
				return err
			}
		}
		return nil
	}
	if items := mapValue(schema["items"]); len(items) > 0 {
		path := prefix + "[]"
		if prefix == "" {
			path = "body[]"
		}
		typeName, format := schemaTypeAndFormat(items)
		if err := addNativeField(walk, RequestFieldSummary{Path: path, Required: parentRequired, Type: typeName, Format: format, Description: stringValue(items["description"])}); err != nil {
			return err
		}
		return walkSchemaFields(items, path, parentRequired, depth+1, walk)
	}
	return nil
}

func addNativeField(walk *nativeSchemaWalk, field RequestFieldSummary) error {
	if len(walk.fields) >= walk.maxFields {
		return fmt.Errorf("schema field summary exceeds the configured %d-field budget", walk.maxFields)
	}
	walk.fields = append(walk.fields, field)
	return nil
}

func smithyShapeFields(model *smithyparser.Model, shape *smithyparser.Shape, budget PromptBudget, location string) ([]RequestFieldSummary, []string, error) {
	if shape == nil {
		return nil, nil, nil
	}
	members := mapValue(shape.Raw["members"])
	if len(members) == 0 {
		return nil, nil, nil
	}
	fields := make([]RequestFieldSummary, 0, len(members))
	var gaps []string
	for _, name := range sortedMapKeys(members) {
		member := mapValue(members[name])
		target := stringValue(member["target"])
		typeName := ""
		format := ""
		if targetShape, ok := model.Shape(target); ok && targetShape != nil {
			typeName, format = smithyType(targetShape)
			switch targetShape.Type {
			case "structure", "union", "map", "list", "set":
				gaps = append(gaps, "Nested Smithy shape details for member "+name+" are not expanded.")
			}
		} else if target != "" {
			typeName = target
			gaps = append(gaps, "Smithy target shape "+target+" is not available for field "+name+".")
		}
		if looksLikeCredentialName(name) {
			gaps = append(gaps, "Credential-shaped Smithy member "+name+" is omitted from workflow data inputs and requires separate auth review.")
			continue
		}
		traits := mapValue(member["traits"])
		_, required := traits["smithy.api#required"]
		fields = append(fields, RequestFieldSummary{Path: name, Type: typeName, Format: format, Required: required})
		if len(fields) > budget.MaxFields {
			return nil, gaps, fmt.Errorf("Smithy shape exceeds the configured %d-field budget", budget.MaxFields)
		}
	}
	sort.SliceStable(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })
	_ = location
	return fields, gaps, nil
}

func smithyType(shape *smithyparser.Shape) (string, string) {
	if shape == nil {
		return "", ""
	}
	if strings.HasPrefix(shape.ID, "smithy.api#") {
		switch strings.TrimPrefix(shape.ID, "smithy.api#") {
		case "String", "Enum", "Unit":
			return "string", ""
		case "Boolean":
			return "boolean", ""
		case "Byte", "Short", "Integer", "Long", "BigInteger":
			return "integer", ""
		case "Float", "Double", "BigDecimal":
			return "number", ""
		case "Blob", "PrimitiveBlob":
			return "string", "byte"
		case "Timestamp":
			return "string", "date-time"
		case "List", "Set":
			return "array", ""
		case "Structure", "Union", "Map", "Document":
			return "object", ""
		}
	}
	switch shape.Type {
	case "string", "enum":
		return "string", ""
	case "boolean":
		return "boolean", ""
	case "byte", "short", "integer", "long", "bigInteger":
		return "integer", ""
	case "float", "double", "bigDecimal":
		return "number", ""
	case "blob":
		return "string", "byte"
	case "timestamp":
		return "string", "date-time"
	case "list", "set":
		return "array", ""
	case "structure", "union", "map", "document":
		return "object", ""
	default:
		return shape.ID, ""
	}
}

func adaptAsyncAPISource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	if err := sourceguard.CheckDocument("asyncapi", source.content); err != nil {
		return operationSourceAdapterResult{}, err
	}
	trimmed := strings.TrimSpace(string(source.content))
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		if err := sourceguard.CheckJSON("asyncapi", source.content); err != nil {
			return operationSourceAdapterResult{}, err
		}
	} else if err := sourceguard.CheckYAML("asyncapi", source.content); err != nil {
		return operationSourceAdapterResult{}, err
	}
	model, err := asyncapiparser.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	result := operationSourceAdapterResult{title: model.Title, operationCount: len(model.Operations)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{"Only locally available message payload metadata is summarized; channel bindings and external references are not resolved."}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{"Only locally available message payload metadata is summarized; channel bindings and external references are not resolved."}},
		{Dimension: "auth", Status: OperationCapabilityPartial, Gaps: []string{"AsyncAPI security requirements are not translated into workflow credential policy."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"AsyncAPI send/receive direction does not establish business read/write effect."}},
	}
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native == nil || strings.TrimSpace(native.ID) == "" {
			continue
		}
		selector := "#/operations/" + escapeJSONPointer(native.ID)
		operation := nativeOperationSummary(source, native.ID, selector, native.Action, native.ChannelRef, native.Summary, native.Description)
		var operationGaps []string
		if len(native.MessageRefs) > 1 {
			operationGaps = append(operationGaps, "This AsyncAPI operation declares multiple message payload alternatives; only the first available payload is summarized.")
		}
		for _, ref := range native.MessageRefs {
			target, ok := model.ResolveSelector(ref)
			if !ok || target.Message == nil || target.Message.Payload == nil {
				operationGaps = append(operationGaps, "AsyncAPI message payload reference "+ref+" is unavailable locally or has no payload schema.")
				continue
			}
			fields, gaps, fieldErr := nativeSchemaFields(target.Message.Payload, budget, func(localRef string) (map[string]any, bool) {
				return localJSONPointerMap(model.Raw, localRef)
			})
			if fieldErr != nil {
				return result, fieldErr
			}
			operationGaps = append(operationGaps, gaps...)
			body := &RequestBodySummary{Required: true, Schema: schemaPointer(target.Message.Payload), Fields: fields, Ref: ref}
			switch strings.ToLower(strings.TrimSpace(native.Action)) {
			case "send":
				if operation.RequestBody == nil {
					operation.RequestBody = body
				}
			case "receive":
				if operation.ResponseBody == nil {
					operation.ResponseBody = &ResponseBodySummary{Schema: body.Schema, Fields: body.Fields, Ref: ref}
				}
			default:
				operationGaps = append(operationGaps, "AsyncAPI operation action is not recognized as send or receive; payload direction is unresolved.")
			}
		}
		capabilities := nativeCapabilities(OperationCapabilityPartial, OperationCapabilityPartial, OperationCapabilityPartial, OperationCapabilityPartial, selector)
		capabilities[0].Gaps = append(capabilities[0].Gaps, operationGaps...)
		capabilities[1].Gaps = append(capabilities[1].Gaps, operationGaps...)
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffect(operation, selector), capabilities, budget)
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, operationGaps...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func adaptGraphQLSource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := graphql.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	result := operationSourceAdapterResult{title: sourceName(source.input), operationCount: len(model.Operations)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{"Only declared variables or root-field arguments are summarized; literal arguments and unresolved input objects are not inferred."}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{"Only selected root fields with locally inspectable schema types are summarized."}},
		{Dimension: "auth", Status: OperationCapabilityUnsupported, Gaps: []string{"GraphQL schema and operation documents do not establish endpoint authentication requirements."}},
		{Dimension: "effect", Status: OperationCapabilitySupported},
	}
	rootTypes := map[string]string{"query": model.Schema.QueryType, "mutation": model.Schema.MutationType, "subscription": model.Schema.SubscriptionType}
	types := make(map[string]*graphql.Type, len(model.Types))
	for _, item := range model.Types {
		if item != nil {
			types[item.Name] = item
		}
	}
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native == nil || strings.TrimSpace(native.ID) == "" {
			continue
		}
		selector := firstNonEmpty(native.Selector, graphql.OperationSelector(native.ID))
		operation := nativeOperationSummary(source, native.ID, selector, "GRAPHQL "+native.Kind, selector, native.Summary, native.Description)
		for _, variable := range native.Variables {
			if variable == nil {
				continue
			}
			typeName, format := graphqlContractType(variable.Type)
			operation.Parameters = append(operation.Parameters, ParameterSummary{Name: variable.Name, In: "variable", Required: variable.Required, Type: typeName, Format: format, Description: variable.Type.Display})
		}
		rootName := rootTypes[strings.ToLower(native.Kind)]
		root := types[rootName]
		if native.FieldName != "" {
			if field := graphqlRootField(root, native.FieldName); field != nil {
				for _, argument := range field.Args {
					if argument == nil {
						continue
					}
					typeName, format := graphqlContractType(argument.Type)
					operation.Parameters = append(operation.Parameters, ParameterSummary{Name: argument.Name, In: "argument", Required: argument.Required, Type: typeName, Format: format, Description: argument.Description})
				}
			}
		}
		for _, selected := range native.SelectionNames {
			field := graphqlRootField(root, selected)
			typeName, format := "", ""
			required := false
			if field != nil {
				typeName, format = graphqlContractType(field.Type)
				required = field.Type.Required
			}
			if operation.ResponseBody == nil {
				operation.ResponseBody = &ResponseBodySummary{Schema: &SchemaSummary{Type: "object"}}
			}
			operation.ResponseBody.Fields = append(operation.ResponseBody.Fields, RequestFieldSummary{Path: selected, Type: typeName, Format: format, Required: required})
		}
		if native.FieldName != "" {
			typeName, format := graphqlContractType(native.FieldType)
			operation.ResponseBody = &ResponseBodySummary{Schema: &SchemaSummary{Type: "object"}, Fields: []RequestFieldSummary{{Path: native.FieldName, Type: typeName, Format: format, Required: native.FieldType.Required}}}
		}
		nativeSignals := make([]effectSignal, 0, 1)
		switch strings.ToLower(native.Kind) {
		case "query":
			nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectRead, kind: "graphql.operation_kind", reference: selector, word: "query"})
		case "mutation":
			nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectWrite, kind: "graphql.operation_kind", reference: selector, word: "mutation"})
		}
		inputStatus := OperationCapabilityPartial
		if model.SourceKind == graphql.ArtifactKindIntrospectionJSON || model.SourceKind == graphql.ArtifactKindSDL {
			inputStatus = OperationCapabilitySupported
		}
		outputStatus := OperationCapabilityPartial
		if root != nil && (len(native.SelectionNames) > 0 || native.FieldName != "") {
			outputStatus = OperationCapabilitySupported
		}
		capabilities := nativeCapabilities(inputStatus, outputStatus, OperationCapabilityUnsupported, OperationCapabilitySupported, selector)
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffectWithNative(operation, selector, nativeSignals), capabilities, budget)
		for i := range candidate.Summary.Outputs {
			candidate.Summary.Outputs[i].Required = nil
		}
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		if model.SourceKind == graphql.ArtifactKindOperationDocument {
			candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, "GraphQL operation document alone does not provide endpoint auth policy or complete response schema."))
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func graphqlRootField(root *graphql.Type, name string) *graphql.Field {
	if root == nil {
		return nil
	}
	for _, field := range root.Fields {
		if field != nil && field.Name == name {
			return field
		}
	}
	return nil
}

func graphqlContractType(value graphql.TypeRef) (string, string) {
	name := strings.TrimSpace(value.Name)
	if name == "" {
		name = strings.Trim(value.Display, "![]")
	}
	if value.List || strings.Contains(value.Display, "[") {
		return "array", ""
	}
	switch strings.ToLower(name) {
	case "string", "id":
		return "string", ""
	case "int":
		return "integer", ""
	case "float":
		return "number", ""
	case "boolean":
		return "boolean", ""
	case "":
		return "", ""
	default:
		return name, ""
	}
}

func adaptOpenRPCSource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := openrpc.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	result := operationSourceAdapterResult{title: model.Info.Title, operationCount: len(model.Methods)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{"Only locally declared JSON-RPC parameter descriptors and supported JSON Schema fields are summarized."}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{"Only locally declared JSON-RPC result descriptors and supported JSON Schema fields are summarized."}},
		{Dimension: "auth", Status: OperationCapabilityUnsupported, Gaps: []string{"OpenRPC method metadata does not establish endpoint authentication requirements."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"JSON-RPC method names and schemas do not establish business read/write effect."}},
	}
	for _, native := range model.Methods {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native == nil || strings.TrimSpace(native.Name) == "" {
			continue
		}
		selector := openrpc.MethodSelector(native.Name)
		operation := nativeOperationSummary(source, native.Name, selector, "JSON-RPC", selector, native.Summary, native.Description)
		var operationGaps []string
		for _, parameter := range native.Params {
			if parameter == nil {
				continue
			}
			descriptor := openrpcDescriptor(model, parameter)
			typeName, format := schemaTypeAndFormat(descriptor.Schema)
			name := firstNonEmpty(descriptor.Name, parameter.Name, localReferenceName(descriptor.Ref))
			operation.Parameters = append(operation.Parameters, ParameterSummary{Name: name, In: "param", Required: descriptor.Required, Type: typeName, Format: format, Description: firstNonEmpty(descriptor.Summary, descriptor.Description)})
			if descriptor.Schema == nil {
				operationGaps = append(operationGaps, "JSON-RPC parameter "+name+" has no local schema descriptor.")
			}
			if descriptor.Ref != "" && descriptor.Schema == nil {
				operationGaps = append(operationGaps, "JSON-RPC parameter descriptor reference "+descriptor.Ref+" is unresolved locally.")
			}
		}
		if native.Result != nil {
			descriptor := openrpcDescriptor(model, native.Result)
			if descriptor.Schema != nil {
				fields, gaps, fieldErr := nativeSchemaFields(descriptor.Schema, budget, func(ref string) (map[string]any, bool) {
					return localJSONPointerMap(model.Raw, ref)
				})
				if fieldErr != nil {
					return result, fieldErr
				}
				operationGaps = append(operationGaps, gaps...)
				operation.ResponseBody = &ResponseBodySummary{Schema: schemaPointer(descriptor.Schema), Fields: fields, Ref: descriptor.Ref}
			} else {
				operation.ResponseBody = &ResponseBodySummary{Description: descriptor.Description, Ref: descriptor.Ref}
				operationGaps = append(operationGaps, "JSON-RPC result descriptor has no local result schema.")
			}
		}
		capabilities := nativeCapabilities(OperationCapabilityPartial, OperationCapabilityPartial, OperationCapabilityUnsupported, OperationCapabilityPartial, selector)
		capabilities[0].Gaps = append(capabilities[0].Gaps, operationGaps...)
		capabilities[1].Gaps = append(capabilities[1].Gaps, operationGaps...)
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffect(operation, selector), capabilities, budget)
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, operationGaps...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func openrpcDescriptor(model *openrpc.Model, descriptor *openrpc.ContentDescriptor) *openrpc.ContentDescriptor {
	if descriptor == nil {
		return &openrpc.ContentDescriptor{}
	}
	if descriptor.Ref != "" {
		name := localReferenceName(descriptor.Ref)
		if resolved := model.Components.ContentDescriptors[name]; resolved != nil {
			return resolved
		}
	}
	return descriptor
}

func localReferenceName(reference string) string {
	reference = strings.TrimSpace(reference)
	if index := strings.LastIndex(reference, "/"); index >= 0 && index < len(reference)-1 {
		reference = reference[index+1:]
	}
	return strings.ReplaceAll(strings.ReplaceAll(reference, "~1", "/"), "~0", "~")
}

func localJSONPointerMap(root map[string]any, reference string) (map[string]any, bool) {
	if !strings.HasPrefix(reference, "#/") {
		return nil, false
	}
	var current any = root
	for _, part := range strings.Split(strings.TrimPrefix(reference, "#/"), "/") {
		part = strings.ReplaceAll(strings.ReplaceAll(part, "~1", "/"), "~0", "~")
		object, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		current, ok = object[part]
		if !ok {
			return nil, false
		}
	}
	value, ok := current.(map[string]any)
	return value, ok
}

func adaptGRPCSource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := grpcproto.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	methods := model.MethodSummaries()
	result := operationSourceAdapterResult{title: sourceName(source.input), operationCount: len(methods)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{"Only top-level protobuf request fields are summarized; nested message shapes are opaque."}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{"Only top-level protobuf response fields are summarized; nested message shapes are opaque."}},
		{Dimension: "auth", Status: OperationCapabilityUnsupported, Gaps: []string{"Protobuf service metadata does not establish endpoint authentication requirements."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"gRPC method names and message shapes do not establish business read/write effect."}},
	}
	for _, nativeSummary := range methods {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		target, resolveErr := model.ResolveSelector(nativeSummary.Selector)
		if resolveErr != nil || target.Method == nil {
			continue
		}
		native := target.Method
		selector := nativeSummary.Selector
		operation := nativeOperationSummary(source, native.SourceOperationID, selector, "gRPC", nativeSummary.FullMethod, "", "")
		var operationGaps []string
		request := findProtoMessage(model, native.RequestType, native.Package)
		if request != nil {
			fields, fieldGaps, fieldErr := grpcMessageFields(request, budget)
			if fieldErr != nil {
				return result, fieldErr
			}
			operationGaps = append(operationGaps, fieldGaps...)
			operation.RequestBody = &RequestBodySummary{Required: true, Schema: &SchemaSummary{Type: "object", Ref: request.Selector}, Fields: fields, Ref: request.Selector}
		} else {
			operationGaps = append(operationGaps, "The gRPC request message shape is not available locally.")
		}
		response := findProtoMessage(model, native.ResponseType, native.Package)
		if response != nil {
			fields, fieldGaps, fieldErr := grpcMessageFields(response, budget)
			if fieldErr != nil {
				return result, fieldErr
			}
			operationGaps = append(operationGaps, fieldGaps...)
			operation.ResponseBody = &ResponseBodySummary{Schema: &SchemaSummary{Type: "object", Ref: response.Selector}, Fields: fields, Ref: response.Selector}
		} else {
			operationGaps = append(operationGaps, "The gRPC response message shape is not available locally.")
		}
		if native.ClientStreaming || native.ServerStreaming {
			operationGaps = append(operationGaps, "Streaming request or response semantics are not modeled by the step value contract.")
		}
		capabilities := nativeCapabilities(OperationCapabilityPartial, OperationCapabilityPartial, OperationCapabilityUnsupported, OperationCapabilityPartial, selector)
		capabilities[0].Gaps = append(capabilities[0].Gaps, operationGaps...)
		capabilities[1].Gaps = append(capabilities[1].Gaps, operationGaps...)
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffect(operation, selector), capabilities, budget)
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, operationGaps...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func findProtoMessage(model *grpcproto.Model, name, packageName string) *grpcproto.Message {
	name = strings.TrimPrefix(strings.TrimSpace(name), ".")
	packageName = strings.TrimSpace(packageName)
	for _, file := range model.Files {
		for _, message := range flattenProtoMessages(file.Messages) {
			if message == nil {
				continue
			}
			if message.FullName == name || message.Name == name || packageName != "" && message.FullName == packageName+"."+name {
				return message
			}
		}
	}
	return nil
}

func flattenProtoMessages(messages []*grpcproto.Message) []*grpcproto.Message {
	var out []*grpcproto.Message
	var visit func([]*grpcproto.Message)
	visit = func(items []*grpcproto.Message) {
		for _, item := range items {
			if item == nil {
				continue
			}
			out = append(out, item)
			visit(item.Messages)
		}
	}
	visit(messages)
	return out
}

func grpcMessageFields(message *grpcproto.Message, budget PromptBudget) ([]RequestFieldSummary, []string, error) {
	fields := make([]RequestFieldSummary, 0, len(message.Fields))
	var gaps []string
	for _, field := range message.Fields {
		if field == nil {
			continue
		}
		if looksLikeCredentialName(firstNonEmpty(field.JSONName, field.Name)) {
			gaps = append(gaps, "Credential-shaped protobuf fields were omitted from workflow data metadata and need separate authentication review.")
			continue
		}
		typeName, format := grpcType(field.Type, field.TypeName)
		if field.Repeated {
			typeName = "array"
		}
		fields = append(fields, RequestFieldSummary{Path: firstNonEmpty(field.JSONName, field.Name), Type: typeName, Format: format, Required: field.Required})
		if len(fields) > budget.MaxFields {
			return nil, gaps, fmt.Errorf("protobuf message exceeds the configured %d-field budget", budget.MaxFields)
		}
	}
	sort.SliceStable(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })
	return fields, uniqueSortedStrings(gaps), nil
}

func grpcType(value, typeName string) (string, string) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "string", "bytes":
		if strings.EqualFold(strings.TrimSpace(value), "bytes") {
			return "string", "byte"
		}
		return "string", ""
	case "bool", "boolean":
		return "boolean", ""
	case "double", "float":
		return "number", ""
	case "int32", "int64", "uint32", "uint64", "sint32", "sint64", "fixed32", "fixed64", "sfixed32", "sfixed64", "integer":
		return "integer", ""
	case "enum":
		return "string", ""
	case "message", "group":
		return "object", ""
	default:
		if typeName != "" {
			return "object", ""
		}
		return strings.TrimSpace(value), ""
	}
}

func adaptODataSource(ctx context.Context, source loadedOperationSource, budget PromptBudget) (operationSourceAdapterResult, error) {
	model, err := odata.Parse(source.content)
	if err != nil {
		return operationSourceAdapterResult{}, err
	}
	operations := model.OperationSummaries()
	result := operationSourceAdapterResult{title: "OData " + model.Version, operationCount: len(operations)}
	if err := enforceNativeOperationBudget(&result, budget); err != nil {
		return result, err
	}
	const inputRequirednessGap = "OData parameter nullability does not establish whether a call argument is required."
	const outputRequirednessGap = "OData nullable metadata does not prove that a property will be present in every response."
	result.capabilities = []OperationCapability{
		{Dimension: "inputs", Status: OperationCapabilityPartial, Gaps: []string{inputRequirednessGap}},
		{Dimension: "outputs", Status: OperationCapabilityPartial, Gaps: []string{outputRequirednessGap}},
		{Dimension: "auth", Status: OperationCapabilityUnsupported, Gaps: []string{"OData CSDL does not establish endpoint authentication requirements."}},
		{Dimension: "effect", Status: OperationCapabilityPartial, Gaps: []string{"OData entity-set, singleton, and navigation selectors identify resources, not a read-only HTTP operation."}},
	}
	for _, native := range operations {
		if err := ctx.Err(); err != nil {
			return result, err
		}
		if native.Selector == "" || native.ID == "" {
			continue
		}
		selector := native.Selector
		operation := nativeOperationSummary(source, native.ID, selector, "OData "+native.Kind, selector, "", "")
		for _, parameter := range native.Parameters {
			typeName, format := odataContractType(parameter.Type)
			if parameter.Collection {
				typeName = "array"
			}
			operation.Parameters = append(operation.Parameters, ParameterSummary{Name: parameter.Name, In: "parameter", Type: typeName, Format: format, Description: parameter.Type})
		}
		outputType := firstNonEmpty(native.ReturnType, native.EntityType)
		var operationGaps []string
		if outputType != "" {
			typeName, format := odataContractType(outputType)
			if native.Collection || strings.HasPrefix(outputType, "Collection(") {
				typeName = "array"
			}
			fields, gaps, fieldErr := odataEntityFields(model, outputType, budget)
			if fieldErr != nil {
				return result, fieldErr
			}
			operationGaps = append(operationGaps, gaps...)
			operation.ResponseBody = &ResponseBodySummary{Schema: &SchemaSummary{Type: typeName, Format: format}, Fields: fields}
		}
		var nativeSignals []effectSignal
		effectCapability := OperationCapabilityPartial
		switch strings.ToLower(native.Kind) {
		case "action", "action-import":
			effectCapability = OperationCapabilitySupported
			nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectWrite, kind: "odata.operation_kind", reference: selector, word: "action"})
		case "function", "function-import":
			effectCapability = OperationCapabilitySupported
			nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectRead, kind: "odata.operation_kind", reference: selector, word: "function"})
		default:
			nativeSignals = append(nativeSignals, effectSignal{effect: OperationEffectUnknown, kind: "odata.resource_kind", reference: selector, word: native.Kind})
		}
		capabilities := nativeCapabilities(OperationCapabilityPartial, OperationCapabilityPartial, OperationCapabilityUnsupported, effectCapability, selector)
		capabilities[0].Gaps = append(capabilities[0].Gaps, inputRequirednessGap)
		capabilities[1].Gaps = append(capabilities[1].Gaps, outputRequirednessGap)
		if effectCapability != OperationCapabilitySupported {
			capabilities[3].Gaps = append(capabilities[3].Gaps, "OData resource metadata does not establish whether a selected operation reads or mutates data; the selector is not a read-only HTTP operation.")
			capabilities[3].Evidence = []OperationEvidence{{Kind: "odata.resource_kind", Reference: selector}}
		}
		candidate, diagnostics, candidateErr := candidateFromOperation(source, operation, selector, assessOperationEffectWithNative(operation, selector, nativeSignals), capabilities, budget)
		for i := range candidate.Summary.Inputs {
			candidate.Summary.Inputs[i].Required = nil
		}
		for i := range candidate.Summary.Outputs {
			candidate.Summary.Outputs[i].Required = nil
		}
		candidate.Summary.Gaps = uniqueSortedStrings(append(candidate.Summary.Gaps, append(operationGaps, inputRequirednessGap, outputRequirednessGap)...))
		result.diagnostics = append(result.diagnostics, diagnostics...)
		if candidateErr != nil {
			return result, candidateErr
		}
		result.candidates = append(result.candidates, candidate)
	}
	return result, nil
}

func odataContractType(value string) (string, string) {
	value = strings.TrimSpace(value)
	if strings.HasPrefix(value, "Collection(") && strings.HasSuffix(value, ")") {
		return "array", ""
	}
	switch strings.TrimPrefix(value, "Edm.") {
	case "String", "Guid", "Binary":
		if strings.HasSuffix(value, "Binary") {
			return "string", "byte"
		}
		return "string", ""
	case "Boolean":
		return "boolean", ""
	case "Byte", "SByte", "Int16", "Int32", "Int64":
		return "integer", ""
	case "Decimal", "Single", "Double":
		return "number", ""
	case "Date", "DateTimeOffset", "TimeOfDay":
		return "string", "date-time"
	default:
		return value, ""
	}
}

func odataEntityFields(model *odata.Model, typeName string, budget PromptBudget) ([]RequestFieldSummary, []string, error) {
	typeName = strings.TrimPrefix(strings.TrimSuffix(strings.TrimPrefix(typeName, "Collection("), ")"), ".")
	var gaps []string
	for _, schema := range model.Schemas {
		if schema == nil {
			continue
		}
		for _, entity := range append(append([]*odata.EntityType(nil), schema.EntityTypes...), schema.ComplexTypes...) {
			if entity == nil || entity.FullName != typeName && entity.Name != typeName && schema.Namespace+"."+entity.Name != typeName {
				continue
			}
			var fields []RequestFieldSummary
			for _, property := range entity.Properties {
				if property == nil {
					continue
				}
				if looksLikeCredentialName(property.Name) {
					gaps = append(gaps, "Credential-shaped OData properties were omitted from workflow data metadata and need separate authentication review.")
					continue
				}
				fieldType, format := odataContractType(property.Type)
				if property.Collection {
					fieldType = "array"
				}
				fields = append(fields, RequestFieldSummary{Path: property.Name, Type: fieldType, Format: format})
				if len(fields) > budget.MaxFields {
					return nil, gaps, fmt.Errorf("OData entity type exceeds the configured %d-field budget", budget.MaxFields)
				}
			}
			sort.SliceStable(fields, func(i, j int) bool { return fields[i].Path < fields[j].Path })
			return fields, uniqueSortedStrings(gaps), nil
		}
	}
	return nil, gaps, nil
}
