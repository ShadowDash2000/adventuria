package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionActions)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	err = ActionToRecord(action, record)
	if err != nil {
		return nil, err
	}

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToAction(record)
}

func (r *Repository) Update(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionActions, action.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActionNotFound
		}
		return nil, err
	}

	err = ActionToRecord(action, record)
	if err != nil {
		return nil, err
	}

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToAction(record)
}

func (r *Repository) Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error) {
	if action.IsNew() {
		return r.Create(ctx, action)
	}

	return r.Update(ctx, action)
}

func (r *Repository) GetLastActionByPlayerId(ctx context.Context, playerId string, timeFrom time.Time) (*model.ActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	dateFrom, err := types.ParseDateTime(timeFrom)
	if err != nil {
		return nil, err
	}

	var record core.Record
	err = pb.RecordQuery(schema.CollectionActions).
		WithContext(ctx).
		Where(dbx.HashExp{schema.ActionSchema.Player: playerId}).
		AndWhere(dbx.NewExp("created > {:date}", dbx.Params{
			"date": dateFrom,
		})).
		OrderBy("created DESC", "rowid DESC").
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrActionNotFound
		}
		return nil, err
	}

	action, err := RecordToAction(&record)
	if err != nil {
		return nil, err
	}

	return action, nil
}
