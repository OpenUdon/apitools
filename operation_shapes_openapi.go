package apitools

import (
	"context"
	"strings"

	"github.com/OpenUdon/uws/binding"
)

func openAPIShapes(ctx context.Context, root map[string]any, source binding.Source, budget *operationShapeBudget) ([]binding.OperationShape, error) {
	version := stringValue(root["openapi"])
	swagger := stringValue(root["swagger"]) == "2.0"
	if !swagger && !strings.HasPrefix(version, "3.0.") && !strings.HasPrefix(version, "3.1.") {
		return nil, ErrOperationShapeTable
	}
	info, ok := root["info"].(map[string]any)
	if !ok || stringValue(info["title"]) == "" || stringValue(info["version"]) == "" {
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
			if method != strings.ToLower(method) {
				return nil, ErrOperationShapeTable
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
			var serversKnown bool
			shape.Servers, serversKnown = shapeHTTPServers(root, pathItem, native)
			var err error
			shape.Security, err = openAPIShapeSecurity(root, native)
			if err != nil {
				return nil, err
			}
			complete := serversKnown && shape.Security.Known && strings.HasPrefix(path, "/")
			dialectUnknown := false
			if _, present := root["jsonSchemaDialect"]; present {
				complete = false
				dialectUnknown = true
			}
			if _, present := native["callbacks"]; present {
				complete = false
			}
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
					switch location {
					case "path", "query", "header":
					case "cookie":
						if swagger {
							return nil, ErrOperationShapeTable
						}
					case "body", "formData":
						if !swagger {
							return nil, ErrOperationShapeTable
						}
					default:
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
					if swagger && location == "formData" {
						complete = false
					}
					if _, present := parameter["content"]; present {
						complete = false
					}
					for _, key := range []string{"style", "explode", "collectionFormat"} {
						if _, present := parameter[key]; present {
							complete = false
						}
					}
					if _, present := parameter["allowReserved"]; present {
						complete = false
					}
					required, isBool := parameter["required"].(bool)
					if _, declared := parameter["required"]; declared && !isBool {
						return nil, ErrOperationShapeTable
					}
					if location == "path" {
						if !required {
							return nil, ErrOperationShapeTable
						}
					}
					inputs[key] = binding.Input{Location: location, Name: name, Required: required, Schema: budget.schema(value, shapeLocalResolver(root))}
				}
			}
			for _, in := range inputs {
				shape.Inputs = append(shape.Inputs, in)
			}
			pathNames := map[string]bool{}
			for remaining := path; remaining != ""; {
				start := strings.IndexAny(remaining, "{}")
				if start < 0 {
					break
				}
				if remaining[start] != '{' {
					return nil, ErrOperationShapeTable
				}
				end := strings.IndexByte(remaining[start+1:], '}')
				if end < 0 {
					return nil, ErrOperationShapeTable
				}
				name := remaining[start+1 : start+1+end]
				if name == "" || strings.ContainsAny(name, "{}") {
					return nil, ErrOperationShapeTable
				}
				pathNames[name] = true
				remaining = remaining[start+end+2:]
			}
			for name := range pathNames {
				if _, declared := inputs["path\x00"+name]; !declared {
					return nil, ErrOperationShapeTable
				}
			}
			for _, in := range inputs {
				if in.Location == "path" && !pathNames[in.Name] {
					return nil, ErrOperationShapeTable
				}
			}
			if raw, exists := native["requestBody"]; exists {
				body, resolved := shapeLocalObject(root, raw)
				if !resolved {
					return nil, ErrOperationShapeTable
				}
				schema, _ := shapeBodyContent(root, body, budget)
				required, isBool := body["required"].(bool)
				if _, declared := body["required"]; declared && !isBool {
					return nil, ErrOperationShapeTable
				}
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
				schema, hasBody := shapeBodyContent(root, success[0], budget)
				if hasBody {
					shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body", Schema: schema})
				} else {
					complete = false
				}
				for _, name := range sortedMapKeys(mapValue(success[0]["headers"])) {
					header, resolved := shapeLocalObject(root, mapValue(success[0]["headers"])[name])
					if !resolved {
						return nil, ErrOperationShapeTable
					}
					value := header["schema"]
					if swagger {
						value = header
					}
					for _, key := range []string{"style", "explode", "collectionFormat", "content"} {
						if _, present := header[key]; present {
							complete = false
						}
					}
					shape.Outputs = append(shape.Outputs, binding.Output{Location: "header", Name: name, Schema: budget.schema(value, shapeLocalResolver(root))})
				}
			} else if len(success) > 1 {
				shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body"})
				complete = false
			} else {
				complete = false
			}
			if swagger {
				for _, key := range []string{"consumes", "produces"} {
					value, declared := native[key]
					if !declared {
						value = root[key]
					}
					media, ok := value.([]any)
					supported := ok && len(media) == 1 && shapeJSONMediaType(stringValue(media[0]))
					if key == "consumes" {
						for i := range shape.Inputs {
							if shape.Inputs[i].Location == "body" && !supported {
								shape.Inputs[i].Schema.Known = false
							}
						}
					} else {
						for i := range shape.Outputs {
							if shape.Outputs[i].Location == "body" && !supported {
								shape.Outputs[i].Schema.Known = false
							}
						}
					}
				}
			}
			for i := range shape.Inputs {
				if dialectUnknown {
					shape.Inputs[i].Schema.Known = false
				}
				complete = complete && shape.Inputs[i].Schema.Known
			}
			for i := range shape.Outputs {
				if dialectUnknown {
					shape.Outputs[i].Schema.Known = false
				}
				complete = complete && shape.Outputs[i].Schema.Known
			}
			shape.Complete = complete
			if err := budget.append(&out, shape); err != nil {
				return nil, err
			}
			if len(out) > binding.MaxOperations {
				return nil, ErrOperationShapeTable
			}
		}
	}
	return out, nil
}
