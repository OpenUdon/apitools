# Status M01 - Private Harness Bootstrap

Status markers are `[ ]` pending, `[+]` completed, `[~]` in progress, `[!]`
blocked, and `[X]` cancelled. Table markers remain wrapped in backticks.

| Item | State | Notes |
|---|---|---|
| Canonical harness directory created under `../tofu/apitools` | `[+]` | Initial private harness files added. |
| `../apitools/AGENTS.md` symlinked to private harness | `[+]` | Public tracked file removed from index; local symlink points to `../tofu/apitools/AGENTS.md`. |
| `../apitools/memory-bank` symlinked to private harness | `[+]` | New local symlink points to `../tofu/apitools/memory-bank`. |
| `../apitools/evolution` symlinked to private harness | `[+]` | New local symlink points to `../tofu/apitools/evolution`. |
| Public repo ignores harness paths | `[+]` | Ignore entries cover `AGENTS.md`, `memory-bank`, and `evolution` symlink paths. |
| Verification completed | `[+]` | `go test ./...`, `go vet ./...`, `git diff --check`, `go run ./cmd/apitools search --help`, and `go run ./cmd/apitools import --help` passed in `../apitools`. |
