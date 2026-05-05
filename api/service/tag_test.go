package service

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/moduleforge/core-api/observer"
	coreservice "github.com/moduleforge/core-api/service"
	coredb "github.com/moduleforge/core-model/db"
	tagsdb "github.com/moduleforge/tags-model/db"
)

// adminPrincipal returns an admin Principal.
func adminPrincipal(entityID int64) coreservice.Principal {
	return coreservice.Principal{EntityID: entityID, UserID: entityID, IsAdmin: true}
}

// userPrincipal returns a non-admin Principal.
func userPrincipal(entityID int64) coreservice.Principal {
	return coreservice.Principal{EntityID: entityID, UserID: entityID, IsAdmin: false}
}

// buildService constructs a TagService with allow-all authz and a no-op observer,
// wired to the provided mocks. Returns a recording observer for assertion.
func buildService(coreQ *mockCoreQuerier, tagQ *mockTagQuerier) (*TagService, *recordingObserver) {
	rec := &recordingObserver{}
	obs := observer.NewObserverGroup(rec)
	svc := &TagService{
		db:             newFakeDB(),
		az:             allowAllAuthz{},
		obs:            obs,
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	return svc, rec
}

// buildServiceWithTag sets up a mock environment with one seeded tag and returns
// the queriers, the entity UUIDs, and the service.
func buildServiceWithTag(t *testing.T) (svc *TagService, coreQ *mockCoreQuerier, tagQ *mockTagQuerier, tagUUID uuid.UUID, ownerID, subjectID, entityID int64) {
	t.Helper()
	coreQ = newMockCoreQuerier()
	tagQ = newMockTagQuerier()

	svc, _ = buildService(coreQ, tagQ)

	// Seed owner, subject entities in coreQ.
	_, ownerID = coreQ.seedEntity("natural_person")
	_, subjectID = coreQ.seedEntity("natural_person")

	// Seed tag entity.
	_, entityID = coreQ.seedEntity("tag")

	// Link entity UUID → tag ID mapping in tagQ.
	tagUUID, _ = tagQ.seedTag(entityID, ownerID, subjectID, "label", "urgent", nil)

	// Ensure coreQ knows about the tagUUID → entityID mapping.
	row := coreQ.entitiesByID[entityID]
	coreQ.entities[tagUUID] = row

	return
}

// --- Create validation tests ---

func TestTagService_Create_MissingPurpose(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := userPrincipal(1)

	_, err := svc.Create(context.Background(), coreQ, actor, CreateTagInput{
		SubjectEntityUUID: uuid.New(),
		Purpose:           "",
		Value:             "x",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("want ErrInvalidInput, got %v", err)
	}
}

func TestTagService_Create_MissingValue(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := userPrincipal(1)

	_, err := svc.Create(context.Background(), coreQ, actor, CreateTagInput{
		SubjectEntityUUID: uuid.New(),
		Purpose:           "label",
		Value:             "",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("want ErrInvalidInput, got %v", err)
	}
}

func TestTagService_Create_InvalidColor(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := userPrincipal(1)
	badColor := "#ZZZZZZZZ"

	_, err := svc.Create(context.Background(), coreQ, actor, CreateTagInput{
		SubjectEntityUUID: uuid.New(),
		Purpose:           "label",
		Value:             "v",
		Color:             &badColor,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("want ErrInvalidInput, got %v", err)
	}
}

func TestTagService_Create_SubjectNotFound(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)

	// Seed owner entity.
	_, ownerID := coreQ.seedEntity("natural_person")
	actor := userPrincipal(ownerID)

	// Unknown subject UUID.
	_, err := svc.Create(context.Background(), coreQ, actor, CreateTagInput{
		SubjectEntityUUID: uuid.New(), // not seeded
		Purpose:           "label",
		Value:             "v",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

// --- Authorize-denied tests ---

func TestTagService_Create_AuthorizeDenied(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	authzErr := errors.New("not authorized")
	svc := &TagService{
		db:             newFakeDB(),
		az:             denyAllAuthz{err: authzErr},
		obs:            observer.NewObserverGroup(),
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := userPrincipal(1)

	_, err := svc.Create(context.Background(), coreQ, actor, CreateTagInput{
		SubjectEntityUUID: uuid.New(),
		Purpose:           "label",
		Value:             "v",
	})
	if !errors.Is(err, authzErr) {
		t.Errorf("want authz error, got %v", err)
	}
	// No DB calls should have been made: tagQ should be empty.
	if len(tagQ.tags) != 0 {
		t.Error("expected no tags created when authz denied")
	}
}

func TestTagService_Update_AuthorizeDenied(t *testing.T) {
	_, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	authzErr := errors.New("not authorized")
	svc := &TagService{
		db:             newFakeDB(),
		az:             denyAllAuthz{err: authzErr},
		obs:            observer.NewObserverGroup(),
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := adminPrincipal(999)
	color := "#FF0000FF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &color})
	if !errors.Is(err, authzErr) {
		t.Errorf("want authz error, got %v", err)
	}
}

func TestTagService_Delete_AuthorizeDenied(t *testing.T) {
	_, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	authzErr := errors.New("not authorized")
	svc := &TagService{
		db:             newFakeDB(),
		az:             denyAllAuthz{err: authzErr},
		obs:            observer.NewObserverGroup(),
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := adminPrincipal(999)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if !errors.Is(err, authzErr) {
		t.Errorf("want authz error, got %v", err)
	}
}

// --- In-tx observer error causes rollback ---

func TestTagService_Update_InTxObserverError_Propagates(t *testing.T) {
	_, coreQ, tagQ, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	obsErr := errors.New("observer failure")
	rec := &recordingObserver{observeErr: obsErr}
	obs := observer.NewObserverGroup(rec) // default policy = PolicyPropagate
	svc := &TagService{
		db:             newFakeDB(),
		az:             allowAllAuthz{},
		obs:            obs,
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := userPrincipal(ownerID)
	color := "#FF0000FF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &color})
	if !errors.Is(err, obsErr) {
		t.Errorf("want observer error propagated, got %v", err)
	}
	if len(rec.observeCalls) != 1 {
		t.Errorf("expected 1 observe call, got %d", len(rec.observeCalls))
	}
}

// --- Observer records correct op/resource ---

func TestTagService_Update_ObserverCalled(t *testing.T) {
	_, coreQ, tagQ, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	rec := &recordingObserver{}
	obs := observer.NewObserverGroup(rec)
	svc := &TagService{
		db:             newFakeDB(),
		az:             allowAllAuthz{},
		obs:            obs,
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := userPrincipal(ownerID)
	color := "#FF0000FF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &color})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.observeCalls) != 1 {
		t.Fatalf("expected 1 observe call, got %d", len(rec.observeCalls))
	}
	call := rec.observeCalls[0]
	if call.op != "update" {
		t.Errorf("op: want update, got %q", call.op)
	}
	if call.resource != "tag" {
		t.Errorf("resource: want tag, got %q", call.resource)
	}
}

// --- GetByUUID authz tests ---

func TestTagService_Get_AdminSeesAll(t *testing.T) {
	svc, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := adminPrincipal(999) // not owner or subject

	tag, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.EntityUUID != tagUUID {
		t.Errorf("uuid mismatch")
	}
}

func TestTagService_Get_OwnerSeesOwn(t *testing.T) {
	svc, coreQ, tagQ, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(ownerID)

	_, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Get_SubjectSeesAsSubject(t *testing.T) {
	svc, coreQ, tagQ, tagUUID, _, subjectID, _ := buildServiceWithTag(t)
	actor := userPrincipal(subjectID)

	_, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Get_OtherGets404(t *testing.T) {
	svc, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(777) // unrelated

	_, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, tagUUID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestTagService_Get_NotFound(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := adminPrincipal(1)

	_, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestTagService_Get_AuthorizeDenied(t *testing.T) {
	_, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	authzErr := errors.New("not authorized")
	svc := &TagService{
		db:             newFakeDB(),
		az:             denyAllAuthz{err: authzErr},
		obs:            observer.NewObserverGroup(),
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := adminPrincipal(999)

	_, err := svc.GetByUUID(context.Background(), coreQ, tagQ, actor, tagUUID)
	if !errors.Is(err, authzErr) {
		t.Errorf("want authz error, got %v", err)
	}
}

// --- Search authz tests ---

func TestTagService_Search_NoFilter_Returns400Analog(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := userPrincipal(1)

	_, err := svc.Search(context.Background(), coreQ, tagQ, actor, SearchTagsFilter{})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("want ErrInvalidInput, got %v", err)
	}
}

func TestTagService_Search_AdminSeesAll(t *testing.T) {
	svc, coreQ, tagQ, _, ownerID, _, _ := buildServiceWithTag(t)
	actor := adminPrincipal(999)

	ownerUUID := coreQ.entitiesByID[ownerID].Uuid
	tags, err := svc.Search(context.Background(), coreQ, tagQ, actor, SearchTagsFilter{
		OwnerEntityUUID: &ownerUUID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("want 1 tag, got %d", len(tags))
	}
}

func TestTagService_Search_NonAdminFilteredToOwned(t *testing.T) {
	svc, coreQ, tagQ, _, ownerID, subjectID, _ := buildServiceWithTag(t)

	ownerUUID := coreQ.entitiesByID[ownerID].Uuid

	// actor is the subject: should see tags where they're subject.
	actor := userPrincipal(subjectID)
	tags, err := svc.Search(context.Background(), coreQ, tagQ, actor, SearchTagsFilter{
		OwnerEntityUUID: &ownerUUID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("want 1 tag (as subject), got %d", len(tags))
	}

	// actor is unrelated: should see nothing.
	actor2 := userPrincipal(888)
	tags2, err := svc.Search(context.Background(), coreQ, tagQ, actor2, SearchTagsFilter{
		OwnerEntityUUID: &ownerUUID,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags2) != 0 {
		t.Errorf("want 0 tags for unrelated actor, got %d", len(tags2))
	}
}

// --- ListBySubject tests ---

func TestTagService_ListBySubject_SubjectSeesAll(t *testing.T) {
	svc, coreQ, tagQ, _, _, subjectID, _ := buildServiceWithTag(t)
	actor := userPrincipal(subjectID)

	subjectUUID := coreQ.entitiesByID[subjectID].Uuid
	tags, err := svc.ListBySubject(context.Background(), coreQ, tagQ, actor, subjectUUID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("want 1 tag, got %d", len(tags))
	}
}

func TestTagService_ListBySubject_ThirdPartySeeOnlyOwned(t *testing.T) {
	svc, coreQ, tagQ, _, ownerID, subjectID, _ := buildServiceWithTag(t)

	subjectUUID := coreQ.entitiesByID[subjectID].Uuid

	// Owner is also owner of the tag — should see their own tag.
	actor := userPrincipal(ownerID)
	tags, err := svc.ListBySubject(context.Background(), coreQ, tagQ, actor, subjectUUID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags) != 1 {
		t.Errorf("want 1 tag (owned), got %d", len(tags))
	}

	// Unrelated actor sees nothing.
	actor2 := userPrincipal(999)
	tags2, err := svc.ListBySubject(context.Background(), coreQ, tagQ, actor2, subjectUUID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tags2) != 0 {
		t.Errorf("want 0 tags for unrelated, got %d", len(tags2))
	}
}

func TestTagService_ListBySubject_NotFound(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := userPrincipal(1)

	_, err := svc.ListBySubject(context.Background(), coreQ, tagQ, actor, uuid.New(), nil)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

// --- UpdateByUUID authz tests ---

func TestTagService_Update_OwnerCanChangeColor(t *testing.T) {
	svc, coreQ, _, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(ownerID)
	newColor := "#FF0000FF"

	tag, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &newColor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Color == nil || *tag.Color != newColor {
		t.Errorf("want color %q, got %v", newColor, tag.Color)
	}
}

func TestTagService_Update_AdminCanChangeColor(t *testing.T) {
	svc, coreQ, _, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := adminPrincipal(999)
	newColor := "#00FF00FF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &newColor})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Update_SubjectGetsForbidden(t *testing.T) {
	svc, coreQ, _, tagUUID, _, subjectID, _ := buildServiceWithTag(t)
	actor := userPrincipal(subjectID)
	newColor := "#0000FFFF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &newColor})
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("want ErrForbidden, got %v", err)
	}
}

func TestTagService_Update_StrangerGets404(t *testing.T) {
	svc, coreQ, _, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(888)
	newColor := "#FFFFFFFF"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &newColor})
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestTagService_Update_InvalidColor(t *testing.T) {
	svc, coreQ, _, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(ownerID)
	bad := "notacolor"

	_, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: &bad})
	if !errors.Is(err, ErrInvalidInput) {
		t.Errorf("want ErrInvalidInput, got %v", err)
	}
}

func TestTagService_Update_NilColorClears(t *testing.T) {
	svc, coreQ, _, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(ownerID)

	tag, err := svc.UpdateByUUID(context.Background(), coreQ, actor, tagUUID, UpdateTagInput{Color: nil})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Color != nil {
		t.Errorf("want nil color after clear, got %v", tag.Color)
	}
}

// --- DeleteByUUID authz tests ---

func TestTagService_Delete_OwnerOK(t *testing.T) {
	svc, coreQ, _, tagUUID, ownerID, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(ownerID)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Delete_AdminOK(t *testing.T) {
	svc, coreQ, _, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := adminPrincipal(999)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTagService_Delete_SubjectGetsForbidden(t *testing.T) {
	svc, coreQ, _, tagUUID, _, subjectID, _ := buildServiceWithTag(t)
	actor := userPrincipal(subjectID)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if !errors.Is(err, ErrForbidden) {
		t.Errorf("want ErrForbidden, got %v", err)
	}
}

func TestTagService_Delete_StrangerGets404(t *testing.T) {
	svc, coreQ, _, tagUUID, _, _, _ := buildServiceWithTag(t)
	actor := userPrincipal(888)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestTagService_Delete_NotFound(t *testing.T) {
	coreQ := newMockCoreQuerier()
	tagQ := newMockTagQuerier()
	svc, _ := buildService(coreQ, tagQ)
	actor := adminPrincipal(1)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("want ErrNotFound, got %v", err)
	}
}

func TestTagService_Delete_ObserverCalled(t *testing.T) {
	_, coreQ, tagQ, tagUUID, _, _, _ := buildServiceWithTag(t)
	rec := &recordingObserver{}
	obs := observer.NewObserverGroup(rec)
	svc := &TagService{
		db:             newFakeDB(),
		az:             allowAllAuthz{},
		obs:            obs,
		newCoreQuerier: func(_ pgx.Tx) coredb.Querier { return coreQ },
		newTagQuerier:  func(_ pgx.Tx) tagsdb.Querier { return tagQ },
	}
	actor := adminPrincipal(999)

	err := svc.DeleteByUUID(context.Background(), coreQ, actor, tagUUID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.observeCalls) != 1 {
		t.Fatalf("expected 1 observe call, got %d", len(rec.observeCalls))
	}
	call := rec.observeCalls[0]
	if call.op != "delete" {
		t.Errorf("op: want delete, got %q", call.op)
	}
	if call.resource != "tag" {
		t.Errorf("resource: want tag, got %q", call.resource)
	}
}
