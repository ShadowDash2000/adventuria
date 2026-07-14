package themes

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByIdDb(ctx context.Context, idDb string) (*model.Theme, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, theme *model.Theme) (*model.Theme, error)
}

type Themes struct {
	repository repository
}

func NewThemes(repo repository) *Themes {
	return &Themes{
		repository: repo,
	}
}

func (t *Themes) GetOrCreate(ctx context.Context, data model.ThemeCreate) (*model.Theme, error) {
	theme, err := t.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrThemeNotFound) {
			return model.NewTheme(data)
		}
		return nil, err
	}

	return theme, nil
}

func (t *Themes) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *Themes) Save(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	return t.repository.Save(ctx, theme)
}
