# M82 source publication evidence

The confirmed STG11_SOURCE_PUBLICATION grant published reviewed M82 source and
retirement closure through normal fast-forward pushes to the exact APItools
origin/main. No tag, force push, branch deletion or deployment was performed.

- Accepted implementation: `54583f9b2f452b7cc522360c5aeeff29ca22f96c`, whole review 3/10.
- Resolved module: `v0.0.0-20261006210844-54583f9b2f45`; configured registry origin hash matches accepted implementation.
- Independently observed retirement closure: `fae9982e42d6b16fe7a5ebfd342a016613a62adb`.
- `git ls-remote origin refs/heads/main` returned that exact closure after push.
- Accepted implementation is an ancestor of the observed closure.
- Exact destination: `git@github.com-tabilet:OpenUdon/apitools.git`, `refs/heads/main`.

Consumers pin accepted source and retain observed closure/publication evidence.
Producer metadata remains advisory until independently reproduced against exact
source bytes, and never supplies credential readiness, approval or execution.
[Qualification](m82-qualification.md), [contract](operation-shapes.md) and the
[frozen retirement record](../tabilet/docs/history/status-M82.md) preserve scope,
verification and resolved findings. Later documentation-only descendants do not
change the accepted code. Kinet/OpenUdon/Udon adoption remains package-local work.
