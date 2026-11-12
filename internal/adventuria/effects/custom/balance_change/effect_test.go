package balance_change

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestBalanceChange_CanUse(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)
	eff := &BalanceChange{
		EffectBase: effects.NewEffectBase(
			*model.RestoreEffectInfo(model.EffectData{
				Id:   "eff1",
				Type: Type,
			}),
		),
	}

	t.Run("success", func(t *testing.T) {
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil))

		if !eff.CanUse(ctx, actor) {
			t.Error("CanUse should return true")
		}
	})
}

func TestBalanceChange_Subscribe(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

	setup := func() (*player.Actor, *bool, func(context.Context)) {
		events := model.NewEvents()
		progress := model.RestorePlayerProgress(model.PlayerProgressData{
			Balance: 50,
		})
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			progress,
			nil,
			nil,
		)
		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}
		return factory.New(events, player), &callbackCalled, callback
	}

	effectCtx := effects.Context{
		InvItemID: "trigger_item",
		Priority:  10,
	}

	t.Run("successful activation", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "100",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "trigger_item"})
		item := model.RestoreItem(model.ItemData{Id: "item1"})
		invItem := model.RestoreInventoryItem(inv, item)

		err = actor.Events().OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if actor.Player().Progress().Balance() != 150 {
			t.Errorf("Expected balance 150, got %d", actor.Player().Progress().Balance())
		}
	})

	t.Run("negative change", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff2",
					Type:  Type,
					Value: "-30",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "trigger_item"})
		invItem := model.RestoreInventoryItem(inv, nil)

		err = actor.Events().OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if actor.Player().Progress().Balance() != 20 { // 50 - 30
			t.Errorf("Expected balance 20, got %d", actor.Player().Progress().Balance())
		}
	})

	t.Run("wrong trigger item", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &BalanceChange{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "100",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		inv := model.RestoreInventory(model.InventoryData{Id: "other_item"})
		invItem := model.RestoreInventoryItem(inv, nil)

		err = actor.Events().OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: invItem,
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if *called {
			t.Error("Callback should not be called")
		}

		if actor.Player().Progress().Balance() != 50 {
			t.Errorf("Balance should not change, got %d", actor.Player().Progress().Balance())
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
