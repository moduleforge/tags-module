package httpapi

import (
	"context"
	"io"
	"log/slog"

	"github.com/google/uuid"

	coreservice "github.com/moduleforge/core-api/service"
	coredb "github.com/moduleforge/core-model/db"
	"github.com/moduleforge/tags-api/service"
	tagsdb "github.com/moduleforge/tags-model/db"
)

// --- fake PrincipalExtractor ---

type fakePrincipalExtractor struct {
	p  *coreservice.Principal
	ok bool
}

func (f *fakePrincipalExtractor) FromContext(_ context.Context) (*coreservice.Principal, bool) {
	return f.p, f.ok
}

// --- fake TagServicer ---

type fakeTagService struct {
	tag  service.Tag
	tags []service.Tag
	err  error
}

func (f *fakeTagService) Create(_ context.Context, _ coredb.Querier, _ coreservice.Principal, _ service.CreateTagInput) (service.Tag, error) {
	return f.tag, f.err
}

func (f *fakeTagService) GetByUUID(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ uuid.UUID) (service.Tag, error) {
	return f.tag, f.err
}

func (f *fakeTagService) Search(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ service.SearchTagsFilter) ([]service.Tag, error) {
	return f.tags, f.err
}

func (f *fakeTagService) ListBySubject(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ uuid.UUID, _ *string) ([]service.Tag, error) {
	return f.tags, f.err
}

func (f *fakeTagService) UpdateByUUID(_ context.Context, _ coredb.Querier, _ coreservice.Principal, _ uuid.UUID, _ service.UpdateTagInput) (service.Tag, error) {
	return f.tag, f.err
}

func (f *fakeTagService) DeleteByUUID(_ context.Context, _ coredb.Querier, _ coreservice.Principal, _ uuid.UUID) error {
	return f.err
}

var _ service.TagServicer = (*fakeTagService)(nil)

// --- fake coredb.Querier (pass-through; real work done by fakeTagService) ---

type fakeCoreQuerier struct{}

func (f *fakeCoreQuerier) ArchiveEntity(_ context.Context, _ uuid.UUID) error { return nil }
func (f *fakeCoreQuerier) CreateCorporation(_ context.Context, _ coredb.CreateCorporationParams) (coredb.Corporation, error) {
	return coredb.Corporation{}, nil
}
func (f *fakeCoreQuerier) CreateEntity(_ context.Context, _ int64) (coredb.Entity, error) {
	return coredb.Entity{}, nil
}
func (f *fakeCoreQuerier) CreateLegalEntity(_ context.Context, _ int64) (int64, error) { return 0, nil }
func (f *fakeCoreQuerier) CreateNaturalPerson(_ context.Context, _ coredb.CreateNaturalPersonParams) (coredb.NaturalPerson, error) {
	return coredb.NaturalPerson{}, nil
}
func (f *fakeCoreQuerier) CreateServiceAccount(_ context.Context, _ coredb.CreateServiceAccountParams) (coredb.ServiceAccount, error) {
	return coredb.ServiceAccount{}, nil
}
func (f *fakeCoreQuerier) GetCorporationByEntityID(_ context.Context, _ int64) (coredb.Corporation, error) {
	return coredb.Corporation{}, nil
}
func (f *fakeCoreQuerier) GetEntityByID(_ context.Context, _ int64) (coredb.GetEntityByIDRow, error) {
	return coredb.GetEntityByIDRow{}, nil
}
func (f *fakeCoreQuerier) GetEntityByUUID(_ context.Context, _ uuid.UUID) (coredb.GetEntityByUUIDRow, error) {
	return coredb.GetEntityByUUIDRow{}, nil
}
func (f *fakeCoreQuerier) GetLegalEntityByEntityID(_ context.Context, _ int64) (int64, error) {
	return 0, nil
}
func (f *fakeCoreQuerier) GetNaturalPersonByEntityID(_ context.Context, _ int64) (coredb.NaturalPerson, error) {
	return coredb.NaturalPerson{}, nil
}
func (f *fakeCoreQuerier) GetServiceAccountByEntityID(_ context.Context, _ int64) (coredb.ServiceAccount, error) {
	return coredb.ServiceAccount{}, nil
}
func (f *fakeCoreQuerier) GetTypeBySlug(_ context.Context, _ string) (coredb.Type, error) {
	return coredb.Type{}, nil
}
func (f *fakeCoreQuerier) GetTypeByID(_ context.Context, _ int64) (coredb.Type, error) {
	return coredb.Type{}, nil
}
func (f *fakeCoreQuerier) UnarchiveEntity(_ context.Context, _ uuid.UUID) error { return nil }
func (f *fakeCoreQuerier) UpdateCorporation(_ context.Context, _ coredb.UpdateCorporationParams) error {
	return nil
}
func (f *fakeCoreQuerier) UpdateNaturalPerson(_ context.Context, _ coredb.UpdateNaturalPersonParams) error {
	return nil
}

var _ coredb.Querier = (*fakeCoreQuerier)(nil)

// noopLogger returns a slog.Logger that discards all output.
func noopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 1}))
}

// buildTestDeps builds a Deps with mocked services.
// The fakeTagService receives all handler calls; queriers are no-ops.
func buildTestDeps(ext *fakePrincipalExtractor, tagSvc *fakeTagService) Deps {
	if tagSvc == nil {
		tagSvc = &fakeTagService{}
	}
	svcs := &service.Services{}
	svcs.Tag = tagSvc

	return Deps{
		CoreQuerier: &fakeCoreQuerier{},
		Services:    svcs,
		Principal:   ext,
		Logger:      noopLogger(),
	}
}
