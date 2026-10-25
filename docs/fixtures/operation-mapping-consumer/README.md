# Ordinary combined producer consumer

Copy .txt inputs to unsuffixed names in a disposable standalone module. Use
cached Go1.26.6, GOWORK=off, GOPROXY=off, GOSUMDB=off, GOAUTH=off,
GOTOOLCHAIN=local and GOFLAGS=-mod=readonly. Set M91_MODULE_DIR to acquired
APItools Dir solely to read exact fixtures, never to replace imports. Run
go test -json ./... -count=1, go vet ./... and go build -o probe .
Pass the probe APItools testdata/operation-shapes/m83/eight-families.json and
OpenUdon docs/fixtures/authoring-card-v1/complete.json. It calls both metadata
helpers on13operations/eightfamilies and typed card Parse/RenderText, returning
BuildInfo. Missing required mappings may refuse as incompatible; incomplete
eligibility stays indeterminate. Tests include19actual source positive/refusal
cases. No credential, admission, approval or operation authority follows. Exact
receipts, graph and retained failed checks: ../../m91-published-proof.json.
