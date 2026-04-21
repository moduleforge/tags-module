package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/moduleforge/core-api/display"
	tagsdb "github.com/moduleforge/tags-model/db"
)

// RegisterBuiltins wires the default display renderer for the "tag" type into reg.
// Call this once at service startup after constructing the registry.
//
// Renderer behaviour:
//   - tag name: "<purpose>:<value>".
func RegisterBuiltins(reg *display.Registry, q tagsdb.Querier) {
	reg.Register("tag", display.FieldName, func(ctx context.Context, tx pgx.Tx, entityID int64) (string, error) {
		var querier tagsdb.Querier
		if tx != nil {
			querier = tagsdb.New(tx)
		} else {
			querier = q
		}
		tag, err := querier.GetTagByEntityID(ctx, entityID)
		if err != nil {
			return "", fmt.Errorf("display tag name: %w", err)
		}
		return tag.Purpose + ":" + tag.Value, nil
	})
}
