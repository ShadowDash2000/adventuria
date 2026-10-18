package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Players struct {
	pb core.App
}

func NewPlayers(pb core.App) *Players {
	return &Players{pb: pb}
}

func (p *Players) GetIDByName(ctx context.Context, name string) (id string, err error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, p.pb)

	err = pb.DB().
		Select(schema.PlayerSchema.Id).
		From(schema.CollectionPlayers).
		Where(dbx.HashExp{
			schema.PlayerSchema.Name: name,
		}).
		Limit(1).
		Row(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrPlayerNotFound
		}
		return "", err
	}

	return id, nil
}
