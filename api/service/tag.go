package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/moduleforge/core-api/audit"
	coreservice "github.com/moduleforge/core-api/service"
	coredb "github.com/moduleforge/core-model/db"
	tagsdb "github.com/moduleforge/tags-model/db"
)

// colorRe validates the color format "#RRGGBBAA" (8 hex digits after #).
var colorRe = regexp.MustCompile(`^#[0-9A-Fa-f]{8}$`)

// pgUniqueViolation is the Postgres error code for unique_violation.
const pgUniqueViolation = "23505"

// TxBeginner abstracts transaction creation for testability.
// *pgxpool.Pool satisfies this interface.
type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// CreateTagInput carries the fields required to create a tag.
type CreateTagInput struct {
	SubjectEntityUUID uuid.UUID
	Purpose           string
	Value             string
	Color             *string // optional; must match "#RRGGBBAA" if set
}

// UpdateTagInput carries the only mutable field on a tag.
type UpdateTagInput struct {
	// Color = nil clears the color (sets DB column to NULL).
	// Color = &"" is invalid (rejected by regex).
	// Color = &"#RRGGBBAA" sets the color.
	Color *string
}

// SearchTagsFilter filters the tag search. At least one of OwnerEntityUUID /
// SubjectEntityUUID is required; the rest are optional.
type SearchTagsFilter struct {
	OwnerEntityUUID   *uuid.UUID
	SubjectEntityUUID *uuid.UUID
	Purpose           *string
	Value             *string
}

// Tag is the service-layer representation of a tag, using public UUIDs.
type Tag struct {
	EntityUUID  uuid.UUID
	OwnerUUID   uuid.UUID
	SubjectUUID uuid.UUID
	Purpose     string
	Value       string
	Color       *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TagServicer defines tag CRUD operations available to httpapi handlers.
type TagServicer interface {
	Create(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, tx TxBeginner, in CreateTagInput) (Tag, error)
	GetByUUID(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, entityUUID uuid.UUID) (Tag, error)
	Search(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, filter SearchTagsFilter) ([]Tag, error)
	ListBySubject(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, subjectUUID uuid.UUID, purposeFilter *string) ([]Tag, error)
	UpdateByUUID(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, entityUUID uuid.UUID, in UpdateTagInput) (Tag, error)
	DeleteByUUID(ctx context.Context, coreQ coredb.Querier, tagQ tagsdb.Querier, actor coreservice.Principal, entityUUID uuid.UUID, tx TxBeginner) error
}

// TagService implements TagServicer with audit logging.
type TagService struct {
	aw audit.Writer
}

// Compile-time assertion.
var _ TagServicer = (*TagService)(nil)

// Create inserts an entity row and a tags row in a single transaction.
// Owner is always set server-side to actor.EntityID — any client-supplied
// owner is ignored. The subject entity must already exist.
func (s *TagService) Create(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	txb TxBeginner,
	in CreateTagInput,
) (Tag, error) {
	// Input validation.
	in.Purpose = strings.TrimSpace(in.Purpose)
	in.Value = strings.TrimSpace(in.Value)
	if in.Purpose == "" {
		return Tag{}, fmt.Errorf("%w: purpose is required", ErrInvalidInput)
	}
	if len(in.Purpose) > 512 {
		return Tag{}, fmt.Errorf("%w: purpose exceeds 512 characters", ErrInvalidInput)
	}
	if in.Value == "" {
		return Tag{}, fmt.Errorf("%w: value is required", ErrInvalidInput)
	}
	if len(in.Value) > 512 {
		return Tag{}, fmt.Errorf("%w: value exceeds 512 characters", ErrInvalidInput)
	}
	if in.Color != nil && !colorRe.MatchString(*in.Color) {
		return Tag{}, fmt.Errorf("%w: color must match #RRGGBBAA", ErrInvalidInput)
	}

	// Resolve subject entity.
	subjectEntity, err := coreQ.GetEntityByUUID(ctx, in.SubjectEntityUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tag{}, fmt.Errorf("%w: subject entity not found", ErrNotFound)
		}
		return Tag{}, fmt.Errorf("tag.Create resolve subject: %w", err)
	}

	// Resolve owner entity (actor).
	ownerEntity, err := coreQ.GetEntityByID(ctx, actor.EntityID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tag{}, fmt.Errorf("%w: owner entity not found", ErrNotFound)
		}
		return Tag{}, fmt.Errorf("tag.Create resolve owner: %w", err)
	}

	// Resolve the type ID for 'tag'.
	tagType, err := coreQ.GetTypeBySlug(ctx, "tag")
	if err != nil {
		return Tag{}, fmt.Errorf("tag.Create resolve type: %w", err)
	}

	// Open transaction.
	dbTx, err := txb.Begin(ctx)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.Create begin tx: %w", err)
	}
	defer dbTx.Rollback(ctx) //nolint:errcheck

	txCoreQ := coredb.New(dbTx)
	txTagQ := tagsdb.New(dbTx)

	// Insert entity row.
	entity, err := txCoreQ.CreateEntity(ctx, tagType.ID)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.Create entity: %w", err)
	}

	// Build color param.
	colorParam := pgtype.Text{}
	if in.Color != nil {
		colorParam = pgtype.Text{String: *in.Color, Valid: true}
	}

	// Insert tags row.
	tag, err := txTagQ.CreateTag(ctx, tagsdb.CreateTagParams{
		EntityID:  entity.ID,
		OwnerID:   actor.EntityID,
		SubjectID: subjectEntity.ID,
		Purpose:   in.Purpose,
		Value:     in.Value,
		Color:     colorParam,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return Tag{}, fmt.Errorf("%w: tag already exists", ErrConflict)
		}
		return Tag{}, fmt.Errorf("tag.Create insert: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return Tag{}, fmt.Errorf("tag.Create commit: %w", err)
	}

	result := hydrateTag(entity.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, tag)

	eid := entity.ID
	_ = s.aw.Write(ctx, "create", "tag", &eid, nil, tagSnapshot(result))

	return result, nil
}

