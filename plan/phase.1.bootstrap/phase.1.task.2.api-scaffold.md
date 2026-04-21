# Phase 1, Task 2 — Scaffold tags-module/api

## Context

tags-module/api holds the Go service + HTTP router for tags. This task sets up the empty package skeleton; real code lands in Phase 3.

## Acceptance

- `tags-module/api/go.mod` — module `github.com/moduleforge/tags-api`, same Go version as `core-module/api/go.mod`.
- go.mod requires the same third-party deps core-module/api uses that tags-module will also need: `github.com/jackc/pgx/v5`, `github.com/google/uuid`, `github.com/go-chi/chi/v5`, `github.com/moduleforge/core-api` (latest local via go.work), `github.com/moduleforge/core-model` (latest local via go.work), `github.com/moduleforge/tags-model` (latest local via go.work). Leave versions unpinned where go.work will redirect to the workspace.
- `tags-module/api/service/doc.go` — one-line package doc stub: `// Package service exposes tx-aware tag CRUD for consumer apps.`
- `tags-module/api/httpapi/doc.go` — one-line stub: `// Package httpapi provides the chi subrouter for tag endpoints.`
- `tags-module/api/Makefile` — canonical targets (`build`, `test`, `lint`, `lint-fix`, `clean`). Mirror `core-module/api/Makefile`.
- `tags-module/api/README.md` — one paragraph.

## How to verify

- `cd tags-module/api && go build ./...` exits 0 (empty packages build).
- `cd tags-module/api && go vet ./...` exits 0.

## Reference

- Template: `core-module/api/` — copy structure, rename module path.
- Leave go.sum creation until after go.work is updated in Task 1.4 — `go mod tidy` in this directory should then succeed.
