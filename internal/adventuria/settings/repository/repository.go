package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
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

	errs := pb.ExpandRecord(record, []string{
		schema.SettingsSchema.IgdbFilterGameTypes,
		schema.SettingsSchema.IgdbFilterPlatforms,
	}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
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

	errs := pb.ExpandRecord(&record, []string{
		schema.SettingsSchema.IgdbFilterGameTypes,
		schema.SettingsSchema.IgdbFilterPlatforms,
	}, nil)
	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to expand records: %v", errs)
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

func (r *Repository) CurrentSeason(ctx context.Context) (string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var currentSeason string
	err := pb.RecordQuery(schema.CollectionSettings).
		WithContext(ctx).
		Select(schema.SettingsSchema.CurrentSeason).
		OrderBy("updated DESC").
		Limit(1).
		Row(&currentSeason)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrSettingsNotFound
		}
		return "", err
	}

	return currentSeason, nil
}

func (r *Repository) IsEventEnded(ctx context.Context) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var eventEnded bool
	err := pb.RecordQuery(schema.CollectionSettings).
		WithContext(ctx).
		Select(schema.SettingsSchema.EventEnded).
		OrderBy("updated DESC").
		Limit(1).
		Row(&eventEnded)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errs.ErrSettingsNotFound
		}
		return false, err
	}

	return eventEnded, nil
}

func (r *Repository) UpdateIGDBGamesParsedByID(ctx context.Context, id string, amount int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionSettings,
			dbx.Params{
				schema.SettingsSchema.IgdbGamesParsed: amount,
			},
			dbx.HashExp{
				schema.SettingsSchema.Id: id,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrSettingsNotFound
	}

	return nil
}
