package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
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

func (r *Repository) GetByID(ctx context.Context, id string) (*model.EffectInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionEffects).
		WithContext(ctx).
		Where(dbx.HashExp{schema.EffectSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrEffectNotFound
		}
		return nil, err
	}

	return RecordToEffect(&record), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []string) ([]*model.EffectInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionEffects).
		WithContext(ctx).
		Where(dbx.In(
			schema.EffectSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToEffects(records), nil
}

func (r *Repository) GetAllByItemID(ctx context.Context, itemId string) ([]*model.EffectInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var itemRecord core.Record
	err := pb.RecordQuery(schema.CollectionItems).
		WithContext(ctx).
		Where(dbx.HashExp{schema.ItemSchema.Id: itemId}).
		Limit(1).
		One(&itemRecord)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrItemNotFound
		}
		return nil, err
	}

	var records []*core.Record
	err = pb.RecordQuery(schema.CollectionEffects).
		WithContext(ctx).
		Where(dbx.In(
			schema.EffectSchema.Id,
			pbhelper.SliceToAny(
				itemRecord.GetStringSlice(schema.ItemSchema.Effects))...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToEffects(records), nil
}
