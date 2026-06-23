package effects

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.EffectInfo, error)
	GetByIDs(ctx context.Context, ids []string) ([]*model.EffectInfo, error)
	GetAllByItemID(ctx context.Context, itemId string) ([]*model.EffectInfo, error)
}

type inventories interface {
	DeleteByID(ctx context.Context, id string) error
}

type Effects struct {
	repository  repository
	inventories inventories
}

func NewEffects(repository repository, inventories inventories) *Effects {
	return &Effects{
		repository:  repository,
		inventories: inventories,
	}
}

func (e *Effects) SubscribeActiveEffects(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	items []*model.InventoryItem,
) error {
	for _, item := range items {
		if !item.Inventory().IsActive() {
			continue
		}

		effects, err := e.GetByIDs(ctx, item.UnappliedEffects())
		if err != nil {
			return err
		}

		for _, effect := range effects {
			unsubKey := item.Inventory().ID() + ":" + effect.Data().ID()

			unsubs, err := effect.Subscribe(
				ctx,
				events,
				player,
				model.EffectContext{
					InvItemID: item.Inventory().ID(),
					Priority:  0,
				},
				func(ctx context.Context) {
					events.Unsubscribe(unsubKey)
					e.effectCallback(ctx, events, player, item, effect)
				},
			)
			if err != nil {
				return err
			}

			events.AddUnsubs(unsubKey, unsubs...)
		}
	}

	return nil
}

func (e *Effects) effectCallback(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	item *model.InventoryItem,
	effect model.Effect,
) {
	item.Inventory().AddAppliedEffects(effect.Data().ID())

	if item.Inventory().AppliedEffectsCount() < item.Item().EffectsCount() {
		return
	}

	player.LastAction().AddUsedItems(item.Item().ID())
	for _, effectId := range item.Item().Effects() {
		events.Unsubscribe(item.Inventory().ID() + ":" + effectId)
	}

	err := e.inventories.DeleteByID(ctx, item.Inventory().ID())
	if err != nil {
		// TODO
	}
}

func (e *Effects) SubscribePersistentEffects(ctx context.Context, events *model.Events, player *model.Player) error {
	for _, effectDef := range GetAllPersistent() {
		unsubs, err := effectDef.Effect().Subscribe(ctx, events, player)
		if err != nil {
			return err
		}

		events.AddUnsubs(player.ID()+":"+string(effectDef.Type()), unsubs...)
	}
	return nil
}

func (e *Effects) UnsubscribeActiveEffects(events *model.Events, items []*model.InventoryItem) {
	for _, item := range items {
		for _, effectId := range item.Item().Effects() {
			events.Unsubscribe(item.Inventory().ID() + ":" + effectId)
		}
	}
}

func (e *Effects) GetView(ctx context.Context, events *model.Events, player *model.Player, id string) (any, error) {
	effect, err := e.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	effectWithView, ok := effect.(model.WithView)
	if !ok {
		return nil, nil
	}

	return effectWithView.GetView(ctx, events, player)
}

func (e *Effects) GetByID(ctx context.Context, id string) (model.Effect, error) {
	effectInfo, err := e.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return effectInfoToEffect(effectInfo)
}

func (e *Effects) GetByIDs(ctx context.Context, ids []string) ([]model.Effect, error) {
	effects, err := e.repository.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	return effectInfosToEffects(effects)
}

func (e *Effects) GetAllByItemID(ctx context.Context, itemId string) ([]model.Effect, error) {
	effectInfos, err := e.repository.GetAllByItemID(ctx, itemId)
	if err != nil {
		return nil, err
	}
	return effectInfosToEffects(effectInfos)
}
