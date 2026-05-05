// Package service exposes tx-aware tag CRUD for consumer apps.
// Consumer modules wire in their own auth and audit implementations
// via the authz.Authorizer and observer.MutationObserver interfaces.
package service

import (
	"github.com/jackc/pgx/v5"

	"github.com/moduleforge/core-api/authz"
	"github.com/moduleforge/core-api/observer"
	"github.com/moduleforge/core-api/txhelper"
	coredb "github.com/moduleforge/core-model/db"
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

// New constructs a Services aggregate.
//
// coreQ is typically coredb.New(pool) and is the base querier for pre-tx
// entity resolution reads.
//
// tagQ is typically tagsdb.New(pool) and is the base querier for pre-tx
// tag reads.
//
// db is the connection pool (or any txhelper.DB) used to open transactions
// for mutating operations.
//
// az gates every operation; a non-nil error from az.Authorize aborts the
// operation immediately.
//
// obs receives in-tx and post-commit notifications for every mutation;
// pass observer.NewObserverGroup() for a no-op group.
func New(coreQ coredb.Querier, tagQ tagsdb.Querier, db txhelper.DB, az authz.Authorizer, obs *observer.ObserverGroup) *Services {
	newCoreQ := func(tx pgx.Tx) coredb.Querier { return coredb.New(tx) }
	newTagQ := func(tx pgx.Tx) tagsdb.Querier { return tagsdb.New(tx) }
	return &Services{
		Tag: &TagService{
			db:             db,
			az:             az,
			obs:            obs,
			newCoreQuerier: newCoreQ,
			newTagQuerier:  newTagQ,
		},
		q: tagQ,
	}
}

// Querier returns the base Querier so handlers can derive tx-scoped variants.
func (s *Services) Querier() tagsdb.Querier {
	return s.q
}
