# Phase 2, Task 2 — Migration 0201: `tags` table + triggers + indexes

## Context

The core of the tags-module schema: the `tags` subtype table plus its constraints and triggers.

## Acceptance

Migration file `tags-module/model/migrations/0201_tags.sql` must:

### Table

```sql
CREATE TABLE tags (
  entity_id   BIGINT PRIMARY KEY REFERENCES entities(id) ON DELETE RESTRICT,
  owner_id    BIGINT NOT NULL REFERENCES entities(id),
  subject_id  BIGINT NOT NULL REFERENCES entities(id),
  purpose     TEXT NOT NULL CHECK (char_length(purpose) <= 512),
  value       TEXT NOT NULL CHECK (char_length(value) <= 512),
  color       TEXT CHECK (color SIMILAR TO '#[0-9A-Fa-f]{8}'),
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

### Indexes

```sql
CREATE UNIQUE INDEX tags_owner_subject_purpose_idx
  ON tags (owner_id, subject_id, purpose);

CREATE INDEX tags_subject_id_idx ON tags (subject_id);
-- owner_id alone is already covered by the leading column of the unique index,
-- no separate index needed.
```

### Type-descent trigger

Mirror `legal_entities_check_type` from `core-module/model/migrations/0008_legal_entities.sql`, renamed for `tags`, asserting that `NEW.entity_id`'s fundamental type descends from `'tag'`. Fire `BEFORE INSERT`.

### Immutability trigger

One function `tags_reject_immutable_changes()` that raises on any `DISTINCT FROM` to `owner_id`, `subject_id`, `purpose`, or `value`. Each check should produce a clear error message naming the offending column. Fire `BEFORE UPDATE` on `tags`. Model up the shape of `entities_immutable_type()` but extend for four columns.

### updated_at trigger

```sql
CREATE TRIGGER tags_set_updated_at
  BEFORE UPDATE ON tags
  FOR EACH ROW EXECUTE FUNCTION set_updated_at();
```

### Comments

Add one top-of-file comment explaining the table is the single-subject tag subtype and that all identity/lifecycle columns (uuid, archived_at) live on the parent `entities` row.

## How to verify

- `atlas migrate validate --dir file://migrations` exits 0.
- In a scratch DB with core migrations applied, running `atlas migrate apply` against the composed dir (core 0000–0011 + tags 0200–0201) applies cleanly.
- A manual SQL scratch test:
  1. Insert a type `tag`, an entity with that type, an owner entity (any existing type), a subject entity.
  2. Insert a `tags` row — succeeds.
  3. Insert a second row with identical `(owner_id, subject_id, purpose)` — rejected by unique index.
  4. Insert a row with `color = 'notacolor'` — rejected by CHECK.
  5. Insert a row with `color = '#FF00FFAA'` — accepted.
  6. UPDATE any of `owner_id/subject_id/purpose/value` — rejected by trigger.
  7. UPDATE `color` — accepted; `updated_at` bumps.

These manual checks are not required in an automated test, but the implementer should do at least one round in psql against a scratch DB and note results in the PR.

## Reference

- `core-module/model/migrations/0007_entities.sql` — immutability trigger pattern.
- `core-module/model/migrations/0008_legal_entities.sql` — type-descent trigger pattern.
- `core-module/model/migrations/0000_helpers.sql` — `set_updated_at` + `type_is_or_descends_from`.
