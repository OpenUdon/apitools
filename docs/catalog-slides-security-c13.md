# Google Slides advisory scope evidence

The `google-slides-discovery-auth-overlay` advisory covers five Google Slides
methods. It selects `presentations.readonly` for the three reads and
`presentations` for the two writes. Each rule uses the exact native method
and path. The advisory has no root-level scope, so its read-only scope
cannot carry over to a write or to another method.

| Discovery method | HTTP method and native path | Selected scope suffix |
| --- | --- | --- |
| `slides.presentations.get` | `GET /v1/presentations/{+presentationId}` | `presentations.readonly` |
| `slides.presentations.pages.get` | `GET /v1/presentations/{presentationId}/pages/{pageObjectId}` | `presentations.readonly` |
| `slides.presentations.pages.getThumbnail` | `GET /v1/presentations/{presentationId}/pages/{pageObjectId}/thumbnail` | `presentations.readonly` |
| `slides.presentations.create` | `POST /v1/presentations` | `presentations` |
| `slides.presentations.batchUpdate` | `POST /v1/presentations/{presentationId}:batchUpdate` | `presentations` |

The retained publisher snapshots list these scope alternatives under
`https://www.googleapis.com/auth/`.

| Methods | Advertised scope suffixes in both snapshots |
| --- | --- |
| The three reads | `drive`, `drive.file`, `drive.readonly`, `presentations`, `presentations.readonly` |
| `presentations.create` | `drive`, `drive.file`, `presentations` |
| `presentations.batchUpdate` | `drive`, `drive.file`, `drive.readonly`, `presentations`, `spreadsheets`, `spreadsheets.readonly` |

The `batchUpdate` list includes Drive and Sheets read-only alternatives.
The snapshots do not establish how those alternatives apply across request
types. This advisory selects the Slides write scope and makes no claim
about effective permissions for the other alternatives. The metadata
covers the selected methods and leaves the wider root-scope overlay class
unreviewed.

## Retained source identity

| Publisher SDK | Discovery revision | SHA-256 of `slides/v1/slides-api.json` |
| --- | --- | --- |
| `google.golang.org/api@v0.221.0` | `20240305` | `c56a87d562396319c1af8a2a2068ed50a31fe1e34437dd3ba593a57990d268d9` |
| `google.golang.org/api@v0.259.0` | `20250401` | `f45d36026cc10cbb790a83405dc774c5815d9287604bdfd335475fd1cae0b022` |

Both files matched the entries in the retained Google SDK module ZIPs.
The `v0.259.0` ZIP checksum in the module cache matches the Udon pin
`h1:90TaGVIxScrh1Vn/XI2426kRpBqHwWIzVBzJsVZ5XrQ=`.
The canonical catalog reference separately records revision `20260513`.
That revision was unavailable in the inspected cache, and the supervised
C07 check retained no downloadable artifact for this reference. The scope
evidence above is bound to the two SDK revisions. [Publisher SDK locator](https://github.com/googleapis/google-api-go-client/blob/v0.259.0/slides/v1/slides-api.json) · [C07 record](../tabilet/docs/history/status-C07.md#status-record)

To reproduce the offline proof, put the two retained copies in a directory
as `v0.221.0.json` and `v0.259.0.json`, then run the focused test.

```sh
APITOOLS_C13_SLIDES_DISCOVERY_DIR="$SLIDES_SOURCE_DIR" \
  go test -count=1 -run '^TestGoogleSlidesCatalogScopeEvidence$' .
```

This checks both digests, parses the native methods and verifies the
selected advisory paths and scopes.

`apitools catalog overlay-view google-slides` lists the five advisory rules.
The standalone view does not load the Discovery operations, so it displays
the method matches as unresolved. The retained-source check above verifies
each path against the native methods.
