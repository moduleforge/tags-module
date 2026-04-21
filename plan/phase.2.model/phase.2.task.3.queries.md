# Phase 2, Task 3 — sqlc queries for tags

## Context

Define the sqlc queries the service layer needs. Each query corresponds to an operation in the Phase 3 TagService.

## Acceptance

Create `tags-module/model/queries/tags.sql` with the following named queries (sqlc annotations shown; follow the exact format used in `core-module/model/queries/service_accounts.sql`):

- **CreateTag** `:one` — INSERT INTO tags, returns all columns.
- **GetTagByEntityID** `:one` — SELECT by `entity_id`.
- **GetTagByEntityUUID** `:one` — JOIN entities; select by `entities.uuid`; returns tag columns + the entity uuid.
- **UpdateTagColor** `:exec` — UPDATE color WHERE entity_id = $1.
- **DeleteTag** `:exec` — DELETE FROM tags WHERE entity_id = $1.
- **ListTagsBySubjectEntityID** `:many` — SELECT tags WHERE subject_id = $1, optional purpose filter via `sqlc.narg('purpose')`.
- **SearchTags** `:many` — SELECT tags with optional filters: `sqlc.narg('owner_id')`, `sqlc.narg('subject_id')`, `sqlc.narg('purpose')`, `sqlc.narg('value')`. Use the sqlc `COALESCE(sqlc.narg('x'), column)` pattern for optional equality filters (see how users-module handles similar optional filters if present; otherwise use `WHERE (owner_id = sqlc.narg('owner_id') OR sqlc.narg('owner_id') IS NULL)` repeated per filter).
- **CountTagsBySubjectEntityID** `:one` — COUNT(*) WHERE subject_id = $1 — optional, but useful for pagination later; include it.

All queries return columns via SELECT list, not `*`, for sqlc stability.

## How to verify

- `cd tags-module/model && make build` produces `db/tags.sql.go` without errors.
- Generated Querier has methods matching the names above.
- `go vet ./...` in tags-module/model exits 0.

## Reference

- `core-module/model/queries/natural_persons.sql` — most similar query surface (CRUD + narg usage if present).
- `core-module/model/queries/service_accounts.sql` — simpler reference.

## Open question during implementation

If `sqlc.narg` pattern isn't already used in this project, default to the `WHERE (col = $N OR $N IS NULL)` pattern and pass `pgtype.Int8`/`pgtype.Text` for the optional params. Document whichever approach is chosen in a one-line comment at the top of the file. If unsure, check how any existing query handles optional filters before inventing a new style.
