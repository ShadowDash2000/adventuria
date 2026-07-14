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

func (r *Repository) GetOrCreate(ctx context.Context, data model.GameTypeCreate) (*model.GameType, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionGameTypes).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.GameTypeSchema.IdDb: data.IdDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.NewGameType(data)
		}

		return nil, err
	}

	return RecordToGameType(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionGameTypes).
		WithContext(ctx).
		Select(
			schema.GameTypeSchema.Id,
			schema.GameTypeSchema.Checksum,
		).
		Where(dbx.In(
			schema.GameTypeSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.GameTypeSchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionGameTypes)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	GameTypeToRecord(gameType, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToGameType(record), nil
}

func (r *Repository) Update(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionGameTypes, gameType.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGameTypeNotFound
		}
		return nil, err
	}

	GameTypeToRecord(gameType, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToGameType(record), nil
}

func (r *Repository) Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	if gameType.IsNew() {
		return r.Create(ctx, gameType)
	}

	return r.Update(ctx, gameType)
}
