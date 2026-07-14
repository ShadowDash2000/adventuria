package tags

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByIdDb(ctx context.Context, idDb string) (*model.Tag, error)
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
	tag, err := t.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrTagNotFound) {
			return model.NewTag(data)
		}
		return nil, err
	}

	return tag, nil
}

func (t *Tags) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *Tags) Save(ctx context.Context, tag *model.Tag) (*model.Tag, error) {
	return t.repository.Save(ctx, tag)
}
