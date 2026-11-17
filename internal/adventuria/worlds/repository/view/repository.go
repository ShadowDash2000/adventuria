package view

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria/worlds/service/view"
	"adventuria/pkg/pbtransaction"
	"context"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetAllBoard(ctx context.Context) ([]*view.WorldBoard, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var worldRows []*worldBoardRow
	err := pb.DB().
		Select(
			schema.WorldSchema.Id,
			schema.WorldSchema.Sort,
		).
		From(schema.CollectionWorlds).
		WithContext(ctx).
		All(&worldRows)
	if err != nil {
		return nil, err
	}

	if len(worldRows) == 0 {
		return nil, nil
	}

	res := make([]*view.WorldBoard, len(worldRows))
	for i, row := range worldRows {
		res[i] = row.toWorldBoard()
	}

	return res, nil
}
