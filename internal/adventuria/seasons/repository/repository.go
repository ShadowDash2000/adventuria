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

func (r *Repository) Create(ctx context.Context, season *model.Season) (*model.Season, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionSeasons)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	err = SeasonToRecord(season, record)
	if err != nil {
		return nil, err
	}

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToSeason(record), nil
}

func (r *Repository) GetFirst(ctx context.Context) (*model.Season, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionSeasons).
		WithContext(ctx).
		OrderBy("created DESC").
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSeasonNotFound
		}
		return nil, err
	}

	return RecordToSeason(&record), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.Season, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionSeasons).
		WithContext(ctx).
		Where(dbx.HashExp{schema.SeasonSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSeasonNotFound
		}
		return nil, err
	}

	return RecordToSeason(&record), nil
}
