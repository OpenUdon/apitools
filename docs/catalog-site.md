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

## Review and preview monthly changes

Build `catalog.json` from each clean reviewed revision and check both
`SHA256SUMS` files. Compare the two catalogs using the installed maintainer
binary. Use `-first-publication` once, when the site has no predecessor.

```sh
/tmp/catalogsite compare -from /tmp/previous-site/catalog.json -to /tmp/reviewed-site/catalog.json \
  -period YYYY-MM > /tmp/YYYY-MM.json
```

Review the proposal and commit it as `cmd/catalogsite/changes/YYYY-MM.json`.
It names the two reviewed commits and sorted provider IDs added, changed, or
removed in the public projection. The first record compares the reviewed C12.2
and C12.3 builds and records that the site has no previous publication. The
site renders these committed records. The maintenance cron job never invokes
the site generator or Pages workflow.

For preview, build the reviewed commit as above in a clean checkout and serve
its output on loopback, for example with `python3 -m http.server 8765 --bind
127.0.0.1 --directory /tmp/apitools-site`. Use the backup preview directory
for rollback. Verify its manifest, stop the current loopback server, and
serve the backup directory on the same local address. Check the source commit
in `/catalog.json` and the page content before closing the preview.

## Publish and roll back

The GitHub Pages workflow runs only on manual dispatch from `main`. It builds
and uploads an artifact from the full approved source commit on `main`'s
history. The workflow creates no commit. Each action below needs its own
approval. Keep the preceding DNS values and the last verified published
commit and manifest as rollback evidence.

1. List and approve the exact outgoing commits, then push to `main`.
2. Verify the domain under the GitHub organization’s Pages settings. Treat
   any DNS verification TXT record that GitHub requests as a separate DNS
   change. Retain domain verification while the domain is assigned to Pages.
3. Set the repository Pages source to GitHub Actions and configure the
   custom domain `apis.kinet.com`.
4. Approve and apply the routing record. For the subdomain GitHub documents a
   CNAME to the organization Pages host `openudon.github.io`. Confirm the
   exact existing values, target and Cloudflare proxy mode with the DNS
   owner before changing the record. [GitHub custom-domain setup](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/managing-a-custom-domain-for-your-github-pages-site).
5. Dispatch `Publish reviewed catalog` with the approved full commit SHA.
   Save the run and artifact manifest. Compare served bytes to that manifest
   and record the certificate and HTTP redirect checks in the deployment
   attempt. [GitHub Pages workflows](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages).
6. Once the certificate is issued, enable Enforce HTTPS under its separate
   approval and verify the redirect and certificate.

To roll back content, separately approve a dispatch for the last verified
published source commit. Compare the regenerated manifest with the retained
manifest for that commit, then validate the served files again.
If the first release needs rollback, restore or remove the site routing DNS
before disabling Pages or clearing the custom domain, under separate
approvals for each action. Do the same before retiring the site. Record the
DNS and Pages outcomes and verify that the removed route has expired. Leave
domain verification in place while the organization retains the domain.
[GitHub custom-domain verification](https://docs.github.com/en/pages/configuring-a-custom-domain-for-your-github-pages-site/verifying-your-custom-domain-for-github-pages).
