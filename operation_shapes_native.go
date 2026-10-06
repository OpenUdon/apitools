package apitools

import (
	"bytes"
	"context"
	"strings"

	"github.com/OpenUdon/apitools/graphql"
	"github.com/OpenUdon/apitools/grpcproto"
	"github.com/OpenUdon/apitools/internal/sourceguard"
	"github.com/OpenUdon/apitools/odata"
	"github.com/OpenUdon/apitools/openrpc"
	asyncparser "github.com/OpenUdon/asyncapi"
	smithyparser "github.com/OpenUdon/awssmithy"
	googleparser "github.com/OpenUdon/googlediscovery"
	"github.com/OpenUdon/uws/binding"
)

func sourceOperationShapes(ctx context.Context, kind OperationSourceKind, data []byte, source binding.Source) ([]binding.OperationShape, error) {
	if sourceguard.CheckDocument(string(kind), data) != nil {
		return nil, ErrOperationShapeTable
	}
	switch kind {
	case OperationSourceOpenAPI, OperationSourceGoogleDiscovery, OperationSourceAWSSmithy, OperationSourceAsyncAPI:
		root, err := decodeShapeDocument(ctx, data)
		if err != nil {
			return nil, err
		}
		switch kind {
		case OperationSourceOpenAPI:
			return openAPIShapes(ctx, root, source)
		case OperationSourceGoogleDiscovery:
			return discoveryShapes(ctx, root, source)
		case OperationSourceAWSSmithy:
			return smithyShapes(ctx, root, source)
		case OperationSourceAsyncAPI:
			return asyncShapes(ctx, root, source)
		}
	case OperationSourceOpenRPC:
		if _, err := decodeShapeDocument(ctx, data); err != nil {
			return nil, err
		}
		model, err := openrpc.Parse(data)
		if err != nil {
			return nil, err
		}
		return openRPCShapes(ctx, model, source)
	case OperationSourceGraphQL:
		if bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			if _, err := decodeShapeDocument(ctx, data); err != nil {
				return nil, err
			}
		}
		model, err := graphql.Parse(data)
		if err != nil {
			return nil, err
		}
		return graphQLShapes(ctx, model, source)
	case OperationSourceGRPCProtobuf:
		if bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			if _, err := decodeShapeDocument(ctx, data); err != nil {
				return nil, err
			}
		}
		model, err := grpcproto.Parse(data)
		if err != nil {
			return nil, err
		}
		return grpcShapes(ctx, model, source)
	case OperationSourceOData:
		if bytes.HasPrefix(bytes.TrimSpace(data), []byte("{")) {
			if _, err := decodeShapeDocument(ctx, data); err != nil {
				return nil, err
			}
		}
		model, err := odata.Parse(data)
		if err != nil {
			return nil, err
		}
		return odataShapes(ctx, model, source)
	}
	return nil, ErrOperationShapeTable
}

