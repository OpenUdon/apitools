package apitools

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/OpenUdon/apitools/internal/sourceguard"
	"github.com/OpenUdon/uws/binding"
	"gopkg.in/yaml.v3"
)

// ErrOperationShapeTable reports invalid, unsafe or over-budget source metadata.
// It deliberately contains no source values, filesystem paths or parser excerpts.
var ErrOperationShapeTable = errors.New("invalid or bounded operation-shape source contract")

// ShapeSourceInput assigns an explicit workflow source ID to local API bytes.
// URL remains provenance only; Content takes precedence over Path.
type ShapeSourceInput struct {
	ID string
	OperationSourceInput
}

// OperationShapeOptions selects local sources and optionally tighter limits.
// Defaults/ceilings: 32 sources, 20 MiB per source, 64 MiB aggregate raw bytes,
// 10,000 operations, and the UWS 8 MiB table/256 KiB schema contracts.
// Limits cannot raise these ceilings. Parsing is cooperatively cancellable;
// consumers requiring hard CPU/RSS/deadline bounds must use an isolated worker.
type OperationShapeOptions struct {
	Sources       []ShapeSourceInput
	MaxBytes      int64
	MaxOperations int
}

// BuildOperationShapeTable reproduces source-neutral UWS shapes from explicit
// local artifacts through the existing native parsers. No provenance URL,
// external reference, credential or provider is contacted. Missing/unsupported
// evidence stays unknown; malformed sources and limits return no partial table.
// Returned tables are metadata claims, never workflow or execution authority.
func BuildOperationShapeTable(ctx context.Context, options OperationShapeOptions) (binding.ShapeTable, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return binding.ShapeTable{}, err
	}
	if len(options.Sources) == 0 || len(options.Sources) > 32 || options.MaxBytes < 0 || options.MaxBytes > DefaultMaxBytes || options.MaxOperations < 0 || options.MaxOperations > binding.MaxOperations {
		return binding.ShapeTable{}, ErrOperationShapeTable
	}
	maxBytes := resolvedLocalMaxBytes(options.MaxBytes)
	maxOperations := options.MaxOperations
	if maxOperations == 0 {
		maxOperations = binding.MaxOperations
	}
	table := binding.ShapeTable{Version: binding.TableVersion}
	budget := newOperationShapeBudget(ctx, maxOperations)
	ids := map[string]bool{}
	var totalBytes int64
	for _, input := range options.Sources {
		if err := ctx.Err(); err != nil {
			return binding.ShapeTable{}, err
		}
		if !shapeText(input.ID, 256) || ids[input.ID] || !validOperationSourceKind(input.Kind) {
			return binding.ShapeTable{}, ErrOperationShapeTable
		}
		ids[input.ID] = true
		provenance, _, err := sanitizeOperationSourceURL(input.URL)
		if err != nil {
			return binding.ShapeTable{}, ErrOperationShapeTable
		}
		var content []byte
		if input.Content != nil {
			if int64(len(input.Content)) > maxBytes {
				return binding.ShapeTable{}, ErrOperationShapeTable
			}
			content = append([]byte(nil), input.Content...)
		} else if input.Path != "" {
			content, err = readLocalSpecFile(input.Path, maxBytes)
			if err != nil {
				return binding.ShapeTable{}, ErrOperationShapeTable
			}
		} else {
			return binding.ShapeTable{}, ErrOperationShapeTable
		}
		totalBytes += int64(len(content))
		if totalBytes > 64<<20 {
			return binding.ShapeTable{}, ErrOperationShapeTable
		}
		digest := sha256.Sum256(content)
		source := binding.Source{ID: input.ID, Kind: string(input.Kind), SHA256: hex.EncodeToString(digest[:]), URL: provenance}
		if err := budget.source(source); err != nil {
			return binding.ShapeTable{}, err
		}
		operations, err := sourceOperationShapes(ctx, input.Kind, content, source, budget)
		if ctxErr := ctx.Err(); ctxErr != nil {
			return binding.ShapeTable{}, ctxErr
		}
		if err != nil || len(operations) > maxOperations-len(table.Operations) {
			return binding.ShapeTable{}, ErrOperationShapeTable
		}
		table.Sources = append(table.Sources, source)
		table.Operations = append(table.Operations, operations...)
	}
	data, err := table.Marshal()
	if err != nil {
		return binding.ShapeTable{}, ErrOperationShapeTable
	}
	if err := ctx.Err(); err != nil {
		return binding.ShapeTable{}, err
	}
	return binding.ParseTable(data)
}

// VerifyOperationShapeTable independently reproduces a claimed table against
// the caller's exact local source set. Structural validity or a matching digest
// alone cannot establish producer correctness. Reproduction grants no authority.
func VerifyOperationShapeTable(ctx context.Context, options OperationShapeOptions, claimed binding.ShapeTable) error {
	if ctx == nil {
		ctx = context.Background()
	}
	data, err := boundedShapeTableBytes(ctx, claimed)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return ErrOperationShapeTable
	}
	reproduced, err := BuildOperationShapeTable(ctx, options)
	if err != nil {
		return err
	}
	expected, err := reproduced.Marshal()
	if err != nil || !bytes.Equal(data, expected) {
		return ErrOperationShapeTable
	}
	return nil
}

func shapeText(value string, max int) bool {
	if value == "" || len(value) > max || !utf8.ValidString(value) {
		return false
	}
	clean, changed := sanitizePromptString(value, max)
	return !changed && clean == value
}

func nativeShape(source binding.Source, id, ref, protocol string) binding.OperationShape {
	shape := binding.OperationShape{Source: source, Protocol: protocol}
	shape.Selector = binding.Selector{Kind: "ref", Value: ref, Key: ref}
	if id != "" {
		shape.Aliases = []binding.Selector{{Kind: "id", Value: id, Key: ref}}
	}
	return shape
}

