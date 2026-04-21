# Phase 4, Task 3 — `<TagList>` component

## Context

Read-only rendering of a subject's tags. Fetches on mount and on `subject`/`purposes` change.

## Acceptance

Create `tags-module/gui/src/TagList.tsx` exporting `TagList`:

```ts
export interface TagListProps {
  subject: string;                      // subject entity UUID
  purposes?: string[];                  // undefined or [] = free-form (all)
  noPurpose?: boolean;
  client: ReturnType<typeof createTagsClient>;
  className?: string;                   // passthrough for layout
}
```

Behavior:

- Treats `purposes === undefined` and `purposes.length === 0` as equivalent — both mean "all purposes".
- On mount and on `[subject, purposesKey]` change (where `purposesKey = JSON.stringify(purposes ?? [])`), calls `client.listBySubject(subject, purposes)`.
- If `purposes.length > 1`, filter the result on the client to tags whose `purpose` is in the set (redundancy with server for robustness).
- Renders a flex container of `<TagChip tag={t} noPurpose={noPurpose} />` for each tag.
- While loading: render a small skeleton (one or two placeholder chips).
- On error: render a small inline error message; do not throw.
- Empty state: render nothing (no "no tags" message) — keeps the component unobtrusive when embedded next to other content.

## How to verify

- `npm run typecheck` exits 0.
- Visual in Phase 5.
