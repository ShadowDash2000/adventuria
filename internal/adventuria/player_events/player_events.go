package player_events

import (
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
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
	playerEventDef, ok := Get(playerEvent.Type())
	if !ok {
		return nil, fmt.Errorf("unknown player_event type: %s", playerEvent.Type())
	}

	err := playerEventDef.New().Verify(ctx, playerEvent.Payload())
	if err != nil {
		return nil, err
	}

	if playerEvent.IsNew() {
		return p.repository.Create(ctx, playerEvent)
	}

	return p.repository.Update(ctx, playerEvent)
}