// decodeShapeDocument shares the structural guards, but retains JSON numeric
// lexemes instead of passing schema constraints through float64 summaries.
func decodeShapeDocument(ctx context.Context, data []byte) (map[string]any, error) {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return nil, ErrOperationShapeTable
	}
	var value any
	if trimmed[0] == '{' {
		if err := sourceguard.CheckJSONWithLimits(ctx, "operation-shapes", data, sourceguard.DefaultLimits()); err != nil {
			return nil, err
		}
		d := json.NewDecoder(bytes.NewReader(data))
		d.UseNumber()
		var err error
		value, err = shapeJSONValue(ctx, d)
		if err != nil {
			return nil, err
		}
		if _, err := d.Token(); err != io.EOF {
			return nil, ErrOperationShapeTable
		}
	} else {
		node, err := sourceguard.YAMLDocument(ctx, "operation-shapes", data, sourceguard.DefaultLimits())
		if err != nil {
			return nil, err
		}
		value, err = shapeYAMLValue(node, nil)
		if err != nil {
			return nil, err
		}
	}
	root, ok := value.(map[string]any)
	if !ok || sourceguard.CheckValue("operation-shapes", root) != nil {
		return nil, ErrOperationShapeTable
	}
	return root, nil
}

func shapeJSONValue(ctx context.Context, d *json.Decoder) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	token, err := d.Token()
	if err != nil {
		return nil, err
	}
	delimiter, container := token.(json.Delim)
	if !container {
		return token, nil
	}
	switch delimiter {
	case '{':
		out := map[string]any{}
		for d.More() {
			token, err := d.Token()
			if err != nil {
				return nil, err
			}
			key, ok := token.(string)
			if !ok {
				return nil, ErrOperationShapeTable
			}
			if _, duplicate := out[key]; duplicate {
				return nil, ErrOperationShapeTable
			}
			value, err := shapeJSONValue(ctx, d)
			if err != nil {
				return nil, err
			}
			out[key] = value
		}
		if token, err := d.Token(); err != nil || token != json.Delim('}') {
			return nil, ErrOperationShapeTable
		}
		return out, nil
	case '[':
		out := []any{}
		for d.More() {
			value, err := shapeJSONValue(ctx, d)
			if err != nil {
				return nil, err
			}
			out = append(out, value)
		}
		if token, err := d.Token(); err != nil || token != json.Delim(']') {
			return nil, ErrOperationShapeTable
		}
		return out, nil
	}
	return nil, ErrOperationShapeTable
}

func shapeYAMLValue(n *yaml.Node, path []string) (any, error) {
	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) != 1 {
			return nil, ErrOperationShapeTable
		}
		return shapeYAMLValue(n.Content[0], path)
	case yaml.MappingNode:
		out := map[string]any{}
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i]
			if key.Kind != yaml.ScalarNode {
				return nil, ErrOperationShapeTable
			}
			name := key.Value
			if key.Tag != "!!str" {
				// Numeric response-code keys are the only accepted non-string
				// mapping keys. Never coerce arbitrary schema or extension keys.
				responseCodes := len(path) == 4 && path[0] == "paths" || len(path) == 5 && path[0] == "components" && path[1] == "pathItems"
				if !responseCodes || path[len(path)-1] != "responses" || !isHTTPMethod(path[len(path)-2]) && path[len(path)-2] != "trace" || key.Tag != "!!int" && key.Tag != "!!float" {
					return nil, ErrOperationShapeTable
				}
				value, err := shapeYAMLValue(key, nil)
				if err != nil {
					return nil, err
				}
				number, ok := value.(json.Number)
				code, valid := new(big.Rat).SetString(number.String())
				if !ok || !valid || !code.IsInt() || !code.Num().IsInt64() || code.Num().Int64() < 100 || code.Num().Int64() > 599 {
					return nil, ErrOperationShapeTable
				}
				name = code.Num().String()
			}
			if _, exists := out[name]; exists {
				return nil, ErrOperationShapeTable
			}
			value, err := shapeYAMLValue(n.Content[i+1], append(path, name))
			if err != nil {
				return nil, err
			}
			out[name] = value
		}
		return out, nil
	case yaml.SequenceNode:
		out := make([]any, 0, len(n.Content))
		for _, child := range n.Content {
			value, err := shapeYAMLValue(child, append(path, ""))
			if err != nil {
				return nil, err
			}
			out = append(out, value)
		}
		return out, nil
	case yaml.ScalarNode:
		switch n.Tag {
		case "!!str":
			return n.Value, nil
		case "!!null":
			return nil, nil
		case "!!bool":
			return n.Value == "true" || n.Value == "True" || n.Value == "TRUE", nil
		case "!!int", "!!float":
			value := strings.TrimPrefix(strings.ReplaceAll(n.Value, "_", ""), "+")
			if strings.HasPrefix(value, ".") {
				value = "0" + value
			}
			if strings.HasPrefix(value, "-.") {
				value = "-0" + value[1:]
			}
			if strings.HasSuffix(value, ".") {
				value += "0"
			}
			if !json.Valid([]byte(value)) {
				return nil, ErrOperationShapeTable
			}
			return json.Number(value), nil
		}
	}
	return nil, ErrOperationShapeTable
}

func sortShapeInputs(inputs []binding.Input) {
	sort.Slice(inputs, func(i, j int) bool {
		if inputs[i].Location != inputs[j].Location {
			return inputs[i].Location < inputs[j].Location
		}
		return inputs[i].Name < inputs[j].Name
	})
}
