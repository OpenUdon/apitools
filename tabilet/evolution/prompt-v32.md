# v32 — Public catalog website at apis.kinet.com

On 2026-11-25, the owner requested a public website for the APItools catalog
alongside the monthly updates and approved the revised `memory-bank-propose`
plan at 2026-11-25T20:00:04Z, when `HEAD` was `8e76fba`.
`468fef327154cf248fb53ab5e0b765898ce46067` was committed at 20:00:58Z and is
where the planning edits were written. The owner approved corrections to
those edits later the same day. The edits were then rebased onto
`513c6acfe3774cc8dcd915104fc3dd9458139617`, after C08 retired and the workflow
upgrade, with the resulting changes approved the same day. Recorded
assumption: the audience is developers and agents choosing an API.

The owner selected a static, independent site built only from committed,
reviewed catalog data for every built-in provider: provider and category pages
with coverage, provenance and verification dates, in-browser keyword search,
browsing that works without JavaScript, committed monthly change records, and
versioned JSON in APItools' own small format. Spec or overlay bytes, auth
guidance while [FND-05](../memory-bank/milestone.md#open-review-findings)
remains open, hosted semantic search, accounts, user submissions, analytics,
logos, third-party assets, and `"current"` or `"latest"` claims are excluded.

The owner selected GitHub Pages with custom domain `apis.kinet.com`, deployed
from a manual workflow artifact after the monthly review so no generated pages
are committed to any branch and the monthly job cannot publish. The in-page
CSP limitation (Pages cannot block framing via response headers) was disclosed
to the owner before host selection. C12.1 owns the narrower publication
sign-off, takedown contact and read-only DNS/Pages readiness check; C12.5 owns
separately approved going-live steps and horizon closure of the open review
findings register, including FND-04, FND-05, FND-07 and FND-08, while findings
owned by C09 rows stay with C09. C09's rows and its integrated C06–C09/C10
review are unchanged; only the register's horizon closure moves from C09 to C12. The
broader demand-grown directory, APIs.guru `list.json` compatibility,
unreviewed seeds, crawling and user submissions stay in
[Candidate Directions](../memory-bank/milestone.md#candidate-directions).
Planning only: no implementation, launch, push, GitHub settings, DNS or
deployment is authorized.
