package balance_change

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestBalanceChange_CanUse(t *testing.T) {
	ctx := t.Context()
	eff := &BalanceChange{
		EffectBase: effects.NewEffectBase(
			*model.RestoreEffectInfo(model.EffectData{
				Id:   "eff1",
				Type: Type,
			}),
		),
	}

	t.Run("success", func(t *testing.T) {
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil)

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})
}

func TestBalanceChange_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (*model.Events, *model.Player, *bool, func(context.Context)) {
		events := model.NewEvents()
		progress := model.RestorePlayerProgress(model.PlayerProgressData{
			Balance: 50,
		})
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			progress,
			nil,
		)
		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}
		return events, player, &callbackCalled, callback
	}

	effectCtx := model.EffectContext{
		InvItemID: "trigger_item",
		Priority:  10,
	}

	t.Run("successful activation", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "100",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "trigger_item"})
		item := model.RestoreItem(model.ItemData{Id: "item1"})
		invItem := model.RestoreInventoryItem(inv, item)

		err = events.OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if player.Progress().Balance() != 150 {
			t.Errorf("Expected balance 150, got %d", player.Progress().Balance())
		}
	})

	t.Run("negative change", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff2",
					Type:  Type,
					Value: "-30",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "trigger_item"})
		invItem := model.RestoreInventoryItem(inv, nil)

		err = events.OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if player.Progress().Balance() != 20 { // 50 - 30
			t.Errorf("Expected balance 20, got %d", player.Progress().Balance())
		}
	})

	t.Run("wrong trigger item", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "100",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "other_item"})
		invItem := model.RestoreInventoryItem(inv, nil)

		err = events.OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if *called {
			t.Error("Callback should not be called")
		}

		if player.Progress().Balance() != 50 {
			t.Errorf("Balance should not change, got %d", player.Progress().Balance())
		}
	})
}

func TestBalanceChange_Verify(t *testing.T) {
	ctx := t.Context()
	eff := &BalanceChange{}

	t.Run("success", func(t *testing.T) {
		err := eff.Verify(ctx, "100")
		if err != nil {
			t.Errorf("Verify failed for valid value: %v", err)
		}
	})

	t.Run("invalid value", func(t *testing.T) {
		err := eff.Verify(ctx, "abc")
		if err == nil {
			t.Error("Verify should fail for non-integer value")
		}
	})

	t.Run("empty value", func(t *testing.T) {
		err := eff.Verify(ctx, "")
		if err == nil {
			t.Error("Verify should fail for empty value")
		}
	})
}
