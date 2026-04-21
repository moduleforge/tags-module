# Phase 2, Task 5 — Hash migrations

## Context

Atlas tracks migration integrity via `atlas.sum`. Regenerate it after adding 0200 and 0201.

## Acceptance

- `cd tags-module/model && make migrate.hash` succeeds.
- `tags-module/model/migrations/atlas.sum` is present and up to date.
- `atlas migrate validate --dir file://migrations` exits 0.

## How to verify

- Delete and regenerate atlas.sum; confirm the result is byte-identical to what was committed.
