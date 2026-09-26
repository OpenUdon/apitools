package apitools

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestLocalFilesFindsValidOpenAPIDocuments(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(dir, "support.yaml"), `openapi: 3.0.0
info:
  title: Support Ticket API
  version: 1.0.0
  description: Fetch and update support tickets.
paths:
  /tickets:
    get:
      responses:
        "200":
          description: ok
`)
	writeLocalFile(t, filepath.Join(dir, "notes.yaml"), `not: openapi`)

	got, err := LocalFiles(context.Background(), LocalOptions{
		Dir:     dir,
		BaseDir: base,
		Query:   "support ticket",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1: %#v", len(got), got)
	}
	if got[0].RelativePath != "openapi/support.yaml" {
		t.Fatalf("relative path = %q", got[0].RelativePath)
	}
	if got[0].Title != "Support Ticket API" || got[0].Metadata.OperationCount != 1 {
		t.Fatalf("metadata = %#v", got[0])
	}
	if got[0].Score == 0 {
		t.Fatalf("score should be positive")
	}
}

func TestLocalFilesResolvesSymlinkedAncestorsOfScanRoot(t *testing.T) {
	realBase := t.TempDir()
	dir := filepath.Join(realBase, "openapi")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(dir, "support.yaml"), validLocalSpec("Support API", "Support"))
	aliasParent := filepath.Join(t.TempDir(), "workspace")
	symlinkOrSkip(t, realBase, aliasParent)

	got, err := LocalFiles(context.Background(), LocalOptions{
		Dir:     filepath.Join(aliasParent, "openapi"),
		BaseDir: aliasParent,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelativePath != "openapi/support.yaml" || got[0].Path != filepath.Join(dir, "support.yaml") {
		t.Fatalf("local results = %#v", got)
	}
}

func TestLocalFilesFindsDraftOpenAPIDocuments(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(dir, "draft.yaml"), `openapi: 3.0.0
info:
  title: Draft API
  version: 1.0.0
paths:
  /items:
    get:
      operationId: listItems
`)

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base, Query: "draft"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Title != "Draft API" || got[0].Metadata.OpenAPI != "3.0.0" {
		t.Fatalf("results = %#v", got)
	}
}

func TestLocalFilesSortsByScoreThenPath(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(dir, "z.yaml"), validLocalSpec("Billing API", "Invoices"))
	writeLocalFile(t, filepath.Join(dir, "a.yaml"), validLocalSpec("Mail API", "Send mail"))

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base, Query: "mail"})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].RelativePath != "openapi/a.yaml" {
		t.Fatalf("first = %#v", got)
	}
}

func TestLocalFilesRejectsMissingDir(t *testing.T) {
	_, err := LocalFiles(context.Background(), LocalOptions{})
	if err == nil || !strings.Contains(err.Error(), "directory is required") {
		t.Fatalf("expected directory error, got %v", err)
	}
}

func TestLocalFilesRejectsSymlinkedRoot(t *testing.T) {
	base := t.TempDir()
	realDir := filepath.Join(base, "real-openapi")
	if err := os.MkdirAll(realDir, 0o755); err != nil {
		t.Fatal(err)
	}
	linkDir := filepath.Join(base, "openapi")
	symlinkOrSkip(t, realDir, linkDir)

	_, err := LocalFiles(context.Background(), LocalOptions{Dir: linkDir, BaseDir: base})
	if err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink error, got %v", err)
	}
}

func TestLocalFilesSkipsSymlinkedCandidatesAndKeepsScanning(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(base, "target.yaml")
	writeLocalFile(t, target, validLocalSpec("Target API", "Target"))
	symlinkOrSkip(t, target, filepath.Join(dir, "support.yaml"))

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("symlink candidate was included: %#v", got)
	}
}

func TestLocalFilesKeepsRegularCandidateBesideSymlink(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "target.yaml")
	writeLocalFile(t, target, validLocalSpec("Target API", "Target"))
	symlinkOrSkip(t, target, filepath.Join(dir, "support.yaml"))

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Path != target {
		t.Fatalf("results = %#v, want the regular target only", got)
	}
}

