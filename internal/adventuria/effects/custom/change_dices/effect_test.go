package change_dices

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChangeDices_CanUse(t *testing.T) {
	ctx := t.Context()
	eff := &ChangeDices{
		EffectBase: effects.NewEffectBase(
			*model.RestoreEffectInfo(model.EffectData{
				Id:   "eff1",
				Type: Type,
			}),
		),
	}

	t.Run("success", func(t *testing.T) {
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})
}

func TestChangeDices_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (*model.Events, *model.Player, *bool, func(context.Context)) {
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
		return events, player, &callbackCalled, callback
	}

	effectCtx := model.EffectContext{
		Priority: 10,
	}

	t.Run("successful dice change", func(t *testing.T) {
		events, player, called, callback := setup()
		eff := &ChangeDices{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "d4;d4",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		e := &model.OnBeforeRollEvent{
			Dices: []model.Dice{model.DiceD6()},
		}
		err = events.OnBeforeRoll().Trigger(ctx, e)
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		if len(e.Dices) != 2 {
			t.Errorf("Expected 2 dices, got %d", len(e.Dices))
		}

		for _, d := range e.Dices {
			if d != model.DiceD4() {
				t.Errorf("Expected d4, got %v", d)
			}
		}
	})

	t.Run("invalid value in subscribe", func(t *testing.T) {
		events, player, _, callback := setup()
		eff := &ChangeDices{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: "invalid",
				}),
			),
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		e := &model.OnBeforeRollEvent{}
		err = events.OnBeforeRoll().Trigger(ctx, e)
		if err == nil {
			t.Error("Trigger should fail for invalid dice type")
		}
	})
}

func TestChangeDices_Verify(t *testing.T) {
	ctx := t.Context()
	eff := &ChangeDices{}

	t.Run("success single dice", func(t *testing.T) {
		err := eff.Verify(ctx, "d6")
		if err != nil {
			t.Errorf("Verify failed for valid value: %v", err)
		}
	})

	t.Run("success multiple dices", func(t *testing.T) {
		err := eff.Verify(ctx, "d4;d6;d4")
		if err != nil {
			t.Errorf("Verify failed for valid value: %v", err)
		}
	})

	t.Run("invalid dice type", func(t *testing.T) {
		err := eff.Verify(ctx, "d10")
		if err == nil {
			t.Error("Verify should fail for unknown dice type")
		}
	})

	t.Run("empty value", func(t *testing.T) {
		err := eff.Verify(ctx, "")
		if err == nil {
			t.Error("Verify should fail for empty value")
		}
	})
}
