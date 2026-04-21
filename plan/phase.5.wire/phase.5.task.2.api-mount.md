# Phase 5, Task 2 — users-module/api consumes tags-api

## Context

Wire the tags router into users-module's existing server.

## Acceptance

- `users-module/api/go.mod` adds `require github.com/moduleforge/tags-api v0.0.0` (go.work redirects to the workspace).
- `users-module/api/cmd/server/main.go` changes:
  - Import `tagsapi "github.com/moduleforge/tags-api/httpapi"` and `tagsservice "github.com/moduleforge/tags-api/service"` and `tagsdb "github.com/moduleforge/tags-model/db"`.
  - After `coreSvcs := coreservice.New(...)` block, add:
    ```go
    tagsSvcs := tagsservice.New(tagsdb.New(pool), auditWriter)
    tagsRouter := tagsapi.NewRouter(tagsapi.Deps{
        Pool:      pool,
        Services:  tagsSvcs,
        Audit:     auditWriter,
        Principal: auth.CorePrincipalAdapter{}, // same adapter satisfies both interfaces
        Logger:    logger,
    })
    ```
  - Inside the authenticated `/v1` group (where coreRouter is mounted), mount `tagsRouter` alongside: `r.Mount("/", tagsRouter)` — or whatever the symmetrical pattern is. Verify coreRouter's mount path pattern and mirror it.
- `auth.CorePrincipalAdapter` already satisfies core's `PrincipalExtractor`; confirm it also satisfies tags-api's (which is the same interface imported from core-api — so it does). If tags-api imports `PrincipalExtractor` from core-api rather than redeclaring it, the adapter works directly. This is the expected design.

## How to verify

- `cd users-module/api && go build ./...` exits 0.
- `go vet ./...` exits 0.
- `cd users-module/api && make test` still exits 0.
- A manual `grep -n tagsRouter main.go` shows the router is mounted.

## Reference

- `users-module/api/cmd/server/main.go` — existing coreRouter mount, around lines 200–230.
- `users-module/api/internal/auth/core_adapter.go` — PrincipalExtractor adapter.
