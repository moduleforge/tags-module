# Phase 3, Task 2 — Display renderer registration

## Context

core-module/api/display.Registry is the project-wide mechanism for rendering human-readable strings from entities. tags-module registers its own renderer for the `tag` type.

## Acceptance

Create `tags-module/api/service/display.go` exporting `RegisterBuiltins(reg *display.Registry, q tagsdb.Querier)` (or similar signature mirroring `core-module/api/service/display_builtins.go`). The function registers:

- `("tag", display.FieldName)` → renderer that loads the tag by entity id and returns `purpose + ":" + value`.

Do not register a `FieldDescription` renderer for v1; it's not needed and deferred.

Callers (users-module's main.go in Phase 5) will call this after calling core's `service.RegisterBuiltins` on the same registry.

## How to verify

- `go build ./...` in tags-module/api.
- Add a unit test in Task 3.4 that constructs a Registry, registers tags' builtin, and asserts `Render(ctx, tx, entityID, display.FieldName)` returns `purpose:value`.

## Reference

- `core-module/api/service/display_builtins.go` — the pattern.
- `core-module/api/display/registry.go` — the Registry type.
