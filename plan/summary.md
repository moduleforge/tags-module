# tags-module — Plan summary

## Purpose

Add a third portable module, sibling to `core-module` and `users-module`, that models user-applied Tags. Tags are Entities (per the core entity-typing pattern), so tags-module registers its own concrete type under `entity` and stands up its own `tags` subtype table, Go service + HTTP router, and React component library. Downstream apps mount the router next to core's and drop the React components into their UI to get tagging end-to-end without re-implementing anything.

## Scope

**Model:**
- New concrete `type` registered under `entity`: `tag` (no parent subtype — Tags are not Legal Entities).
- New subtype table `tags`:
  - `entity_id BIGINT PRIMARY KEY REFERENCES entities(id) ON DELETE RESTRICT` — tag's own entity anchor.
  - `owner_id BIGINT NOT NULL REFERENCES entities(id)` — set server-side from principal, immutable.
  - `subject_id BIGINT NOT NULL REFERENCES entities(id)` — the tagged thing, immutable.
  - `purpose TEXT NOT NULL CHECK (char_length(purpose) <= 512)` — immutable.
  - `value TEXT NOT NULL CHECK (char_length(value) <= 512)` — immutable.
  - `color TEXT CHECK (color SIMILAR TO '#[0-9A-Fa-f]{8}')` — `#RRGGBBAA` hex, nullable, mutable.
  - `created_at`, `updated_at` — managed by `set_updated_at` trigger.
  - `UNIQUE (owner_id, subject_id, purpose)` — one purpose per (owner, subject).
