# Phase 4, Task 4 — `<TagEditor>` component

## Context

Full CRUD editor for tags on a single subject. Superset of `<TagList>` behavior plus add/remove + color-edit.

## Acceptance

Create `tags-module/gui/src/TagEditor.tsx` exporting `TagEditor`:

```ts
export interface TagEditorProps {
  subject: string;
  purposes?: string[];          // undefined or [] = free-form
  noPurpose?: boolean;
  client: ReturnType<typeof createTagsClient>;
  className?: string;
  onChange?: (tags: Tag[]) => void;  // fired after each successful mutation
}
```

Behavior:

- Same fetch-on-mount as `<TagList>`.
- Renders existing tags as `<TagChip tag noPurpose onRemove onColorChange />` — so chips are removable and their color editable.
  - `onRemove` calls `client.remove(tag.uuid)` then refetches.
  - `onColorChange` calls `client.updateColor(tag.uuid, color)` then updates state.
- Renders an "add" form:
  - **purpose handling:**
    - `purposes === undefined` or `purposes.length === 0` → free-form `<input>`.
    - `purposes.length === 1` → purpose is fixed to `purposes[0]`, no input shown (display as a label next to the value input).
    - `purposes.length > 1` → `<select>` populated from `purposes`.
  - **value input** — always shown, required, maxlength 512.
  - **color picker** — optional; reuse the chip's color-edit UI if simple, else a plain `<input type="color">` + alpha slider.
  - **submit button** — disabled when submitting; error shown inline on failure.
- On successful Create, refetch and clear inputs. Call `onChange` with the new list.
- Deduplicate locally: if the user tries to create a (purpose) that already exists among the fetched tags and owner is self, show a friendly inline error ("A tag with this purpose already exists") before calling the server. The server will reject with 409 anyway; this is UX-only.

## How to verify

- `npm run typecheck` exits 0.
- Visual in Phase 5.

## Notes

- The editor is not "locked" to single-owner semantics at the UI level — if the current user edits a subject that already has someone else's tag with the same purpose, they can still create their own (different owner → unique). The "already exists" check should compare against tags the current user owns, not all tags in the list. For v1, simplest heuristic: don't attempt the client-side dedupe and just let the server 409 surface through an inline error message. Ship the simple version.
