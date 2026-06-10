package reviews

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

type repository interface {
	Save(ctx context.Context, review *model.Review) (*model.Review, error)
}

type Reviews struct {
	repository repository
}

func NewReviews(repository repository) *Reviews {
	return &Reviews{repository: repository}
}

func (r *Reviews) Save(ctx context.Context, review *model.Review) (*model.Review, error) {
	return r.repository.Save(ctx, review)
}
