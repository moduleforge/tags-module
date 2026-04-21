# Phase 5 — Wire into users-module

## Goal

Connect tags-module into the users-module runtime: compose migrations, mount the tags router inside the authenticated `/v1` group, register the tag display renderer, and link the gui into users-module/gui.

## Preconditions

- Phase 2, 3, 4 complete.

## Outputs

- `users-module/model` compose pipeline extended to include tags migrations + queries.
- `users-module/api/go.mod` requires `github.com/moduleforge/tags-api` (and transitively `github.com/moduleforge/tags-model`).
- `users-module/api/cmd/server/main.go` constructs TagService, mounts tagsRouter at `/v1` inside the authed group, registers display renderer.
- `users-module/gui` consumes `@moduleforge/tags-gui` via yalc.
- Optionally one small end-to-end integration test inside users-module/api that creates → reads → deletes a tag.

## Hard rules

- users-module code does NOT import `tags-module/model` directly. The Querier is exposed indirectly via tags-module/api's Services aggregate — identical to how users-module consumes core-model (only core-api types cross the boundary).
- Actually: users-module's auditWriter satisfies core's `audit.Writer` interface; since tags-module uses the same interface, the same writer is reused.
- Mount point: `r.Mount("/", tagsRouter)` inside the authenticated `/v1` subtree, next to `r.Mount("/", coreRouter)`. All tags paths (`/tags`, `/entities/{uuid}/tags`) end up at `/v1/tags` and `/v1/entities/{uuid}/tags`.

## Tasks

- 5.1 Extend users-module/model compose pipeline
- 5.2 Update users-module/api to require and mount tags-api
- 5.3 Register display renderer alongside core's
- 5.4 yalc link tags-module/gui
- 5.5 users-module/api integration test

## How to verify

- `cd users-module/api && make test` still green.
- `atlas migrate validate` on composed dir succeeds with 0000–0011 + 0100–01NN + 0200–0201 in order.
- `make dev.start` in users-module brings up a stack where `POST /v1/tags` (with a valid bearer) creates a tag (manual, Phase 6 smoke).
