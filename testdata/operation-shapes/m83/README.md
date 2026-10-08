# M83 corrected projection fixture

This successor fixture preserves the frozen M82 v1 fixture and all original
source/candidate/schema/grammar bytes. It uses the five corrected source files
under `sources/` plus the unchanged OpenRPC, gRPC and OData source bytes from
`../../operation-candidates/v1/sources/`. Source IDs remain the eight kinds.

The same `uws.shape-table.v1` wire describes corrected derived metadata: AsyncAPI
2 directions, optional GraphQL defaults and distinct response aliases, Smithy
literal omission, normal Discovery endpoint/optional request and numeric YAML
response/extension/TRACE/unknown int32 evidence. Native security remains symbolic.
This changes neither authority nor runtime support. Consumers require fresh
assessment/confirmation/grants for new derived table/package/worker identities.

`TestOperationShapesRemediationGolden` checks exact bytes, raw SHA-256, native
selectors, schema/security limitations, source-order independence and independent
reproduction. Separate M83 regressions verify equivalent versions and refusals.
Table size is 9,615 bytes; SHA-256 is `cfe13bcd44f62b0b650a20833984231426bd0a83cee91538b96d5332d1107fad`. The whole local code review passed at iteration 4; named publication and ordinary consumer proof remain separate.
