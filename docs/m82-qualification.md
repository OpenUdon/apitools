# M82 consumer qualification — in progress

M82.1–M82.3 are committed. M82.4 is selected and in progress. This record does
not establish whole-milestone acceptance or source publication. The closing
whole code review passed at 3/10. Publication and final milestone acceptance
remain pending. The confirmed Stage 11 goal remains active.

## Source and public contract

Implementation baseline: `18fdf1fc8063384e637143d2c39a322f404463f6`, with
uncommitted reviewed M82.4 HTTP completeness, aggregate resource and source-dialect/ref-sibling guards plus
their negative fixtures. That baseline does not contain the later guard changes.
The producer/verifier APIs and public UWS wire are documented in
[operation-shapes.md](operation-shapes.md).

UWS C09 is pinned to accepted source
`6a267306032edc687a298cefc8bba7019d3ad059`, pseudo-version
`v0.0.0-20261006181058-6a267306032e`. Its accepted source and retirement closure
`8e5be730aa68aa4cb4f6c591a3a9425c6b8bd55c` are ancestors of independently
observed authorized UWS origin/main
`0e1e10d0e3e39b20768f1b4862c51bbb1c6cf9cf`. Observed APItools origin/main is
still `9364afac98c97b79c9e6b9335bd9c71d4e8722d5`; no M82 source was published.
The destination remains `git@github.com-tabilet:OpenUdon/apitools.git`, main.
Publication must follow final qualification/review and use the existing normal
fast-forward grant; consumers must wait for exact accepted publication evidence.

The frozen eight-family fixture is 9,829 bytes, SHA-256
`dc20d2287322a4b20d5d92b1a0d1ec8036f1bec853b642d7df8797ed096c3ba1`.
It contains eight sources/twelve native operations. Prior candidate v1 wires,
source artifacts and existing exported metadata APIs remain unchanged. Source
tooling still performs no credential resolution or workflow/API execution.
Native protocol/schema/security limits remain explicit. The standalone module
closure retains Horizon `v1.14.5` / HCL `v2.24.0`.

## Verification obtained so far

Go1.26.6 standalone/offline APItools full tests and vet passed after the M82.4
guards. Focused shape race checks passed, including all eight source families,
supported advisory compatibility, OR/AND security, forged/stale claims, exact
JSON/YAML numeric lexemes, ambiguous encodings/IDs/native collisions, missing/
cyclic/external refs, protocol limitations, resolver isolation and bounded wire.
Catalog generated outputs are current; gofmt and patch checks pass. Compatible
retained staticcheck2026.2.1 reported six pre-existing diagnostics: four SA1019
Discovery-wrapper uses in api_versions.go/api_versions_compare.go, one ST1013
status literal in catalog_discovery_remote_test.go, and one ST1005 error string
in operation_source_native.go. All four files are byte-identical to the pre-M82
baseline `ab2e944`; no diagnostic concerns the new shape files. This is not a
zero-diagnostic result. APItools's required milestone commands do not mandate
zero staticcheck; no suppression, waiver or unrelated legacy cleanup was added.

Workspace full OpenUdon/Udon regressions passed during M82.1/M82.2. OpenUdon
standalone `GOWORK=off GOPROXY=off go test ./...` passed against its retained
module pins at source `77b400d778ea1bba112ded74e1f787d19c38e4c4`. This proves
retained compatibility, not adoption of the new APItools pin.

Udon source `2b022c49be0bd269e59a09cabf310cba5fca79cb` retains local sibling
replacements even with GOWORK off. Its original readonly module invocation
reports that go.mod needs updating because the new APItools dependency raises
UWS to accepted C09. No tracked Udon file was changed. A disposable modfile
changed only that UWS require to
`v0.0.0-20261006181058-6a267306032e`; offline `go test -mod=mod -modfile=FILE
./...` passed. The temporary modfile SHA-256 was
`8eb4ec8b4590b37aa0708dbaa96d6de1bf4dd1da26da08625b95465228a8099e`, sum file
`14d9783a41313dc7a6cd82fd72b7da844b94f775a8e1f527b2c6d4240990b03e`.
Actual consumer dependency adoption remains with the pending owning milestone;
this qualification does not silently change it.

## Remaining acceptance work

Finish M82.4 evidence/verification, persist and perform the full bounded closing
review, fix every blocking finding, publish/independently verify exact accepted
source, reconcile UWS:M08/Kinet:M46/OpenUdon:P09 and retire the complete local
record. Preserve all frozen history and consumer/browser pins. No downstream
milestone may start from this draft evidence. Deployment and live migration/
provider/model/mail/registration operations remain outside the confirmed goal.

## Whole code review

Three persisted iterations completed the code gate. Review 1 fixed aggregate
schema/reference amplification and source-dialect keyword proof; review 2 fixed
OpenAPI 3.0 boolean/union/null evidence; review 3 found no remaining P1/P2.
The final full standalone APItools tests/vet and focused races passed, as did
current full workspace OpenUdon/Udon regression suites. The frozen fixture and
earlier candidate/source bytes remain unchanged. No evolution bump: implementation
follows the approved Stage 11 direction. Exact source publication and completed
task/retirement metadata remain necessary before downstream execution.
