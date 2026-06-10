package inventories

import (
	"adventuria/internal/adventuria_new/model"
	"adventuria/internal/adventuria_new/scope"
	"adventuria/pkg/helper"
	"context"
	"time"

	"github.com/google/uuid"
)

type repository interface {
	Save(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	GetByID(ctx context.Context, id string) (*model.Inventory, error)
	GetPlayerInventoryItemByID(ctx context.Context, playerId, itemId string) (*model.InventoryItem, error)
	GetAllByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error)
	DeleteByID(ctx context.Context, id string) error
	GetPlayerUsedSlots(ctx context.Context, playerId string) (int, error)
	GetAllDroppableUsingSlotByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error)
}

type effectsService interface {
	GetAllByItemID(ctx context.Context, itemId string) ([]model.Effect, error)
	SubscribeActiveEffects(ctx context.Context, events *model.Events, player *model.Player, items []*model.InventoryItem) error
	UnsubscribeActiveEffects(events *model.Events, items []*model.InventoryItem)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
}

type Inventories struct {
	repository repository
	effects    effectsService
	items      items
}

func NewInventories(repository repository, effects effectsService, items items) *Inventories {
	return &Inventories{
		repository: repository,
		effects:    effects,
		items:      items,
	}
}

func (i *Inventories) Save(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	return i.repository.Save(ctx, inventory)
}

func (i *Inventories) GetByID(ctx context.Context, id string) (*model.Inventory, error) {
	return i.repository.GetByID(ctx, id)
}

func (i *Inventories) GetPlayerInventoryItemByID(ctx context.Context, playerId, itemId string) (*model.InventoryItem, error) {
	return i.repository.GetPlayerInventoryItemByID(ctx, playerId, itemId)
}

func (i *Inventories) GetAllByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error) {
	return i.repository.GetAllByPlayerID(ctx, playerId)
}

func (i *Inventories) AddItem(
	ctx context.Context,
	events *model.Events,
	playerId string,
	item *model.Item,
) (*model.InventoryItem, error) {
	onBeforeItemAddEvent := model.OnBeforeItemAddEvent{
		ItemRecord:    item,
		ShouldAddItem: true,
	}

	err := events.OnBeforeItemAdd().Trigger(&onBeforeItemAddEvent)
	if err != nil {
		return nil, err
	}

	if !onBeforeItemAddEvent.ShouldAddItem {
		return nil, nil
	}

	inventoryCreate := model.InventoryCreate{
		Player: playerId,
		Item:   item.ID(),
	}
	if item.IsActiveByDefault() {
		inventoryCreate.Activated = time.Now()
		inventoryCreate.IsActive = true
	}
	inventory, err := model.NewInventory(uuid.New(), inventoryCreate)
	if err != nil {
		return nil, err
	}

	inventory, err = i.repository.Save(ctx, inventory)
	if err != nil {
		return nil, err
	}

	inventoryItem := model.RestoreInventoryItem(inventory, item)

	err = events.OnAfterItemAdd().Trigger(&model.OnAfterItemAddEvent{
		Item: inventoryItem,
	})
	if err != nil {
		return nil, err
	}

	return inventoryItem, nil
}

func (i *Inventories) GetPlayerUsedSlots(ctx context.Context, playerId string) (int, error) {
	return i.repository.GetPlayerUsedSlots(ctx, playerId)
}

func (i *Inventories) HasEmptySlots(ctx context.Context, player *model.Player) (bool, error) {
	usedSlots, err := i.GetPlayerUsedSlots(ctx, player.ID())
	if err != nil {
		return false, err
	}
	return usedSlots < player.Progress().MaxInventorySlots(), nil
}

func (i *Inventories) AddItemByID(
	ctx context.Context,
	events *model.Events,
	playerId string,
	itemId string,
) (*model.InventoryItem, error) {
	item, err := i.items.GetByID(ctx, itemId)
	if err != nil {
		return nil, err
	}
	return i.AddItem(ctx, events, playerId, item)
}

func (i *Inventories) CanUseItem(ctx context.Context, scope *scope.Scope, itemId string) (bool, error) {
	inventory, err := i.GetByID(ctx, itemId)
	if err != nil {
		return false, err
	}

	if inventory.IsActive() {
		return false, nil
	}

	effs, err := i.effects.GetAllByItemID(ctx, inventory.Item())
	if err != nil {
		return false, err
	}

	for _, effect := range effs {
		if !effect.CanUse(ctx, scope.Events(), scope.Player()) {
			return false, nil
		}
	}

	return true, nil
}

func (i *Inventories) UseItem(ctx context.Context, events *model.Events, player *model.Player, itemId string) error {
	item, err := i.GetPlayerInventoryItemByID(ctx, player.ID(), itemId)
	if err != nil {
		return err
	}

	item.Inventory().SetIsActive(true)
	item.Inventory().SetActivated(time.Now())

	err = i.effects.SubscribeActiveEffects(ctx, events, player, []*model.InventoryItem{item})
	if err != nil {
		return err
	}

	_, err = i.Save(ctx, item.Inventory())
	if err != nil {
		return err
	}

	return nil
}

func (i *Inventories) DropItem(ctx context.Context, events *model.Events, player *model.Player, item *model.InventoryItem) error {
	err := i.repository.DeleteByID(ctx, item.Inventory().ID())
	if err != nil {
		return err
	}

	i.effects.UnsubscribeActiveEffects(events, []*model.InventoryItem{item})

	if price := item.Item().Price(); price > 0 {
		err = player.Progress().BalanceChange(price / 2)
		if err != nil {
			return err
		}
	}

	return nil
}

func (i *Inventories) DropPlayerItemByID(ctx context.Context, events *model.Events, player *model.Player, itemId string) error {
	item, err := i.GetPlayerInventoryItemByID(ctx, player.ID(), itemId)
	if err != nil {
		return err
	}

	return i.DropItem(ctx, events, player, item)
}

func (i *Inventories) DropRandomItem(ctx context.Context, events *model.Events, player *model.Player) error {
	items, err := i.repository.GetAllDroppableUsingSlotByPlayerID(ctx, player.ID())
	if err != nil {
		return err
	}

	if len(items) == 0 {
		return nil
	}

	return i.DropItem(ctx, events, player, helper.RandomItemFromSlice(items))
}
