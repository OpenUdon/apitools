package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/OpenUdon/apitools"
	"github.com/OpenUdon/apitools/internal/catalogupdate"
)

func TestCatalogMaintenancePreflightAndDryRun(t *testing.T) {
	dir := catalogCLIStateFixture(t)
	config := filepath.Join(dir, "maintenance.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"cache","state_db":"state.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if code := runCatalogMaintenance([]string{"run", "--config", config, "--dry-run", "--json"}, &out, &errOut); code != 0 {
		t.Fatalf("dry run exited %d: %s", code, errOut.String())
	}
	var report apitools.CatalogUpdateReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || !report.DryRun || len(report.Providers) != 334 {
		t.Fatalf("bad dry-run scope: %d providers, %v", len(report.Providers), err)
	}
	for _, path := range []string{"state.sqlite", "state.sqlite.monthly.lock", "cache"} {
		if _, err := os.Stat(filepath.Join(dir, path)); !os.IsNotExist(err) {
			t.Fatalf("dry run touched %s: %v", path, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "cache"), 0o700); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	args := []string{"preflight", "--config", config, "--json"}
	if code := runCatalogMaintenance(args, &out, &errOut); code != 0 {
		t.Fatalf("preflight exited %d: %s", code, errOut.String())
	}
	var readiness catalogMaintenanceReadiness
	if err := json.Unmarshal(out.Bytes(), &readiness); err != nil || !filepath.IsAbs(readiness.CacheDir) || !filepath.IsAbs(readiness.Binary) {
		t.Fatalf("preflight did not resolve absolute paths: %+v, %v", readiness, err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) != 2 {
		t.Fatalf("preflight left files: %v, %v", entries, err)
	}
	if err := os.WriteFile(config, []byte(`{"cache_dir":"cache","state_db":"state.sqlite","expected_volume_uuid":"apitools-fixture-lost-volume"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := runCatalogMaintenance(args, &out, &errOut); code != exitRuntime || !strings.Contains(errOut.String(), "apitools-fixture-lost-volume") {
		t.Fatalf("missing mount passed preflight: %d, %s", code, errOut.String())
	}
}

func TestCatalogMaintenanceScratchScheduleAndReports(t *testing.T) {
	dir := catalogCLIStateFixture(t)
	programConfig := filepath.Join(dir, "user's maintenance.json")
	if err := os.WriteFile(programConfig, []byte(`{"cache_dir":"cache","state_db":"state.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "cache"), 0o700); err != nil {
		t.Fatal(err)
	}
	cron := filepath.Join(dir, "scratch.crontab")
	const existing = "CRON_TZ=Etc/UTC\n# User job\n0 8 * * * /opt/bin/report\n"
	if err := os.WriteFile(cron, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	var calls int
	client := &apitools.Client{AllowUnsafeHosts: true, HTTPClient: &http.Client{Transport: catalogFixtureTransport(func(r *http.Request) (*http.Response, error) {
		calls++
		return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("")), Request: r}, nil
	})}}
	var out, errOut bytes.Buffer
	command := func(want int, args ...string) []byte {
		t.Helper()
		out.Reset()
		errOut.Reset()
		if code := runCatalogMaintenanceWithClient(args, &out, &errOut, client); code != want {
			t.Fatalf("%v exited %d, expected %d: %s", args, code, want, errOut.String())
		}
		return bytes.Clone(out.Bytes())
	}
	readCron := func() string {
		t.Helper()
		contents, err := os.ReadFile(cron)
		if err != nil {
			t.Fatal(err)
		}
		return string(contents)
	}
	command(exitRuntime, "enable", "--config", programConfig, "--after-run", "unknown", "--crontab-file", cron)
	if readCron() != existing || calls != 0 {
		t.Fatal("enable without first-run proof touched the scratch schedule or publishers")
	}
	command(exitRuntime, "enable", "--config", programConfig, "--after-run", "unknown", "--crontab-file", programConfig)
	if !strings.Contains(errOut.String(), "must be separate") {
		t.Fatalf("configuration collision was not refused: %s", errOut.String())
	}
	command(exitRuntime, "enable", "--config", programConfig, "--after-run", "unknown", "--crontab-file", programConfig+" ")
	if !strings.Contains(errOut.String(), "trailing whitespace") {
		t.Fatalf("normalization alias was not refused: %s", errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := runCatalogUpdateWithClient([]string{"--all", "--config", programConfig, "--json"}, &out, &errOut, client); code != exitSuccess {
		t.Fatalf("operator's first run exited %d: %s", code, errOut.String())
	}
	var first apitools.CatalogUpdateReport
	if err := json.Unmarshal(out.Bytes(), &first); err != nil || !first.Completed || len(first.Providers) != 334 {
		t.Fatalf("operator's first report has wrong scope: %+v, %v", first, err)
	}
	stateAlias := filepath.Join(dir, "state.alias")
	if err := os.Link(filepath.Join(dir, "state.sqlite"), stateAlias); err != nil {
		t.Fatal(err)
	}
	command(exitRuntime, "enable", "--config", programConfig, "--after-run", first.RunID, "--crontab-file", stateAlias)
	if !strings.Contains(errOut.String(), "aliases a maintenance") {
		t.Fatalf("state alias was not refused: %s", errOut.String())
	}
	enable := []string{"enable", "--config", programConfig, "--after-run", first.RunID, "--crontab-file", cron}
	command(exitSuccess, enable...)
	installed := readCron()
	unmanaged, entry, err := splitCatalogCronTable(installed)
	if err != nil || unmanaged != existing || !strings.HasPrefix(entry, "17 3 1 * * '/") || !strings.Contains(entry, "user'\\''s maintenance.json") || !strings.Contains(entry, first.RunID) {
		t.Fatalf("unexpected managed schedule: %s (%v)", installed, err)
	}
	command(exitSuccess, enable...)
	if readCron() != installed {
		t.Fatal("enable changed an already installed entry")
	}
	beforeRefused := calls
	command(exitRuntime, "run", "--config", programConfig, "--after-run", "unknown", "--json")
	if calls != beforeRefused {
		t.Fatal("scheduled invocation contacted publishers before checking the first-run proof")
	}
	run := command(exitSuccess, "run", "--config", programConfig, "--after-run", first.RunID, "--json")
	var scheduled catalogMaintenanceReport
	if err := json.Unmarshal(run, &scheduled); err != nil || !scheduled.Run.Completed || scheduled.Run.RunID == first.RunID ||
		scheduled.Summary.ProvidersChecked != 334 || scheduled.Summary.Failures == 0 || scheduled.CacheUsage == nil ||
		scheduled.EmbeddingUsage.Status != "not-requested" || scheduled.EmbeddingUsage.Requests != 0 {
		t.Fatalf("wrong scheduled report: %s, %v", run[:min(len(run), 500)], err)
	}
	status := command(exitSuccess, "status", "--config", programConfig, "--crontab-file", cron, "--json")
	var observed catalogCronStatus
	if err := json.Unmarshal(status, &observed); err != nil || !observed.Enabled || observed.LatestRunID != scheduled.Run.RunID {
		t.Fatalf("wrong schedule status: %s, %v", status, err)
	}
	latest := command(exitSuccess, "report", "--config", programConfig, "--json")
	if !bytes.Equal(latest, run) {
		t.Fatal("latest report differs from the saved scheduled report")
	}
	text := command(exitSuccess, "report", "--config", programConfig, "--run", first.RunID)
	if !bytes.Contains(text, []byte("Run "+first.RunID+" is complete")) || !bytes.Contains(text, []byte("Providers checked 334/334")) {
		t.Fatalf("saved text report missing operator evidence: %s", text[:min(len(text), 500)])
	}
	cache, err := catalogupdate.OpenArtifactCache(filepath.Join(dir, "cache"), "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer cache.Close()
	store, err := catalogupdate.OpenStore(t.Context(), filepath.Join(dir, "state.sqlite"), cache)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	now := time.Now()
	previousPeriod := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0).Format("2006-01")
	pending, err := store.MonthlyRun(t.Context(), previousPeriod, first.RunIdentity, first.CatalogSHA256, []string{"attio"})
	if err != nil {
		t.Fatal(err)
	}
	obsoleteProof, err := store.StartRunSnapshot(t.Context(), first.RunIdentity, first.CatalogSHA256, []string{"attio"})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Record(t.Context(), obsoleteProof.ID, catalogupdate.Outcome{ProviderID: "attio", Status: "checked"}); err != nil {
		t.Fatal(err)
	}
	if err := store.FinishRun(t.Context(), obsoleteProof.ID); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	beforeSupersede := calls
	command(exitRuntime, "supersede", "--config", programConfig, "--run", scheduled.Run.RunID)
	superseded := command(exitSuccess, "supersede", "--config", programConfig, "--run", pending.ID, "--json")
	var archived catalogMaintenanceReport
	if err := json.Unmarshal(superseded, &archived); err != nil || archived.Run.RunID != pending.ID ||
		archived.Run.MonthlyPeriod != previousPeriod || archived.Run.SupersededAt.IsZero() || archived.Run.Completed || calls != beforeSupersede || readCron() != installed {
		t.Fatalf("supersede did not preserve the run or isolate cron and publishers: %+v, %v", archived.Run, err)
	}
	if text := command(exitSuccess, "report", "--config", programConfig, "--run", pending.ID); !bytes.Contains(text, []byte("is superseded")) {
		t.Fatalf("historical report dropped supersession: %s", text)
	}
	if err := os.WriteFile(cron, []byte(strings.Replace(installed, first.RunID, obsoleteProof.ID, 1)), 0o600); err != nil {
		t.Fatal(err)
	}
	stale := command(exitRuntime, "status", "--config", programConfig, "--crontab-file", cron, "--json")
	if err := json.Unmarshal(stale, &observed); err != nil || !observed.Enabled || !strings.Contains(observed.StateError, "Complete and inspect a full reviewed run") {
		t.Fatalf("status hid stale proof or the installed schedule: %s, %v", stale, err)
	}
	command(exitSuccess, "disable", "--crontab-file", cron)
	command(exitSuccess, "disable", "--crontab-file", cron)
	if readCron() != existing {
		t.Fatal("disable altered unrelated entries")
	}
}

func TestCatalogMaintenanceRejectsUnsafeCrontabEdits(t *testing.T) {
	for _, raw := range []string{catalogCronStart + "\n", catalogCronEnd + "\n",
		catalogCronStart + "\n17 3 1 * * echo one\n17 3 1 * * echo two\n" + catalogCronEnd + "\n"} {
		if _, _, err := splitCatalogCronTable(raw); err == nil {
			t.Fatalf("accepted malformed block %q", raw)
		}
	}
	for _, path := range []string{"/config%other", "/config\n17 3 * * * evil"} {
		if _, err := makeCatalogCronEntry("/bin/apitools", path, "fixture"); err == nil {
			t.Fatalf("accepted unsafe cron config path %q", path)
		}
	}
}

func TestCatalogMaintenanceRejectsMissingSidecarThroughDirectoryAlias(t *testing.T) {
	root := t.TempDir()
	stateDir := filepath.Join(root, "state")
	if err := os.Mkdir(stateDir, 0o700); err != nil {
		t.Fatal(err)
	}
	alias := filepath.Join(root, "alias")
	if err := os.Symlink(stateDir, alias); err != nil {
		t.Fatal(err)
	}
	actual := filepath.Join(stateDir, "state.sqlite-journal")
	protected := filepath.Join(alias, "state.sqlite-journal")
	if _, err := resolveCatalogCronScratchPath(actual, protected); err == nil || !strings.Contains(err.Error(), "must be separate") {
		t.Fatalf("allowed alternate path to a pending journal: %v", err)
	}
}

func TestCatalogMaintenanceBoundsCronOutputWhileCopying(t *testing.T) {
	output := catalogCronOutput{limit: 8}
	if _, err := io.WriteString(&output, "abcd"); err != nil {
		t.Fatal(err)
	}
	if _, err := io.Copy(&output, strings.NewReader("too long")); err == nil || output.buffer.String() != "abcd" {
		t.Fatalf("capture ignored the limit: %q, %v", output.buffer.String(), err)
	}
}

func TestCatalogMaintenanceRejectsEmptySensitiveFlags(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("PATH", dir) // Even a broken selector cannot reach the host crontab.
	config := filepath.Join(dir, "maintenance.json")
	if err := os.WriteFile(config, []byte(`{"cache_dir":"cache","state_db":"state.sqlite"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"disable", "--crontab-file="}, {"run", "--config", config, "--after-run", ""}, {"report", "--config", config, "--run="}} {
		var out, errOut bytes.Buffer
		if code := runCatalogMaintenance(args, &out, &errOut); code != exitUsage || !strings.Contains(errOut.String(), "requires a nonempty value") {
			t.Fatalf("%v did not fail before accessing cron or state: %d, %s", args, code, errOut.String())
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "cache")); !os.IsNotExist(err) {
		t.Fatalf("malformed command touched storage: %v", err)
	}
}

func TestCatalogMaintenanceReportCountsPublisherHeadAndBaselineIndependently(t *testing.T) {
	const baseline = `{"reviewed_baseline":{"version":{"version":"1.0"}},"diagnostic":"Fixture head needs review"}`
	run := apitools.CatalogUpdateReport{RunID: "fixture", Providers: []string{"aws-ec2", "datadog", "stripe", "hubspot"},
		Outcomes: []catalogupdate.Outcome{
			{ProviderID: "aws-ec2", ReferenceID: "aws-ec2-smithy-model", Status: "publisher-spec-unexamined"},
			{ProviderID: "datadog", ReferenceID: "datadog-api-v2", Status: "publisher-spec-invalid", Evidence: []byte(baseline)},
			{ProviderID: "stripe", ReferenceID: "stripe-latest-openapi-spec3", Status: "reviewed-pin-reproduced", Evidence: []byte(baseline)},
			{ProviderID: "hubspot", ReferenceID: "hubspot-public-api-spec-index", Status: "publisher-index-inventoried",
				Evidence: []byte(`{"inventory":{"Sources":[{"Version":"2026-09-beta","Stage":"LATEST","Stable":false,"VersionLabelConflict":true}]}}`)},
		}}
	var out, errOut bytes.Buffer
	if code := writeCatalogMaintenanceReport(t.Context(), &out, &errOut, run, nil, apitools.CatalogUpdateConfig{}, true); code != exitSuccess {
		t.Fatalf("JSON report exited %d: %s", code, errOut.String())
	}
	var report catalogMaintenanceReport
	if err := json.Unmarshal(out.Bytes(), &report); err != nil || report.Summary.ReviewedPins != 2 || report.Summary.PendingReview != 3 || report.Summary.Failures != 1 {
		t.Fatalf("head and baseline counts lost evidence: %+v, %v", report.Summary, err)
	}
	out.Reset()
	if code := writeCatalogMaintenanceReport(t.Context(), &out, &errOut, run, nil, apitools.CatalogUpdateConfig{}, false); code != exitSuccess {
		t.Fatalf("text report exited %d: %s", code, errOut.String())
	}
	for _, line := range []string{
		"aws-ec2/aws-ec2-smithy-model publisher-spec-unexamined (publisher head needs review)",
		"datadog/datadog-api-v2 publisher-spec-invalid (reviewed baseline adopted)",
		"stripe/stripe-latest-openapi-spec3 reviewed-pin-reproduced (publisher head needs review)",
		"hubspot/hubspot-public-api-spec-index publisher-index-inventoried (conflicting publisher release labels need review)",
	} {
		if !strings.Contains(out.String(), line) {
			t.Fatalf("text report omitted %q: %s", line, out.String())
		}
	}
}
