package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) IsValidFileURL(ctx context.Context, rawURL string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionReviewsFiles)
	if err != nil {
		return false, err
	}

	return strings.Contains(rawURL, "/api/files/"+collection.Id), nil
}

func (r *Repository) GetSize(ctx context.Context, id string) (int64, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionReviewsFiles).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ReviewsFilesSchema.Id: id,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errs.ErrReviewFileNotFound
		}

		return 0, err
	}

	fsys, err := pb.NewFilesystem()
	if err != nil {
		return 0, err
	}
	defer fsys.Close()

	reader, err := fsys.GetReader(record.BaseFilesPath() + "/" + record.GetString(schema.ReviewsFilesSchema.File))
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	return reader.Size(), nil
}

func (r *Repository) SaveFromBytes(ctx context.Context, bytes []byte, name string) (string, string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	f, err := filesystem.NewFileFromBytes(bytes, name)
	if err != nil {
		return "", "", err
	}

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionReviewsFiles)
	if err != nil {
		return "", "", err
	}

	record := core.NewRecord(collection)
	record.Set(schema.ReviewsFilesSchema.File, f)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return "", "", err
	}

	return record.Id,
		"/api/files/" + record.BaseFilesPath() + "/" + record.GetString(schema.ReviewsFilesSchema.File),
		nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionReviewsFiles).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.ReviewsFilesSchema.Id: id,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrReviewFileNotFound
		}

		return err
	}

	return pb.DeleteWithContext(ctx, &record)
}
