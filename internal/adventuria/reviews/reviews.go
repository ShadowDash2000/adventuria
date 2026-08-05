package reviews

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	Create(ctx context.Context, review *model.Review) (*model.Review, error)
	Update(ctx context.Context, review *model.Review) (*model.Review, error)
	GetByActionID(ctx context.Context, actionId string) (*model.Review, error)
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

func (r *Reviews) GetByActionID(ctx context.Context, actionId string) (*model.Review, error) {
	return r.repository.GetByActionID(ctx, actionId)
}
