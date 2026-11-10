package inventories

import (
	"adventuria/internal/adventuria/errs"
	inventoryModel "adventuria/internal/adventuria/inventories/model"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"adventuria/pkg/helper"
	"context"
	"slices"
)

type repository interface {
	Create(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	Update(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error)
	GetByID(ctx context.Context, id string) (*model.Inventory, error)
	GetPlayerInventoryItemByID(ctx context.Context, playerId, itemId string) (*model.InventoryItem, error)
	GetPlayerInventoryItemWithEffectsByID(ctx context.Context, playerId, itemId string) (*inventoryModel.InventoryItemDetails, error)
	GetAllByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error)
	GetAllDetailedByPlayerID(ctx context.Context, playerId string) ([]*inventoryModel.InventoryItemDetails, error)
	DeleteByID(ctx context.Context, id string) error
	DeleteByIDs(ctx context.Context, ids []string) error
	GetPlayerUsedSlots(ctx context.Context, playerId string) (int, error)
	GetAllDroppableUsingSlotByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItem, error)
}

type effectsService interface {
	CreateFromEffectInfos(infos []*model.EffectInfo) ([]model.Effect, error)
	GetAllByItemID(ctx context.Context, itemId string) ([]model.Effect, error)
	SubscribeInventoryEffects(ctx context.Context, events *model.Events, player *model.Player, items []*model.InventoryItemWithEffects) error
	UnsubscribeActiveEffects(events *model.Events, items []*model.InventoryItem)
}

type items interface {
	GetByID(ctx context.Context, id string) (*model.Item, error)
	GetBySlug(ctx context.Context, slug string) (*model.Item, error)
}

type players interface {
	GetByID(ctx context.Context, playerId, seasonId string) (*model.Player, error)
}

type playerInfo interface {
	GetNameByID(ctx context.Context, id string) (string, error)
}

type playerEvents interface {
	Save(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error)
}

type eventsSubscriber interface {
	Subscribe(ctx context.Context, events *model.Events, player *model.Player, inventory []*model.InventoryItemWithEffects) error
}

type Inventories struct {
	repository       repository
	effects          effectsService
	items            items
	players          players
	playerInfo       playerInfo
	playerEvents     playerEvents
	eventsSubscriber eventsSubscriber
}

func NewInventories(
	repository repository,
	effects effectsService,
	items items,
	players players,
	playerInfo playerInfo,
	playerEvents playerEvents,
	eventsSubscriber eventsSubscriber,
) *Inventories {
	return &Inventories{
		repository:       repository,
		effects:          effects,
		items:            items,
		players:          players,
		playerInfo:       playerInfo,
		playerEvents:     playerEvents,
		eventsSubscriber: eventsSubscriber,
	}
}

