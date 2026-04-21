# Phase 4, Task 1 — API client helpers

## Context

A tiny wrapper over `fetch` that the three components use. Must be injectable so consumers can pass their own auth headers.

## Acceptance

Create `tags-module/gui/src/lib/api.ts` exporting:

```ts
export interface Tag {
  uuid: string;
  ownerUuid: string;
  subjectUuid: string;
  purpose: string;
  value: string;
  color?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TagsClientOptions {
  // Base URL for the API (e.g. "/v1"). Required.
  baseUrl: string;
  // Fetch implementation; defaults to global fetch.
  fetchImpl?: typeof fetch;
  // Called before every request; return Headers (or an object) to merge.
  // Typically used by consumers to add Authorization: Bearer.
  headers?: () => Record<string, string> | Promise<Record<string, string>>;
}

export function createTagsClient(opts: TagsClientOptions): {
  listBySubject(subjectUuid: string, purposes?: string[]): Promise<Tag[]>;
  create(input: { subject: string; purpose: string; value: string; color?: string }): Promise<Tag>;
  updateColor(uuid: string, color: string | null): Promise<Tag>;
  remove(uuid: string): Promise<void>;
};
```

Implementation notes:

- `listBySubject` calls `GET {baseUrl}/entities/{subjectUuid}/tags` optionally repeated per purpose filter on the client side if the server's single-purpose query param requires it (server accepts a single `purpose=` — the client issues N calls and concatenates if `purposes.length > 1`, OR simply requests all and filters client-side when multiple purposes are provided; prefer the latter for simplicity).
- `create` calls `POST {baseUrl}/tags`.
- `updateColor` calls `PUT {baseUrl}/tags/{uuid}` with body `{color}`.
- `remove` calls `DELETE {baseUrl}/tags/{uuid}`.
- All methods throw an `Error` with the server's message on non-2xx.
- The client should be dependency-free (no axios / tanstack-query) — just fetch.

## How to verify

- `npm run typecheck` exits 0.
- A quick scratch `.test.ts` (no runner needed; don't commit) that constructs the client with a mock fetch and asserts the URLs and bodies are correct.

## Notes

- We intentionally pass `Authorization` via the injected `headers()` callback rather than hardcoding bearer-token handling. Consumers (like users-module/gui) already have auth context and know how to produce headers.
