# @moduleforge/tags-gui

React component library for tag management UIs. Will ship `<TagChip>`, `<TagList>`, and `<TagEditor>` presentational components built on top of `@moduleforge/core-gui` primitives. This package is scaffolded in Phase 1; components land in Phase 4.

## Build

```bash
npm run build
```

Outputs `dist/index.js` (CJS), `dist/index.mjs` (ESM), and `dist/index.d.ts` (types) via tsup.

## Fresh checkout note

`@moduleforge/core-gui` is resolved via yalc in a local development setup. On a fresh checkout, run `make link-tags` from the repo root before `npm install` to publish and link `@moduleforge/core-gui` into this package's `.yalc/` store. Until that step is performed, `@moduleforge/core-gui` resolves to `"*"` and may not install correctly without the yalc link.
