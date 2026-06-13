package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.Item, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.HashExp{schema.ItemSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrItemNotFound
		}
		return nil, err
	}

	return RecordToItem(&record), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]*model.Item, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.In(
			schema.ItemSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToItems(records), nil
}

func (r *Repository) GetAllRollable(ctx context.Context) ([]*model.Item, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.ItemSchema.IsRollable: true},
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToItems(records), nil
}

func (r *Repository) GetAllRollableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.ItemSchema.Type: t},
			dbx.HashExp{schema.ItemSchema.IsRollable: true},
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToItems(records), nil
}

func (r *Repository) GetAllBuyableByType(ctx context.Context, t model.ItemType) ([]*model.Item, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.ItemSchema.Type: t},
			dbx.HashExp{schema.ItemSchema.IsRollable: true},
			dbx.NewExp(pbhelper.GreaterThan(schema.ItemSchema.Price, "0")),
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToItems(records), nil
}
