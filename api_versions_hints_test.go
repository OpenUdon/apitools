package apitools

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type versionAdapterFunc func(context.Context, APIVersionHintContext) (APIVersionHints, error)

func (f versionAdapterFunc) DiscoverHints(ctx context.Context, input APIVersionHintContext) (APIVersionHints, error) {
	return f(ctx, input)
}

func TestAPIVersionHintsBudgetAndUntrustedClaims(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		v := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/v"), ".json")
		_, _ = w.Write([]byte(versionSpec(v)))
	}))
	defer server.Close()
	req := versionTestRequest(server.URL)
	hints := APIVersionHints{SchemaVersion: APIVersionHintsSchemaVersion}
	for i := 2; i < 10; i++ {
		hints.Hints = append(hints.Hints, APIVersionHint{URL: fmt.Sprintf("%s/v%d.json", server.URL, i), Version: fmt.Sprint(i), Producer: "fake-adapter"})
	}
	report, err := (&Client{AllowUnsafeHosts: true}).VerifyHints(context.Background(), req, hints, APIVersionDiscoveryOptions{Network: true})
	if err != nil || calls.Load() != 3 || report.Status != "newer_found" || len(report.Versions) != 8 {
		t.Fatalf("report=%#v calls=%d err=%v", report, calls.Load(), err)
	}
	unexamined := 0
	for _, r := range report.Versions {
		if r.Validation == "unexamined" {
			unexamined++
		}
	}
	if unexamined != 5 {
		t.Fatalf("omitted excess hints: %#v", report)
	}
	hints.Hints = []APIVersionHint{{URL: server.URL + "/v2.json", Version: "v99"}, {URL: server.URL + "/v2.json", Version: "v99"}}
	before := calls.Load()
	report, err = (&Client{AllowUnsafeHosts: true}).VerifyHints(context.Background(), req, hints, APIVersionDiscoveryOptions{Network: true})
	if err != nil || calls.Load() != before+1 || report.Preferred != "v2" || report.Status != "newer_found" {
		t.Fatalf("hint claim upgraded source: %#v %v", report, err)
	}
}

func TestAPIVersionAdapterFailureAndLateHints(t *testing.T) {
	var calls atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { calls.Add(1); http.NotFound(w, r) }))
	defer server.Close()
	req := versionTestRequest(server.URL)
	c := &Client{AllowUnsafeHosts: true}
	opts := APIVersionDiscoveryOptions{Network: true, Timeout: 20 * time.Millisecond, AdapterEnabled: true, AdapterTimeout: 60 * time.Millisecond}
	adapter := versionAdapterFunc(func(ctx context.Context, input APIVersionHintContext) (APIVersionHints, error) {
		input.Report.Known.Version = "mutated"
		<-ctx.Done()
		return APIVersionHints{}, ctx.Err()
	})
	start := time.Now()
	report, err := c.DiscoverAPIVersionsWithAdapter(context.Background(), req, opts, adapter)
	if err != nil || report.Known.Version != req.Known.Version || time.Since(start) > 310*time.Millisecond {
		t.Fatalf("adapter failure %#v %v", report, err)
	}
	adapter = versionAdapterFunc(func(context.Context, APIVersionHintContext) (APIVersionHints, error) {
		time.Sleep(35 * time.Millisecond)
		return APIVersionHints{SchemaVersion: APIVersionHintsSchemaVersion, Hints: []APIVersionHint{{URL: server.URL + "/v2.json"}}}, nil
	})
	before := calls.Load()
	report, err = c.DiscoverAPIVersionsWithAdapter(context.Background(), req, opts, adapter)
	if err != nil || calls.Load() != before+1 {
		t.Fatalf("late hints reset network budget: %#v %v count=%d", report, err, calls.Load())
	}
}
