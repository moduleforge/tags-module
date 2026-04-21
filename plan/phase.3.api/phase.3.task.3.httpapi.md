# Phase 3, Task 3 — httpapi router + handlers

## Context

Thin chi subrouter mountable at any prefix, exposing the six routes in `phase.3.api.md`. Follows core-module/api/httpapi conventions exactly.

## Acceptance

Create these files under `tags-module/api/httpapi/`:

### router.go

```go
type Deps struct {
    Pool      *pgxpool.Pool
    Services  *service.Services
    Audit     audit.Writer
    Principal service.PrincipalExtractor
    Logger    *slog.Logger
}

func NewRouter(deps Deps) chi.Router { … }
```

Routes registered:

```
POST   /tags
GET    /tags
GET    /tags/{uuid}
PUT    /tags/{uuid}
DELETE /tags/{uuid}
GET    /entities/{uuid}/tags
```

### tags.go

Handlers for `/tags` and `/tags/{uuid}`:

- `handleCreateTag` — decode body `{subject, purpose, value, color?}` (where `subject` is the subject entity UUID); extract principal; call `services.Tag.Create`; write 201 + full tag.
- `handleSearchTags` — parse query string (`owner`, `subject`, `purpose`, `value` all optional); extract principal; call `services.Tag.Search`; write 200 + array. 400 if no owner/subject filter. 400 with clear body if owner/subject aren't valid UUIDs.
- `handleGetTag` — path param uuid; call `services.Tag.GetByUUID`; 200 + tag; 404 on ErrNotFound.
- `handlePutTag` — decode body, reject any key other than `color`; call `services.Tag.UpdateByUUID`; 200 + updated tag.
- `handleDeleteTag` — call `services.Tag.DeleteByUUID`; 204 No Content on success.

### subject_tags.go

- `handleSubjectTags` — path param uuid = subject entity UUID; optional `purpose` query filter; call `services.Tag.ListBySubject`; return `{tags: [...]}` (wrap in an object for forward compat).

### response.go

Local JSON + error helpers, or thin re-export of core-module's if they're exported. Keep it tiny — decode/encode/writeError should be <40 lines total.

### Principal extraction

Use `deps.Principal.FromContext(r.Context())` (the same pattern core-module uses). If `ok == false`, return 401.

### Body-key rejection for PUT

Implement as a "strict" decoder: use `json.NewDecoder(r.Body).DisallowUnknownFields()` and decode into a typed struct with only `Color *string`. Any other key causes decode failure → 400.

## How to verify

- `go build ./...` exits 0 in tags-module/api.
- `go vet ./...` exits 0.
- All route wiring can be sanity-checked by starting a minimal main in a _test.go (handled in Task 3.4).

## Reference

- `core-module/api/httpapi/router.go` — the Deps shape.
- `core-module/api/httpapi/service_accounts.go` — simplest CRUD handler set to model up.
- `core-module/api/httpapi/response.go` — response helpers.
