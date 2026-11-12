package scope

import (
	"adventuria/internal/adventuria/board/actors/player"
	effects2 "adventuria/internal/adventuria/effects"
	"context"
)

type effects interface {
	SubscribeInventoryEffects(ctx context.Context, actor *player.Actor, items []*effects2.InventoryItemWithEffects) error
	SubscribePersistentEffects(ctx context.Context, actor *player.Actor) error
}

type cellEvents interface {
	ListenToCellEvents(ctx context.Context, actor *player.Actor) error
}

type worlds interface {
	SubscribeEffects(ctx context.Context, actor *player.Actor, worldId string) error
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

func (e *EventsSubscriber) Subscribe(ctx context.Context, actor *player.Actor, inventory []*effects2.InventoryItemWithEffects) error {
	err := e.effects.SubscribeInventoryEffects(ctx, actor, inventory)
	if err != nil {
		return err
	}

	err = e.effects.SubscribePersistentEffects(ctx, actor)
	if err != nil {
		return err
	}

	err = e.cellEvents.ListenToCellEvents(ctx, actor)
	if err != nil {
		return err
	}

	err = e.worlds.SubscribeEffects(ctx, actor, actor.Player().Progress().CurrentWorld())
	if err != nil {
		return err
	}

	return nil
}
