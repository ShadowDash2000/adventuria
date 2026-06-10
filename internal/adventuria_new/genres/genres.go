package genres

import "context"

type repository interface {
	Exists(ctx context.Context, id string) (bool, error)
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
