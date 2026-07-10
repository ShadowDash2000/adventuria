package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
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

func (r *Repository) Create(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionSteamSpy)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	SteamSpyToRecord(steamSpy, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToSteamSpy(record), nil
}

func (r *Repository) Update(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionSteamSpy, steamSpy.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSteamSpyNotFound
		}
		return nil, err
	}

	SteamSpyToRecord(steamSpy, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToSteamSpy(record), nil
}

func (r *Repository) Save(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error) {
	if steamSpy.IsNew() {
		return r.Create(ctx, steamSpy)
	}

	return r.Update(ctx, steamSpy)
}

func (r *Repository) ExistsByIdDb(ctx context.Context, idDb int) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionSteamSpy).
		WithContext(ctx).
		Select(schema.SteamSpySchema.Id).
		Where(dbx.HashExp{schema.SteamSpySchema.IdDb: idDb}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repository) GetByAppID(ctx context.Context, id int) (*model.SteamSpy, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record *core.Record
	err := pb.RecordQuery(schema.CollectionSteamSpy).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.SteamSpySchema.IdDb: id,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSteamSpyNotFound
		}
		return nil, err
	}

	return RecordToSteamSpy(record), nil
}
