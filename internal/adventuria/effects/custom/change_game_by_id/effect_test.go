package change_game_by_id

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	activitiesMocks "adventuria/internal/adventuria/activities/mocks"
	"adventuria/internal/adventuria/board/actors/player"
	cellsMocks "adventuria/internal/adventuria/cells/mocks"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChangeGameById_CanUse(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

	setup := func() (*ChangeGameById, *actionsMocks.Actions, *cellsMocks.Cells) {
		mActions := &actionsMocks.Actions{}
		mCells := &cellsMocks.Cells{}
		mActivities := &activitiesMocks.Activities{}
		eff := &ChangeGameById{
			EffectBase: effects.NewEffectBase(*model.RestoreEffectInfo(model.EffectData{
				Id:   "eff1",
				Type: Type,
			})),
			actions:    mActions,
			cells:      mCells,
			activities: mActivities,
		}
		return eff, mActions, mCells
	}

	t.Run("success", func(t *testing.T) {
		eff, mActions, mCells := setup()
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil))

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				IsChangeGameNotAllowed: false,
			}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, actor *player.Actor, t model.ActionType) bool {
			return true
		}

		if !eff.CanUse(ctx, actor) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("change game not allowed on cell", func(t *testing.T) {
		eff, _, mCells := setup()
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil))

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				IsChangeGameNotAllowed: false,
			}), nil
		}

		if eff.CanUse(ctx, actor) {
			t.Error("CanUse should return false when change game is not allowed on cell")
		}
	})

	t.Run("cannot drop", func(t *testing.T) {
		eff, mActions, mCells := setup()
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil))

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, actor *player.Actor, actionType model.ActionType) bool {
			if actionType == actions.ActionTypeDrop {
				return false
			}
			return true
		}

		if eff.CanUse(ctx, actor) {
			t.Error("CanUse should return false when cannot drop")
		}
	})

	t.Run("cannot done", func(t *testing.T) {
		eff, mActions, mCells := setup()
		actor := factory.New(model.NewEvents(), model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil))

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, actor *player.Actor, actionType model.ActionType) bool {
			if actionType == actions.ActionTypeDone {
				return false
			}
			return true
		}

		if eff.CanUse(ctx, actor) {
			t.Error("CanUse should return false when cannot done")
		}
	})
}

func TestChangeGameById_Subscribe(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

	setup := func() (*player.Actor, *bool, func(context.Context)) {
		events := model.NewEvents()
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			&model.PlayerProgress{},
			model.RestoreAction(model.ActionData{}),
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

	t.Run("OnAfterItemUse - successful activity change", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &ChangeGameById{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "new_game_id",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "trigger_item",
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if actor.Player().LastAction().Activity() != "new_game_id" {
			t.Errorf("Expected activity 'new_game_id', got '%s'", actor.Player().LastAction().Activity())
		}
	})

	t.Run("OnAfterMove - call callback", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &ChangeGameById{}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = actor.Events().OnAfterMove().Trigger(ctx, &model.OnAfterMoveEvent{})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}
	})
}
