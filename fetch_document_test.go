package apitools

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchDocumentRejectsUnsafeHosts(t *testing.T) {
	c := &Client{}
	for _, host := range []string{"localhost", "127.0.0.1", "[::1]", "10.0.0.1", "192.168.1.1", "169.254.169.254", "169.254.0.1", "224.0.0.1", "0.0.0.0", "[fe80::1%25eth0]"} {
		doc, err := c.FetchDocument(context.Background(), "http://"+host+"/openapi.yaml")
		if err == nil {
			t.Fatalf("host %q: expected rejection, got %d bytes", host, len(doc.Content))
		}
		// An unreachable address also errors; only the guard's refusal counts.
		if msg := err.Error(); !strings.Contains(msg, "private") && !strings.Contains(msg, "localhost") && !strings.Contains(msg, "scoped") {
			t.Fatalf("host %q: error %q is not the guard's refusal", host, msg)
		}
		if len(doc.Content) != 0 || doc.FinalURL != nil || doc.ContentType != "" {
			t.Fatalf("host %q: a refused fetch returned data: %+v", host, doc)
		}
	}
}

func TestFetchDocumentRejectsBadURLsAndNilClient(t *testing.T) {
	var nilClient *Client
	for _, raw := range []string{"", "   ", "not a url", "ftp://example.com/x", "http://user:pass@example.com/x"} {
		if _, err := (&Client{}).FetchDocument(context.Background(), raw); err == nil {
			t.Fatalf("input %q: expected an error", raw)
		}
	}
	if _, err := nilClient.FetchDocument(context.Background(), "http://127.0.0.1/x"); err == nil {
		t.Fatal("a nil client must still apply the host guard")
	}
}

func TestFetchDocumentReturnsContentFinalURLAndContentType(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/start":
			http.Redirect(w, r, "/final.yaml", http.StatusFound)
		case "/final.yaml":
			w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
			_, _ = w.Write([]byte("not an OpenAPI document, and that is fine"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	c := &Client{AllowUnsafeHosts: true}
	doc, err := c.FetchDocument(context.Background(), server.URL+"/start")
	if err != nil {
		t.Fatal(err)
	}
	if string(doc.Content) != "not an OpenAPI document, and that is fine" {
		t.Fatalf("content = %q: the fetch must not validate or alter the document", doc.Content)
	}
	if doc.FinalURL == nil || doc.FinalURL.Path != "/final.yaml" {
		t.Fatalf("final URL = %v, want the redirect target", doc.FinalURL)
	}
	if doc.ContentType != "application/yaml; charset=utf-8" {
		t.Fatalf("content type = %q", doc.ContentType)
	}
}

func TestFetchDocumentDecodesGzipWithinTheBudget(t *testing.T) {
	encoded := gzipBytes(t, []byte(strings.Repeat("a", 512)))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(encoded)
	}))
	defer server.Close()

	doc, err := (&Client{AllowUnsafeHosts: true, MaxBytes: 1024}).FetchDocument(context.Background(), server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if string(doc.Content) != strings.Repeat("a", 512) {
		t.Fatalf("decoded content has %d bytes", len(doc.Content))
	}
}

func TestFetchDocumentReturnsHTTPStatusError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "gone", http.StatusGone)
	}))
	defer server.Close()

	_, err := (&Client{AllowUnsafeHosts: true}).FetchDocument(context.Background(), server.URL)
	var statusErr HTTPStatusError
	if !errors.As(err, &statusErr) || statusErr.Code != http.StatusGone {
		t.Fatalf("err = %v, want an HTTPStatusError with status 410", err)
	}
}

func TestFetchDocumentEnforcesTheByteLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(strings.Repeat("a", 4096)))
	}))
	defer server.Close()

	doc, err := (&Client{AllowUnsafeHosts: true, MaxBytes: 1024}).FetchDocument(context.Background(), server.URL)
	if err == nil || !strings.Contains(err.Error(), "larger than") {
		t.Fatalf("err = %v, want a size error", err)
	}
	if len(doc.Content) != 0 {
		t.Fatalf("an oversize fetch returned %d bytes", len(doc.Content))
	}
}

func TestFetchDocumentEnforcesTheTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(5 * time.Second):
		}
	}))
	defer server.Close()

	start := time.Now()
	_, err := (&Client{AllowUnsafeHosts: true, Timeout: 100 * time.Millisecond}).FetchDocument(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected the request to time out")
	}
	if elapsed := time.Since(start); elapsed > 3*time.Second {
		t.Fatalf("the timeout was not applied: returned after %v", elapsed)
	}
}

func TestFetchDocumentReachesLoopbackOnlyWhenAllowed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer server.Close()

	if _, err := (&Client{}).FetchDocument(context.Background(), server.URL); err == nil {
		t.Fatal("a default client must refuse a loopback server")
	}
	if doc, err := (&Client{AllowUnsafeHosts: true}).FetchDocument(context.Background(), server.URL); err != nil || string(doc.Content) != "ok" {
		t.Fatalf("an allowing client: content %q, err %v", doc.Content, err)
	}
}

func TestFetchDocumentRefusesAnUnlistedPort(t *testing.T) {
	// A public literal address, so the refusal comes from the port policy and no
	// DNS lookup or connection is involved.
	for _, raw := range []string{"http://93.184.216.34:8080/openapi.yaml", "https://93.184.216.34:8443/openapi.yaml"} {
		_, err := (&Client{}).FetchDocument(context.Background(), raw)
		if err == nil || !strings.Contains(err.Error(), "refusing URL port") {
			t.Fatalf("%s: err = %v, want the port policy's refusal", raw, err)
		}
	}
}