func TestLocalFilesSkipsOversizedCandidatesAndKeepsValid(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	valid := validLocalSpec("Support API", "Support")
	writeLocalFile(t, filepath.Join(dir, "support.yaml"), valid)
	maxBytes := int64(len(valid) + 1)
	writeLocalFile(t, filepath.Join(dir, "oversized.yaml"), strings.Repeat("x", int(maxBytes)+10))

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base, MaxBytes: maxBytes})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Path != filepath.Join(dir, "support.yaml") {
		t.Fatalf("results = %#v, want the valid candidate only", got)
	}
}

func TestLocalFilesDeduplicatesAndEnforcesTraversalBounds(t *testing.T) {
	base := t.TempDir()
	dir := filepath.Join(base, "openapi")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := validLocalSpec("Support API", "Support")
	writeLocalFile(t, filepath.Join(dir, "a.yaml"), content)
	writeLocalFile(t, filepath.Join(dir, "b.yaml"), content)
	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].RelativePath != "openapi/a.yaml" {
		t.Fatalf("deduplicated results = %#v", got)
	}

	writeLocalFile(t, filepath.Join(dir, "c.yaml"), validLocalSpec("Mail API", "Mail"))
	visitBound, err := discoverLocalSources(context.Background(), LocalSourceDiscoveryOptions{
		Roots:             []string{dir},
		MaxVisitedEntries: 2,
	}, true, base)
	if err != nil {
		t.Fatal(err)
	}
	if !visitBound.report.Truncated || len(visitBound.localResults) > 1 || len(visitBound.report.Diagnostics) == 0 || !strings.Contains(visitBound.report.Diagnostics[0].Message, "visit limit") {
		t.Fatalf("bounded traversal = %#v, results %#v", visitBound.report, visitBound.localResults)
	}
	candidateBound, err := discoverLocalSources(context.Background(), LocalSourceDiscoveryOptions{
		Roots:         []string{dir},
		MaxCandidates: 1,
	}, true, base)
	if err != nil {
		t.Fatal(err)
	}
	if !candidateBound.report.Truncated || len(candidateBound.localResults) != 1 || len(candidateBound.report.Diagnostics) == 0 || !strings.Contains(candidateBound.report.Diagnostics[0].Message, "candidate acceptance limit") {
		t.Fatalf("bounded candidates = %#v, results %#v", candidateBound.report, candidateBound.localResults)
	}
}

// TestLocalFilesFailsOnUnreadableRoot proves that a scan root that cannot be
// read (for example, permission denied) is a scan failure, not an empty
// result: descendant read errors are recorded as rejections so the rest of a
// scan can continue, but the root itself failing means the scan made no
// progress at all.
func TestLocalFilesFailsOnUnreadableRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("directory permission bits do not restrict reads the same way on Windows")
	}
	base := t.TempDir()
	dir := filepath.Join(base, "locked")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLocalFile(t, filepath.Join(dir, "spec.yaml"), validLocalSpec("Locked API", "Locked"))
	if err := os.Chmod(dir, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(dir, 0o755)
	if _, err := os.ReadDir(dir); err == nil {
		t.Skip("this process can read a mode-000 directory (likely running as root); the permission check cannot be exercised")
	}

	got, err := LocalFiles(context.Background(), LocalOptions{Dir: dir, BaseDir: base})
	if err == nil {
		t.Fatalf("LocalFiles() on an unreadable root = (%#v, nil), want an error", got)
	}
	if len(got) != 0 {
		t.Fatalf("LocalFiles() on an unreadable root returned %d results, want none", len(got))
	}
	if !strings.Contains(err.Error(), "cannot read local source root") {
		t.Fatalf("err = %v, want a root-read failure", err)
	}
}

func writeLocalFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func symlinkOrSkip(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
}

func validLocalSpec(title, description string) string {
	return `openapi: 3.0.0
info:
  title: ` + title + `
  version: 1.0.0
  description: ` + description + `
paths:
  /items:
    get:
      responses:
        "200":
          description: ok
`
}
