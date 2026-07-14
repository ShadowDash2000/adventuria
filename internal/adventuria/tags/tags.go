package tags

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetOrCreate(ctx context.Context, data model.TagCreate) (*model.Tag, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, tag *model.Tag) (*model.Tag, error)
}

type Tags struct {
	repository repository
}

func NewTags(repo repository) *Tags {
	return &Tags{
		repository: repo,
	}
}

func (t *Tags) GetOrCreate(ctx context.Context, data model.TagCreate) (*model.Tag, error) {
	return t.repository.GetOrCreate(ctx, data)
}

func (t *Tags) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *Tags) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	return t.repository.Save(ctx, tag)
}
