# Synthetic catalog root

All content here was written for tests and may be copied and redistributed.
It contains no third-party specifications, credentials, live accounts or real
service endpoints. `.invalid` URLs are metadata only and are never fetched.

`catalog.json` describes two providers, including multiword lookup keys, sharing
one OpenAPI document. `registrations.json` binds both links to the document's
exact SHA-256 and byte count. The source has one read and one write operation.
There is no committed SQLite database or generated operation index.

From the APItools repository, prepare a **new** absolute disposable destination
whose parent already exists:

```bash
go run testdata/catalog-root/prepare.go /tmp/my-new-synthetic-catalog
```

The helper refuses an existing destination, copies only the synthetic data
and registers it using the existing SQLite cache API. Load `catalog.json` into
`catalog.Catalog` and pass it explicitly in tests; it is not added to the
built-in provider catalog. Configure `catalog.RootOptions.Directory` with the
new directory. `sqlitecache.ReadCatalogArtifacts` reads its registrations
without migration, pruning or access-time changes.

M81.4 will add the offline index command and its exact invocation here. This
fixture currently proves root/registration preparation, not discovery or
index acceptance. Delete only the disposable root you created after checking.
