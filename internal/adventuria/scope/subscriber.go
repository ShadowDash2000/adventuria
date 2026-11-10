package scope

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type effects interface {
	SubscribeInventoryEffects(ctx context.Context, events *model.Events, player *model.Player, items []*model.InventoryItemWithEffects) error
	SubscribePersistentEffects(ctx context.Context, events *model.Events, player *model.Player) error
}

type cellEvents interface {
	ListenToCellEvents(ctx context.Context, events *model.Events, player *model.Player) error
}

type worlds interface {
	SubscribeEffects(ctx context.Context, events *model.Events, player *model.Player, worldId string) error
}

type EventsSubscriber struct {
	effects    effects
	cellEvents cellEvents
	worlds     worlds
}

func NewEventsSubscriber(effects effects, cellEvents cellEvents, worlds worlds) *EventsSubscriber {
	return &EventsSubscriber{
		effects:    effects,
		cellEvents: cellEvents,
		worlds:     worlds,
	}
}

func (e *EventsSubscriber) Subscribe(ctx context.Context, events *model.Events, player *model.Player, inventory []*model.InventoryItemWithEffects) error {
	err := e.effects.SubscribeInventoryEffects(ctx, events, player, inventory)
	if err != nil {
		return err
	}

	err = e.effects.SubscribePersistentEffects(ctx, events, player)
	if err != nil {
		return err
	}

	err = e.cellEvents.ListenToCellEvents(ctx, events, player)
	if err != nil {
		return err
	}

	err = e.worlds.SubscribeEffects(ctx, events, player, player.Progress().CurrentWorld())
	if err != nil {
		return err
	}

	return nil
}
