# C12 — Public catalog website

**Owner.** APItools. **Dependency.** Accepted [C13](../docs/history/status-C13.md). **Parallel-safe: no.**
C11 is skipped because the [M86](status-M86.md) records use an unqualified `C11` for `UWS:C11`. C12 needs its own launch after C13. See the [milestone specification](milestone.md#c12--public-catalog-website) for scope, exclusions, acceptance, class O/R verification and `AUTHORIZATION_REQUIREMENTS`.

**Closing duty.** C12 is the last active catalog milestone. C12.5 closes the [open review findings](milestone.md#open-review-findings) register under its existing rules and verifies the accepted C13 dispositions for FND-05, FND-07 and FND-08 and the owner-approved FND-04 deferral.

## Tasks

| Item | State | Notes |
| --- | --- | --- |
| C12.1 — Sign-off and readiness | `[+]` | Before sign-off, resolve how operation labels such as Jenkins "Get CSRF crumb" and FileMaker "Delete session" are published; these examples come from C08 whole review round 1 L-2 ([C08 record](../docs/history/status-C08.md#status-record)). Record the owner's decision on publishing provider names, descriptions, operation labels, license notes and links, plus a correction/takedown contact. Record hostname `apis.kinet.com`, its DNS owner and the GitHub Pages host. Under explicit authority, read current DNS and host settings without changing them, and confirm read-only that the hostname is not already in use and that nothing found (such as a CAA record) would prevent a certificate; the certificate itself is proved in C12.5. Class R check under the [verification-attempt gate](milestone.md#verification-attempt-gate); keep DNS/host credentials out of the repository and registrar account details out of planning records. Blocks C12.2–C12.5. |
| C12.2 — Data and generator | `[ ]` | Add maintainer-only `cmd/catalogsite` using the Go standard library, with no new public Go API. Build deterministically from an identified clean commit of reviewed catalog data with a pinned Go toolchain; the generator makes no network request, while a hosted runner may download pinned modules verified by `go.sum`. Cover every built-in provider at the built commit (334 at planning); emit no spec or overlay bytes, no auth guidance, and no claim that a provider or version is current or latest. Take the base URL as an explicit build input and use relative internal links. Emit versioned JSON in APItools' own small format with a documented schema version and the source commit, and publish a SHA-256 manifest of the output. |
| C12.3 — Pages and search | `[ ]` | Build provider and category pages with coverage, provenance and verification dates, plus in-browser keyword search by name, alias, category and keyword, with every provider and category page reachable without JavaScript and no third-party assets, logos, analytics, accounts or user submissions. Search matches words in public profile text and must not present an incidental operation-label match as a supported capability. Follow the evidence rules accepted in [C08](../docs/history/status-C08.md#status-record), whose whole review round 1 P1-1 closed in round 2. C13.1 delivered the CLI search correction for [FND-07](milestone.md#open-review-findings). Verify the website's own public-profile search behavior. Escape all catalog text and allow only HTTPS source links. Add a restrictive in-page CSP (Pages cannot set response headers to block framing) and accessible markup with named Go-test checks for document language, title, landmarks, form labels and link text. |
| C12.4 — Monthly changes and publishing | `[ ]` | Compare two reviewed revisions and commit a small monthly record; the monthly changes page matches it, including the first publication with no predecessor. Add a manual-only GitHub Pages artifact workflow with every GitHub Action pinned by commit SHA and least permissions, committing no generated pages to any branch and keeping the monthly job unable to publish. Add `apitools-catalog-maintenance` skill steps for preview, publish and rollback, and exercise rollback on the loopback preview. Document domain verification, DNS rollback and removal: verify the domain before adding DNS, and remove DNS before retiring the site. |
| C12.5 — First publication and closure | `[ ]` | Ordered going-live steps, each needing separate explicit approval: (a) push the workflow to `main`, which publishes every unpushed local commit, so list them for the owner first; (b) verify the domain with GitHub; (c) enable Pages with GitHub Actions as source and set the custom domain; (d) add the DNS record; (e) dispatch the deployment; (f) enforce HTTPS once the certificate issues. Each deployment is a counted class R attempt under the [verification-attempt gate](milestone.md#verification-attempt-gate) and its existing limit. Check that live files match the published SHA-256 manifest for the recorded commit, that `apis.kinet.com` serves a valid certificate, and that HTTP redirects to HTTPS. Own horizon closure of the open review findings register. Verify the accepted C13 dispositions for FND-05, FND-07 and FND-08 and the owner-approved FND-04 deferral. Run the required class O whole-unit gates before whole review. |

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

## Review

Review gate iterations: 0

Whole review pending, default limit 10. Review and model roles follow
[Monthly reviewed catalog](milestone.md#monthly-reviewed-catalog). Use the
[test policy](test-policy.md), [verification-attempt gate](milestone.md#verification-attempt-gate)
and [milestone review procedure](milestone.md#milestone-review-procedure).
Required-gate records include elapsed wall-clock time. Persist
[required gate recovery](milestone.md#required-gate-recovery) state here before
each interactive Goal/Next required-gate dispatch.
