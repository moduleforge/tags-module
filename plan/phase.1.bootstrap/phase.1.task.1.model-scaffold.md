# Phase 1, Task 1 — Scaffold tags-module/model

## Context

tags-module/model holds the Tag-entity schema (migrations 0200+), sqlc queries, and generated `db/` package. This task creates the skeleton; migrations/queries land in Phase 2.

## Acceptance

- `tags-module/model/go.mod` — module `github.com/moduleforge/tags-model`, same Go version as `core-module/model/go.mod`.
- `tags-module/model/atlas.hcl` — mirror `core-module/model/atlas.hcl`, migration dir `file://migrations`.
- `tags-module/model/sqlc.yaml` — v2 config, engine postgresql, schema `./migrations`, queries `./queries`, output `./db`, package `db`, `emit_interface: true`, `emit_json_tags: true`, `emit_prepared_queries: false`. Package import path in overrides: `github.com/moduleforge/tags-model/db`.
- `tags-module/model/Makefile` — canonical targets (`build`, `test`, `lint`, `lint-fix`, `clean`, `migrate.new`, `migrate.up`, `migrate.status`, `migrate.hash`). Model up `core-module/model/Makefile` exactly.
- `tags-module/model/.gitignore` — match `core-module/model/.gitignore`.
- `tags-module/model/migrations/.keep` and `tags-module/model/queries/.keep` — empty dirs present in git.
- `tags-module/model/README.md` — one-paragraph description, install/build/migrate commands, and a one-line note documenting the reserved migration range `0200–0299`.

## How to verify

- `cd tags-module/model && make build` exits 0. (sqlc produces an empty `db/db.go` since no queries yet.)
- `atlas migrate validate --dir file://migrations` runs (empty dir is acceptable).

## Reference

- Template: `core-module/model/` — copy the structure, rename module path, update import paths in sqlc overrides.
- Make convention memory: user preference favours canonical target naming (build default, test, lint / lint-fix, migrate.* namespace).
