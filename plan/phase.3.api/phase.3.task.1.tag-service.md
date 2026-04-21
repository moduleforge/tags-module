# Phase 3, Task 1 — TagService implementation

## Context

The service layer is the heart of the module: CRUD + authz + audit. Handlers are thin wrappers over it.

## Acceptance

Create `tags-module/api/service/tag.go` containing:

### Types

```go
type CreateTagInput struct {
    SubjectEntityUUID uuid.UUID
    Purpose           string
    Value             string
    Color             *string // optional; must match "#RRGGBBAA" if set
}

type UpdateTagInput struct {
    Color *string // only mutable field; nil = no change
}

type SearchTagsFilter struct {
    OwnerEntityUUID   *uuid.UUID
    SubjectEntityUUID *uuid.UUID
    Purpose           *string
    Value             *string
}

type Tag struct {
    EntityUUID    uuid.UUID
    OwnerUUID     uuid.UUID
    SubjectUUID   uuid.UUID
    Purpose       string
    Value         string
    Color         *string
    CreatedAt     time.Time
    UpdatedAt     time.Time
}

type TagServicer interface {
    Create(ctx, q, actor Principal, in CreateTagInput) (Tag, error)
    GetByUUID(ctx, q, actor Principal, entityUUID uuid.UUID) (Tag, error)
    Search(ctx, q, actor Principal, filter SearchTagsFilter) ([]Tag, error)
    ListBySubject(ctx, q, actor Principal, subjectUUID uuid.UUID, purposeFilter *string) ([]Tag, error)
    UpdateByUUID(ctx, q, actor Principal, entityUUID uuid.UUID, in UpdateTagInput) (Tag, error)
    DeleteByUUID(ctx, q, actor Principal, entityUUID uuid.UUID, pool *pgxpool.Pool) error
}
```

(Exact signature details — receiver, package of `Principal`, etc. — follow core-module/api/service conventions. `q` is `tagsdb.Querier`. `Principal` and `audit.Writer` come from core.)

### Rules

- **Create**
  - Trim + validate `Purpose` and `Value` non-empty, length ≤ 512.
  - If `Color` set, validate the regex `^#[0-9A-Fa-f]{8}$` in Go too (defense in depth; DB also checks).
  - Resolve subject UUID → internal entity id (`GetEntityByUUID` from coredb). 404 if missing.
  - Resolve `tag` type id (`GetTypeBySlug`).
  - Open a tx (caller can pass one; if not, pool-based helper opens it): INSERT entity with `fundamental_type_id = tag`, INSERT tags row, commit.
  - Audit write with `op="create"`, `resource="tag"`, target = new entity id, before=nil, after={uuid, owner, subject, purpose, value, color}.
  - Return hydrated Tag.
- **GetByUUID**
  - Resolve entity by UUID. 404 if missing or its fundamental type isn't `tag`.
  - Load tags row.
  - Authz: admin OR `actor.EntityID == tag.OwnerID` OR `actor.EntityID == tag.SubjectID` → OK. Otherwise return `ErrNotFound` (404).
  - Return Tag.
- **Search**
  - Require actor is authenticated (non-nil).
  - Validate filter: at least one of OwnerEntityUUID / SubjectEntityUUID is set — else `ErrInvalidInput`.
  - Resolve provided UUIDs to internal ids.
  - Call `SearchTags` on the querier.
  - If not admin: post-filter to rows where `owner_id == actor.EntityID` OR `subject_id == actor.EntityID`.
  - Return list (possibly empty).
- **ListBySubject**
  - Resolve subject UUID; 404 if unknown.
  - Fetch via `ListTagsBySubjectEntityID` with optional purpose filter.
  - Authz post-filter: if not admin and `actor.EntityID != subject.ID`, keep only tags the actor owns. (Subject sees all tags about themselves; third parties see only tags they own on that subject.)
  - Return list.
- **UpdateByUUID**
  - Resolve + load tag.
  - Authz: admin OR `actor.EntityID == tag.OwnerID`. Subject-only → 403. Other → 404.
  - If `in.Color != nil`, validate regex, update via `UpdateTagColor`.
  - Audit write with before/after snapshots showing the color change.
  - Return hydrated Tag.
- **DeleteByUUID**
  - Resolve + load tag.
  - Authz: admin OR owner. Else 404 (subject can't delete; nor can strangers).
  - Open a tx: DELETE from tags, DELETE from entities. Commit.
  - Audit write with `op="delete"`, before snapshot, after=nil.
  - Returns nil or `ErrNotFound` / `ErrForbidden`.

### Error surface

Use the same sentinel errors pattern core-module uses (`ErrNotFound`, `ErrForbidden`, `ErrInvalidInput`). If core exports them, import and reuse. Otherwise define locally with identical semantics.

### Services aggregate

Create `tags-module/api/service/service.go`:

```go
type Services struct {
    Tag TagServicer
    q   tagsdb.Querier
}

func New(q tagsdb.Querier, aw audit.Writer) *Services {
    return &Services{
        Tag: &TagService{aw: aw},
        q:   q,
    }
}

func (s *Services) Querier() tagsdb.Querier { return s.q }
```

Plus a `pool` field on Deps or Services if needed by handlers to open tx (see how core-module does this).

## How to verify

- `cd tags-module/api && go build ./...` exits 0.
- `go vet ./...` exits 0.

## Reference

- `core-module/api/service/service_account.go` — most similar shape (create via single subtype).
- `core-module/api/service/natural_person.go` — more elaborate CRUD, useful for UpdateByUUID pattern.
- `core-module/api/service/errors.go` — error-sentinel pattern.
