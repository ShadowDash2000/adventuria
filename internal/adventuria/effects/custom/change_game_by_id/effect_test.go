package change_game_by_id

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	activitiesMocks "adventuria/internal/adventuria/activities/mocks"
	cellsMocks "adventuria/internal/adventuria/cells/mocks"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChangeGameById_CanUse(t *testing.T) {
	ctx := t.Context()

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
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				IsChangeGameNotAllowed: false,
			}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool {
			return true
		}

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("change game not allowed on cell", func(t *testing.T) {
		eff, _, mCells := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				IsChangeGameNotAllowed: false,
			}), nil
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when change game is not allowed on cell")
		}
	})

	t.Run("cannot drop", func(t *testing.T) {
		eff, mActions, mCells := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, player *model.Player, actionType model.ActionType) bool {
			if actionType == actions.ActionTypeDrop {
				return false
			}
			return true
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when cannot drop")
		}
	})

	t.Run("cannot done", func(t *testing.T) {
		eff, mActions, mCells := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{}), nil
		}

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, player *model.Player, actionType model.ActionType) bool {
			if actionType == actions.ActionTypeDone {
				return false
			}
			return true
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when cannot done")
		}
	})
}

func TestChangeGameById_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (*model.Events, *model.Player, *bool, func(context.Context)) {
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
		return events, player, &callbackCalled, callback
	}

	effectCtx := model.EffectContext{
		InvItemID: "trigger_item",
		Priority:  10,
	}

	t.Run("OnAfterItemUse - successful activity change", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &ChangeGameById{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "new_game_id",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
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

		if player.LastAction().Activity() != "new_game_id" {
			t.Errorf("Expected activity 'new_game_id', got '%s'", player.LastAction().Activity())
		}
	})

	t.Run("OnAfterMove - call callback", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &ChangeGameById{}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterMove().Trigger(ctx, &model.OnAfterMoveEvent{})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}
	})
}
