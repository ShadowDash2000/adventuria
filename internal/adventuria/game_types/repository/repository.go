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

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.GameType, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionGameTypes).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.GameTypeSchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGameTypeNotFound
		}

		return nil, err
	}

	return RecordToGameType(&record), nil
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

	var record core.Record
	err := pb.RecordQuery(schema.CollectionGameTypes).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.GameTypeSchema.Id: gameType.ID(),
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGameTypeNotFound
		}
		return nil, err
	}

	GameTypeToRecord(gameType, &record)
	err = pb.SaveWithContext(ctx, &record)
	if err != nil {
		return nil, err
	}

	return RecordToGameType(&record), nil
}
