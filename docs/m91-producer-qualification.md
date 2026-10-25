# M91 local producer qualification

This is local, unpublished producer evidence. M91.1–M91.3 are implemented;
M91.4 remains incomplete for required publication and consumer handoff. Whole
review **1/10 completed with an evidence blocker**. No publication, dependency
acquisition or new consumer adoption is established by these checks.

The implementation was frozen at
`0ff4ef561dcc5f784a630b1559ac688bca35802d`, following integration baseline
`1aa42bbbd6cce01b902e1ce612243ec0b61b928a`. The coordinated detached input is
`/home/peter/Workspace/apitools.goal/m91-qualification`; its copied ledger is
frozen input, not a second execution owner. Primary M91 status owns progress.
The source worktree was clean before and after the full qualification.

## Environment and observed checks

All commands used the explicit cached Go1.26.6 executable at
`/home/peter/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.6.linux-amd64/bin/go`,
with `GOWORK=off GOPROXY=off GOSUMDB=off GOAUTH=off GOTOOLCHAIN=local` and
`GOFLAGS=-mod=readonly`. No workspace/directory replacement or dependency/tool
acquisition was used. Existing ordinary module/build caches were reused.
Temporary output stayed under the approved exclusive UID1002/mode0700
`/tmp/k13-8f7f` with process-local umask0077; logs are mode0600.

Full qualification invocation **1** completed in **30.768 seconds** within its
20-minute hard cap:

| Check | Outcome | Seconds |
|---|---|---|
| `go test ./... -count=1 -timeout=15m` | PASS, 18 packages (one without tests) | 18.775 |
| `go vet ./...` | PASS | 8.809 |
| CLI build | PASS | 1.872 |
| Actual CLI `search --help` / `import --help` | PASS | 0.231 combined |
| `go run ./cmd/cataloggen -check` | PASS | 0.818 |
| Actual CLI `catalog check` | PASS | 0.115 |

The selected graph contains **77 modules including main**. Actual CLI BuildInfo
contains **39 modules including main**, Go1.26.6 and no replacements. Exact
versions and sums, canonical identity hashes, command outcomes and retained log
hashes are in [m91-local-proof.json](m91-local-proof.json). Every compiled
dependency agrees with the selected version/sum. The UWS dependency remains
`v0.0.0-20261008043726-b099f6803277`; Horizon `v1.14.5` and HashiCorp HCL
`v2.24.0` remain present. This adds no HCL-removal claim or module upgrade.

Both the default automatic build and the one supplemental explicit
`go build -buildvcs=true` build succeeded, but actual BuildInfo omitted VCS
fields in both. The explicit build took 2.541 seconds and produced the same
CLI SHA256 as the automatic build:
`263fb9948dfd2cf83de055aebb7e483868e438d921313aeb3662799e0985819b`.
No `vcs.revision` or `vcs.modified=false` observation is claimed. The source
revision/cleanliness evidence comes separately from Git. Both original build
observations are retained; no further stamp experiment was attempted.

Affected race invocation **1** passed in **89.189 seconds** within its five-minute
hard cap: `go test -race . ./operationlifecycle -count=1 -timeout=4m` (root test
execution 37.570s; operationlifecycle 1.328s). Source remained clean at the same
revision after settlement. No full-suite replay was performed.

## Existing checker observations

Affected staticcheck invocation **1** used compatible staticcheck2026.2.1
(`0.8.1`), verified tool SHA256
`c42e607df5bc643e53de8a619f18ec93793fa8d57d6981acbaf4cf28c0372ea7`.
It ran `. ./operationlifecycle` once, exited **1** in **17.660 seconds**, and
reported these six observations:

| Code | Existing locations |
|---|---|
| SA1019 | `api_versions.go:14`, `api_versions.go:418`, `api_versions_compare.go:7`, `api_versions_compare.go:39` |
| ST1013 | `catalog_discovery_remote_test.go:139` |
| ST1005 | `operation_source_native.go:402` |

The current raw log SHA256 is
`fa29556fac17c1cb6f04e1d5788595336fe9892f171141ae4d2b52c4e34f53f5`.
There are no new M91-file diagnostics. The coordinator independently reconciled
the [existing M85 policy](m85-qualification.md): six retained observations,
no suppression, new waiver or zero-result claim. Three diagnostic-bearing
files match the frozen M85 `checkerFiles` hashes exactly. The fourth changed
only by an import and appended tests; its offending `302` expression is
unchanged, with its location shifted by one line. The old raw checker log is
unavailable, so byte-identical historical output is **not** claimed. The actual
current nonzero exit is preserved. No checker retry or unrelated source repair
was performed.

The checker used its own explicit mode0700
`/tmp/k13-8f7f/m91-staticcache`; historical default cache links were preserved.
The older installed staticcheck2024.1.1 was inspected but not run on M91 code.

## Handoff and remaining acceptance

The [SDK contract](operation-mapping-v1.md) and
[versioned fixture corpus](../testdata/operation-mapping/v1/README.md) are the
consumer handoff inputs. The 19 actual-source cases cover business-message and
fake Cloudflare resource mappings, while the frozen eight-family partial
shapes remain unresolved. Ordinary producer external-package tests pass; these
are local API compatibility evidence, not new published-consumer adoption.

All pushes remain denied by the human's “commit in all packages but no push”
instruction and `EXTERNAL_MUTATIONS:none`. Independent ordinary publication,
exact published-module acquisition/selection and actual W22/M136 use remain
unperformed. Existing unchanged consumer pins cannot establish M91 adoption;
no `go.work` or directory replacement was substituted. Udon consumer execution
is outside this executor's four-repository effective verification subset and
was not run. OpenUdon:P14 is independently owned and is not an M91 dependency
or import. APItools:M90 stays reserved and unscheduled.

The source worktree, private raw logs/results/binaries and checker cache remain
retained for coordinator reconciliation under the approved roots. Whole
acceptance, fresh bounded whole review, shared-memory consolidation and
retirement remain coordinator work after missing required evidence is resolved.

The fresh read-only whole reviewer inspected all five lenses and found no
production-code P1/P2 defect or P3 finding. Its P2 M91-W1 retains the required
ordinary publication, reusable consumer proof and available existing-consumer
checks. New W22/M136 adoption remains downstream work, not an additional M91
prerequisite. The full gate has not passed; the next complete review is
iteration 2 after the missing evidence is obtained.

M91.4's narrow DEFAULT task review checked commands/results, source and module
identities, fixture handoff and remaining authority/evidence. This was not a
whole-milestone five-lens review; it consumes no whole-review iteration.
Consumed invocations remain full1 / affected-static1 / race1 / explicit-VCS
build1, with no retries. Format and patch checks pass.

Coordinator documentation updates add the README guide link and align the
owning milestone/status summaries with this local readiness. Production,
fixtures, module pins and test inputs remain the qualified
`0ff4ef561dcc5f784a630b1559ac688bca35802d` bytes;
the passing full/race evidence is reused for this documentation-only delta.
