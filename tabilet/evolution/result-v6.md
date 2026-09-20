# Result V6 - Google Discovery Converter Package

`apitools` now owns reusable Google Discovery conversion through
`github.com/OpenUdon/apitools/googlediscovery`.

- `Convert` and `ConvertMap` produce bounded OpenAPI 3.0.1-shaped metadata
  from Google Discovery JSON bytes or decoded maps.
- The converter preserves Discovery resources, methods, root/resource/method
  parameters, schemas, media upload request bodies, deterministic operation
  IDs, and Google OAuth scope metadata.
- Follow-up review hardening rejects duplicate OpenAPI path/method mappings and
  distinguishes raw simple media uploads from multipart uploads.
- Drive and Gmail fixtures cover conversion behavior in `apitools`.
- Catalog Discovery artifacts remain classified as `google-discovery` and
  `valid-structured` review evidence unless an explicit caller converts them.
- `../udon` delegates reusable conversion to `apitools` and keeps native
  document adaptation plus runtime-plan generation locally.
