# tags-module — TODO

Status: `[ ]` not started · `[~]` in progress · `[x]` done · `[!]` blocked

- [x] **Phase 1 — Bootstrap tags-module skeleton** (depends on: none)
  - [x] 1.1 Scaffold `tags-module/model/` (go.mod, atlas.hcl, sqlc.yaml, Makefile, .gitignore, README)
  - [x] 1.2 Scaffold `tags-module/api/` (go.mod, Makefile, empty service/httpapi packages, go.sum)
  - [x] 1.3 Scaffold `tags-module/gui/` (package.json, tsconfig, tsup config, empty src/index.ts)
  - [x] 1.4 Update top-level `go.work` to include tags-module modules + replace directives
  - [x] 1.5 Update root `Makefile` (add tags-module to GO_PROJECTS; add `link-tags` / `unlink-tags`; wire a `link-all`)
  - [x] 1.6 Update root `README.md` to describe tags-module alongside core/users

- [x] **Phase 2 — Model** (depends on: 1)
  - [x] 2.1 Migration 0200 — register `tag` type under `entity`
  - [x] 2.2 Migration 0201 — `tags` table + triggers (type-descent, immutability, set_updated_at) + unique index
  - [x] 2.3 `queries/tags.sql` — CRUD + search queries
  - [x] 2.4 sqlc build → `tags-module/model/db/`
  - [x] 2.5 atlas migrate hash → `atlas.sum`
  - [x] 2.6 Compile-check `tags-module/model`

- [ ] **Phase 3 — API** (depends on: 2)
  - [ ] 3.1 `tags-module/api/service/tag.go` — TagService (Create, Get, Search, Update, Delete) with authz + audit
  - [ ] 3.2 `tags-module/api/service/display.go` — register `tag → name → "purpose:value"`
  - [ ] 3.3 `tags-module/api/httpapi/router.go` + handlers (`tags.go`, `subject_tags.go`, `response.go`)
  - [ ] 3.4 Unit tests (service authz paths, uniqueness collision, immutable fields, handler auth paths)
  - [ ] 3.5 `tags-module/api/openapi.fragment.yaml`
  - [ ] 3.6 `make test` green in `tags-module/api`

- [ ] **Phase 4 — GUI** (depends on: 1; may run in parallel with 2/3)
  - [ ] 4.1 API client helpers (`src/lib/api.ts`) — `listTagsForSubject`, `createTag`, `updateTagColor`, `deleteTag`
  - [ ] 4.2 `<TagChip>` component (presentational, color-aware)
  - [ ] 4.3 `<TagList>` component (fetch by subject + optional purposes, renders chips)
  - [ ] 4.4 `<TagEditor>` component (add/remove; purpose handling per `purposes` prop)
  - [ ] 4.5 `src/index.ts` re-exports
  - [ ] 4.6 tsup build produces ESM + types
  - [ ] 4.7 `npm run typecheck` clean

- [ ] **Phase 5 — Wire into users-module** (depends on: 2, 3, 4)
  - [ ] 5.1 Extend `users-module/model` compose pipeline to copy tags migrations/queries
  - [ ] 5.2 `atlas migrate hash` on composed dir; regenerate sqlc
  - [ ] 5.3 `users-module/api/go.mod` require `github.com/moduleforge/tags-api`
  - [ ] 5.4 `users-module/api/cmd/server/main.go` construct TagService + mount tags router inside authed `/v1`
  - [ ] 5.5 Tags service registers its display renderer into the same `display.Registry`
  - [ ] 5.6 `make link-tags` (root Makefile target) publishes tags-module/gui via yalc and adds into users-module/gui
  - [ ] 5.7 Sanity: a users-module integration test creates a tag via HTTP, reads it back, deletes it

- [ ] **Phase 6 — Verification + cleanup** (depends on: all)
  - [ ] 6.1 `make test` green across every sub-project
  - [ ] 6.2 `make dev.start` smoke (manual — hand to user)
  - [ ] 6.3 `atlas migrate status` shows 0000–00NN, 0100–01NN, 0200–02NN in order
  - [ ] 6.4 grep sanity: no `tags` tables in core-module or users-module/model
  - [ ] 6.5 Audit log entries present for create/update/delete (manual — hand to user)
  - [ ] 6.6 Update users-module summary + root README noting tags-module dependency

## Reports

Drop progress notes into `report.<N>.<topic>.md` in this directory as work proceeds.
