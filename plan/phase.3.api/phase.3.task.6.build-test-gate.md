# Phase 3, Task 6 — Build/test gate

## Context

Final integration sanity for the api package before moving to gui / wiring.

## Acceptance

- `cd tags-module/api && make build` exits 0.
- `cd tags-module/api && make test` exits 0 with coverage numbers reported at expected levels.
- `go vet ./...` exits 0 in tags-module/api.
- Top-level `go work sync` still exits 0.

No new files; this is a gate task.
