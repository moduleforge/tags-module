# Phase 2 — Model

## Goal

Land the SQL schema for tags: register the `tag` concrete type, create the `tags` subtype table with all constraints/triggers, write the sqlc queries, and generate the `db/` package.

## Preconditions

- Phase 1 complete — tags-module/model scaffold exists and builds empty.
- Phase 1 established that migrations 0200–0299 are reserved for tags-module.

## Outputs

- `tags-module/model/migrations/0200_type_tag.sql`
- `tags-module/model/migrations/0201_tags.sql`
- `tags-module/model/queries/tags.sql`
- `tags-module/model/db/*.go` (sqlc output, committed per project convention)
- `tags-module/model/migrations/atlas.sum`

## Hard rules

- Schema file must assume the core-module schema is already applied. Specifically: the `types` table, the `entities` table, the `set_updated_at()` function, and the `type_is_or_descends_from()` helper must all exist. They come from core-module migrations 0000–0011.
- Follow core's trigger patterns exactly (see `core-module/model/migrations/0008_legal_entities.sql` for the type-descent check pattern; see `0007_entities.sql` for `entities_fundamental_type_immutable` pattern for immutability triggers; see `0000_helpers.sql` for `set_updated_at`).
- **No Postgres-specific syntax beyond what core already uses.** `SIMILAR TO` for the color CHECK is acceptable (SQL:1999 standard). Regex `~` is not.
- UUIDs are on `entities`, not on `tags`.

## Schema decisions (recap from plan summary)

- `tags` columns:
  - `entity_id BIGINT PRIMARY KEY REFERENCES entities(id) ON DELETE RESTRICT`
  - `owner_id BIGINT NOT NULL REFERENCES entities(id)`
  - `subject_id BIGINT NOT NULL REFERENCES entities(id)`
  - `purpose TEXT NOT NULL CHECK (char_length(purpose) <= 512)`
  - `value TEXT NOT NULL CHECK (char_length(value) <= 512)`
  - `color TEXT CHECK (color SIMILAR TO '#[0-9A-Fa-f]{8}')`
  - `created_at TIMESTAMPTZ NOT NULL DEFAULT now()`
  - `updated_at TIMESTAMPTZ NOT NULL DEFAULT now()`
- Indexes:
  - `UNIQUE (owner_id, subject_id, purpose)` — the main uniqueness constraint.
  - `INDEX (subject_id)` — speeds `/entities/{uuid}/tags` lookups.
  - `INDEX (owner_id)` — speeds owner-scoped search (the UNIQUE index leading with owner_id also helps but a standalone is clearer for query planner hints; keep it if it appears costly later).
- Triggers:
  - `tags_type_check` BEFORE INSERT — entity's fundamental type must descend from `tag`.
  - `tags_immutable_fields` BEFORE UPDATE — reject any DISTINCT FROM on `owner_id`, `subject_id`, `purpose`, `value`.
  - `tags_set_updated_at` BEFORE UPDATE — reuses the existing `set_updated_at()` function from core.

## Tasks

- 2.1 Write migration 0200 (register `tag` type)
- 2.2 Write migration 0201 (`tags` table + triggers + indexes)
- 2.3 Write `queries/tags.sql`
- 2.4 Regenerate sqlc (`make build`)
- 2.5 Hash migrations (`make migrate.hash`)
- 2.6 Compile-check model

## How to verify

- `cd tags-module/model && make build` exits 0 and produces `db/tags.sql.go` (and friends) with the expected query functions.
- `atlas migrate validate --dir file://migrations` exits 0.
- `make migrate.hash` updates `atlas.sum` cleanly.
