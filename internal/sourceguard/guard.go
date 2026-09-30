// Package sourceguard defines the common resource contract for parsers that
// consume untrusted API-description artifacts.
package sourceguard

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

const (
	// MaxDocumentBytes is the maximum size accepted by direct parser entry
	// points. Discovery callers may impose a smaller limit before parsing.
	MaxDocumentBytes = 20 << 20
	// MaxNestingDepth bounds structural and recursive parser nesting.
	MaxNestingDepth = 100
	// MaxWorkItems bounds tokens, wire fields, or decoded nodes handled by one
	// parse operation.
	MaxWorkItems = 100_000
	// MaxStructuralItems bounds the streaming JSON/YAML preflight. Structural
	// scans are linear in the already-bounded source bytes and need a higher
	// ceiling than semantic parser work so valid large OpenAPI documents (for
	// example Kubernetes Swagger) remain inspectable.
	MaxStructuralItems = 1_000_000
)

// Limits are internal, explicit budgets for the registered catalog index path.
// Existing parsers always select DefaultLimits; these are not caller request
// options and cannot raise nesting or remove structural limits.
type Limits struct {
	MaxDocumentBytes   int
	MaxStructuralItems int
	MaxNestingDepth    int
}

func DefaultLimits() Limits {
	return Limits{MaxDocumentBytes: MaxDocumentBytes, MaxStructuralItems: MaxStructuralItems, MaxNestingDepth: MaxNestingDepth}
}

func checkLimits(ctx context.Context, kind string, data []byte, limits Limits) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if limits.MaxDocumentBytes <= 0 || limits.MaxDocumentBytes > 128<<20 || limits.MaxStructuralItems <= 0 || limits.MaxStructuralItems > 8_000_000 || limits.MaxNestingDepth <= 0 || limits.MaxNestingDepth > MaxNestingDepth {
		return fmt.Errorf("%s: invalid explicit parser limits", kind)
	}
	if len(data) > limits.MaxDocumentBytes {
		return fmt.Errorf("%s: document exceeds maximum size %d bytes", kind, limits.MaxDocumentBytes)
	}
	return nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}

// CheckDocument enforces the shared direct-parser byte limit.
func CheckDocument(kind string, data []byte) error {
	if len(data) > MaxDocumentBytes {
		return fmt.Errorf("%s: document exceeds maximum size %d bytes", kind, MaxDocumentBytes)
	}
	return nil
}

// CheckYAML validates YAML structure without expanding aliases into recursive
// Go values. Alias nodes are rejected because expansion can multiply work
// beyond the source byte and node budgets.
func CheckYAML(kind string, data []byte) error {
	_, err := YAMLDocument(context.Background(), kind, data, DefaultLimits())
	return err
}

// YAMLDocument validates one bounded alias-free document, retaining its node
// tree so an index caller can decode it without parsing the source twice.
func YAMLDocument(ctx context.Context, kind string, data []byte, limits Limits) (*yaml.Node, error) {
	if err := checkLimits(ctx, kind, data, limits); err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(contextReader{ctx: ctx, reader: bytes.NewReader(data)})
	var document yaml.Node
	if err := decoder.Decode(&document); err != nil {
		return nil, fmt.Errorf("%s: decode YAML structure: %w", kind, err)
	}
	var trailing yaml.Node
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("%s: YAML contains multiple documents", kind)
		}
		return nil, fmt.Errorf("%s: decode YAML trailing content: %w", kind, err)
	}
	type item struct {
		node  *yaml.Node
		depth int
	}
	stack := []item{{node: &document, depth: 0}}
	work := 0
	for len(stack) > 0 {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		if current.node == nil {
			continue
		}
		work++
		if work > limits.MaxStructuralItems {
			return nil, fmt.Errorf("%s: YAML node count exceeds maximum %d", kind, limits.MaxStructuralItems)
		}
		if current.depth > limits.MaxNestingDepth {
			return nil, fmt.Errorf("%s: YAML nesting exceeds maximum depth %d", kind, limits.MaxNestingDepth)
		}
		if current.node.Kind == yaml.AliasNode {
			return nil, fmt.Errorf("%s: YAML aliases are not supported in untrusted source metadata", kind)
		}
		for i := len(current.node.Content) - 1; i >= 0; i-- {
			stack = append(stack, item{node: current.node.Content[i], depth: current.depth + 1})
		}
	}
	return &document, ctx.Err()
}

