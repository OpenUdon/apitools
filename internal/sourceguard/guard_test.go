package sourceguard

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestExplicitIndexLimitsRemainBounded(t *testing.T) {
	limits := Limits{MaxDocumentBytes: 128 << 20, MaxStructuralItems: 4, MaxNestingDepth: 100}
	if err := CheckJSONWithLimits(context.Background(), "index", []byte("[0,1,2,3]"), limits); err == nil {
		t.Fatal("accepted structural overflow")
	}
	if _, err := YAMLDocument(context.Background(), "index", []byte("a: b\nc: d\n"), limits); err == nil {
		t.Fatal("accepted YAML structural overflow")
	}
	limits.MaxStructuralItems = 0
	if err := CheckJSONWithLimits(context.Background(), "index", []byte("{}"), limits); err == nil {
		t.Fatal("accepted unbounded limits")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := CheckJSONWithLimits(ctx, "index", []byte("{}"), DefaultLimits()); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
}

func TestCheckJSONUsesStructuralBudgetSeparateFromSemanticWork(t *testing.T) {
	data := []byte("[" + strings.Repeat("0,", MaxWorkItems) + "0]")
	if err := CheckJSON("openapi", data); err != nil {
		t.Fatalf("linear JSON above the semantic work budget was rejected: %v", err)
	}
}

func TestCheckJSONRejectsStructuralBudgetOverflow(t *testing.T) {
	data := []byte("[" + strings.Repeat("0,", MaxStructuralItems) + "0]")
	err := CheckJSON("openapi", data)
	if err == nil || !strings.Contains(err.Error(), "JSON token count exceeds maximum") {
		t.Fatalf("structural overflow error = %v", err)
	}
}

func TestCheckValueRejectsDepthWorkAndRetainedBytes(t *testing.T) {
	deep := map[string]any{}
	cursor := deep
	for range MaxNestingDepth + 1 {
		next := map[string]any{}
		cursor["next"] = next
		cursor = next
	}
	if err := CheckValue("decoded", deep); err == nil || !strings.Contains(err.Error(), "nesting exceeds") {
		t.Fatalf("depth error = %v", err)
	}

	wide := make([]any, MaxWorkItems)
	if err := CheckValue("decoded", wide); err == nil || !strings.Contains(err.Error(), "work exceeds") {
		t.Fatalf("work error = %v", err)
	}

	large := map[string]any{"value": strings.Repeat("x", MaxDocumentBytes+1)}
	if err := CheckValue("decoded", large); err == nil || !strings.Contains(err.Error(), "string data exceeds") {
		t.Fatalf("retained-byte error = %v", err)
	}
}
