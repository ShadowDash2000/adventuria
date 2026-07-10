package themes

import (
	"adventuria/internal/adventuria/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	GetOrCreate(ctx context.Context, id uuid.UUID, data model.ThemeCreate) (*model.Theme, error)
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

func (t *Themes) GetOrCreate(ctx context.Context, id uuid.UUID, data model.ThemeCreate) (*model.Theme, error) {
	return t.repository.GetOrCreate(ctx, id, data)
}

func (t *Themes) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *Themes) Save(ctx context.Context, theme *model.Theme) (*model.Theme, error) {
	return t.repository.Save(ctx, theme)
}
