package apitools

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/OpenUdon/uws/binding"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Shape schemas retain validation constraints, never examples/defaults or
// descriptions. Unsupported keywords/dialects and recursive/external refs
// make the projection unknown rather than silently weakening the contract.
func (budget *operationShapeBudget) schema(value any, resolver nativeSchemaResolver) binding.Schema {
	if budget.err != nil {
		return binding.Schema{}
	}
	if err := budget.ctx.Err(); err != nil {
		budget.err = err
		return binding.Schema{}
	}
	work := 0
	projected, known := projectShapeSchema(value, resolver, map[string]bool{}, 0, &work, budget.schemaDialect)
	budget.schemaWork += work
	if work > 10000 || budget.schemaWork > 100000 {
		budget.err = ErrOperationShapeTable
		return binding.Schema{}
	}
	if projected == nil {
		return binding.Schema{}
	}
	data, err := json.Marshal(projected)
	if err != nil || len(data) > binding.MaxSchemaBytes || len(data) > binding.MaxTableBytes-budget.schemaBytes {
		budget.err = ErrOperationShapeTable
		return binding.Schema{}
	}
	budget.schemaBytes += len(data)
	if known {
		compiler := jsonschema.NewCompiler()
		compiler.UseLoader(shapeSchemaLoader{})
		compiler.AssertFormat()
		if err := compiler.AddResource("https://apitools.invalid/shape", projected); err != nil {
			known = false
		} else if _, err := compiler.Compile("https://apitools.invalid/shape"); err != nil {
			known = false
		}
	}
	return binding.Schema{Known: known, JSON: data}
}

type shapeSchemaLoader struct{}

func (shapeSchemaLoader) Load(string) (any, error) {
	return nil, errors.New("external schema resources are unavailable")
}

func projectShapeSchema(value any, resolver nativeSchemaResolver, active map[string]bool, depth int, work *int, dialect int) (any, bool) {
	*work++
	if depth > 50 || *work > 10000 {
		*work = 10001
		return nil, false
	}
	if boolean, ok := value.(bool); ok {
		return boolean, dialect != shapeSchemaOpenAPI30
	}
	schema, ok := value.(map[string]any)
	if !ok || schema == nil {
		return nil, false
	}
	out := map[string]any{}
	known := true
	if dialect != shapeSchemaOpenAPI30 && schema["$schema"] == "https://json-schema.org/draft/2020-12/schema" {
		dialect = shapeSchema2020
	}
	for _, key := range sortedMapKeys(schema) {
		child := schema[key]
		if dialect != shapeSchema2020 {
			switch key {
			case "prefixItems", "dependentSchemas", "$defs":
				known = false
			}
		}
		if dialect == shapeSchemaOpenAPI30 {
			if key == "type" {
				kind, ok := child.(string)
				if !ok || kind == "null" {
					known = false
				}
			}
			switch key {
			case "const", "contains", "if", "then", "else", "propertyNames", "patternProperties", "definitions", "$schema":
				known = false
			}
		}
		switch key {
		case "description", "title", "examples", "example", "default", "$comment", "deprecated":
			// Annotation values are outside the binding projection.
		case "readOnly", "writeOnly":
			if child != false {
				known = false
			}
		case "$ref":
			// Reference sibling meaning differs across retained OpenAPI/schema
			// dialects. Preserve a projection without claiming containment proof.
			if len(schema) > 1 {
				known = false
			}
			ref, ok := child.(string)
			if !ok || resolver == nil || active[ref] {
				known = false
				continue
			}
			target, found := resolver(ref)
			if !found {
				known = false
				continue
			}
			active[ref] = true
			resolved, complete := projectShapeSchema(target, resolver, active, depth+1, work, dialect)
			delete(active, ref)
			if resolved != nil {
				out["allOf"] = []any{resolved}
			}
			known = known && complete
		case "type", "format", "enum", "const", "required", "minimum", "maximum", "multipleOf", "minLength", "maxLength", "pattern", "minItems", "maxItems", "uniqueItems", "minProperties", "maxProperties":
			out[key] = child
		case "exclusiveMinimum", "exclusiveMaximum":
			if flag, ok := child.(bool); ok {
				known = false
				if flag {
					bound := "minimum"
					if key == "exclusiveMaximum" {
						bound = "maximum"
					}
					if schema[bound] == nil {
						known = false
					} else {
						out[key] = schema[bound]
					}
				}
			} else {
				out[key] = child
			}
		case "properties", "patternProperties", "$defs", "definitions", "dependentSchemas":
			children, ok := child.(map[string]any)
			if !ok {
				known = false
				continue
			}
			items := map[string]any{}
			for _, name := range sortedMapKeys(children) {
				projected, complete := projectShapeSchema(children[name], resolver, active, depth+1, work, dialect)
				if projected == nil {
					projected = map[string]any{}
				}
				items[name] = projected
				known = known && complete
			}
			out[key] = items
		case "items", "additionalProperties", "not", "contains", "if", "then", "else", "propertyNames":
			projected, complete := projectShapeSchema(child, resolver, active, depth+1, work, dialect)
			if projected != nil {
				out[key] = projected
			}
			known = known && complete
		case "anyOf", "oneOf", "allOf", "prefixItems":
			children, ok := child.([]any)
			if !ok {
				known = false
				continue
			}
			items := []any{}
			for _, v := range children {
				projected, complete := projectShapeSchema(v, resolver, active, depth+1, work, dialect)
				if projected != nil {
					items = append(items, projected)
				}
				known = known && complete
			}
			if key == "allOf" {
				prior, _ := out[key].([]any)
				items = append(prior, items...)
			}
			out[key] = items
		case "nullable":
			if nullable, ok := child.(bool); !ok || nullable {
				known = false
			}
		case "$schema":
			// The binding compiler's default is 2020-12. Other dialects may
			// assign different meaning to otherwise familiar keywords.
			if child != "https://json-schema.org/draft/2020-12/schema" {
				known = false
			}
		default:
			// Custom keywords, IDs, dynamic anchors and dialect-specific semantics
			// require reviewed support; no remote loader is supplied here.
			known = false
		}
	}
	return out, known
}

