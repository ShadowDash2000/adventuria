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

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Platform, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlatforms).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.PlatformSchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlatformNotFound
		}

		return nil, err
	}

	return RecordToPlatform(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionPlatforms).
		WithContext(ctx).
		Select(
			schema.PlatformSchema.Id,
			schema.PlatformSchema.Checksum,
		).
		Where(dbx.In(
			schema.PlatformSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.PlatformSchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, platform *model.Platform) (*model.Platform, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlatforms)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlatformToRecord(platform, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlatform(record), nil
}

func (r *Repository) Update(ctx context.Context, platform *model.Platform) (*model.Platform, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlatforms, platform.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlatformNotFound
		}
		return nil, err
	}

	PlatformToRecord(platform, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlatform(record), nil
}

func (r *Repository) Save(ctx context.Context, platform *model.Platform) (*model.Platform, error) {
	if platform.IsNew() {
		return r.Create(ctx, platform)
	}

	return r.Update(ctx, platform)
}
