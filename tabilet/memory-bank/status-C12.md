# C12 — Public catalog website

**Owner.** APItools. **Dependency.** Accepted [C13](../docs/history/status-C13.md). **Parallel-safe: no.**
C11 is skipped because the [M86](status-M86.md) records use an unqualified `C11` for `UWS:C11`. C12 needs its own launch after C13. See the [milestone specification](milestone.md#c12--public-catalog-website) for scope, exclusions, acceptance, class O/R verification and `AUTHORIZATION_REQUIREMENTS`.

**Closing duty.** C12 is the last active catalog milestone. C12.5 closes the [open review findings](milestone.md#open-review-findings) register under its existing rules and verifies the accepted C13 dispositions for FND-05, FND-07 and FND-08 and the owner-approved FND-04 deferral.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C12.1 — Sign-off and readiness | `[+]` | Before sign-off, resolve how operation labels such as Jenkins "Get CSRF crumb" and FileMaker "Delete session" are published; these examples come from C08 whole review round 1 L-2 ([C08 record](../docs/history/status-C08.md#status-record)). Record the owner's decision on publishing provider names, descriptions, operation labels, license notes and links, plus a correction/takedown contact. Record hostname `apis.kinet.com`, its DNS owner and the GitHub Pages host. Under explicit authority, read current DNS and host settings without changing them, and confirm read-only that the hostname is not already in use and that nothing found (such as a CAA record) would prevent a certificate; the certificate itself is proved in C12.5. Class R check under the [verification-attempt gate](milestone.md#verification-attempt-gate); keep DNS/host credentials out of the repository and registrar account details out of planning records. Blocks C12.2–C12.5. |
| C12.2 — Data and generator | `[+]` | Add maintainer-only `cmd/catalogsite` using the Go standard library, with no new public Go API. Build deterministically from an identified clean commit of reviewed catalog data with a pinned Go toolchain; the generator makes no network request, while a hosted runner may download pinned modules verified by `go.sum`. Cover every built-in provider at the built commit (334 at planning); emit no spec or overlay bytes, no auth guidance, and no claim that a provider or version is current or latest. Take the base URL as an explicit build input and use relative internal links. Emit versioned JSON in APItools' own small format with a documented schema version and the source commit, and publish a SHA-256 manifest of the output. |
| C12.3 — Pages and search | `[+]` | Build provider and category pages with coverage, provenance and verification dates, plus in-browser keyword search by name, alias, category and keyword, with every provider and category page reachable without JavaScript and no third-party assets, logos, analytics, accounts or user submissions. Search matches words in public profile text and must not present an incidental operation-label match as a supported capability. Follow the evidence rules accepted in [C08](../docs/history/status-C08.md#status-record), whose whole review round 1 P1-1 closed in round 2. C13.1 delivered the CLI search correction for [FND-07](milestone.md#open-review-findings). Verify the website's own public-profile search behavior. Escape all catalog text and allow only HTTPS source links. Add a restrictive in-page CSP (Pages cannot set response headers to block framing) and accessible markup with named Go-test checks for document language, title, landmarks, form labels and link text. |
| C12.4 — Monthly changes and publishing | `[+]` | Compare two reviewed revisions and commit a small monthly record; the monthly changes page matches it, including the first publication with no predecessor. Add a manual-only GitHub Pages artifact workflow with every GitHub Action pinned by commit SHA and least permissions, committing no generated pages to any branch and keeping the monthly job unable to publish. Add `apitools-catalog-maintenance` skill steps for preview, publish and rollback, and exercise rollback on the loopback preview. Document domain verification, DNS rollback and removal: verify the domain before adding DNS, and remove DNS before retiring the site. |
| C12.5 — First publication and closure | `[+]` | Ordered going-live steps, each needing separate explicit approval: (a) push the workflow to `main`, which publishes every unpushed local commit, so list them for the owner first; (b) verify the domain with GitHub; (c) enable Pages with GitHub Actions as source and set the custom domain; (d) add the DNS record; (e) dispatch the deployment; (f) enforce HTTPS once the certificate issues. Each deployment is a counted class R attempt under the [verification-attempt gate](milestone.md#verification-attempt-gate) and its existing limit. Check that live files match the published SHA-256 manifest for the recorded commit, that `apis.kinet.com` serves a valid certificate, and that HTTP redirects to HTTPS. Own horizon closure of the open review findings register. Verify the accepted C13 dispositions for FND-05, FND-07 and FND-08 and the owner-approved FND-04 deferral. Run the required class O whole-unit gates before whole review. |

## Launch and pending inputs

The owner replied “i approve” to the C12-only launch checksum
`3314da274f273ecc58076e566289f4582b0fc1b0126e012f15b6d9b757e9c94f`.
The raw launch matches that checksum. C12 starts at local `main` commit
`b4f2e3fea3c1794dd5b4c5de6bd0ee4fe4dd9ea7`, with C13 accepted and five
C13 commits awaiting publication. The stock protocol identity is
`4facf0102dd672fc43ff6e4a79384b3a5d0bbfe77a06ed2e55c401a8650d05ae`.
APItools has no project addendum.

After approving the launch, the owner replied “i approve” to the separate
C12.1 publication proposal and named read-only readiness check. The approved
site fields are provider names, descriptions, operation labels, license notes
and source links, excluding the Jenkins “Get CSRF crumb” and FileMaker
“Delete session” operation labels. The owner confirmed ownership of the
website and said they can coordinate DNS changes through Cloudflare. The
owner is the DNS decision owner for C12. The public correction and takedown
contact is the [APItools issue tracker](https://github.com/OpenUdon/apitools/issues),
as specified by the owner. The selected host is GitHub Pages for
`OpenUdon/apitools`.

Required gate `c12.1-readiness` attempt `c12.1-readiness-a1` started from local
commit `b4f2e3fea3c1794dd5b4c5de6bd0ee4fe4dd9ea7` under the owner's direct
approval of one supervised Class R pass. The launch review gate limit and
distinct-key ceiling are both 3; no key was charged before dispatch. The live
scope is bounded `dig` records for `apis.kinet.com` and its parent
`kinet.com`, and one `gh api --method GET repos/OpenUdon/apitools/pages`
using the existing GitHub CLI configuration. No retry, setting change,
publication, or deployment is included.

Attempt `c12.1-readiness-a1` ran once on 2026-11-25 at 23:39:28Z. All
six DNS queries answered. `apis.kinet.com` returned A, AAAA and HTTPS answers
consistent with Cloudflare. CNAME returned no answer. CAA returned no answer
at the hostname or parent zone. These DNS answers do not distinguish a
dedicated hostname from a wildcard or establish the origin behind the proxy.
The GitHub Pages GET returned HTTP 404, leaving Pages state unverified under
the current GitHub CLI access. No settings changed.
[Read-only readiness receipt](../../dist/c12-goal/c12.1-readiness-a1.txt).

The required gate failed on two keys. `dns-host-availability` records
existing A/AAAA/HTTPS answers with hostname use unresolved.
`github-pages-config` records HTTP 404 with Pages state unverified.
Each key used 1 of 3 gate attempts. The distinct failure chain has 2 of 3
signatures. The owner's one-attempt Class R grant is consumed.

The owner subsequently confirmed that `apis.kinet.com` can be repointed
to GitHub Pages without interrupting an existing service. This resolves the
hostname-use question. The owner directly inspected the repository's
[Pages settings](https://github.com/OpenUdon/apitools/settings/pages) and
attested that GitHub Pages is disabled and that the source control showed
“Branch.” This is owner-supplied settings evidence. The original API 404
and charged attempts remain in the record. The owner accepted this evidence
with the earlier DNS results and directed that no further live readiness
check is needed. No external changes were approved. C12.1 resumed for
local verification and task review.

Local C12.1 verification passed at the documented base
`b4f2e3fea3c1794dd5b4c5de6bd0ee4fe4dd9ea7` over the saved
[reviewed ledger diff](../../dist/c12-goal/c12.1-reviewed.diff) with SHA-256
`f6b479ac94d3d3c9aab4e1b0936f79750e0ce3af9b656a90c3873856aa5c046a`.
Using Go 1.26.6 offline, `go build ./...`, `go vet ./...`,
`go run ./cmd/cataloggen -check`, `go run ./cmd/apitools catalog check`
and `git diff --check` passed in about three seconds. The row changes
only evidence in the ledger. No Go package is affected.
[Local gate log](../../dist/c12-goal/c12.1-offline.log).
The separate GPT-6.2 Astra xhigh task review passed with no findings.

C12.2 implemented the clean-source Go generator and public v1 projection.
The first task review found an inherited AWS ELBv2 freshness claim and
unchecked credential-style fragments in HTTPS links. The projection now
uses neutral AWS ELBv2 wording, rejects such fragments, and tests provider
identities and representative coverage labels. The Astra xhigh re-review
passed all three dispositions.

Required C12.2 gate at base `26796ff55521e07ab152434b20d92186193e97e9`
with [candidate diff](../../dist/c12-goal/c12.2-correction.diff) SHA-256
`3190e4c2689f474ab508551f707e8f0720291913b407b3c338d1c350cd5e898f`
passed. Offline Go 1.26.6 `go build ./...`, `go vet ./...`,
`go test -count=1 ./cmd/catalogsite`, and `git diff --check` passed.
Two `go build -trimpath -buildvcs=true` builds of clean candidate snapshot
`64dc18256a92254488bfe6ef526b3aa01bb6af28` produced identical JSON and
manifests. `sha256sum --check SHA256SUMS` passed. Coverage for 334 providers
is 123 official OpenAPI, 55 other official machine-readable, 137 advisory
overlays and 19 docs only. The corrected gate took about five seconds.
[Gate log](../../dist/c12-goal/c12.2-correction.log).

Required C12.3 gate `c12.3-a1` used base
`7f7560e0d8ca94184068fa65f348998b67a1c1ff` and source diff
`9050e1474d4bb71f188ed8cf557048a7ea27fda69e270328096ea483240856d0`.
Build, vet, the whole `cmd/catalogsite` package and the artifact manifest
passed. The isolated Chromium loopback check failed at first navigation
with `ERR_NAME_NOT_RESOLVED`. The browser fixture blocked the numeric
loopback address along with external hosts. Key `chromium-loopback-setup`
used 1 of 3 attempts (one distinct signature of 3). Corrective attempt
`c12.3-a2` starts on the same candidate code and explicitly exempts only
the local address used by the server. [First gate](../../dist/c12-goal/c12.3-row.log).

Attempt `c12.3-a2` passed its offline commands, but independent task review
rejected the scratch Node assertions as browser acceptance evidence. It also
found missing source verification dates on category cards and an incorrect
singular count. The replacement Go test owns the browser assertions and
uses Chromium under test with a disposable profile and loopback only. The
category cards link the first dated source for each provider, and all
directories now render “1 provider.” Corrected required gate `c12.3-a3`
starts on the same base with [saved diff](../../dist/c12-goal/c12.3-review-correction.diff)
SHA-256 `97a84b0d7c05dfb5570fa44df127a56b7225c9db643c375ff63a061c20c5452b`.

Attempt `c12.3-a3` passed in 17 seconds with offline Go 1.26.6
`go build ./...`, `go vet ./...`, `go test -v -count=1 ./cmd/catalogsite`
with `APITOOLS_SITE_TEST_CHROME=/usr/bin/google-chrome`, and
`git diff --check`. The Go test checked eight Chromium 157.0.8103.38
cases against the embedded production script and static HTML. All 388
files in the clean candidate `cf8884527027f7b77a1a87d28c576a0169641358`
manifest verified. The clean pass resets the required gate failure keys.
The separate Astra xhigh re-review passed with no open findings.
[Corrected gate](../../dist/c12-goal/c12.3-review-correction.log).

Required C12.4 gate `c12.4-a1` starts from
`47f0f53b476c2d86281589026585391ef8260777` and [saved diff](../../dist/c12-goal/c12.4-verification.diff)
SHA-256 `2041686bdf721fb6277d669677177110dc5604d1e2143f2f782d03620688e6c2`.
The offline Go gate includes the browser and loopback rollback checks.
The first-record comparison uses clean builds of the two accepted rows.
No live workflow dispatch is part of this gate.

The first C12.4 task review required a stronger same-address rollback
exercise and found that older incomplete ancestors could pass the
workflow's artifact preflight. The corrected Go rollback test verifies all
manifest hashes and restores visibly different content on the original
address. The workflow checks the complete manifest and requires the three
public entry points to be nonempty and covered by it. Corrective gate
`c12.4-a2` starts from the same base with [saved diff](../../dist/c12-goal/c12.4-correction.diff)
SHA-256 `88c922954185d1a594699affdafddd489acbc51132940bc8021f49c4751b689f`.

C12.4 gate `c12.4-a2` passed in 29 seconds with offline Go 1.26.6
`go build ./...`, `go vet ./...`, and `go test -v -count=1 ./cmd/catalogsite`
with Chromium. This included same-address A → B → A loopback rollback of
distinguishable provider pages and complete manifest verification.
`git diff --check` passed. The first record was reproduced byte for byte
using `catalogsite compare` on clean reviewed builds of
`7f7560e0d8ca94184068fa65f348998b67a1c1ff` and
`47f0f53b476c2d86281589026585391ef8260777`, with no previous publication and
zero public provider additions, changes or removals. The exact workflow build
step accepted the clean candidate `1488b88a1a7ad7d213b60761590a4584a5061a14`
and all 389 listed files. It rejected the incomplete C12.2 and C12.3 builds
before upload. The independent Astra xhigh re-review passed with no findings.
[Corrected gate](../../dist/c12-goal/c12.4-correction.log).

## First publication approval

C12.5 proposed one fast-forward push from the locally recorded
`origin/main` at `c5072d43a83e9601a0965ca770669fa8e29c19bb` to
`a4527c10002e77ef2e769c404c8d96d2b227b139` on `refs/heads/main` at
`git@github.com-tabilet:OpenUdon/apitools.git`. The [nine outgoing commits](../../dist/c12-goal/c12.5-proposed-outgoing.txt)
include five from C13 and four from C12. The owner directly approved this
exact push in response to the proposed launch-policy exception. The effective
exception `c12.5-first-push` is limited to this one `git.push` by the
coordinator, with exact-commit content and fast-forward mode. The underlying
launch file remains at its approved checksum.

Preflight matched the entire approved outgoing range, the local repository
and push URL, and the live `refs/heads/main` at `c5072d43a83e9601a0965ca770669fa8e29c19bb`.
[Remote preflight](../../dist/c12-goal/c12.5-push-before.txt).

### Owner-performed changes before the pending push

The owner reported these actions through the Claude pane at their request.
They occurred ahead of the planned C12.5 order and consumed no goal-run
approval or deployment attempt.

- The owner enabled GitHub Pages for `OpenUdon/apitools` with GitHub Actions
  as the source and `apis.kinet.com` as the custom domain. Pages reported
  “DNS check successful.” This supersedes the earlier disabled state for
  current readiness. The C12.1 attestation remains historical evidence.
- The owner added a DNS-only CNAME from `apis.kinet.com` to
  `openudon.github.io` in Cloudflare with TTL Auto. They also added the
  `_github-pages-challenge-openudon.apis` TXT challenge. Certificate issuance
  and HTTPS enforcement are pending.
- C12.5(b) was performed by the owner. In a later report relayed through
  the Claude pane at their request, the owner attested that the OpenUdon
  organization Pages settings showed “Successfully verified apis.kinet.com”
  and listed `apis.kinet.com` as **Verified** under **Verified domains**.
  This completes the domain-verification step as owner-performed evidence.
- At the owner's request the Claude pane reported a read-only
  `git ls-remote origin refs/heads/main` result of
  `c5072d43a83e9601a0965ca770669fa8e29c19bb`. The owner-performed Pages
  setup produced no commit on `main`. The proposed push is still a
  fast-forward and has the same nine outgoing commits.

The owner's relay is evidence only. In a subsequent direct message, the
owner approved deployment and HTTPS enforcement. This applies to C12.5(e),
dispatching `Publish reviewed catalog` on `main` for release commit
`a4527c10002e77ef2e769c404c8d96d2b227b139`, and C12.5(f), enabling Enforce
HTTPS for `apis.kinet.com` after GitHub issues the certificate. The approvals
remain pending execution. The approved launch's external-mutation policy
is still `none`.

C12.5(a) completed on 2026-11-26 at 02:09:59Z. The exact non-force push
reported `c5072d4..a4527c1`. A separate `git ls-remote` read confirmed
`refs/heads/main` at `a4527c10002e77ef2e769c404c8d96d2b227b139` and the
local tracking ref agrees. The one-push exception is consumed.
[Push transcript](../../dist/c12-goal/c12.5-push.txt) and
[independent remote read](../../dist/c12-goal/c12.5-push-after.txt).

The owner replied “yes” directly to the proposal for a scoped launch-policy
exception covering the two previously approved operations. This permits one
manual dispatch of `Publish reviewed catalog` on `OpenUdon/apitools` `main`
for `a4527c10002e77ef2e769c404c8d96d2b227b139`, monitoring and required live
checks, then setting Enforce HTTPS for `apis.kinet.com` after certificate
issuance. The checksum-approved launch bytes remain unchanged.

The C12.5 preflight read GitHub `main` at the exact release commit, workflow
ID `431090433` active, and Pages public with workflow build type and the
approved CNAME. GitHub initially reported no certificate and HTTPS
enforcement disabled. The one supervised Class R deployment attempt
`c12.5-deploy-a1` is reserved before the workflow dispatch POST. It uses
the local release artifact proof below, the `main` ref and `source_commit`
set to the release commit. Monitor the returned run without redispatch.

The dispatch returned [GitHub Actions run 45740811218](https://github.com/OpenUdon/apitools/actions/runs/45740811218).
GitHub identifies attempt 1 as `workflow_dispatch` on `main` at the exact
release commit. Build and deployment succeeded at 02:13:01Z on 2026-11-26.
The dispatch grant is consumed. The public site's HTTP manifest matched the
local release manifest byte for byte, and all 389 HTTP responses were 200 and
matched their declared SHA-256 digests. The live catalog identified the
exact release commit. [Run receipt](../../dist/c12-goal/c12.5-run-45740811218.json)
and [live checksum receipt](../../dist/c12-goal/c12.5-http-checksums.txt).

HTTPS verification remains part of the same attempt. After deployment GitHub
still reported a null certificate and disabled enforcement. An HTTPS probe
failed hostname validation while the public HTTP manifest was available.
The CNAME resolves to `openudon.github.io`. Wait for issuance before using
the one-time enforcement grant. At 02:16:51Z GitHub still had no certificate,
and the live TLS endpoint served a `*.github.io` certificate without
`apis.kinet.com` in its subject alternative names. C12.5 is blocked on GitHub
certificate issuance. The approved enforcement remains unused. Continue the
live checks of this same deployment once the certificate is available.
[Pages pre-enforcement receipt](../../dist/c12-goal/c12.5-pages-before-https.json).

At the owner's request, the 2026-11-26 06:05:40Z recheck confirmed that the
domain remains verified and GitHub still reports a null certificate and
disabled HTTPS enforcement. The live endpoint still serves `*.github.io`.
The CNAME still points to `openudon.github.io`, whose CAA permits Let's
Encrypt. No CAA answer appeared at `kinet.com`. The public HTTP manifest
still matched the release SHA-256. The enforcement grant remains unused.

The next owner-requested recheck at 2026-11-26 13:51:32Z returned the same
verified domain, null certificate and disabled enforcement. TLS still failed
hostname validation. The CNAME and public HTTP manifest matched the approved
deployment. C12.5 still awaits certificate issuance.

The 2026-11-26 21:19:08Z recheck still showed no GitHub certificate and a
live hostname validation failure. GitHub Pages DNS health reported the
domain valid, served by Pages and eligible for HTTPS, with no CAA error.
Its HTTPS result was `peer_failed_verification`. The live HTTP manifest
again matched the release. The enforcement grant remains unused.

The owner-requested 2026-11-27 01:10:08Z recheck returned the same null
certificate, disabled enforcement, hostname validation failure and
`peer_failed_verification`. Pages still reported HTTPS eligibility with
no CAA error. The public HTTP manifest matched the release again.

The owner-requested 2026-11-27 03:20:35Z recheck, about 25 hours after
deployment, still showed no certificate and failed live TLS hostname
validation. The public HTTP manifest matched the release. GitHub's
[certificate troubleshooting guidance](https://docs.github.com/en/pages/getting-started-with-github-pages/securing-your-github-pages-site-with-https#troubleshooting-certificate-provisioning-certificate-not-yet-created-error)
documents removing and re-adding the custom domain to restart provisioning.
This would require separate authority to change Pages settings and could
briefly interrupt the custom-domain site. The enforcement grant remains unused.

At 2026-11-27 08:34:30Z, about 30 hours after deployment, GitHub's Pages
API still showed the verified domain with no certificate and disabled
enforcement. Live TLS still failed hostname validation and the HTTP
manifest still matched the approved release.

The owner replied “yes” directly to the proposal to remove and immediately
re-add `apis.kinet.com` to restart certificate provisioning, including the
disclosed possibility of a short interruption. The scoped exception covers
one paired GitHub Pages custom-domain reset on `OpenUdon/apitools` only.
Preflight at 08:43:06Z still found the verified custom domain on the public
workflow-backed site, with no certificate and HTTPS enforcement disabled.
The reset changes only the Pages `cname` field, from `apis.kinet.com` to null
and back, and is reserved as `c12.5-domain-reset-a1`. Verification continues
on the original deployment run and the existing HTTPS enforcement grant.

GitHub acknowledged removal and restoration, and the Pages API confirmed
the original verified custom domain and public workflow source at
2026-11-27 08:43:39Z. The HTTP manifest still matches the release. The
one-time custom-domain reset is consumed. No workflow was redispatched.
[Reset receipt](../../dist/c12-goal/c12.5-domain-reset.log).

At 2026-11-27 08:44:26Z the Pages API reported an approved certificate for
`apis.kinet.com` expiring 2027-02-24. The live HTTPS manifest request
returned 200 with successful certificate validation. Reserve the owner's
existing one-time C12.5(f) grant before setting `https_enforced` to `true`
on this repository's Pages site.

GitHub confirmed HTTPS enforcement enabled at 08:44:41Z. The root page and
`/catalog.json` redirect from HTTP with 301 to their exact HTTPS URLs.
OpenSSL validated the Let's Encrypt certificate for `apis.kinet.com`.
The HTTPS manifest is byte-identical to the retained release manifest,
SHA-256 `cc504342dde19d78b6a56023a6a0e4e4ee3b1608638d1356a3af9cf0c6fb498e`.
All 389 files matched that manifest over certificate-verified HTTPS. The
first fetch of `providers/gotify/index.html` returned 503, and a single
targeted GET then returned 200 with matching bytes. The live catalog
identifies `a4527c10002e77ef2e769c404c8d96d2b227b139` and 334 providers.
Class R deployment attempt `c12.5-deploy-a1` passed on the original run.
The HTTPS enforcement grant is consumed.
[Enforcement](../../dist/c12-goal/c12.5-enforce-https.log),
[redirect and certificate](../../dist/c12-goal/c12.5-http-redirect.txt),
[HTTPS file checks](../../dist/c12-goal/c12.5-https-checksums.txt) and
[targeted retry](../../dist/c12-goal/c12.5-https-gotify-retry.txt).

The local prepublication suite passed. The corrected row build and vet cover
the identical implementation at the proposed head. The whole
`go test -count=1 ./...` with Chromium, generator freshness, catalog quality,
search/import help and diff checks passed in about 53 seconds.
[Whole offline log](../../dist/c12-goal/c12-local-full-offline.log).
Two local generations from a clean checkout of the exact proposed head were
byte-identical. All 389 entries in `SHA256SUMS` verified, and the manifest
SHA-256 is
`cc504342dde19d78b6a56023a6a0e4e4ee3b1608638d1356a3af9cf0c6fb498e`.
[Local artifact verification](../../dist/c12-goal/release-local.log).

The C12.5 offline gate was run fresh against the exact release implementation
using Go 1.26.6 and Chromium. `go build ./...`, `go vet ./...`,
`go test -count=1 ./...`, generator freshness, catalog quality, search/import
help and the full diff check passed in about 57 seconds. This proof is reused
because the source, tools and offline inputs are unchanged. Live verification
has now passed.
[C12.5 offline gate](../../dist/c12-goal/c12.5-offline.log).

C12.5 reconciled the horizon register against the frozen C06–C10 and C13
records. FND-01, FND-02, FND-03, FND-05, FND-06, FND-07 and FND-08 have
accepted fixes in their recorded owners. The C13 record retains the direct
owner deferral of FND-04. The horizon closure will prune the fixed rows and
keep FND-04 and last-used ID FND-08.

Independent GPT-6.2 Astra xhigh task review passed with no findings. It
verified the exact pre-closure ledger diff over the deployed source commit,
the 389 retained HTTPS file hashes, the TLS and redirect receipts, one
dispatch, the owner-scoped reset and enforcement, and the finding outcomes.
[Task-review diff](../../dist/c12-goal/c12.5-task-review.diff) SHA-256
`f38e11fc6d015ecee72a13c07639d0c14fb280ac1e01b45c5d706817c6698169`.

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Review and model roles follow
[Monthly reviewed catalog](milestone.md#monthly-reviewed-catalog). Use the
[test policy](test-policy.md), [verification-attempt gate](milestone.md#verification-attempt-gate)
and [milestone review procedure](milestone.md#milestone-review-procedure).
Required-gate records include elapsed wall-clock time. Persist
[required gate recovery](milestone.md#required-gate-recovery) state here before
each interactive Goal/Next required-gate dispatch.