// GetByUUID resolves a tag by entity UUID and enforces read authz.
func (s *TagService) GetByUUID(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	entityUUID uuid.UUID,
) (Tag, error) {
	row, err := tagQ.GetTagByEntityUUID(ctx, entityUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tag{}, ErrNotFound
		}
		return Tag{}, fmt.Errorf("tag.GetByUUID: %w", err)
	}

	// Authz: admin OR owner OR subject → OK, else 404.
	if !actor.IsAdmin && actor.EntityID != row.OwnerID && actor.EntityID != row.SubjectID {
		return Tag{}, ErrNotFound
	}

	// Resolve UUIDs for owner and subject.
	ownerEntity, err := coreQ.GetEntityByID(ctx, row.OwnerID)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.GetByUUID resolve owner: %w", err)
	}
	subjectEntity, err := coreQ.GetEntityByID(ctx, row.SubjectID)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.GetByUUID resolve subject: %w", err)
	}

	return hydrateTag(row.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, tagFromUUIDRow(row)), nil
}

// Search finds tags matching filter, post-filtered by authz.
func (s *TagService) Search(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	filter SearchTagsFilter,
) ([]Tag, error) {
	if filter.OwnerEntityUUID == nil && filter.SubjectEntityUUID == nil {
		return nil, fmt.Errorf("%w: at least one of owner or subject is required", ErrInvalidInput)
	}

	params := tagsdb.SearchTagsParams{}

	if filter.OwnerEntityUUID != nil {
		ownerEntity, err := coreQ.GetEntityByUUID(ctx, *filter.OwnerEntityUUID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []Tag{}, nil // no results; don't leak non-existence
			}
			return nil, fmt.Errorf("tag.Search resolve owner: %w", err)
		}
		params.OwnerID = pgtype.Int8{Int64: ownerEntity.ID, Valid: true}
	}

	if filter.SubjectEntityUUID != nil {
		subjectEntity, err := coreQ.GetEntityByUUID(ctx, *filter.SubjectEntityUUID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return []Tag{}, nil
			}
			return nil, fmt.Errorf("tag.Search resolve subject: %w", err)
		}
		params.SubjectID = pgtype.Int8{Int64: subjectEntity.ID, Valid: true}
	}

	if filter.Purpose != nil {
		params.Purpose = pgtype.Text{String: *filter.Purpose, Valid: true}
	}
	if filter.Value != nil {
		params.Value = pgtype.Text{String: *filter.Value, Valid: true}
	}

	rows, err := tagQ.SearchTags(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("tag.Search query: %w", err)
	}

	// Post-filter for non-admins.
	var result []Tag
	for _, row := range rows {
		if !actor.IsAdmin && actor.EntityID != row.OwnerID && actor.EntityID != row.SubjectID {
			continue
		}
		ownerEntity, err := coreQ.GetEntityByID(ctx, row.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("tag.Search resolve owner uuid: %w", err)
		}
		subjectEntity, err := coreQ.GetEntityByID(ctx, row.SubjectID)
		if err != nil {
			return nil, fmt.Errorf("tag.Search resolve subject uuid: %w", err)
		}
		entityRow, err := coreQ.GetEntityByID(ctx, row.EntityID)
		if err != nil {
			return nil, fmt.Errorf("tag.Search resolve entity uuid: %w", err)
		}
		result = append(result, hydrateTag(entityRow.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, row))
	}

	if result == nil {
		result = []Tag{}
	}
	return result, nil
}

