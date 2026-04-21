# Phase 1, Task 3 — Scaffold tags-module/gui

## Context

tags-module/gui will ship `<TagChip>`, `<TagList>`, `<TagEditor>` as a presentational React library built with tsup, published locally via yalc (same pattern as core-module/gui). This task sets up the skeleton; components land in Phase 4.

## Acceptance

- `tags-module/gui/package.json` — identical layout to `core-module/gui/package.json` with these deltas:
  - `"name": "@moduleforge/tags-gui"`
  - `"version": "0.0.1"`
  - `dependencies` includes `"@moduleforge/core-gui": "file:.yalc/@moduleforge/core-gui"` (will be populated by yalc) — so Tag components can reuse core's shadcn Button/Badge primitives.
  - `peerDependencies`: `react`, `react-dom` at the same versions core-module/gui uses.
  - Scripts: `build`, `typecheck`, `clean` — copy from core-module/gui.
- `tags-module/gui/tsconfig.json` — copy from `core-module/gui/tsconfig.json`.
- `tags-module/gui/tsup.config.ts` — copy from `core-module/gui/tsup.config.ts`; entry is `src/index.ts`.
- `tags-module/gui/src/index.ts` — empty placeholder: `export {};`.
- `tags-module/gui/README.md` — one paragraph: what it ships, how to consume.
- `tags-module/gui/.gitignore` — match `core-module/gui/.gitignore`.

## How to verify

- `cd tags-module/gui && npm install && npm run build` produces `dist/index.js` and `dist/index.d.ts`.
- `npm run typecheck` exits 0.

## Notes

- Do NOT add yalc-link wiring in package.json yet; the actual link happens via the root `make link-tags` target defined in Task 1.5. At scaffold time, the `@moduleforge/core-gui` dep may resolve via `npm install` against a not-yet-yalc-linked package — that is OK as long as install + typecheck pass. If npm install chokes because the yalc path isn't there yet, use `"@moduleforge/core-gui": "*"` as a stopgap and document in README that `make link-tags` must be run first in a fresh checkout.
