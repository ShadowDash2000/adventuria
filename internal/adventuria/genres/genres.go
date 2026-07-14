package genres

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Exists(ctx context.Context, id string) (bool, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Genre, error)
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

func (g *Genres) GetOrCreate(ctx context.Context, data model.GenreCreate) (*model.Genre, error) {
	genre, err := g.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrGenreNotFound) {
			return model.NewGenre(data)
		}
		return nil, err
	}

	return genre, nil
}

func (g *Genres) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return g.repository.GetChecksumsByIDs(ctx, ids)
}

func (g *Genres) Save(ctx context.Context, genre *model.Genre) (*model.Genre, error) {
	return g.repository.Save(ctx, genre)
}
