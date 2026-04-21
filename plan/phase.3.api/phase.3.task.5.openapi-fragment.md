# Phase 3, Task 5 — OpenAPI fragment

## Context

Document the six tag endpoints in an OpenAPI 3.0.3 fragment. Consumers merge this fragment into their own spec.

## Acceptance

`tags-module/api/openapi.fragment.yaml` includes:

- `paths`:
  - `/tags` — POST, GET
  - `/tags/{uuid}` — GET, PUT, DELETE
  - `/entities/{uuid}/tags` — GET
- `components.schemas`:
  - `Tag` — uuid, ownerUuid, subjectUuid, purpose, value, color (nullable), createdAt, updatedAt.
  - `CreateTagRequest` — subject (uuid), purpose, value, color (optional).
  - `UpdateTagRequest` — color (required; the only mutable field).
  - `TagListResponse` — `{tags: [Tag]}`.
  - `ErrorResponse` — reuse if core-module defines one; else define locally.
- Every response documented with correct status codes (201/200/204/400/401/403/404).
- Security scheme references `Authorization: Bearer` (matching users-module).

## How to verify

- `npx @redocly/cli lint openapi.fragment.yaml` exits 0.

## Reference

- `core-module/api/openapi.fragment.yaml` — the shape/style to match.
