# v32 result — Public catalog website plan recorded

The serial APItools catalog horizon is now **C09 → C12**, adding five
pending rows in [C12](../memory-bank/status-C12.md) after accepted
[C09](../memory-bank/status-C09.md). `C11` is skipped because
[M86](../memory-bank/status-M86.md) uses an unqualified `C11` for `UWS:C11`.

| Owner | Planned result | Rows |
| --- | --- | ---: |
| [C12](../memory-bank/status-C12.md) | Publication sign-off and read-only `apis.kinet.com` DNS/Pages readiness, reproducible `cmd/catalogsite` generator with versioned JSON and SHA-256 manifest, static pages/search with in-page CSP and accessibility checks, committed monthly change records and a pinned manual Pages workflow/skill steps, and separately approved ordered going-live steps with horizon closure | 5 |

[milestone.md](../memory-bank/milestone.md#c12--public-catalog-website) records
the C12 specification, acceptance criteria, class O/R verification and
`AUTHORIZATION_REQUIREMENTS`, and narrows the
`Published provider version directory` candidate to the remaining broader
directory, unreviewed seeds, crawling, APIs.guru compatibility and user
submissions while keeping its rationale and gates. C12 becomes the last active
milestone, so C12.5 owns horizon closure of the open review findings register,
including FND-04, FND-05, FND-07 and FND-08; findings owned by C09 rows stay
with C09. C09's rows and its integrated C06–C09/C10 review are unchanged; only
the register's horizon closure moves from C09 to C12. Product/architecture/tech-stack records, Go contracts and consumer pins
are unchanged. C12 is outside the C09 launch reference, starts at review 0/10
and needs its own launch after C09.

Only planning records were written at `468fef327154cf248fb53ab5e0b765898ce46067`,
including the owner-approved corrections made later the same day. They were
rebased onto `513c6acfe3774cc8dcd915104fc3dd9458139617`, after C08 retired and the
workflow upgrade, with the resulting changes approved the same day.
Implementation, class O verification, C12.1's class R DNS/host inspection,
GitHub Pages or DNS changes, push, and C12.5's class R deployments require
later execution and explicit authority under the
[test policy](../memory-bank/test-policy.md).
