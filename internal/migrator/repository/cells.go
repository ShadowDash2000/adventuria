package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Cells struct {
	pb core.App
}

func NewCells(pb core.App) *Cells {
	return &Cells{pb: pb}
}

func (c *Cells) GetIDByName(ctx context.Context, name string) (string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, c.pb)

	var id string
	err := pb.DB().
		Select(schema.CellSchema.Id).
		From(schema.CollectionCells).
		Where(dbx.HashExp{
			schema.CellSchema.Name: name,
		}).
		WithContext(ctx).
		Row(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrCellNotFound
		}
		return "", err
	}

	return id, nil
}
