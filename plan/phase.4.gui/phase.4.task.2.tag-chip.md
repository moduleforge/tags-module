# Phase 4, Task 2 — `<TagChip>` component

## Context

Presentational chip that renders a single Tag. Reused by `<TagList>` and `<TagEditor>`. No data fetching here.

## Acceptance

Create `tags-module/gui/src/TagChip.tsx` exporting `TagChip`:

```ts
export interface TagChipProps {
  tag: Tag;                 // from ./lib/api
  noPurpose?: boolean;      // if true, render only value
  onRemove?: () => void;    // if set, chip shows an X button that calls this
  onColorChange?: (color: string | null) => void; // if set, chip is color-editable
}
```

Behavior:

- Renders a small rounded badge (reuse `Badge` from `@moduleforge/core-gui` if exported; otherwise a div styled similarly).
- Background color derives from `tag.color` when present. If absent, use a neutral default (e.g. `var(--muted)`).
- Text: `tag.purpose + ":" + tag.value` unless `noPurpose`, then `tag.value`.
- If `onRemove` is defined, render a small `×` button after the text.
- If `onColorChange` is defined, the chip is clickable and opens a minimal color popover (native `<input type="color">` plus an alpha slider). Emit `#RRGGBBAA` strings.
- Text color should contrast the background. Simple approach: luminance-based black/white decision. Do not over-engineer.

## How to verify

- `npm run typecheck` exits 0.
- Visual check happens in Phase 5 when rendered inside users-module/gui.

## Notes

- Alpha slider: native HTML doesn't provide one directly. Pair `<input type="color">` (for RGB) with a separate `<input type="range" min="0" max="255">` for alpha. Compose into `#RRGGBBAA`. Simple, works everywhere.
- If `Badge` from core-gui isn't sufficient (e.g. it doesn't accept custom colors), render a `<span>` with inline styles. Do not duplicate the Badge primitive here.
