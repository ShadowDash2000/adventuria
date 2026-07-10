package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) UpdateStreamStatusOrSkip(ctx context.Context, playerId string, status bool) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionPlayers,
			dbx.Params{
				schema.PlayerSchema.IsStreamLive: status,
			},
			dbx.HashExp{
				schema.PlayerSchema.Id:           playerId,
				schema.PlayerSchema.IsStreamLive: !status,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return false, err
	}

	rowsAffected, _ := res.RowsAffected()

	return rowsAffected > 0, nil
}