// ListBySubject returns all tags targeting a given subject entity, filtered by authz.
func (s *TagService) ListBySubject(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	subjectUUID uuid.UUID,
	purposeFilter *string,
) ([]Tag, error) {
	subjectEntity, err := coreQ.GetEntityByUUID(ctx, subjectUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("tag.ListBySubject resolve subject: %w", err)
	}

	purposeParam := pgtype.Text{}
	if purposeFilter != nil {
		purposeParam = pgtype.Text{String: *purposeFilter, Valid: true}
	}

	rows, err := tagQ.ListTagsBySubjectEntityID(ctx, tagsdb.ListTagsBySubjectEntityIDParams{
		SubjectID: subjectEntity.ID,
		Purpose:   purposeParam,
	})
	if err != nil {
		return nil, fmt.Errorf("tag.ListBySubject query: %w", err)
	}

	var result []Tag
	for _, row := range rows {
		// Authz post-filter:
		//   admin → all pass
		//   actor == subject → all pass (subject sees all tags about themselves)
		//   otherwise → only rows the actor owns
		if !actor.IsAdmin && actor.EntityID != subjectEntity.ID && actor.EntityID != row.OwnerID {
			continue
		}
		ownerEntity, err := coreQ.GetEntityByID(ctx, row.OwnerID)
		if err != nil {
			return nil, fmt.Errorf("tag.ListBySubject resolve owner uuid: %w", err)
		}
		entityRow, err := coreQ.GetEntityByID(ctx, row.EntityID)
		if err != nil {
			return nil, fmt.Errorf("tag.ListBySubject resolve entity uuid: %w", err)
		}
		result = append(result, hydrateTag(entityRow.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, row))
	}

	if result == nil {
		result = []Tag{}
	}
	return result, nil
}

// UpdateByUUID updates the color of a tag identified by entity UUID.
func (s *TagService) UpdateByUUID(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	entityUUID uuid.UUID,
	in UpdateTagInput,
) (Tag, error) {
	row, err := tagQ.GetTagByEntityUUID(ctx, entityUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Tag{}, ErrNotFound
		}
		return Tag{}, fmt.Errorf("tag.UpdateByUUID fetch: %w", err)
	}

	// Authz: admin OR owner → OK; subject → 403; other → 404.
	if !actor.IsAdmin && actor.EntityID != row.OwnerID {
		if actor.EntityID == row.SubjectID {
			return Tag{}, ErrForbidden
		}
		return Tag{}, ErrNotFound
	}

	// Validate non-nil color before writing.
	if in.Color != nil && !colorRe.MatchString(*in.Color) {
		return Tag{}, fmt.Errorf("%w: color must match #RRGGBBAA", ErrInvalidInput)
	}

	before := colorSnapshot(row.Color)

	colorParam := pgtype.Text{}
	if in.Color != nil {
		colorParam = pgtype.Text{String: *in.Color, Valid: true}
	}

	updated, err := tagQ.UpdateTagColor(ctx, tagsdb.UpdateTagColorParams{
		EntityID: row.EntityID,
		Color:    colorParam,
	})
	if err != nil {
		return Tag{}, fmt.Errorf("tag.UpdateByUUID update: %w", err)
	}

	ownerEntity, err := coreQ.GetEntityByID(ctx, row.OwnerID)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.UpdateByUUID resolve owner: %w", err)
	}
	subjectEntity, err := coreQ.GetEntityByID(ctx, row.SubjectID)
	if err != nil {
		return Tag{}, fmt.Errorf("tag.UpdateByUUID resolve subject: %w", err)
	}

	result := hydrateTag(row.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, updated)

	eid := row.EntityID
	_ = s.aw.Write(ctx, "update", "tag", &eid, before, colorSnapshot(updated.Color))

	return result, nil
}