func shapeLocalObject(root map[string]any, value any) (map[string]any, bool) {
	object, ok := value.(map[string]any)
	if !ok {
		return nil, false
	}
	seen := map[string]bool{}
	for {
		ref, exists := object["$ref"]
		if !exists {
			return object, true
		}
		text, ok := ref.(string)
		if !ok || seen[text] || len(object) != 1 {
			return object, false
		}
		seen[text] = true
		object, ok = localJSONPointerMap(root, text)
		if !ok {
			return nil, false
		}
	}
}

func shapeLocalResolver(root map[string]any) nativeSchemaResolver {
	return func(ref string) (map[string]any, bool) { return localJSONPointerMap(root, ref) }
}

func partialTypeSchema(typeName string, nullable bool) binding.Schema {
	if typeName == "" {
		return binding.Schema{}
	}
	value := any(typeName)
	if nullable {
		value = []string{typeName, "null"}
	}
	data, _ := json.Marshal(map[string]any{"type": value})
	return binding.Schema{JSON: data}
}

func shapeBodyContent(root map[string]any, object map[string]any, budget *operationShapeBudget) (binding.Schema, bool) {
	if schema, exists := object["schema"]; exists {
		return budget.schema(schema, shapeLocalResolver(root)), true
	}
	content := mapValue(object["content"])
	if len(content) != 1 {
		return binding.Schema{}, false
	}
	for _, media := range content {
		item, ok := media.(map[string]any)
		if !ok {
			return binding.Schema{}, false
		}
		return budget.schema(item["schema"], shapeLocalResolver(root)), true
	}
	return binding.Schema{}, false
}

func shapeHTTPServers(root map[string]any, pathItem, operation map[string]any) ([]string, bool) {
	if root["swagger"] == "2.0" {
		host := stringValue(root["host"])
		value, declared := operation["schemes"]
		if !declared {
			value = root["schemes"]
		}
		schemes, _ := value.([]any)
		if host == "" || strings.ContainsAny(host, "/?#@") || len(schemes) == 0 {
			return nil, false
		}
		var servers []string
		for _, scheme := range schemes {
			value, ok := scheme.(string)
			if !ok || value != "https" && value != "http" {
				return nil, false
			}
			raw := value + "://" + host + stringValue(root["basePath"])
			clean, redacted, err := sanitizeOperationSourceURL(raw)
			if err != nil || redacted {
				return nil, false
			}
			servers = append(servers, clean)
		}
		return servers, true
	}
	value, exists := operation["servers"]
	if !exists {
		value, exists = pathItem["servers"]
	}
	if !exists {
		value, exists = root["servers"]
	}
	if !exists {
		return nil, false
	}
	items, ok := value.([]any)
	if !ok {
		return nil, false
	}
	servers := []string{}
	for _, raw := range items {
		server := mapValue(raw)
		url, ok := server["url"].(string)
		if !ok || strings.ContainsAny(url, "{}") {
			return nil, false
		}
		clean, redacted, err := sanitizeOperationSourceURL(url)
		if err != nil || redacted {
			return nil, false
		}
		servers = append(servers, clean)
	}
	return servers, len(servers) > 0
}
