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

func (r *Repository) Create(ctx context.Context, review *model.Review) (*model.Review, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionReviews)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	ReviewToRecord(review, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToReview(record), nil
}

func (r *Repository) Update(ctx context.Context, review *model.Review) (*model.Review, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionReviews, review.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrReviewNotFound
		}
		return nil, err
	}

	ReviewToRecord(review, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToReview(record), nil
}

func (r *Repository) Save(ctx context.Context, review *model.Review) (*model.Review, error) {
	if review.IsNew() {
		return r.Create(ctx, review)
	}

	return r.Update(ctx, review)
}
