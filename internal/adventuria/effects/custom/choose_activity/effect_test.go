package choose_activity

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	activitiesMocks "adventuria/internal/adventuria/activities/mocks"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChooseActivity_CanUse(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

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
		lastAction := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				Activities: &model.ActionActivitiesState{
					Ids: []string{"id1"},
				},
			},
		})
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, lastAction, nil))

		mActions.CanDoFunc = func(ctx context.Context, actor *player.Actor, t model.ActionType) bool {
			return t == actions.ActionTypeCompleteActivity
		}

		if !eff.CanUse(ctx, actor) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("failure", func(t *testing.T) {
		eff, _ := setup()
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, &model.ActionInfo{}, nil))

		if eff.CanUse(ctx, actor) {
			t.Error("CanUse should return false")
		}
	})
}

func TestChooseActivity_Subscribe(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

	setup := func() (*ChooseActivity, *player.Actor, *bool, effects.Callback) {
		mActivities := &activitiesMocks.Activities{}

		eff := &ChooseActivity{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			activities: mActivities,
		}

		events := model.NewEvents()
		action := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				Activities: &model.ActionActivitiesState{
					Ids: []string{"game1", "game2"},
				},
			},
		})
		player := model.RestorePlayer(model.PlayerData{Id: "p1"}, &model.PlayerProgress{}, action, nil)

		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}

		return eff, factory.New(events, player), &callbackCalled, callback
	}

	effectCtx := effects.Context{
		InvItemID: "item1",
		Priority:  10,
	}

	t.Run("success", func(t *testing.T) {
		eff, actor, called, callback := setup()

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
			Data:      map[string]any{"activity_id": "game1"},
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if actor.Player().LastAction().Activity() != "game1" {
			t.Errorf("Expected activity game1, got %s", actor.Player().LastAction().Activity())
		}
	})

	t.Run("wrong item id", func(t *testing.T) {
		eff, actor, called, callback := setup()

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
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
		eff, actor, called, callback := setup()

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
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
		eff, actor, called, callback := setup()

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
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
