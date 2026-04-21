# Phase 3 — API

## Goal

Implement the Go service + HTTP router for tags: tx-aware TagService with authorization and audit, chi subrouter with handlers, display renderer registration, and unit tests.

## Preconditions

- Phase 2 complete: `tags-module/model/db/` is generated and compiling.
- core-module/api is stable (audit.Writer, Principal/PrincipalExtractor, display.Registry all available).

## Outputs

- `tags-module/api/service/tag.go` — `TagServicer` interface + `TagService` implementation.
- `tags-module/api/service/service.go` — aggregate `Services` struct + `New` constructor.
- `tags-module/api/service/display.go` — `RegisterBuiltins(reg *display.Registry, q tagsdb.Querier)`.
- `tags-module/api/service/errors.go` — shared sentinel errors (or re-export core's).
- `tags-module/api/httpapi/router.go` — `NewRouter(Deps) chi.Router`.
- `tags-module/api/httpapi/tags.go` — handlers for `/tags` and `/tags/{uuid}`.
- `tags-module/api/httpapi/subject_tags.go` — handler for `/entities/{uuid}/tags`.
- `tags-module/api/httpapi/response.go` — local JSON/error helpers (may import core's if trivially re-exported).
- Unit tests alongside each source file.
- `tags-module/api/openapi.fragment.yaml` — OpenAPI 3.0.3 fragment for the tag endpoints.

## Routes exposed (final)

| Verb   | Path                               | Purpose                                          | Auth                             |
|--------|-----------------------------------|--------------------------------------------------|----------------------------------|
| POST   | `/tags`                           | Create a tag                                     | authenticated                    |
| GET    | `/tags`                           | Search tags; requires ≥1 of owner/subject        | authenticated                    |
| GET    | `/tags/{uuid}`                    | Read tag by entity UUID                          | owner or subject or admin; else 404 |
| PUT    | `/tags/{uuid}`                    | Update tag `color` (only field accepted)         | owner or admin                   |
| DELETE | `/tags/{uuid}`                    | Hard delete tag (entity + row in one tx)         | owner or admin                   |
| GET    | `/entities/{uuid}/tags`           | List tags whose subject is the given entity      | authenticated                    |

(All mount under whatever prefix the consumer chooses — typically `/v1`.)

## Hard rules

- **No dependency on users-module.** tags-module/api imports only tags-model, core-model, core-api, chi, pgx, uuid, stdlib.
- **No auth middleware in tags-module.** Consumer middleware populates context; tags-module's handlers read via `PrincipalExtractor`.
- **Services accept a `tagsdb.Querier`** so consumers can compose multi-module transactions.
- **Audit is injected** via `audit.Writer`.
- **Authorization in service layer**, not only in handlers. Handlers may short-circuit on obvious failures (unauthenticated → 401) but the authz matrix (owner/subject/admin/404) lives in the service.
- **Create semantics:** owner is set server-side to `actor.EntityID`. Any client-supplied owner field is ignored (or rejected — rejection preferred).
- **PUT body:** the handler should reject any body keys other than `color` (return 400). Do NOT compare old/new values to infer immutability — simply accept only `color` in the decoded shape.
- **DELETE:** tx-open, delete `tags` row, delete `entities` row, commit. Audit write logs the `before` snapshot.
- **Search authz:** `GET /tags` is filtered to the authz-visible set:
  - admins: all results pass.
  - non-admins: row passes if `owner_id == actor.EntityID` OR `subject_id == actor.EntityID`.
  - Enforce ≥1 of owner/subject in the query string (400 otherwise). Authz post-filter still applies even when subject/owner is provided.

## Authorization matrix (recap)

| Operation              | Admin | Owner | Subject | Other |
|------------------------|-------|-------|---------|-------|
| Create                 | ✓     | ✓ (becomes owner) | ✓ (becomes owner) | ✓ (becomes owner) |
| GET detail             | ✓     | ✓     | ✓       | 404   |
| GET list /tags         | ✓ all | ✓ own | ✓ as subject | 400 if no filter; else filtered |
| GET /entities/{uuid}/tags | ✓   | own   | self-subject    | subject=other → filter out tags neither owned nor self-subject (so likely empty) |
| Update (color)         | ✓     | ✓     | ✗ (403) | 404   |
| Delete                 | ✓     | ✓     | ✗ (403) | 404   |

The 403-vs-404 distinction: a non-owner non-admin caller who happens to be the subject gets 403 on mutations (they know the tag exists from GET). A caller who is neither owner nor subject nor admin gets 404 on everything (don't leak existence).

## Tasks

- 3.1 TagService implementation
- 3.2 Display renderer registration
- 3.3 httpapi router + handlers
- 3.4 Unit tests
- 3.5 OpenAPI fragment
- 3.6 Final build/test gate

## How to verify

- `cd tags-module/api && make test` exits 0 with non-trivial coverage (target ≥70% on service, ≥60% on httpapi — mirror core-module/api's numbers).
- `go build ./...` and `go vet ./...` exit 0.
- `npx @redocly/cli lint tags-module/api/openapi.fragment.yaml` exits 0.

## Notes

- Keep handlers thin: decode → call service → encode.
- Use `pgx.Tx` via a `pool.BeginTx` / `coredb.New(tx)` pattern consistent with users-module's existing pattern for cross-table atomic operations in Create and Delete.
- "Can tag any entity you can read" — for v1 this is implemented as "any authenticated user can set subject_id to any existing entity id". If/when a read-authz layer exists on entities, Create can check `service.Entity.CanRead(actor, subject)` or similar. For now just validate that subject exists (`GetEntityByUUID` resolves to internal id; 404 if not found).
