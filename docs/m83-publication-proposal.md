# M83 exact source-publication proposal

This is a concrete proposal, not publication authority. The confirmed goal permits
local implementation/task commits and no external mutations. Prior Udon/UWS grants
are consumed and do not cover APItools. M83.7 remains open; no successor acceptance,
publication, retirement, deployment or consumer adoption is claimed.

The frozen reviewed runtime/source is
`f2c5693ec39ad6981693d2ad126ff26e3fdd564c`, committed at
`2026-10-08T06:14:39Z`; module
`github.com/OpenUdon/apitools@v0.0.0-20261008061439-f2c5693ec39a`.
Whole local code review passed at 4/10 and all owner checks pass. The later
metadata-only evidence/proposal head carries `[skip ci]` and leaves this source
identity unchanged. Its full HEAD is supplied with the handoff; the source is
its ancestor. No `.github` workflow or deployment launcher is tracked in this
repository. The proposal includes no deployment/build-service action.

Local bootstrap module identity:

- Archive Go sum: `h1:+UICyuitKDE3g5rnlCA6Sk+QLIEqqVn8MN0PtOX+NKk=`.
- GoMod sum: `h1:WmUXlfBBoaI6vtv/wzeZfyO8q/pZMhnHiBckkAthYfk=`.
- Raw bootstrap ZIP: 2,202,789 bytes / SHA-256 `0934422193f9bf31b6b98ed917e9c4458ea0d8c4f8ee6de79464b475f8d13cea`.
- All 614 archive files byte-match both frozen source and extracted artifact.
- File-manifest SHA-256: `fe772598292a59f8ff96f3da3f226ff861f62be69e6d75590bc74012a8a86621`.

[Bootstrap proof](m83-bootstrap-proof.json) is explicitly local, without workspace
or directory replacement; it is not independent public resolution. Exact-archive
owner full tests, root/GraphQL/sourceguard races, vet and build pass. An independent
public-API consumer with exact APItools/M09 requirements passes full tests/races/vet/
build, including default/response-key merge, native Unicode, no-fetch and unknown
format/indeterminate binding evidence. Its complete selected graph has 64 modules
including the consumer, 33 compiled modules / 353 packages; all 15,444 selected
artifact files byte-match the cache ZIPs. Owner qualification has 77 selected
modules including owner, 39 compiled modules / 414 packages; all 15,657 dependency
files plus 614 owner files are checked. The [closure](m83-module-closure.json)
records every selected artifact sum and full file-manifest hashes. Extra lower
version graph metadata is cached to allow offline pruning; it is not represented
as selected compiled code. Bootstrap checksum lookup is disabled explicitly;
recorded hashes/content prove the local artifact and cannot replace public proof.

Full retained-pin OpenUdon and Udon suites pass. Candidate-bound Udon full tests
pass through a disposable modfile and the exact bootstrap archive. Candidate-bound
OpenUdon full tests have an unresolved integration failure in
`packagev3.TestSourceBackedChainedInputTypesRemainExact` at unchanged OpenUdon source
`2b4382011fe0f98b52a8f0c892bd64da78bfa82f`. All other packages pass; the full candidate
suite is not passed. The same focused failure reproduces with retained M82 APItools
`v0.0.0-20261006210844-54583f9b2f45` plus corrected M09 root. This isolates it from
an M83-specific exported-API change; it does not establish that the compatibility
expectation is invalid or that a fixture-only repair suffices. Parent revalidation
against accepted M09 and pending M99 scope remains required. The unchanged fixture
uses OpenAPI 3.0.3, a required response integer `n` with minimum 9007199254740993,
`$response.body.n`, then `$steps.fetch.outputs.n`; it expects compatible but reports
`binding.output_field_indeterminate` and `flow.output_unreferenced`. No test expectation,
source, ledger or pin in either sibling is changed. Generated disposable Udon test
trees are retained outside its checkout. Exact failed/passed logs and temporary
modfiles remain under `/home/peter/.cache/apitools-m83-proof/`.

Proposed publication scope after that evidence is reviewed:

1. Normal fast-forward source publication to
   `git@github.com-tabilet:OpenUdon/apitools.git`, `refs/heads/main`, of the exact
   reviewed source and metadata-only evidence head carrying `[skip ci]`.
   Independently observed target baseline is
   `24c36bf40102c2c1d160dc7d0e27fb161e12dbd6`; it is an ancestor of the source.
   No force push, tag, branch deletion or unrelated owner publication is included.
2. Independently resolve the source pseudo-version from the configured ordinary
   origin/module service in fresh caches, verify Origin.Hash, version/time,
   source/ZIP/file content and archive/GoMod sums, and record the ordinary raw ZIP
   identity independently. The local bootstrap ZIP identity is not a remote proof.
3. Reproduce complete owner/public-consumer selected and compiled artifact closures
   and exact files, then run ordinary module-only full tests/races/vet/build.
   Preserve source-isolated unknown/indeterminate cases and no-network metadata behavior.
4. After every acceptance and integration disposition qualifies, record ordinary
   publication evidence, reconcile OpenUdon M99/Kinet M49 to exact accepted/published
   identities, consolidate and retire M83 normally. Any publication of those later
   evidence/retirement commits requires explicit coverage by the named owner grant;
   they cannot change the frozen runtime/source or other owner pins.

This proposal grants no deployment, live user ledger/API/model/mail/host/Cloudflare
operation, registration change, approval upgrade or execution retry. Installed M44,
retained browser/media/Phase A/legacy pins and completed histories stay frozen.