func (i *Inventories) Save(ctx context.Context, inventory *model.Inventory) (*model.Inventory, error) {
	if inventory.IsNew() {
		return i.repository.Create(ctx, inventory)
	}

	return i.repository.Update(ctx, inventory)
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

func (i *Inventories) GetAllWithEffectsByPlayerID(ctx context.Context, playerId string) ([]*model.InventoryItemWithEffects, error) {
	items, err := i.repository.GetAllDetailedByPlayerID(ctx, playerId)
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	itemsWithEffects := make([]*model.InventoryItemWithEffects, len(items))
	for j, item := range items {
		effects, err := i.effects.CreateFromEffectInfos(item.Effects)
		if err != nil {
			return nil, err
		}

		itemsWithEffects[j] = model.RestoreInventoryItemWithEffects(
			item.Inventory,
			item.Item,
			effects,
		)
	}

	return itemsWithEffects, nil
}

func (i *Inventories) GetPlayerInventoryItemWithEffectsByID(ctx context.Context, playerId, itemId string) (*model.InventoryItemWithEffects, error) {
	item, err := i.repository.GetPlayerInventoryItemWithEffectsByID(ctx, playerId, itemId)
	if err != nil {
		return nil, err
	}

	effects, err := i.effects.CreateFromEffectInfos(item.Effects)
	if err != nil {
		return nil, err
	}

	return model.RestoreInventoryItemWithEffects(
		item.Inventory,
		item.Item,
		effects,
	), nil
}

func (i *Inventories) AddItem(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	item *model.Item,
	source item_received.Source,
) (*model.InventoryItemWithEffects, error) {
	onBeforeItemAddEvent := model.OnBeforeItemAddEvent{
		ItemRecord:    item,
		ShouldAddItem: true,
	}

	err := events.OnBeforeItemAdd().Trigger(ctx, &onBeforeItemAddEvent)
	if err != nil {
		return nil, err
	}

	if !onBeforeItemAddEvent.ShouldAddItem {
		return nil, nil
	}

	effects, err := i.effects.GetAllByItemID(ctx, item.ID())
	if err != nil {
		return nil, err
	}

	inventoryCreate := model.InventoryCreate{
		Player: player.ID(),
		Item:   item.ID(),
	}
	if item.IsActiveByDefault() {
		inventoryCreate.IsActive = true
	}
	inventory, err := model.NewInventory(inventoryCreate)
	if err != nil {
		return nil, err
	}

	inventory, err = i.Save(ctx, inventory)
	if err != nil {
		return nil, err
	}

	itemReceivedEvent, err := item_received.NewPlayerEvent(
		model.PlayerEventCreate{
			Player: player.ID(),
			Season: player.Progress().Season(),
			Action: player.LastAction().ID(),
		},
		item_received.Payload{
			ItemId:   item.ID(),
			ItemName: item.Name(),
			ItemType: item.Type(),
			Source:   source,
		},
	)
	if err != nil {
		return nil, err
	}

	_, err = i.playerEvents.Save(ctx, itemReceivedEvent)
	if err != nil {
		return nil, err
	}

	inventoryItem := model.RestoreInventoryItemWithEffects(inventory, item, effects)

	err = i.effects.SubscribeInventoryEffects(ctx, events, player, []*model.InventoryItemWithEffects{inventoryItem})
	if err != nil {
		return nil, err
	}

	err = events.OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
		Item: inventoryItem,
	})
	if err != nil {
		return nil, err
	}

	return inventoryItem, nil
}

func (i *Inventories) TryAddItem(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	item *model.Item,
	source item_received.Source,
) (*model.InventoryItemWithEffects, error) {
	if item.IsUsingSlot() {
		hasEmptySlots, err := i.HasEmptySlots(ctx, player)
		if err != nil {
			return nil, err
		}

		if !hasEmptySlots {
			return nil, errs.ErrNotEnoughInventorySlots
		}
	}

	return i.AddItem(ctx, events, player, item, source)
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
	player *model.Player,
	itemId string,
	source item_received.Source,
) (*model.InventoryItemWithEffects, error) {
	item, err := i.items.GetByID(ctx, itemId)
	if err != nil {
		return nil, err
	}

	return i.AddItem(ctx, events, player, item, source)
}

func (i *Inventories) AddItemBySlug(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	itemSlug string,
	source item_received.Source,
) (*model.InventoryItemWithEffects, error) {
	item, err := i.items.GetBySlug(ctx, itemSlug)
	if err != nil {
		return nil, err
	}

	return i.AddItem(ctx, events, player, item, source)
}

func (i *Inventories) CanUseItem(ctx context.Context, events *model.Events, player *model.Player, itemId string) (bool, error) {
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
		if !effect.CanUse(ctx, events, player) {
			return false, nil
		}
	}

	return true, nil
}

func (i *Inventories) UseItem(ctx context.Context, events *model.Events, player *model.Player, itemId string) error {
	item, err := i.GetPlayerInventoryItemWithEffectsByID(ctx, player.ID(), itemId)
	if err != nil {
		return err
	}

	err = item.Inventory().Activate()
	if err != nil {
		return err
	}

	err = i.effects.SubscribeInventoryEffects(ctx, events, player, []*model.InventoryItemWithEffects{item})
	if err != nil {
		return err
	}

	_, err = i.Save(ctx, item.Inventory())
	if err != nil {
		return err
	}

	return nil
}

func (i *Inventories) CanDropItem(ctx context.Context, playerId string, itemId string) (bool, error) {
	item, err := i.GetPlayerInventoryItemByID(ctx, playerId, itemId)
	if err != nil {
		return false, err
	}

	return item.CanDrop(), nil
}

