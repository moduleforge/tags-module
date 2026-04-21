package httpapi

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moduleforge/core-api/audit"
	coreservice "github.com/moduleforge/core-api/service"
	coredb "github.com/moduleforge/core-model/db"
	"github.com/moduleforge/tags-api/service"
)

// TxBeginner abstracts transaction creation so handlers are testable without
// a real database. *pgxpool.Pool satisfies this interface.
type TxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

// Deps carries the external dependencies that httpapi handlers need.
type Deps struct {
	// Pool is used to open transactions for multi-row creates/deletes.
	Pool *pgxpool.Pool
	// Tx overrides Pool for transaction creation in tests.
	Tx TxBeginner
	// CoreQuerier is the base core-model Querier for entity/type resolution.
	CoreQuerier coredb.Querier
	// Services holds the tag CRUD implementations.
	Services *service.Services
	// Audit is the consumer-provided audit writer.
	Audit audit.Writer
	// Principal extracts caller identity from request context.
	Principal coreservice.PrincipalExtractor
	// Logger is the structured logger for handler-level error messages.
	Logger *slog.Logger
}

// txBeginner returns the TxBeginner to use for starting transactions.
func (d *Deps) txBeginner() TxBeginner {
	if d.Tx != nil {
		return d.Tx
	}
	return d.Pool
}

type handlers struct {
	d Deps
}

// NewRouter wires all tag routes and returns a mountable chi.Router.
// Mount it under any prefix, e.g. r.Mount("/v1", tags.NewRouter(deps)).
func NewRouter(d Deps) chi.Router {
	r := chi.NewRouter()
	h := &handlers{d: d}

	r.Post("/tags", h.handleCreateTag)
	r.Get("/tags", h.handleSearchTags)
	r.Get("/tags/{uuid}", h.handleGetTag)
	r.Put("/tags/{uuid}", h.handlePutTag)
	r.Delete("/tags/{uuid}", h.handleDeleteTag)
	r.Get("/entities/{uuid}/tags", h.handleSubjectTags)

	return r
}
