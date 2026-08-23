package themes

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	Update(ctx context.Context, theme *model.Theme) (*model.Theme, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Theme, error)
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

func (t *Themes) Save(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	if theme.IsNew() {
		return t.repository.Create(ctx, theme)
	}

	return t.repository.Update(ctx, theme)
}