// DeleteByUUID removes a tag. The tags row is deleted; the entity row is
// archived (core-model exposes ArchiveEntity, not a hard DELETE).
func (s *TagService) DeleteByUUID(
	ctx context.Context,
	coreQ coredb.Querier,
	tagQ tagsdb.Querier,
	actor coreservice.Principal,
	entityUUID uuid.UUID,
	txb TxBeginner,
) error {
	row, err := tagQ.GetTagByEntityUUID(ctx, entityUUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("tag.DeleteByUUID fetch: %w", err)
	}

	// Authz: admin OR owner → OK; subject → 403; other → 404.
	if !actor.IsAdmin && actor.EntityID != row.OwnerID {
		if actor.EntityID == row.SubjectID {
			return ErrForbidden
		}
		return ErrNotFound
	}

	// Capture before snapshot for audit.
	ownerEntity, err := coreQ.GetEntityByID(ctx, row.OwnerID)
	if err != nil {
		return fmt.Errorf("tag.DeleteByUUID resolve owner: %w", err)
	}
	subjectEntity, err := coreQ.GetEntityByID(ctx, row.SubjectID)
	if err != nil {
		return fmt.Errorf("tag.DeleteByUUID resolve subject: %w", err)
	}
	beforeSnapshot := tagSnapshot(hydrateTag(row.Uuid, ownerEntity.Uuid, subjectEntity.Uuid, tagFromUUIDRow(row)))

	// Open transaction: delete tags row, archive entity row.
	dbTx, err := txb.Begin(ctx)
	if err != nil {
		return fmt.Errorf("tag.DeleteByUUID begin tx: %w", err)
	}
	defer dbTx.Rollback(ctx) //nolint:errcheck

	txTagQ := tagsdb.New(dbTx)
	txCoreQ := coredb.New(dbTx)

	if err := txTagQ.DeleteTag(ctx, row.EntityID); err != nil {
		return fmt.Errorf("tag.DeleteByUUID delete tag: %w", err)
	}

	if err := txCoreQ.ArchiveEntity(ctx, entityUUID); err != nil {
		return fmt.Errorf("tag.DeleteByUUID archive entity: %w", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return fmt.Errorf("tag.DeleteByUUID commit: %w", err)
	}

	eid := row.EntityID
	_ = s.aw.Write(ctx, "delete", "tag", &eid, beforeSnapshot, nil)

	return nil
}

// --- helpers ---

// hydrateTag converts internal DB rows to the service Tag type.
func hydrateTag(entityUUID, ownerUUID, subjectUUID uuid.UUID, t tagsdb.Tag) Tag {
	tag := Tag{
		EntityUUID:  entityUUID,
		OwnerUUID:   ownerUUID,
		SubjectUUID: subjectUUID,
		Purpose:     t.Purpose,
		Value:       t.Value,
	}
	if t.Color.Valid {
		c := t.Color.String
		tag.Color = &c
	}
	if t.CreatedAt.Valid {
		tag.CreatedAt = t.CreatedAt.Time
	}
	if t.UpdatedAt.Valid {
		tag.UpdatedAt = t.UpdatedAt.Time
	}
	return tag
}

// tagFromUUIDRow converts a GetTagByEntityUUIDRow back into a Tag row shape for hydrateTag.
func tagFromUUIDRow(r tagsdb.GetTagByEntityUUIDRow) tagsdb.Tag {
	return tagsdb.Tag{
		EntityID:  r.EntityID,
		OwnerID:   r.OwnerID,
		SubjectID: r.SubjectID,
		Purpose:   r.Purpose,
		Value:     r.Value,
		Color:     r.Color,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}
}

// tagSnapshot builds an audit snapshot map from a Tag.
func tagSnapshot(t Tag) map[string]any {
	return map[string]any{
		"uuid":         t.EntityUUID.String(),
		"owner_uuid":   t.OwnerUUID.String(),
		"subject_uuid": t.SubjectUUID.String(),
		"purpose":      t.Purpose,
		"value":        t.Value,
		"color":        t.Color,
	}
}

// colorSnapshot builds a before/after snapshot for color-only changes.
func colorSnapshot(c pgtype.Text) map[string]any {
	if c.Valid {
		return map[string]any{"color": c.String}
	}
	return map[string]any{"color": nil}
}
