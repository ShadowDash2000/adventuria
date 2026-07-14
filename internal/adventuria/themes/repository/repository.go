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

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Theme, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionThemes).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ThemeSchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrThemeNotFound
		}

		return nil, err
	}

	return RecordToTheme(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionThemes).
		WithContext(ctx).
		Select(
			schema.ThemeSchema.Id,
			schema.ThemeSchema.Checksum,
		).
		Where(dbx.In(
			schema.ThemeSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.ThemeSchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionThemes)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	ThemeToRecord(theme, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToTheme(record), nil
}

func (r *Repository) Update(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionThemes, theme.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrThemeNotFound
		}
		return nil, err
	}

	ThemeToRecord(theme, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToTheme(record), nil
}

func (r *Repository) Save(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	if theme.IsNew() {
		return r.Create(ctx, theme)
	}

	return r.Update(ctx, theme)
}