- Subtype-table triggers mirror core's existing discipline:
  - Type-descent check on INSERT (entity's fundamental type must descend from `tag`).
  - Immutability triggers on `owner_id`, `subject_id`, `purpose`, `value` (BEFORE UPDATE).
  - Standard `set_updated_at` trigger on UPDATE.

**API (Go, chi subrouter):**
- `POST   /tags` — create; body `{subject, purpose, value, color?}`; owner set from principal.
- `GET    /tags?owner=…&subject=…&purpose=…&value=…` — search; ≥1 of owner/subject required (400 otherwise) to prevent full-table scans.
- `GET    /tags/{uuid}` — detail.
- `PUT    /tags/{uuid}` — body accepts `{color}` only; 400 on anything else.
- `DELETE /tags/{uuid}` — hard delete (tag row + entity row in a tx).
- `GET    /entities/{uuid}/tags?purpose=…` — tags for a subject; empty array when none.
- Display renderer registered: `tag → name → "purpose:value"`. `description` renderer deferred.
- All mutating paths write an audit entry via `audit.Writer` (interface from core-module).

**Authorization:**
- Create: any authenticated user; may tag any entity the caller can read.
- Read (detail, list, subject-listing): owner = full, subject = read-only, admin = full, else = 404 (don't leak existence). Search filters to the authz-visible set.
- Update: owner or admin only.
- Delete: owner or admin only.

**GUI (React library, tsup + yalc, same shape as core-module/gui):**
- `<TagChip tag noPurpose />` — presentational chip coloured from `tag.color`; text is `purpose:value` (or just `value` if `noPurpose`).
- `<TagList subject purposes noPurpose />` — fetches and renders tags for a subject, optionally filtered to a set of purposes.
- `<TagEditor subject purposes noPurpose />` — add/remove; purpose handling:
  - `purposes` undefined or `[]` → free-form input.
  - `purposes.length === 1` → purpose is fixed, user enters value.
  - `purposes.length > 1` → purpose picker, user enters value.

## Locked decisions

1. **Tag is a direct child of `entity`.** No intermediate abstract type.
2. **`owner_id` lives on `tags`** (not on core's `entities`). No core schema changes. Revisit a generic `ownerships` table only when a second consumer needs it.
3. **`UNIQUE (owner_id, subject_id, purpose)`** on the `tags` table — standard DB constraint, no triggers needed for uniqueness.
4. **Immutability enforced by DB triggers**, not API convention.
5. **Color stored as `TEXT` with `SIMILAR TO '#[0-9A-Fa-f]{8}'`.** `SIMILAR TO` is SQL:1999 standard; UI validates independently; color-math helper SQL functions are deferred.
6. **Authorization: owner-full, subject-read, admin-full, else-404.** Encoded explicitly in the service layer.
7. **Audit goes through `core-module/api/audit.Writer`.** users-module's writer already satisfies it.
8. **Migration numbering: `0200–0299`** is reserved for tags-module.
9. **tags-module ships its own compose target.** When users-module composes its migration dir it also pulls from tags-module.
10. **Single subject per tag.** The prior "tag applies to 0+ entities" framing was superseded — two tags with matching `purpose:value` on different subjects are distinct entities.

## Module layout

```
user-components/
├── go.work                                # adds tags-module/model + tags-module/api
├── Makefile                               # gains link-tags, unlink-tags; link-all aggregator
├── core-module/ …                         # unchanged
├── users-module/ …                        # gains tags compose + router mount + gui link
└── tags-module/
    ├── model/      # github.com/moduleforge/tags-model
    │   ├── migrations/0200_type_tag.sql, 0201_tags.sql
    │   ├── queries/tags.sql
    │   └── db/                            # sqlc output
    ├── api/        # github.com/moduleforge/tags-api
    │   ├── service/                       # TagService (tx-aware), DisplayBuiltins
    │   └── httpapi/                       # NewRouter(Deps) chi.Router
    └── gui/        # @moduleforge/tags-gui
        ├── src/{TagChip,TagList,TagEditor,lib/api.ts,index.ts}
        └── dist/                          # tsup output
```

## Architecture notes

- tags-module depends on **core-module** (model: `entities`, `types`, helper functions; api: `audit.Writer`, `Principal`, `PrincipalExtractor`, display registry). Same consumption pattern as users-module.
- tags-module does **not** depend on users-module. A consumer that doesn't want users-module (e.g. a different app) can still use tags-module on top of core-module alone, supplying its own `PrincipalExtractor` and `audit.Writer`.
- Compose pipeline in users-module/model is extended to pull `tags-module/model/migrations/*` alongside core's. Migration numbers don't collide (core 0000–0099, users 0100–0199, tags 0200–0299).
- Display renderer registration happens at tags-module service init and writes into the same core `display.Registry` the consumer already constructs.

## Phases

1. **Bootstrap tags-module skeleton** — dirs, go.mod × 2, package.json, Makefiles, root Makefile wiring, go.work entries.
2. **Model** — migrations 0200–0201, queries, sqlc build, atlas hash.
3. **API** — TagService (tx-aware, authz in service), display renderer, chi subrouter + handlers, tests, OpenAPI fragment.
4. **GUI** — `<TagChip>`, `<TagList>`, `<TagEditor>`, API client, tsup build.
5. **Wire into users-module** — compose tags migrations, go.mod require, main.go mount, gui yalc link.
6. **Verification + cleanup** — full `make test`, smoke, atlas status, docs.

Phases 2–4 can partially overlap once Phase 1 is done. Phase 5 blocks on 2, 3, 4. Phase 6 blocks on all.

## Open questions (v2)

None at plan-write time. Any surfacing during implementation should go into a `report.*.md` in this directory.

## Deferred / noted for later

- SQL helper functions for color math (`color_distance`, `is_light`, etc.) — note in model README.
- Admin-only full search endpoint (unrestricted `GET /tags`) — not needed now.
- Generic `ownerships` table in core — only if a second consumer appears.
- Splitting audit into its own module — users-module's plan already flags this direction.

## Supporting documents

- `TODO.md` — live phase/task checklist.
- `phase.<N>.<title>/phase.<N>.<title>.md` — phase summary + acceptance.
- `phase.<N>.<title>/phase.<N>.task.<M>.<title>.md` — agent-ready task instructions.
