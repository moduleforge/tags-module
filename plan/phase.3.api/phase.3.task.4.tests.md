# Phase 3, Task 4 — Unit tests

## Context

Cover the service-layer authz matrix, immutability error paths, and the handler routing/codes. Use core-module/api tests as reference for mock shape.

## Acceptance

Write table-driven tests covering at minimum:

### Service layer (`tags-module/api/service/tag_test.go`)

- Create: happy path; Create rejects missing purpose; Create rejects invalid color hex.
- Get: admin sees anything; owner sees own; subject sees as-subject; other → `ErrNotFound`.
- Search: 400-analog (`ErrInvalidInput`) when neither owner nor subject supplied; admin sees all; non-admin sees only those where owner==actor OR subject==actor.
- ListBySubject: subject sees all tags about themselves; third parties see only tags they own.
- Update: owner can change color; admin can change color; subject cannot (`ErrForbidden`); stranger → `ErrNotFound`.
- Delete: owner/admin OK; subject → 404 (don't leak: treat as ErrNotFound at service boundary since subject has no mutation right here; pick 404 vs 403 per the table in phase.3.api.md — subject gets 403 for mutations, so use ErrForbidden).
- Uniqueness collision: simulate `pgconn.PgError` with unique_violation code → TagService returns wrapped `ErrConflict` (add this sentinel if not already present).

### Display renderer (`tags-module/api/service/display_test.go`)

- Register into a fresh Registry; resolve a tag entity via the Registry; confirm `FieldName` returns `"<purpose>:<value>"`.

### Handlers (`tags-module/api/httpapi/handlers_test.go`)

Use `httptest.NewServer` + `chi.NewRouter` with mocked services.

- POST /tags: 201 on success; 401 when unauthenticated; 400 on missing subject; 400 on bad color.
- GET /tags: 400 when neither owner nor subject; 200 + [] when no matches.
- GET /tags/{uuid}: 200 for authorized; 404 for unauthorized-stranger.
- PUT /tags/{uuid}: 200 happy; 400 when body includes `purpose`; 400 when body unknown key; 403 when subject tries.
- DELETE /tags/{uuid}: 204 on success; 403 for subject; 404 for stranger.
- GET /entities/{uuid}/tags: 200 + `{tags: [...]}`; 404 when subject entity unknown.

### Mocks

Mirror `core-module/api/service/mock_test.go` and `core-module/api/httpapi/mock_test.go` — hand-rolled mocks backed by maps are sufficient; don't introduce a mocking library.

## How to verify

- `cd tags-module/api && make test` exits 0.
- Coverage ≥ 70% on service, ≥ 60% on httpapi.
- `go vet ./...` exits 0.

## Reference

- `core-module/api/service/service_account_test.go`
- `core-module/api/httpapi/handlers_test.go`
- `core-module/api/httpapi/mock_test.go`
