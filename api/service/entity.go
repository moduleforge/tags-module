package service

import "github.com/moduleforge/core-api/entity"

// Compile-time assertion that TagEntity satisfies the entity.Entity interface.
var _ entity.Entity = TagEntity{}

// TagEntity is the service-layer entity type implementing entity.Entity.
// It carries the internal entity ID so the Authorizer and observers can route
// by resource name and target ID without inspecting concrete types.
type TagEntity struct {
	// ID is the internal entity ID. Nil when the entity has not yet been
	// persisted (pre-create).
	ID *int64
}

// Resource implements entity.Entity.
func (t TagEntity) Resource() string { return "tag" }

// EntityID implements entity.Entity.
func (t TagEntity) EntityID() *int64 { return t.ID }
