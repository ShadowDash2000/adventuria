package tags

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, tag *model.Tag) (*model.Tag, error)
	Update(ctx context.Context, tag *model.Tag) (*model.Tag, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Tag, error)
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
	tag, err := t.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrTagNotFound) {
			return model.NewTag(data)
		}
		return nil, err
	}

	return tag, nil
}

func (t *Tags) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	if tag.IsNew() {
		return t.repository.Create(ctx, tag)
	}

	return t.repository.Update(ctx, tag)
}