func (i *Inventories) DropItem(ctx context.Context, events *model.Events, player *model.Player, item *model.InventoryItem) error {
	err := i.repository.DeleteByID(ctx, item.Inventory().ID())
	if err != nil {
		return err
	}

	i.effects.UnsubscribeActiveEffects(events, []*model.InventoryItem{item})

	err = player.Progress().BalanceChange(item.DropPrice())
	if err != nil {
		return err
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

func (i *Inventories) DropAllItemsByPlayerID(ctx context.Context, events *model.Events, playerId string) error {
	items, err := i.repository.GetAllByPlayerID(ctx, playerId)
	if err != nil {
		return err
	}

	i.effects.UnsubscribeActiveEffects(events, items)

	ids := make([]string, len(items))
	for i, item := range items {
		ids[i] = item.Inventory().ID()
	}

	return i.repository.DeleteByIDs(ctx, ids)
}

func (i *Inventories) GetPlayerInventory(ctx context.Context, playerId, seasonId string) (*inventoryModel.PlayerInventory, error) {
	player, err := i.players.GetByID(ctx, playerId, seasonId)
	if err != nil {
		return nil, err
	}

	playerName, err := i.playerInfo.GetNameByID(ctx, playerId)
	if err != nil {
		return nil, err
	}

	inventoryDetailed, err := i.repository.GetAllDetailedByPlayerID(ctx, playerId)
	if err != nil {
		return nil, err
	}

	playerInventory := &inventoryModel.PlayerInventory{
		PlayerId:           player.ID(),
		PlayerName:         playerName,
		Balance:            player.Progress().Balance(),
		MaxInventorySlots:  player.Progress().MaxInventorySlots(),
		UsedInventorySlots: 0,
		Items:              make([]*inventoryModel.PlayerInventoryItem, len(inventoryDetailed)),
	}

	if len(inventoryDetailed) == 0 {
		return playerInventory, nil
	}

	itemsWithEffects := make([]*model.InventoryItemWithEffects, len(inventoryDetailed))
	for j, inventory := range inventoryDetailed {
		effects, err := i.effects.CreateFromEffectInfos(inventory.Effects)
		if err != nil {
			return nil, err
		}

		itemsWithEffects[j] = model.RestoreInventoryItemWithEffects(
			inventory.Inventory,
			inventory.Item,
			effects,
		)
	}

	events := model.NewEvents()

	err = i.eventsSubscriber.Subscribe(ctx, events, player, itemsWithEffects)
	if err != nil {
		return nil, err
	}

	for i, inventory := range inventoryDetailed {
		itemWithEffects := itemsWithEffects[i]
		unappliedEffectsIds := itemWithEffects.UnappliedEffects()
		unappliedEffectsIdsMap := make(map[string]struct{}, len(unappliedEffectsIds))
		for _, effectId := range unappliedEffectsIds {
			unappliedEffectsIdsMap[effectId] = struct{}{}
		}

		unappliedEffects := make([]*inventoryModel.PlayerInventoryItemEffect, len(unappliedEffectsIdsMap))
		j := 0
		for _, effect := range itemWithEffects.Effects() {
			if _, ok := unappliedEffectsIdsMap[effect.Data().ID()]; !ok {
				continue
			}

			unappliedEffects[j] = &inventoryModel.PlayerInventoryItemEffect{
				Id:   effect.Data().ID(),
				Type: effect.Data().Type(),
			}
			j++
		}

		playerInventory.Items[i] = &inventoryModel.PlayerInventoryItem{
			InventoryId:      inventory.Inventory.ID(),
			ItemId:           inventory.Item.ID(),
			Created:          inventory.Created,
			Name:             inventory.Item.Name(),
			Icon:             inventory.IconPath,
			Description:      inventory.Item.Description(),
			IsActive:         inventory.Inventory.IsActive(),
			CanDrop:          itemWithEffects.CanDrop(),
			CanUse:           itemWithEffects.CanUse(ctx, events, player),
			DropPrice:        itemWithEffects.DropPrice(),
			UnappliedEffects: unappliedEffects,
		}

		if itemWithEffects.Item().IsUsingSlot() {
			playerInventory.UsedInventorySlots++
		}
	}

	slices.SortFunc(playerInventory.Items, func(a, b *inventoryModel.PlayerInventoryItem) int {
		return a.Created.Compare(b.Created)
	})

	return playerInventory, nil
}
