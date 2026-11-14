package actors

import "context"

type repository interface {
	GetActiveActorsBySeasonID(ctx context.Context, seasonId string) ([]*ActiveActor, error)
}

type Actors struct {
	repository repository
}

func NewActors(repository repository) *Actors {
	return &Actors{
		repository: repository,
	}
}

func (a *Actors) GetActiveActorsBySeasonID(ctx context.Context, seasonId string) ([]*ActiveActor, error) {
	return a.repository.GetActiveActorsBySeasonID(ctx, seasonId)
}