// CheckJSON performs a streaming structural pass before callers decode JSON
// into recursive maps. It limits nesting and total tokens without retaining
// attacker-controlled values.
func CheckJSON(kind string, data []byte) error {
	return CheckJSONWithLimits(context.Background(), kind, data, DefaultLimits())
}

// CheckJSONWithLimits is the explicit internal index preflight. Direct parsers
// continue using CheckJSON and its original limits.
func CheckJSONWithLimits(ctx context.Context, kind string, data []byte, limits Limits) error {
	if err := checkLimits(ctx, kind, data, limits); err != nil {
		return err
	}
	decoder := json.NewDecoder(contextReader{ctx: ctx, reader: bytes.NewReader(data)})
	decoder.UseNumber()
	depth := 0
	work := 0
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		token, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("%s: decode JSON structure: %w", kind, err)
		}
		work++
		if work > limits.MaxStructuralItems {
			return fmt.Errorf("%s: JSON token count exceeds maximum %d", kind, limits.MaxStructuralItems)
		}
		if delimiter, ok := token.(json.Delim); ok {
			switch delimiter {
			case '{', '[':
				depth++
				if depth > limits.MaxNestingDepth {
					return fmt.Errorf("%s: JSON nesting exceeds maximum depth %d", kind, limits.MaxNestingDepth)
				}
			case '}', ']':
				depth--
			}
		}
	}
	return nil
}

// CheckValue bounds an already-decoded JSON-shaped value before a parser
// traverses it. ParseMap-style compatibility APIs cannot apply a source-byte
// limit, so this guard accounts for container/scalar work, nesting, and the
// aggregate bytes retained in string keys and values. Cyclic maps or slices
// terminate at the nesting limit instead of recursing indefinitely.
func CheckValue(kind string, value any) error {
	type item struct {
		value any
		depth int
	}
	stack := []item{{value: value}}
	work := 0
	retainedBytes := 0
	addBytes := func(n int) error {
		if n > MaxDocumentBytes-retainedBytes {
			return fmt.Errorf("%s: decoded string data exceeds maximum size %d bytes", kind, MaxDocumentBytes)
		}
		retainedBytes += n
		return nil
	}
	for len(stack) > 0 {
		last := len(stack) - 1
		current := stack[last]
		stack = stack[:last]
		work++
		if work > MaxWorkItems {
			return fmt.Errorf("%s: decoded value work exceeds maximum %d items", kind, MaxWorkItems)
		}
		if current.depth > MaxNestingDepth {
			return fmt.Errorf("%s: decoded value nesting exceeds maximum depth %d", kind, MaxNestingDepth)
		}
		switch typed := current.value.(type) {
		case map[string]any:
			if len(typed) > MaxWorkItems-work-len(stack) {
				return fmt.Errorf("%s: decoded value work exceeds maximum %d items", kind, MaxWorkItems)
			}
			for key, child := range typed {
				if err := addBytes(len(key)); err != nil {
					return err
				}
				stack = append(stack, item{value: child, depth: current.depth + 1})
			}
		case []any:
			if len(typed) > MaxWorkItems-work-len(stack) {
				return fmt.Errorf("%s: decoded value work exceeds maximum %d items", kind, MaxWorkItems)
			}
			for _, child := range typed {
				stack = append(stack, item{value: child, depth: current.depth + 1})
			}
		case string:
			if err := addBytes(len(typed)); err != nil {
				return err
			}
		case []byte:
			if err := addBytes(len(typed)); err != nil {
				return err
			}
		}
	}
	return nil
}

// Budget tracks work shared by recursive or nested binary decoders.
type Budget struct {
	kind  string
	used  int
	limit int
}

// NewBudget returns a work budget using the shared default limit.
func NewBudget(kind string) *Budget {
	return &Budget{kind: kind, limit: MaxWorkItems}
}

// Use accounts for n work items.
func (b *Budget) Use(n int) error {
	if b == nil || n <= 0 {
		return nil
	}
	b.used += n
	if b.used > b.limit {
		return fmt.Errorf("%s: parser work exceeds maximum %d items", b.kind, b.limit)
	}
	return nil
}
