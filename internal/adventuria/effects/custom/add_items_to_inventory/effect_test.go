package add_items_to_inventory

import (
	"adventuria/internal/adventuria/effects"
	inventoriesMocks "adventuria/internal/adventuria/inventories/mocks"
	itemsMocks "adventuria/internal/adventuria/items/mocks"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestAddItemsToInventory_CanUse(t *testing.T) {
	ctx := t.Context()

	setup := func() (*AddItemsToInventory, *inventoriesMocks.Inventories, *itemsMocks.Items) {
		mInventories := &inventoriesMocks.Inventories{}
		mItems := &itemsMocks.Items{}

		eff := &AddItemsToInventory{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			inventories: mInventories,
			items:       mItems,
		}

		return eff, mInventories, mItems
	}

	t.Run("success", func(t *testing.T) {
		eff, _, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil)

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})
}

func TestAddItemsToInventory_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (
		*AddItemsToInventory,
		*model.Events,
		*model.Player,
		*bool,
		func(context.Context),
		*inventoriesMocks.Inventories,
	) {
		mInventories := &inventoriesMocks.Inventories{}
		mItems := &itemsMocks.Items{}

		eff := &AddItemsToInventory{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "item1;item2",
				}),
			),
			inventories: mInventories,
			items:       mItems,
		}

		events := model.NewEvents()
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			&model.PlayerProgress{},
			nil,
		)

		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}

		return eff, events, player, &callbackCalled, callback, mInventories
	}

	effectCtx := model.EffectContext{
		InvItemID: "trigger_item",
		Priority:  10,
	}

	t.Run("successful activation", func(t *testing.T) {
		eff, events, player, called, callback, mInventories := setup()
		addedItems := make([]string, 0)

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		mInventories.AddItemByIDFunc = func(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error) {
			addedItems = append(addedItems, itemId)
			return &model.InventoryItem{}, nil
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "trigger_item",
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if len(addedItems) != 2 {
			t.Errorf("Expected 2 items to be added, got %d", len(addedItems))
		}

		if addedItems[0] != "item1" || addedItems[1] != "item2" {
			t.Errorf("Unexpected items added: %v", addedItems)
		}
	})

	t.Run("wrong trigger item", func(t *testing.T) {
		eff, events, player, called, callback, mInventories := setup()
		addedItemsCount := 0

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		mInventories.AddItemByIDFunc = func(ctx context.Context, events *model.Events, player *model.Player, itemId string) (*model.InventoryItem, error) {
			addedItemsCount++
			return &model.InventoryItem{}, nil
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "other_item",
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if *called {
			t.Error("Callback should not be called")
		}

		if addedItemsCount != 0 {
			t.Error("No items should be added")
		}
	})
}

func TestAddItemsToInventory_Verify(t *testing.T) {
	ctx := t.Context()

	setup := func() (*AddItemsToInventory, *itemsMocks.Items) {
		mItems := &itemsMocks.Items{}
		eff := &AddItemsToInventory{
			items: mItems,
		}
		return eff, mItems
	}

	t.Run("success", func(t *testing.T) {
		eff, mItems := setup()

		mItems.GetByIDsFunc = func(ctx context.Context, ids []string) ([]*model.Item, error) {
			return []*model.Item{{}, {}}, nil
		}

		err := eff.Verify(ctx, "item1;item2")
		if err != nil {
			t.Errorf("Verify failed: %v", err)
		}
	})

	t.Run("empty value", func(t *testing.T) {
		eff, _ := setup()

		err := eff.Verify(ctx, "")
		if err == nil {
			t.Error("Verify should fail for empty value")
		}
	})

	t.Run("item not found", func(t *testing.T) {
		eff, mItems := setup()

		mItems.GetByIDsFunc = func(ctx context.Context, ids []string) ([]*model.Item, error) {
			return []*model.Item{{}}, nil
		}

		err := eff.Verify(ctx, "item1;item2")
		if err == nil {
			t.Error("Verify should fail when items are missing")
		}
	})
}
