// Package service exposes tx-aware tag CRUD for consumer apps.
// Consumer modules wire in their own auth and audit implementations
// via the PrincipalExtractor and audit.Writer interfaces.
package service

import (
	"github.com/moduleforge/core-api/audit"
	tagsdb "github.com/moduleforge/tags-model/db"
)

// Services is the aggregate of all tag service implementations. Consumers
// construct this once at startup and pass it into httpapi.NewRouter.
type Services struct {
	Tag TagServicer

	// q is the base Querier backed by the pool, exposed so handlers can
	// derive tx-scoped queriers via tagsdb.New(tx).
	q tagsdb.Querier
}

// New constructs a Services aggregate. q is typically tagsdb.New(pool).
func New(q tagsdb.Querier, aw audit.Writer) *Services {
	return &Services{
		Tag: &TagService{aw: aw},
		q:   q,
	}
}

// Querier returns the base Querier so handlers can derive tx-scoped variants.
func (s *Services) Querier() tagsdb.Querier {
	return s.q
}
