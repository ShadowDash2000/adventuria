package repository

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
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

func (r *Repository) Create(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionHowLongToBeat)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	HowLongToBeatToRecord(hltb, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToHowLongToBeat(record), nil
}

func (r *Repository) Update(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionHowLongToBeat, hltb.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrHowLongToBeatNotFound
		}
		return nil, err
	}

	HowLongToBeatToRecord(hltb, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToHowLongToBeat(record), nil
}

func (r *Repository) Save(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error) {
	if hltb.IsNew() {
		return r.Create(ctx, hltb)
	}

	return r.Update(ctx, hltb)
}

func (r *Repository) ExistsByIdDb(ctx context.Context, idDb int) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionHowLongToBeat).
		WithContext(ctx).
		Select(schema.HowLongToBeatSchema.Id).
		Where(dbx.HashExp{schema.HowLongToBeatSchema.IdDb: idDb}).
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
