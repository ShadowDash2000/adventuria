package reviews

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, review *model.Review) (*model.Review, error)
	GetByActionAndPlayerID(ctx context.Context, actionId, playerId string) (*model.Review, error)
}

type Reviews struct {
	repository repository
}

func NewReviews(repository repository) *Reviews {
	return &Reviews{repository: repository}
}

func (r *Reviews) Save(ctx context.Context, review *model.Review) (*model.Review, error) {
	if review.IsNew() {
		return r.repository.Create(ctx, review)
	}

	return r.repository.Update(ctx, review)
}

type UpdateInput struct {
	Comment *string
	Score   *float64
}

func (r *Reviews) UpdateByActionAndPlayerID(ctx context.Context, actionId, playerId string, input UpdateInput) (*model.Review, error) {
	review, err := r.repository.GetByActionAndPlayerID(ctx, actionId, playerId)
	if err != nil {
		return nil, err
	}

	nothingToUpdate := true

	if input.Comment != nil {
		comment, err := model.NewReviewComment(*input.Comment)
		if err != nil {
			return nil, err
		}
		review.SetComment(comment)
		nothingToUpdate = false
	}

	if input.Score != nil {
		score, err := model.NewReviewScore(*input.Score)
		if err != nil {
			return nil, err
		}
		review.SetScore(score)
		nothingToUpdate = false
	}

	if nothingToUpdate {
		return nil, errs.ErrNothingToUpdate
	}

	review, err = r.Save(ctx, review)
	if err != nil {
		return nil, err
	}

	return review, nil
}
