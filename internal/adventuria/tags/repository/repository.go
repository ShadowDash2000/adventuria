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

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Tag, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionTags).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.TagSchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTagNotFound
		}

		return nil, err
	}

	return RecordToTag(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionTags).
		WithContext(ctx).
		Select(
			schema.TagSchema.Id,
			schema.TagSchema.Checksum,
		).
		Where(dbx.In(
			schema.TagSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.TagSchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionTags)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	TagToRecord(tag, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToTag(record), nil
}

func (r *Repository) Update(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionTags, tag.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrTagNotFound
		}
		return nil, err
	}

	TagToRecord(tag, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToTag(record), nil
}

func (r *Repository) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if tag.IsNew() {
		return r.Create(ctx, tag)
	}

	return r.Update(ctx, tag)
}
