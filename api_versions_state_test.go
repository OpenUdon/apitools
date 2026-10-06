package apitools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAPIVersionStateIdentityAndConfinement(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "state.json")
	req := versionTestRequest("https://example.com")
	identity := versionRequestIdentity(req, APIVersionDiscoveryOptions{})
	state := apiVersionState{SchemaVersion: APIVersionStateSchemaVersion, Identity: identity, Report: APIVersionDiscoveryReport{SchemaVersion: APIVersionDiscoverySchemaVersion, Known: req.Known, CheckedAt: time.Now().UTC(), Status: "none_found_in_scope"}}
	data, _ := json.Marshal(state)
	if err := writeVersionMetadata(context.Background(), file, data); err != nil {
		t.Fatal(err)
	}
	if _, ok, err := loadVersionState(file, identity); err != nil || !ok {
		t.Fatalf("fresh state ok=%v err=%v", ok, err)
	}
	if _, ok, err := loadVersionState(file, "different"); err != nil || ok {
		t.Fatalf("foreign state ok=%v err=%v", ok, err)
	}
	link := filepath.Join(dir, "link.json")
	if err := os.Symlink(file, link); err != nil {
		t.Fatal(err)
	}
	if _, _, err := loadVersionState(link, identity); err == nil {
		t.Fatal("symlink state accepted")
	}
	if err := writeVersionMetadata(context.Background(), link, data); err == nil {
		t.Fatal("symlink write accepted")
	}
	if err := os.WriteFile(filepath.Join(dir, ".git"), []byte("gitdir: fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := writeVersionMetadata(context.Background(), file, data); err == nil {
		t.Fatal("repository state write accepted")
	}
}

func TestAPIVersionDirectoryOriginsAndAge(t *testing.T) {
	req := versionTestRequest("https://example.com")
	req.Known.ProviderKey = "example.com:mail"
	data := []byte(`{"example.com:mail":{"preferred":"v2","versions":{"v2":{"swaggerUrl":"https://catalog.example/v2.json","updated":"2023-04-21T00:00:00Z","info":{"x-origin":[{"url":"https://example.com/v2.json"}]}},"v10":{"swaggerUrl":"https://catalog.example/v10.json","info":{"x-origin":{"url":"https://example.com/v10.json"}}}}}}`)
	rows, err := versionListRecords(context.Background(), data, req)
	if err != nil || len(rows) != 2 || rows[0].VersionClaims.URLToken != "v10" || rows[1].Evidence != "catalog-preferred" || rows[1].CatalogUpdated == nil || rows[0].SourceURL != "https://example.com/v10.json" {
		t.Fatalf("rows=%#v err=%v", rows, err)
	}
	if got := versionCatalogURL(&Client{}, req); got != "https://api.apis.guru/v2/example.com.json" {
		t.Fatal(got)
	}
	if got := versionCatalogURL(&Client{APIsGuruListURL: "https://mirror.example/list"}, req); got != "https://mirror.example/list" {
		t.Fatal(got)
	}
	bad := req
	bad.OfficialScopes[0].PathPrefix = "/../secret"
	if _, err := validateVersionRequest(&Client{}, bad, APIVersionDiscoveryOptions{}); err == nil {
		t.Fatal("scope escape accepted")
	}
}
