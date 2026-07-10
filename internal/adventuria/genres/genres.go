package genres

import (
	"adventuria/internal/adventuria/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	Exists(ctx context.Context, id string) (bool, error)
	GetOrCreate(ctx context.Context, id uuid.UUID, data model.GenreCreate) (*model.Genre, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, genre *model.Genre) (*model.Genre, error)
}

type Genres struct {
	repository repository
}

func NewGenres(repository repository) *Genres {
	return &Genres{repository: repository}
}

func (g *Genres) Exists(ctx context.Context, id string) (bool, error) {
	return g.repository.Exists(ctx, id)
}

func (g *Genres) GetOrCreate(ctx context.Context, id uuid.UUID, data model.GenreCreate) (*model.Genre, error) {
	return g.repository.GetOrCreate(ctx, id, data)
}

func (g *Genres) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return g.repository.GetChecksumsByIDs(ctx, ids)
}

func (g *Genres) Save(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	return g.repository.Save(ctx, genre)
}
