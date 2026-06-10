package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/model"
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

func (r *Repository) GetByID(ctx context.Context, id string) (*model.World, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionWorlds).
		WithContext(ctx).
		Where(dbx.HashExp{schema.WorldSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return RecordToWorld(&record), nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.World, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionWorlds).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToWorlds(records), nil
}

func (r *Repository) GetDefault(ctx context.Context) (*model.World, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionWorlds).
		WithContext(ctx).
		Where(dbx.HashExp{schema.WorldSchema.IsDefaultWorld: true}).
		Limit(1).
		One(&record)
	if err != nil {
		return nil, err
	}

	return RecordToWorld(&record), nil
}
