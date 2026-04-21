# Phase 4 — GUI

## Goal

Ship `<TagChip>`, `<TagList>`, `<TagEditor>` as a presentational React component library packaged with tsup, consumable via yalc by users-module/gui.

## Preconditions

- Phase 1 complete: `tags-module/gui` scaffold builds empty.
- Phase 3 not strictly required but its endpoint shape / OpenAPI fragment should be referenced when writing the API client.

## Outputs

- `tags-module/gui/src/lib/api.ts` — API client helpers (fetch wrappers).
- `tags-module/gui/src/TagChip.tsx`
- `tags-module/gui/src/TagList.tsx`
- `tags-module/gui/src/TagEditor.tsx`
- `tags-module/gui/src/types.ts` — shared Tag type, prop types.
- `tags-module/gui/src/index.ts` — re-exports public surface.

## Component contracts (summary — see task files for detail)

All components accept the common props:

- `subject: string` — the subject entity's UUID.
- `purposes?: string[]` — restrict to these purposes. `undefined` or `[]` = free-form (no restriction).
- `noPurpose?: boolean` — if true, chips render only the `value` (not `purpose:value`).

`<TagEditor>` adds:

- Ability to add new tags (form with value input, and purpose picker/input depending on `purposes.length`).
- Ability to remove existing tags (clicking an X on the chip).
- Ability to edit color via a color picker on the chip (optional in v1 if scope tightens — see task file).

## Hard rules

- No direct DOM manipulation beyond what React provides.
- No state management library (useState/useEffect are sufficient).
- No router dependency (pure library).
- Reuse shadcn primitives from `@moduleforge/core-gui` (Button, Badge, Input, Label, Card) rather than vendoring them a second time.
- API calls go through an injectable `fetch`-compatible function (default: global fetch). This keeps the library testable and deploy-target agnostic.
- Color input in the editor accepts `#RRGGBBAA` hex; native `<input type="color">` returns `#RRGGBB` without alpha, so handle alpha separately (either a second slider 0–255, or accept `#RRGGBB` and pad with `FF`). Task file proposes a simple default — review and adjust.

## Tasks

- 4.1 API client helpers
- 4.2 `<TagChip>` presentational component
- 4.3 `<TagList>` read-only list
- 4.4 `<TagEditor>` full editor
- 4.5 `src/index.ts` exports
- 4.6 tsup build + typecheck

## How to verify

- `cd tags-module/gui && npm install && npm run build` produces `dist/index.js`, `dist/index.d.ts`.
- `npm run typecheck` exits 0.
- Smoke import in users-module/gui (after Phase 5 link): `import { TagList } from '@moduleforge/tags-gui'` resolves types.

## Notes

- UI validation of color format mirrors the DB/API rules; return friendly inline errors.
- "Can tag any entity you can read" is enforced by the API. The UI does not need to gate the add form — if Create returns 403/404, show the error inline.
