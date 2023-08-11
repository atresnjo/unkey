package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	gen "github.com/unkeyed/unkey/apps/agent/gen/database"
	"github.com/unkeyed/unkey/apps/agent/pkg/entities"
)

func (db *database) CreateKey(ctx context.Context, key entities.Key) error {
	params, err := transformKeyEntitytoCreateKeyParams(key)
	if err != nil {
		return fmt.Errorf("unable to transform key entity to CreateKeyParams")
	}
	return db.write().CreateKey(ctx, params)
}

func transformKeyEntitytoCreateKeyParams(key entities.Key) (gen.CreateKeyParams, error) {
	params := gen.CreateKeyParams{
		ID:             key.Id,
		Hash:           key.Hash,
		Start:          key.Start,
		OwnerID:        sql.NullString{String: key.OwnerId, Valid: key.OwnerId != ""},
		CreatedAt:      key.CreatedAt,
		Expires:        sql.NullTime{Time: key.Expires, Valid: !key.Expires.IsZero()},
		WorkspaceID:    key.WorkspaceId,
		ForWorkspaceID: sql.NullString{String: key.ForWorkspaceId, Valid: key.ForWorkspaceId != ""},
		Name:           sql.NullString{String: key.Name, Valid: key.Name != ""},
		KeyAuthID:      sql.NullString{String: key.KeyAuthId, Valid: true},
	}

	metaJson, err := json.Marshal(key.Meta)
	if err != nil {
		return gen.CreateKeyParams{}, fmt.Errorf("unable to marshal meta: %w", err)
	}
	params.Meta = sql.NullString{String: string(metaJson), Valid: true}

	if key.Ratelimit != nil {
		params.RatelimitType = sql.NullString{String: key.Ratelimit.Type, Valid: true}
		params.RatelimitLimit = sql.NullInt32{Int32: int32(key.Ratelimit.Limit), Valid: true}
		params.RatelimitRefillRate = sql.NullInt32{Int32: key.Ratelimit.RefillRate, Valid: true}
		params.RatelimitRefillInterval = sql.NullInt32{Int32: key.Ratelimit.RefillInterval, Valid: true}
	}

	if key.Remaining != nil {
		params.RemainingRequests = sql.NullInt32{Int32: *key.Remaining, Valid: true}

	}

	return params, nil

}