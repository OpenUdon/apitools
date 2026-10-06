package apitools

import (
	"context"
	"strings"

	"github.com/OpenUdon/uws/binding"
)

func openAPIShapes(ctx context.Context, root map[string]any, source binding.Source) ([]binding.OperationShape, error) {
	version := stringValue(root["openapi"])
	swagger := stringValue(root["swagger"]) == "2.0"
	if !swagger && !strings.HasPrefix(version, "3.0.") && !strings.HasPrefix(version, "3.1.") {
		return nil, ErrOperationShapeTable
	}
	paths, ok := root["paths"].(map[string]any)
	if !ok {
		return nil, ErrOperationShapeTable
	}
	var out []binding.OperationShape
	for _, path := range sortedMapKeys(paths) {
		pathItem, resolved := shapeLocalObject(root, paths[path])
		if !resolved {
			return nil, ErrOperationShapeTable
		}
		for _, method := range sortedMapKeys(pathItem) {
			if !isHTTPMethod(method) {
				continue
			}
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			native, ok := pathItem[method].(map[string]any)
			if !ok {
				return nil, ErrOperationShapeTable
			}
			ref := openAPIOperationSelector(method, path)
			shape := nativeShape(source, stringValue(native["operationId"]), ref, "http")
			shape.Method, shape.Path = strings.ToUpper(method), path
			shape.Servers, _ = shapeHTTPServers(root, pathItem, native)
			inputs := map[string]binding.Input{}
			for _, owner := range []map[string]any{pathItem, native} {
				value, declared := owner["parameters"]
				if !declared {
					continue
				}
				parameters, ok := value.([]any)
				if !ok {
					return nil, ErrOperationShapeTable
				}
				seen := map[string]bool{}
				for _, raw := range parameters {
					parameter, resolved := shapeLocalObject(root, raw)
					if !resolved {
						return nil, ErrOperationShapeTable
					}
					name, location := stringValue(parameter["name"]), stringValue(parameter["in"])
					key := location + "\x00" + name
					if name == "" || location == "" || seen[key] {
						return nil, ErrOperationShapeTable
					}
					seen[key] = true
					value := parameter["schema"]
					if swagger && location != "body" {
						projection := map[string]any{}
						for _, k := range []string{"type", "format", "items", "enum", "minimum", "maximum", "exclusiveMinimum", "exclusiveMaximum", "multipleOf", "minLength", "maxLength", "pattern", "minItems", "maxItems", "uniqueItems"} {
							if v, exists := parameter[k]; exists {
								projection[k] = v
							}
						}
						value = projection
					}
					required, _ := parameter["required"].(bool)
					if location == "path" {
						required = true
					}
					inputs[key] = binding.Input{Location: location, Name: name, Required: required, Schema: projectedShapeSchema(value, shapeLocalResolver(root))}
				}
			}
			for _, in := range inputs {
				shape.Inputs = append(shape.Inputs, in)
			}
			if raw, exists := native["requestBody"]; exists {
				body, resolved := shapeLocalObject(root, raw)
				if !resolved {
					return nil, ErrOperationShapeTable
				}
				schema, _ := shapeBodyContent(root, body)
				required, _ := body["required"].(bool)
				shape.Inputs = append(shape.Inputs, binding.Input{Location: "body", Name: "body", Required: required, Schema: schema})
			}
			sortShapeInputs(shape.Inputs)
			responses, ok := native["responses"].(map[string]any)
			if !ok {
				return nil, ErrOperationShapeTable
			}
			var success []map[string]any
			for _, status := range sortedMapKeys(responses) {
				if len(status) != 3 || status[0] != '2' {
					continue
				}
				response, resolved := shapeLocalObject(root, responses[status])
				if !resolved {
					return nil, ErrOperationShapeTable
				}
				success = append(success, response)
			}
			if len(success) == 1 {
				schema, hasBody := shapeBodyContent(root, success[0])
				if hasBody {
					shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body", Schema: schema})
				}
			} else if len(success) > 1 {
				shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body"})
			}
			// Full transport/security completeness is qualified separately from
			// source-native operation and schema projection.
			out = append(out, shape)
			if len(out) > binding.MaxOperations {
				return nil, ErrOperationShapeTable
			}
		}
	}
	return out, nil
}
