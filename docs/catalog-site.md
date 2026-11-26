# Build the APItools catalog site

`cmd/catalogsite` is a maintainer command. Build it with Go 1.26.6 from a
clean checkout of the reviewed source commit. The command reads the built-in
catalog and search profiles and makes no network requests. Use an output
directory outside the source checkout. Its parent must exist.

```sh
source_commit=$(git rev-parse HEAD)
GOWORK=off GOTOOLCHAIN=local go build -trimpath -buildvcs=true -o /tmp/catalogsite ./cmd/catalogsite
/tmp/catalogsite -source-commit "$source_commit" -base-url https://apis.kinet.com/ -out /tmp/apitools-site
(cd /tmp/apitools-site && sha256sum --check SHA256SUMS)
```

The `-out` directory must be new. `-source-commit` must match the binary's
clean Git build stamp. Compare builds using the same source commit, base URL,
and pinned toolchain. The generator uses no system time or maintenance cache.

## `catalog.json` schema version 1

The top-level keys are `schema_version`, `source_commit` (the full Git
commit), `base_url`, `contact_url`, and `providers`. Each provider appears
once in ID order. Its fields are `id`, `name`, `category`, `description`,
`aliases`, `keywords`, `coverage`, `url`, `sources`, and
`reviewed_operations`.

`coverage` is exactly one of `official-openapi` (including Swagger and
official OpenAPI indexes), `official-machine-readable`,
`docs-with-advisory-overlay`, and `docs-only`, in that precedence order.
An advisory can coexist with an official spec. Each `sources` entry gives
its `id`, `kind`, `authority`, HTTPS `url`, and, when recorded, `verified_at`
(UTC calendar date) and `license_note`. Verification dates refer to the
individual source references. Each `reviewed_operations` entry gives its
`label`, a `source_id` in the same provider, and its `evidence` class
(`pinned-source` or `advisory-overlay`). These labels cover selected
operations. The public projection omits the two labels excluded in C12.1.

`SHA256SUMS` lists SHA-256 hashes and relative paths for every generated file
other than the manifest itself. It is written after the listed files.

## Browse and test the site

The provider and category lists work without JavaScript. When scripts are
enabled, search matches whole words across the public provider profile. The
listings show one dated source for each category entry when available.
Provider pages list the full source set and selected operation labels.

For the browser gate, point `APITOOLS_SITE_TEST_CHROME` at the absolute path
of an installed Chromium executable. `go test -count=1 ./cmd/catalogsite`
runs the assertions with a local server and disposable browser profiles.
The test suite otherwise runs the projection, escaping and link checks
without launching a browser.
