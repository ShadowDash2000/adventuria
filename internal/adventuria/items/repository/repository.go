package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *Repository) GetIDsByFilter(ctx context.Context, filter items.Filter) ([]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var ids []string
	q := pb.DB().
		Select(schema.ItemSchema.Id).
		From(schema.CollectionItems).
		WithContext(ctx)
	buildQueryFromFilter(filter, q)
	err := q.Column(&ids)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func buildQueryFromFilter(filter items.Filter, q *dbx.SelectQuery) {
	if len(filter.Ids) > 0 {
		q.AndWhere(dbx.In(
			schema.ItemSchema.Id,
			pbhelper.SliceToAny(filter.Ids)...,
		))
		return
	}

	if filter.ItemType != "" {
		q.AndWhere(dbx.HashExp{
			schema.ItemSchema.Type: filter.ItemType,
		})
	}

	if filter.IsRollable != nil {
		q.AndWhere(dbx.HashExp{
			schema.ItemSchema.IsRollable: *filter.IsRollable,
		})
	}

	if filter.PriceGreaterThan != nil {
		q.AndWhere(dbx.NewExp(
			fmt.Sprintf("%s > {:price_greater_than}", schema.ItemSchema.Price),
			dbx.Params{"price_greater_than": *filter.PriceGreaterThan},
		))
	}
}
