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

func (r *Repository) GetByIdDb(ctx context.Context, idDb string) (*model.Genre, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionGenres).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.GenreSchema.IdDb: idDb,
		}).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrGenreNotFound
		}

		return nil, err
	}

	return RecordToGenre(&record), nil
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
