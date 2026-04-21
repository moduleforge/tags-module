# Phase 4, Task 6 — Build + typecheck gate

## Context

Final gate before wiring phase.

## Acceptance

- `cd tags-module/gui && npm run build` produces `dist/index.js`, `dist/index.mjs` (if dual-format), and `dist/index.d.ts`.
- `npm run typecheck` exits 0.
- `dist/` is gitignored (per scaffold).

No new files.
