package httpapi

import (
	"context"
	"errors"
	"io"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

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

// --- fake audit.Writer ---

type fakeAuditWriter struct{}

func (f *fakeAuditWriter) Write(_ context.Context, _, _ string, _ *int64, _, _ any) error {
	return nil
}

// --- fake TagServicer ---

type fakeTagService struct {
	tag  service.Tag
	tags []service.Tag
	err  error
}

func (f *fakeTagService) Create(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ service.TxBeginner, _ service.CreateTagInput) (service.Tag, error) {
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

func (f *fakeTagService) UpdateByUUID(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ uuid.UUID, _ service.UpdateTagInput) (service.Tag, error) {
	return f.tag, f.err
}

func (f *fakeTagService) DeleteByUUID(_ context.Context, _ coredb.Querier, _ tagsdb.Querier, _ coreservice.Principal, _ uuid.UUID, _ service.TxBeginner) error {
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

// --- fakeTx for handler tests that exercise the tx path (Create / Delete) ---

type fakeTx struct{}

func (f *fakeTx) Begin(_ context.Context) (pgx.Tx, error) { return f, nil }
func (f *fakeTx) Commit(_ context.Context) error          { return nil }
func (f *fakeTx) Rollback(_ context.Context) error        { return nil }
func (f *fakeTx) CopyFrom(_ context.Context, _ pgx.Identifier, _ []string, _ pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not implemented")
}
func (f *fakeTx) SendBatch(_ context.Context, _ *pgx.Batch) pgx.BatchResults { return nil }
func (f *fakeTx) LargeObjects() pgx.LargeObjects                             { return pgx.LargeObjects{} }
func (f *fakeTx) Prepare(_ context.Context, _, _ string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeTx) Exec(_ context.Context, _ string, _ ...any) (pgconn.CommandTag, error) {
	return pgconn.CommandTag{}, nil
}
func (f *fakeTx) Query(_ context.Context, _ string, _ ...any) (pgx.Rows, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeTx) QueryRow(_ context.Context, _ string, _ ...any) pgx.Row { return nil }
func (f *fakeTx) Conn() *pgx.Conn                                        { return nil }

type fakeTxBeginner struct {
	err error
}

func (f *fakeTxBeginner) Begin(_ context.Context) (pgx.Tx, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &fakeTx{}, nil
}

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
		Pool:        nil,
		Tx:          &fakeTxBeginner{},
		CoreQuerier: &fakeCoreQuerier{},
		Services:    svcs,
		Audit:       &fakeAuditWriter{},
		Principal:   ext,
		Logger:      noopLogger(),
	}
}
