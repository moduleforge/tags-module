# Phase 1, Task 4 — Update top-level go.work

## Context

The root `user-components/go.work` stitches all Go modules together for local dev. Adding tags-module requires two `use` entries and a `replace` directive so core-api's import of `github.com/moduleforge/tags-model` (which will land in Phase 3) resolves locally.

## Acceptance

- `user-components/go.work` includes `./tags-module/api` and `./tags-module/model` in the `use (…)` block.
- `user-components/go.work` adds a `replace github.com/moduleforge/tags-model v0.0.0 => ./tags-module/model` directive alongside the existing core-model replace.
- `go work sync` at repo root exits 0.
- `cd tags-module/api && go mod tidy` produces a clean go.sum that pulls core-api, core-model, tags-model, pgx, chi, and uuid.

## How to verify

- `go work sync` at repo root — no errors.
- `cd tags-module/api && go build ./...` exits 0.
- `cd core-module/api && go build ./...` still exits 0 (we didn't break it).
- `cd users-module/api && go build ./...` still exits 0 (we haven't touched users-module yet).
