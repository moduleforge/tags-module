# Phase 2, Task 6 — Compile-check model package

## Context

Final smoke after generation: tags-module/model should build clean and the generated code should compile against the workspace.

## Acceptance

- `cd tags-module/model && go build ./...` exits 0.
- `go vet ./...` exits 0.
- `make test` runs (may be a no-op if no _test.go files exist — that is acceptable for this phase).
- `go work sync` at repo root still exits 0.
