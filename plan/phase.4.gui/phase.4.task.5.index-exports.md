# Phase 4, Task 5 — src/index.ts re-exports

## Context

Public surface of the library.

## Acceptance

`tags-module/gui/src/index.ts` exports:

```ts
export { TagChip } from './TagChip';
export type { TagChipProps } from './TagChip';
export { TagList } from './TagList';
export type { TagListProps } from './TagList';
export { TagEditor } from './TagEditor';
export type { TagEditorProps } from './TagEditor';
export { createTagsClient } from './lib/api';
export type { Tag, TagsClientOptions } from './lib/api';
```

No other exports (types.ts is internal if used).

## How to verify

- `npm run build` produces a `dist/index.d.ts` that declares exactly the surface above.
- A consumer TS file importing from `@moduleforge/tags-gui` gets types for Tag, the three components, and the client factory.
