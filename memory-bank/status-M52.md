# Status M52 - n8n De-Emphasis And License-Risk Reduction

| Item | State | Notes |
|---|---|---|
| n8n report API removed | `[+]` | Removed the exported catalog n8n gap-report API, tests, CLI command, README usage, and tech-stack smoke command. |
| n8n provider rows removed | `[+]` | Removed n8n Public API and TimeSaved candidate/provider/auth metadata, alias expectations, and artifact registration. |
| Candidate evidence cleaned | `[+]` | Removed built-in n8n node-directory priority evidence and local `try-n8n` fixture seeding from candidates. Built-in candidate evidence is now source-first. |
| Public docs updated | `[+]` | Rewrote catalog queue and non-OpenAPI protocol notes to avoid n8n-specific provider prioritization and to require provider-owned or protocol-owned source artifacts. |
| Active memory superseded | `[+]` | Added this milestone and updated active product, architecture, tech-stack, and milestone notes so M49-M51 are historical only. |
| Verification | `[+]` | Passed full Go, catalog, consumer, diff, and text-audit checks. |

## Boundary Notes

- Do not copy third-party integration code, generated node metadata,
  credential behavior, descriptions, icons, or runtime workflow behavior.
- Runtime connector catalogs may suggest user demand, but they are not provider
  truth and must not drive catalog rows without provider-owned source evidence.
- Historical memory-bank and evolution records may mention n8n as chronology;
  active public docs, code, and release guidance should not depend on it.
