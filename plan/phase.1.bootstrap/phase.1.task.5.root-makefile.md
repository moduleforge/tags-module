# Phase 1, Task 5 — Update root Makefile

## Context

The root `user-components/Makefile` aggregates `build`/`test`/`clean` across all sub-projects and provides yalc link targets for gui packages. Extend it to include tags-module.

## Acceptance

- Add `tags-module/model` and `tags-module/api` to `GO_PROJECTS`.
- Introduce a new variable `TAGS_GUI_DIR := tags-module/gui` and include it in the same build/test/clean loops used for `CORE_GUI_DIR` (i.e. a second `cd $(TAGS_GUI_DIR) && npm run …` block per target, or refactor into a list). Keep it simple — duplicated block is fine, matches the existing style.
- Add `link-tags` target: build + publish `@moduleforge/tags-gui` via `yalc publish --push` from `tags-module/gui`, then `yalc add @moduleforge/tags-gui` in `users-module/gui`, then `yalc update`, then `go work sync`. Model exactly on the existing `link-core` target.
- Add `unlink-tags` target (mirrors `unlink-core`).
- Add a convenience `link-all` target that depends on `link-core` and `link-tags`.
- Update `help` output so the new targets are listed (they auto-list via the `awk` pattern; just make sure the `##` comment is present on each).

## How to verify

- `make help` shows `link-tags`, `unlink-tags`, `link-all`.
- `make build` builds all five sub-projects (core-model, core-api, users-model, users-api, tags-model, tags-api) + both gui libraries (core-gui, tags-gui) + users-module/gui.
- `make link-tags` exits 0 (assumes Phase 1 Task 3 has produced a buildable tags-module/gui).
- `make link-all` exits 0.

## Reference

- Existing `link-core` / `unlink-core` in `user-components/Makefile` — clone and rename.
