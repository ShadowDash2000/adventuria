package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, settings *model.Settings) (*model.Settings, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionSettings)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	SettingsToRecord(settings, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToSettings(record), nil
}

func (r *Repository) GetFirst(ctx context.Context) (*model.Settings, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionSettings).
		WithContext(ctx).
		OrderBy("updated DESC").
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrSettingsNotFound
		}
		return nil, err
	}

	return RecordToSettings(&record), nil
}

func (r *Repository) IsActionsBlocked(ctx context.Context) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var isBlocked bool
	err := pb.RecordQuery(schema.CollectionSettings).
		WithContext(ctx).
		Select(schema.SettingsSchema.BlockAllActions).
		OrderBy("updated DESC").
		Limit(1).
		Row(&isBlocked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.ErrSettingsNotFound
		}
		return false, err
	}

	return isBlocked, nil
}
