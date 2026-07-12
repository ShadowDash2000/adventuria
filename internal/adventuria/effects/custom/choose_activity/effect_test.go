package choose_activity

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	activitiesMocks "adventuria/internal/adventuria/activities/mocks"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChooseActivity_CanUse(t *testing.T) {
	ctx := t.Context()

	setup := func() (*ChooseActivity, *actionsMocks.Actions) {
		mActions := &actionsMocks.Actions{}
		mActivities := &activitiesMocks.Activities{}

		eff := &ChooseActivity{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:    mActions,
			activities: mActivities,
		}

		return eff, mActions
	}

	t.Run("success", func(t *testing.T) {
		eff, mActions := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return t == actions.ActionTypeDone
		}

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("failure", func(t *testing.T) {
		eff, mActions := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return false
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})
}

func TestChooseActivity_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (*ChooseActivity, *model.Events, *model.Player, *bool, model.EffectCallback) {
		mActions := &actionsMocks.Actions{}
		mActivities := &activitiesMocks.Activities{}

		eff := &ChooseActivity{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:    mActions,
			activities: mActivities,
		}

		events := model.NewEvents()
		action := model.RestoreAction(model.ActionData{
			ItemsList: []string{"game1", "game2"},
		})
		player := model.RestorePlayer(model.PlayerData{Id: "p1"}, &model.PlayerProgress{}, action)

		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}

		return eff, events, player, &callbackCalled, callback
	}

	effectCtx := model.EffectContext{
		InvItemID: "item1",
		Priority:  10,
	}

	t.Run("success", func(t *testing.T) {
		eff, events, player, called, callback := setup()

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
			Data:      map[string]any{"activity_id": "game1"},
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if player.LastAction().Activity() != "game1" {
			t.Errorf("Expected activity game1, got %s", player.LastAction().Activity())
		}
	})

	t.Run("wrong item id", func(t *testing.T) {
		eff, events, player, called, callback := setup()

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "wrong_item",
			Data:      map[string]any{"activity_id": "game1"},
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if *called {
			t.Error("Callback should not be called")
		}
	})

	t.Run("missing activity_id in data", func(t *testing.T) {
		eff, events, player, called, callback := setup()

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
			Data:      map[string]any{},
		})
		if err == nil {
			t.Error("Trigger should fail with invalid activity_id")
		}

		if *called {
			t.Error("Callback should not be called")
		}
	})

	t.Run("activity_id not in items list", func(t *testing.T) {
		eff, events, player, called, callback := setup()

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
			Data:      map[string]any{"activity_id": "unknown_game"},
		})
		if err == nil {
			t.Error("Trigger should fail when activity not in items list")
		}

		if *called {
			t.Error("Callback should not be called")
		}
	})
}
