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

func (r *Repository) Create(ctx context.Context, cheapShark *model.CheapShark) (*model.CheapShark, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionCheapshark)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	CheapSharkToRecord(cheapShark, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToCheapShark(record), nil
}

func (r *Repository) Update(ctx context.Context, cheapShark *model.CheapShark) (*model.CheapShark, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionCheapshark, cheapShark.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCheapSharkNotFound
		}
		return nil, err
	}

	CheapSharkToRecord(cheapShark, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToCheapShark(record), nil
}

func (r *Repository) Save(ctx context.Context, cheapShark *model.CheapShark) (*model.CheapShark, error) {
	if cheapShark.IsNew() {
		return r.Create(ctx, cheapShark)
	}

	return r.Update(ctx, cheapShark)
}

func (r *Repository) ExistsByIdDb(ctx context.Context, idDb int) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionCheapshark).
		WithContext(ctx).
		Select(schema.CheapSharkSchema.Id).
		Where(dbx.HashExp{schema.CheapSharkSchema.IdDb: idDb}).
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

func (r *Repository) GetByAppID(ctx context.Context, id int) (*model.CheapShark, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record *core.Record
	err := pb.RecordQuery(schema.CollectionCheapshark).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.CheapSharkSchema.IdDb: id,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrCheapSharkNotFound
		}
		return nil, err
	}

	return RecordToCheapShark(record), nil
}
