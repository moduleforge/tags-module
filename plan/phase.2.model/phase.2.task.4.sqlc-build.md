# Phase 2, Task 4 — Generate sqlc db/ package

## Context

Run sqlc to produce the typed Querier + models from the migrations + queries written in tasks 2.1–2.3.

## Acceptance

- `cd tags-module/model && make build` succeeds.
- Output files present under `tags-module/model/db/`:
  - `db.go` (base DBTX / Queries struct)
  - `models.go` (the `Tag` struct)
  - `tags.sql.go` (method implementations)
  - `querier.go` (Querier interface if `emit_interface: true`)
- `Tag` model includes all tag columns with correct Go types (`pgtype.Text` for nullable `color`, `time.Time` for timestamps, etc.) and JSON tags matching the column names.
- Generated code committed to git (per project convention — check core-module/model to confirm `db/` is not gitignored).

## How to verify

- `go build ./...` exits 0 in tags-module/model.
- `go vet ./...` exits 0.
- Test a trivial instantiation: `var q coredb.Querier = &coredb.Queries{}` compiles in a scratch file (optional).

## Reference

- `core-module/model/db/` — output shape to match.
