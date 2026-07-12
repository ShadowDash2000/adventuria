package cell_points_divide

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestCellPointsDivide_CanUse(t *testing.T) {
	ctx := t.Context()
	eff := &CellPointsDivide{
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

func TestCellPointsDivide_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (*model.Events, *model.Player, *bool, func(context.Context)) {
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
		return events, player, &callbackCalled, callback
	}

	effectCtx := model.EffectContext{
		Priority: 10,
	}

	t.Run("on before done - divide points", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &CellPointsDivide{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "2",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		e := &model.OnBeforeDoneEvent{
			CellPoints: 100,
		}
		err = events.OnBeforeDone().Trigger(ctx, e)
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if e.CellPoints != 50 {
			t.Errorf("Expected CellPoints 50, got %d", e.CellPoints)
		}
	})

	t.Run("on after move - call callback", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &CellPointsDivide{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "2",
				}),
			),
		}

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
