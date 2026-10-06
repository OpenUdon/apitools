package apitools

import (
	"strings"

	"github.com/OpenUdon/uws/binding"
)

func openAPIShapeSecurity(root, operation map[string]any) (binding.Security, error) {
	raw, declared := operation["security"]
	if !declared {
		raw, declared = root["security"]
	}
	if !declared {
		return binding.Security{}, nil
	}
	alternatives, ok := raw.([]any)
	if !ok {
		return binding.Security{}, ErrOperationShapeTable
	}
	security := binding.Security{Known: true}
	if len(alternatives) == 0 {
		security.Alternatives = []binding.SecurityAlternative{{Requirements: []binding.SecurityRequirement{}}}
		return security, nil
	}
	definitions := mapValue(mapValue(root["components"])["securitySchemes"])
	if root["swagger"] == "2.0" {
		definitions = mapValue(root["securityDefinitions"])
	}
	for _, rawAlternative := range alternatives {
		requirements, ok := rawAlternative.(map[string]any)
		if !ok {
			return binding.Security{}, ErrOperationShapeTable
		}
		alternative := binding.SecurityAlternative{Requirements: []binding.SecurityRequirement{}}
		for _, name := range sortedMapKeys(requirements) {
			scopes, ok := requirements[name].([]any)
			if !ok {
				return binding.Security{}, ErrOperationShapeTable
			}
			requirement := binding.SecurityRequirement{Scheme: name, Type: "unknown"}
			seen := map[string]bool{}
			for _, scope := range scopes {
				value, ok := scope.(string)
				if !ok || seen[value] {
					return binding.Security{}, ErrOperationShapeTable
				}
				seen[value] = true
				requirement.Scopes = append(requirement.Scopes, value)
			}
			definition, resolved := shapeLocalObject(root, definitions[name])
			if !resolved {
				security.Known = false
			} else {
				requirement.Type = stringValue(definition["type"])
				switch requirement.Type {
				case "apiKey":
					location := stringValue(definition["in"])
					key := stringValue(definition["name"])
					requirement.Location, requirement.Name = location, key
					if location != "header" && location != "query" && location != "cookie" || key == "" || len(scopes) > 0 {
						security.Known = false
					}
				case "http":
					scheme := strings.ToLower(stringValue(definition["scheme"]))
					if scheme != "basic" && scheme != "bearer" || len(scopes) > 0 {
						security.Known = false
					}
				case "basic":
					if root["swagger"] != "2.0" || len(scopes) > 0 {
						security.Known = false
					}
				case "oauth2":
					declaredScopes := map[string]bool{}
					flows := mapValue(definition["flows"])
					if root["swagger"] == "2.0" {
						flows = map[string]any{"legacy": definition}
					}
					if len(flows) == 0 {
						security.Known = false
					}
					for _, rawFlow := range flows {
						flow, ok := rawFlow.(map[string]any)
						if !ok {
							security.Known = false
							continue
						}
						for scope := range mapValue(flow["scopes"]) {
							declaredScopes[scope] = true
						}
					}
					for _, scope := range requirement.Scopes {
						if !declaredScopes[scope] {
							security.Known = false
						}
					}
				default:
					security.Known = false
					if requirement.Type == "" {
						requirement.Type = "unknown"
					}
				}
			}
			alternative.Requirements = append(alternative.Requirements, requirement)
		}
		security.Alternatives = append(security.Alternatives, alternative)
	}
	return security, nil
}

// Scope-only Discovery metadata is partial authentication evidence. Preserve
// the documented alternatives without claiming credential/account readiness.
func discoveryShapeSecurity(scopes []string) binding.Security {
	security := binding.Security{}
	for _, scope := range scopes {
		security.Alternatives = append(security.Alternatives, binding.SecurityAlternative{Requirements: []binding.SecurityRequirement{{Scheme: "google.oauth2", Type: "oauth2", Scopes: []string{scope}}}})
	}
	return security
}
