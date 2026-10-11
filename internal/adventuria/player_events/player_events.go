package player_events

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	Create(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error)
	Update(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error)
}

type PlayerEvents struct {
	repository repository
}

func NewPlayerEvents(repository repository) *PlayerEvents {
	return &PlayerEvents{
		repository: repository,
	}
}

func (p *PlayerEvents) Save(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error) {
	if playerEvent.IsNew() {
		return p.repository.Create(ctx, playerEvent)
	}

	return p.repository.Update(ctx, playerEvent)
}
