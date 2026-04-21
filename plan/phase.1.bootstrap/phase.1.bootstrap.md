# Phase 1 — Bootstrap tags-module skeleton

## Goal

Stand up empty `tags-module/{model,api,gui}` scaffolds, stitch them into the top-level `go.work`, add yalc plumbing for the gui, and update the root Makefile + README. After this phase all sub-projects compile/build clean but are empty of business logic.

## Outputs

- `tags-module/model/{go.mod, atlas.hcl, sqlc.yaml, Makefile, README.md, .gitignore, migrations/.keep, queries/.keep}`
- `tags-module/api/{go.mod, Makefile, service/doc.go, httpapi/doc.go, README.md}`
- `tags-module/gui/{package.json, tsconfig.json, tsup.config.ts, src/index.ts, README.md, .gitignore}`
- `user-components/go.work` updated to include tags-module paths and `replace` directive
- `user-components/Makefile` gains tags-module entries + `link-tags` / `unlink-tags` targets
- `user-components/README.md` mentions tags-module in the module-layout section

## Hard rules

- Module paths:
  - Go: `github.com/moduleforge/tags-model`, `github.com/moduleforge/tags-api`
  - Node: `@moduleforge/tags-gui`
- Go version pinned to match `core-module/api/go.mod`.
- `tags-module/api` declares dependencies on `github.com/moduleforge/core-api` and `github.com/moduleforge/core-model` (empty usage at this stage — may need a `//go:build ignore` placeholder import if the build complains; otherwise wire in Phase 3).
- `tags-module/gui/package.json` declares `peerDependencies` on `react`, `react-dom` (versions matching `core-module/gui`). Declares `dependencies` on `@moduleforge/core-gui` so consumers get the shared shadcn primitives.
- No runtime logic in this phase — only scaffolding + doc.go / index.ts stubs.

## Reference templates

- `core-module/model/` for model scaffold (atlas.hcl, sqlc.yaml, Makefile conventions)
- `core-module/api/` for api scaffold
- `core-module/gui/` for gui scaffold (tsup, package.json, tsconfig)

## Tasks

- 1.1 Scaffold tags-module/model
- 1.2 Scaffold tags-module/api
- 1.3 Scaffold tags-module/gui
- 1.4 Update top-level go.work
- 1.5 Update root Makefile
- 1.6 Update root README

## How to verify

- `go work sync` at repo root succeeds.
- `cd tags-module/model && make build` exits 0 (sqlc produces empty db package).
- `cd tags-module/api && go build ./...` exits 0.
- `cd tags-module/gui && npm install && npm run build` produces `dist/index.js`.
- `make link-tags` at repo root succeeds (yalc publish + add into users-module/gui).
