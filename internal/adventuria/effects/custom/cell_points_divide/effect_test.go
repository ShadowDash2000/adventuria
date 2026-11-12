package cell_points_divide

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestCellPointsDivide_CanUse(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)
	eff := &CellPointsDivide{
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

func TestCellPointsDivide_Subscribe(t *testing.T) {
	ctx := t.Context()
	factory := player.NewFactory(nil, nil)

	setup := func() (*player.Actor, *bool, func(context.Context)) {
		events := model.NewEvents()
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			&model.PlayerProgress{},
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
		Priority: 10,
	}

	t.Run("on done - divide points", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &CellPointsDivide{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "2",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, actor, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		e := &model.OnActivityDoneEvent{
			Mode: model.EffectRunModeApply,
			Result: model.NewActivityDoneResult(model.ActivityDoneResultData{
				Points: 100,
			}),
		}
		err = actor.Events().OnActivityDone().Trigger(ctx, e)
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if e.Result.Points() != 50 {
			t.Errorf("Expected CellPoints 50, got %d", e.Result.Points())
		}
	})

	t.Run("on after move - call callback", func(t *testing.T) {
		actor, called, callback := setup()
		eff := &CellPointsDivide{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "2",
				}),
			),
		}

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

func TestCellPointsDivide_Verify(t *testing.T) {
	ctx := t.Context()
	eff := &CellPointsDivide{}

	t.Run("success", func(t *testing.T) {
		err := eff.Verify(ctx, "2")
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

	t.Run("zero value", func(t *testing.T) {
		err := eff.Verify(ctx, "0")
		if err == nil {
			t.Error("Verify should fail for zero value")
		}
	})

	t.Run("negative value", func(t *testing.T) {
		err := eff.Verify(ctx, "-1")
		if err == nil {
			t.Error("Verify should fail for negative value")
		}
	})
}
