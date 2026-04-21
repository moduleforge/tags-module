# Phase 5, Task 4 — yalc link tags-module/gui into users-module/gui

## Context

Make `@moduleforge/tags-gui` available to `users-module/gui` via yalc, same pattern as core-gui.

## Acceptance

- Run `make link-tags` from repo root (target defined in Phase 1 Task 1.5).
- `users-module/gui/package.json` now has a `"@moduleforge/tags-gui": "file:.yalc/@moduleforge/tags-gui"` entry.
- `users-module/gui/yalc.lock` updated.
- `cd users-module/gui && npm run typecheck` exits 0 (no users-module code imports tags-gui yet, but the package resolves).

## Optional smoke

If there's an obvious low-stakes place to drop a `<TagEditor subject={user.uuid} />` in a user admin page for manual testing, do it. Otherwise defer to a later task — the module is consumable from this point forward.

## How to verify

- `cd users-module/gui && ls node_modules/@moduleforge/tags-gui/dist` — shows the built artifacts.
- `import { TagEditor } from '@moduleforge/tags-gui'` in a scratch TS file resolves types.
