# Phase 5, Task 3 — Register tag display renderer

## Context

users-module constructs the `display.Registry` and currently registers core's builtins. Extend to also register tags-module's builtin.

## Acceptance

- In `users-module/api/cmd/server/main.go` (or wherever the display registry is built, likely near coreSvcs construction):
  ```go
  tagsservice.RegisterBuiltins(displayReg, tagsdb.New(pool))
  ```
  called immediately after the equivalent `coreservice.RegisterBuiltins(...)` call.
- If users-module doesn't currently expose the display.Registry at all (i.e. core does its renderer registration internally), audit and confirm — the hook point may be inside `coreservice.New` or in a separate init path. Adjust accordingly to plug in tags' renderer through the same mechanism.

## How to verify

- `go build ./...` exits 0.
- (Integration in Phase 5.5): a call to the Registry's `Render(ctx, tx, entityID, display.FieldName)` on a tag entity returns `"<purpose>:<value>"`.

## Reference

- Core's existing display registration in main.go (find by `RegisterBuiltins`).
- `tags-module/api/service/display.go` (Phase 3 Task 2).
