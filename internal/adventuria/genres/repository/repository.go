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

	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionGenres).
		WithContext(ctx).
		Select(schema.GenreSchema.Id).
		Where(dbx.HashExp{schema.GenreSchema.Id: id}).
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

func (r *Repository) GetOrCreate(ctx context.Context, id uuid.UUID, data model.GenreCreate) (*model.Genre, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionGenres).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.GenreSchema.IdDb: data.IdDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.NewGenre(id, data)
		}

		return nil, err
	}

	return RecordToGenre(&record), nil
}

func (r *Repository) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionGenres).
		WithContext(ctx).
		Select(
			schema.GenreSchema.Id,
			schema.GenreSchema.Checksum,
		).
		Where(dbx.In(
			schema.GenreSchema.Id,
			pbhelper.SliceToAny(ids)...,
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	checksums := make(map[string]string, len(records))
	for _, record := range records {
		checksums[record.Id] = record.GetString(schema.GenreSchema.Checksum)
	}

	return checksums, nil
}

func (r *Repository) Create(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionGenres)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	GenreToRecord(genre, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToGenre(record), nil
}

func (r *Repository) Update(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionGenres, genre.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGenreNotFound
		}
		return nil, err
	}

	GenreToRecord(genre, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToGenre(record), nil
}

func (r *Repository) Save(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	if genre.IsNew() {
		return r.Create(ctx, genre)
	}

	return r.Update(ctx, genre)
}