func discoveryShapes(ctx context.Context, root map[string]any, source binding.Source) ([]binding.OperationShape, error) {
	if root["discoveryVersion"] != "v1" {
		return nil, ErrOperationShapeTable
	}
	model, err := googleparser.ParseMap(root)
	if err != nil {
		return nil, err
	}
	var out []binding.OperationShape
	resolver := func(ref string) (map[string]any, bool) {
		ref = strings.TrimPrefix(ref, "#/components/schemas/")
		target, ok := model.Schemas[ref]
		return target, ok
	}
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if native == nil {
			return nil, ErrOperationShapeTable
		}
		shape := nativeShape(source, native.ID, "#/methods/"+escapeJSONPointer(native.ID), "http")
		shape.Method, shape.Path = native.HTTPMethod, native.Path
		shape.Security = discoveryShapeSecurity(native.Scopes)
		// Only declared provenance is emitted; the parser's default Google
		// endpoint is not new source evidence.
		if root["rootUrl"] != nil || root["baseUrl"] != nil {
			if server, redacted, err := sanitizeOperationSourceURL(model.ServerURL); err == nil && !redacted {
				shape.Servers = []string{server}
			}
		}
		for _, parameter := range native.Parameters {
			if parameter == nil {
				return nil, ErrOperationShapeTable
			}
			shape.Inputs = append(shape.Inputs, binding.Input{Location: parameter.Location, Name: firstNonEmpty(parameter.OriginalName, parameter.Name), Required: parameter.Required, Schema: projectedShapeSchema(parameter.Schema, resolver)})
		}
		if native.RequestRef != "" {
			shape.Inputs = append(shape.Inputs, binding.Input{Location: "body", Name: "body", Required: true, Schema: projectedShapeSchema(model.Schemas[native.RequestRef], resolver)})
		}
		if native.ResponseRef != "" {
			shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body", Schema: projectedShapeSchema(model.Schemas[native.ResponseRef], resolver)})
		}
		// The Discovery parser normalizes a dialect-specific schema subset.
		// Do not label that lossy projection a complete JSON Schema contract.
		for i := range shape.Inputs {
			shape.Inputs[i].Schema.Known = false
		}
		for i := range shape.Outputs {
			shape.Outputs[i].Schema.Known = false
		}
		sortShapeInputs(shape.Inputs)
		out = append(out, shape)
	}
	return out, nil
}

func smithyShapes(ctx context.Context, root map[string]any, source binding.Source) ([]binding.OperationShape, error) {
	model, err := smithyparser.ParseMap(root)
	if err != nil {
		return nil, err
	}
	var out []binding.OperationShape
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if native == nil {
			return nil, ErrOperationShapeTable
		}
		shape := nativeShape(source, native.ID, "#/shapes/"+escapeJSONPointer(native.ID), "aws-smithy:"+model.Protocol)
		// Preserve native protocol/member locations; a modeled AWS operation
		// is not an invented generic HTTP operation/server.
		for _, member := range append(append([]*smithyparser.MemberBinding{}, native.InputBindings...), native.UnboundInput...) {
			if member == nil {
				return nil, ErrOperationShapeTable
			}
			location := member.Location
			if location == "" {
				location = "body"
			}
			shape.Inputs = append(shape.Inputs, binding.Input{Location: location, Name: member.MemberName, Required: member.Required, Schema: smithyMemberSchema(model, member.Target)})
		}
		if native.Payload != nil {
			shape.Inputs = append(shape.Inputs, binding.Input{Location: "payload", Name: native.Payload.MemberName, Required: native.Payload.Required, Schema: smithyMemberSchema(model, native.Payload.Target)})
		}
		if native.Output != "" {
			shape.Outputs = append(shape.Outputs, binding.Output{Location: "body", Name: "body", Schema: smithyMemberSchema(model, native.Output)})
		}
		sortShapeInputs(shape.Inputs)
		out = append(out, shape)
	}
	return out, nil
}

func smithyMemberSchema(model *smithyparser.Model, target string) binding.Schema {
	shape := model.Shapes[target]
	if shape == nil {
		kind := strings.TrimPrefix(target, "smithy.api#")
		typeName, _ := smithyType(&smithyparser.Shape{Type: strings.ToLower(kind)})
		return partialTypeSchema(typeName, false)
	}
	typeName, _ := smithyType(shape)
	return partialTypeSchema(typeName, false)
}

func asyncShapes(ctx context.Context, root map[string]any, source binding.Source) ([]binding.OperationShape, error) {
	model, err := asyncparser.ParseMap(root)
	if err != nil {
		return nil, err
	}
	var out []binding.OperationShape
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if native == nil {
			return nil, ErrOperationShapeTable
		}
		shape := nativeShape(source, native.ID, "#/operations/"+escapeJSONPointer(native.ID), "asyncapi")
		var schema binding.Schema
		if len(native.MessageRefs) == 1 {
			target, ok := model.ResolveSelector(native.MessageRefs[0])
			if ok && target.Message != nil {
				schema = projectedShapeSchema(target.Message.Payload, shapeLocalResolver(root))
			}
		}
		switch native.Action {
		case "send", "publish":
			shape.Inputs = []binding.Input{{Location: "payload", Name: "payload", Required: true, Schema: schema}}
		case "receive", "subscribe":
			shape.Outputs = []binding.Output{{Location: "payload", Name: "payload", Schema: schema}}
		}
		out = append(out, shape)
	}
	return out, nil
}

