package reviews_files

import (
	"context"
)

type repository interface {
	IsValidFileURL(ctx context.Context, rawURL string) (bool, error)
	GetSize(ctx context.Context, id string) (int64, error)
	SaveFromBytes(ctx context.Context, bytes []byte, name string) (string, string, error)
	Delete(ctx context.Context, id string) error
}

type ReviewsFiles struct {
	repository repository
}

func NewReviewsFiles(repository repository) *ReviewsFiles {
	return &ReviewsFiles{
		repository: repository,
	}
}

func (r *ReviewsFiles) IsValidFileURL(ctx context.Context, rawURL string) (bool, error) {
	return r.repository.IsValidFileURL(ctx, rawURL)
}

func (r *ReviewsFiles) GetSize(ctx context.Context, id string) (int64, error) {
	return r.repository.GetSize(ctx, id)
}

func (r *ReviewsFiles) SaveFromBytes(ctx context.Context, bytes []byte, name string) (string, string, error) {
	return r.repository.SaveFromBytes(ctx, bytes, name)
}

func (r *ReviewsFiles) Delete(ctx context.Context, id string) error {
	return r.repository.Delete(ctx, id)
}
