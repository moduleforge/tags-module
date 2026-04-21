# Phase 2, Task 1 — Migration 0200: register `tag` type

## Context

tags-module registers its own concrete type `tag` under `entity` via its own migration, per core-module's entity-typing pattern. No core files are edited.

## Acceptance

- `tags-module/model/migrations/0200_type_tag.sql` exists.
- Contents model up `core-module/model/migrations/0006_type_service_account.sql` exactly, with these substitutions:
  - slug: `'tag'`
  - name: `'Tag'`
  - description: `'A labelled annotation applied to another entity.'`
- Includes the guard `DO $$ ... IF NOT EXISTS (SELECT 1 FROM types WHERE slug = 'entity') THEN RAISE EXCEPTION ...` preamble.
- Parent is `entity` (not `legal_entity`).
- `concrete: true`.

## How to verify

- `atlas migrate validate --dir file://migrations` exits 0.
- `make migrate.hash` — atlas.sum updates with the new file.

## Reference

- `core-module/model/migrations/0006_type_service_account.sql` — copy-edit template.