func openRPCShapes(ctx context.Context, model *openrpc.Model, source binding.Source) ([]binding.OperationShape, error) {
	var out []binding.OperationShape
	for _, native := range model.Methods {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if native == nil {
			return nil, ErrOperationShapeTable
		}
		shape := nativeShape(source, native.Name, openrpc.MethodSelector(native.Name), "json-rpc")
		for _, parameter := range native.Params {
			if parameter == nil {
				return nil, ErrOperationShapeTable
			}
			descriptor := openrpcDescriptor(model, parameter)
			name := firstNonEmpty(descriptor.Name, parameter.Name, localReferenceName(parameter.Ref))
			shape.Inputs = append(shape.Inputs, binding.Input{Location: "param", Name: name, Required: descriptor.Required, Schema: projectedShapeSchema(descriptor.Schema, shapeLocalResolver(model.Raw))})
		}
		if native.Result != nil {
			descriptor := openrpcDescriptor(model, native.Result)
			shape.Outputs = []binding.Output{{Location: "result", Name: firstNonEmpty(descriptor.Name, "result"), Schema: projectedShapeSchema(descriptor.Schema, shapeLocalResolver(model.Raw))}}
		}
		// Positional JSON-RPC parameters retain declaration order.
		out = append(out, shape)
	}
	return out, nil
}

func graphQLShapes(ctx context.Context, model *graphql.Model, source binding.Source) ([]binding.OperationShape, error) {
	var out []binding.OperationShape
	for _, native := range model.Operations {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if native == nil {
			return nil, ErrOperationShapeTable
		}
		shape := nativeShape(source, native.ID, native.Selector, "graphql")
		addShapeAliases(&shape, []string{native.Name, "operation:" + native.ID, native.SourceRef})
		for _, variable := range native.Variables {
			if variable == nil {
				return nil, ErrOperationShapeTable
			}
			shape.Inputs = append(shape.Inputs, binding.Input{Location: "variable", Name: variable.Name, Required: variable.Required, Schema: graphQLTypeSchema(variable.Type)})
		}
		if native.FieldName != "" {
			var root *graphql.Type
			for _, candidate := range model.Types {
				if candidate != nil && candidate.Name == native.RootType {
					root = candidate
					break
				}
			}
			if field := graphqlRootField(root, native.FieldName); field != nil {
				for _, argument := range field.Args {
					if argument == nil {
						return nil, ErrOperationShapeTable
					}
					shape.Inputs = append(shape.Inputs, binding.Input{Location: "argument", Name: argument.Name, Required: argument.Required, Schema: graphQLTypeSchema(argument.Type)})
				}
				shape.Outputs = []binding.Output{{Location: "data", Name: native.FieldName, Schema: graphQLTypeSchema(field.Type)}}
			}
		} else {
			for _, name := range native.SelectionNames {
				shape.Outputs = append(shape.Outputs, binding.Output{Location: "data", Name: name})
			}
		}
		sortShapeInputs(shape.Inputs)
		out = append(out, shape)
	}
	return out, nil
}

func graphQLTypeSchema(ref graphql.TypeRef) binding.Schema {
	typeName, _ := graphqlContractType(ref)
	// Lists, custom scalars, selection sets and server auth need additional
	// protocol evidence; the shallow type is explicitly a partial projection.
	return partialTypeSchema(typeName, !ref.Required)
}

