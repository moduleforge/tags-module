# Phase 5, Task 1 — Extend users-module/model compose pipeline

## Context

`users-module/model` composes migrations from core + own files into a single flat `schema/migrations/` dir before running atlas/sqlc. Extend to pull tags-module migrations too.

## Acceptance

- `users-module/model/Makefile` `compose` target copies:
  - `core-module/model/migrations/*.sql`
  - `users-module/model/migrations/*.sql` (the 0100–01NN files)
  - `tags-module/model/migrations/*.sql` (new)
  - Same for queries under `users-module/model/queries` if that's how the pipeline works — follow existing convention. If users-module/model doesn't compose queries (only migrations), leave queries alone.
- `users-module/model/schema/migrations/` (gitignored) rebuilt cleanly.
- `make migrate.hash` on the composed dir updates cleanly.
- sqlc regeneration in users-module/model includes the tags tables if it was configured to — verify whether users-module/model needs its own Querier over tags. If users-module/api only calls tags-api's service interface (which handles its own DB access via tags-model.Querier), then users-module/model need NOT compose tags queries; only the migrations are needed for atlas to build the full schema.

## How to verify

- `cd users-module/model && make build` exits 0.
- `atlas migrate validate` on the composed dir succeeds.
- Full migration list in order: 0000–0011 (core), 0100–01NN (users), 0200–0201 (tags).

## Reference

- `users-module/model/Makefile` — existing compose target.
- Existing compose pipeline for core-module is the model to extend.
