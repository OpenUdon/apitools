# Prompt V10 - Residual Google Discovery Coverage

After M47, close the remaining high-value Google Discovery source coverage gaps
identified from n8n-visible providers.

Add native Google Discovery catalog coverage for:

- Google Contacts / People API.
- Google Perspective / Perspective Comment Analyzer API.
- Google Firebase Realtime Database Management API where Discovery covers
  management resources.

Record Google Ads as a reviewed candidate-only decision if no official
Discovery, OpenAPI, or Smithy source exists. Keep Discovery native, preserve
OAuth/security metadata as catalog metadata only, and do not add API execution,
credential resolution, token fetching, project/account selection, or protobuf
source-family support.