func grpcShapes(ctx context.Context, model *grpcproto.Model, source binding.Source) ([]binding.OperationShape, error) {
	var out []binding.OperationShape
	for _, native := range model.MethodSummaries() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		shape := nativeShape(source, native.SourceOperationID, native.Selector, "grpc-protobuf")
		addShapeAliases(&shape, []string{native.FullMethod, "rpc:" + native.SourceOperationID, native.SourceOperationRef})
		inputSchema := grpcPartialMessageSchema(findProtoMessage(model, native.RequestType, native.Package))
		outputSchema := grpcPartialMessageSchema(findProtoMessage(model, native.ResponseType, native.Package))
		if native.ClientStreaming {
			inputSchema = binding.Schema{}
		}
		if native.ServerStreaming {
			outputSchema = binding.Schema{}
		}
		shape.Inputs = []binding.Input{{Location: "message", Name: native.RequestType, Required: true, Schema: inputSchema}}
		shape.Outputs = []binding.Output{{Location: "message", Name: native.ResponseType, Schema: outputSchema}}
		out = append(out, shape)
	}
	return out, nil
}

func odataShapes(ctx context.Context, model *odata.Model, source binding.Source) ([]binding.OperationShape, error) {
	var out []binding.OperationShape
	for _, native := range model.OperationSummaries() {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		shape := nativeShape(source, native.SourceOperationID, native.Selector, "odata")
		addShapeAliases(&shape, []string{native.Name, "odata:" + native.ID, native.SourceOperationRef})
		for _, parameter := range native.Parameters {
			typeName, _ := odataContractType(parameter.Type)
			if parameter.Collection {
				typeName = "array"
			}
			// CSDL nullability does not prove call-argument requiredness. The
			// operation stays incomplete, so false cannot become absence proof.
			shape.Inputs = append(shape.Inputs, binding.Input{Location: "param", Name: parameter.Name, Schema: partialTypeSchema(typeName, parameter.Nullable)})
		}
		if native.ReturnType != "" || native.EntityType != "" {
			typeName, _ := odataContractType(firstNonEmpty(native.ReturnType, native.EntityType))
			if native.Collection {
				typeName = "array"
			}
			if typeName == "" {
				typeName = "object"
			}
			shape.Outputs = []binding.Output{{Location: "body", Name: "body", Schema: partialTypeSchema(typeName, true)}}
		}
		sortShapeInputs(shape.Inputs)
		out = append(out, shape)
	}
	return out, nil
}

func addShapeAliases(shape *binding.OperationShape, values []string) {
	seen := map[string]bool{shape.Selector.Kind + "\x00" + shape.Selector.Value: true}
	for _, alias := range shape.Aliases {
		seen[alias.Kind+"\x00"+alias.Value] = true
	}
	for _, value := range values {
		if value == "" {
			continue
		}
		kind := "id"
		if strings.HasPrefix(value, "#/") {
			kind = "ref"
		}
		if !seen[kind+"\x00"+value] {
			shape.Aliases = append(shape.Aliases, binding.Selector{Kind: kind, Value: value, Key: shape.Selector.Key})
			seen[kind+"\x00"+value] = true
		}
	}
}

func grpcPartialMessageSchema(message *grpcproto.Message) binding.Schema {
	if message == nil {
		return binding.Schema{}
	}
	properties := map[string]any{}
	var required []string
	for _, field := range message.Fields {
		if field == nil {
			return binding.Schema{}
		}
		name := firstNonEmpty(field.JSONName, field.Name)
		typeName, _ := grpcType(field.Type, field.TypeName)
		child := map[string]any{}
		if typeName != "" {
			child["type"] = typeName
		}
		if field.Repeated {
			child = map[string]any{"type": "array", "items": child}
		}
		properties[name] = child
		if field.Required {
			required = append(required, name)
		}
	}
	value := map[string]any{"type": "object", "properties": properties}
	if len(required) > 0 {
		value["required"] = required
	}
	schema := projectedShapeSchema(value, nil)
	// Proto JSON wire presence, integer/string encodings, oneof, custom
	// options and nested/enum types are not all represented by the parser.
	schema.Known = false
	return schema
}
