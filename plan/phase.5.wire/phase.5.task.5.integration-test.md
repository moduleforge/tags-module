# Phase 5, Task 5 — users-module/api integration test for tags

## Context

Round-trip an auth'd HTTP call sequence through the composed server: create a tag, read it, list by subject, update color, delete.

## Acceptance

Add an integration-level test under `users-module/api/internal/handlers` (or wherever existing integration tests live) that:

1. Spins up the test server (uses existing test infra — don't invent a new one).
2. Registers + authenticates as a non-admin user (owner).
3. Uses a second entity (could be the same user's entity or an admin-created corporation) as the subject.
4. POSTs `/v1/tags` with `{subject, purpose: "rating", value: "5", color: "#FF0000FF"}` — expects 201.
5. GETs `/v1/tags/{uuid}` — expects 200 with matching fields.
6. GETs `/v1/entities/{subject_uuid}/tags` — expects 200 with the new tag in the list.
7. PUTs `/v1/tags/{uuid}` with `{color: "#00FF00FF"}` — expects 200.
8. PUTs `/v1/tags/{uuid}` with `{purpose: "something"}` — expects 400.
9. DELETEs `/v1/tags/{uuid}` — expects 204.
10. GETs `/v1/tags/{uuid}` — expects 404.

Prefer the existing test harness for auth/user setup; don't duplicate a login stack.

If users-module doesn't currently run DB-backed integration tests (some projects only unit-test handlers), skip this task and document the skip in a `report.5.5.integration-test-skipped.md` explaining what was tried and why.

## How to verify

- `cd users-module/api && make test` exits 0 and the new test runs successfully.
- Audit log (`SELECT * FROM audit_log WHERE resource = 'tag' ORDER BY id DESC LIMIT 10`) shows entries for create/update/delete.
