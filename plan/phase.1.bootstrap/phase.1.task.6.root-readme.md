# Phase 1, Task 6 — Update root README

## Context

`user-components/README.md` describes the repo layout and onboarding flow. Update it to include tags-module.

## Acceptance

- The "Module layout" section adds a `tags-module/` entry with a one-line description ("tags on entities — API + React components") and the same sub-dir breakdown style used for core-module and users-module.
- The "First-time setup" section mentions `make link-tags` alongside `make link-core` (or points to `make link-all`).
- The "When you change core-module" sibling section gets a short "When you change tags-module" block noting that Go changes flow through `go work sync` and GUI changes need `make link-tags`.
- No other README content changes.

## How to verify

- Open the README and confirm the three module-layout entries now show core/users/tags.
- `make link-all` is documented.
