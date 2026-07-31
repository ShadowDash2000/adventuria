package add_game_genre

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	activitiesMocks "adventuria/internal/adventuria/activities/mocks"
	"adventuria/internal/adventuria/effects"
	genresMocks "adventuria/internal/adventuria/genres/mocks"
	"adventuria/internal/adventuria/model"
	"context"
	"slices"
	"testing"
)

func TestAddGameGenre_CanUse(t *testing.T) {
	ctx := t.Context()

	setup := func() (*AddGameGenre, *actionsMocks.Actions, *activitiesMocks.Activities) {
		mActions := &actionsMocks.Actions{}
		mGenres := &genresMocks.Genres{}
		mActivities := &activitiesMocks.Activities{}

		eff := &AddGameGenre{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:    mActions,
			genres:     mGenres,
			activities: mActivities,
		}

		return eff, mActions, mActivities
	}

	t.Run("success", func(t *testing.T) {
		eff, mActions, _ := setup()
		lastAction := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				ActivityFilter: &model.ActivityFilter{
					Type: model.ActivityTypeGame,
				},
			},
		})
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, lastAction, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return t == actions.ActionTypeRollWheel
		}

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("cannot roll wheel", func(t *testing.T) {
		eff, mActions, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return false
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})

	t.Run("wrong activity filter type", func(t *testing.T) {
		eff, mActions, _ := setup()
		lastAction := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				ActivityFilter: &model.ActivityFilter{
					Type: model.ActivityTypeGym,
				},
			},
		})
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, lastAction, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when activity filter is not a 'game' type")
		}
	})

	t.Run("has developers filter", func(t *testing.T) {
		eff, mActions, _ := setup()
		lastAction := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				ActivityFilter: &model.ActivityFilter{
					Type:       model.ActivityTypeGame,
					Developers: []string{"dev1"},
				},
			},
		})
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, lastAction, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when cell has developers filter")
		}
	})
}

func TestAddGameGenre_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func() (
		*AddGameGenre,
		*model.Events,
		*model.Player,
		*bool,
		func(context.Context),
		*genresMocks.Genres,
	) {
		mActions := &actionsMocks.Actions{}
		mGenres := &genresMocks.Genres{}
		mActivities := &activitiesMocks.Activities{}

		eff := &AddGameGenre{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:    mActions,
			genres:     mGenres,
			activities: mActivities,
		}

		events := model.NewEvents()
		action := model.RestoreAction(model.ActionData{
			State: model.ActionState{
				ActivityFilter: &model.ActivityFilter{},
			},
		})
		player := model.RestorePlayer(
			model.PlayerData{Id: "p1"},
			&model.PlayerProgress{},
			action,
			nil,
		)

		var callbackCalled bool
		callback := func(ctx context.Context) {
			callbackCalled = true
		}

		return eff, events, player, &callbackCalled, callback, mGenres
	}

	effectCtx := model.EffectContext{
		InvItemID: "item1",
		Priority:  10,
	}

	t.Run("successful activation", func(t *testing.T) {
		eff, events, player, called, callback, mGenres := setup()
		genreID := "rpg"

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		mGenres.ExistsFunc = func(ctx context.Context, id string) (bool, error) {
			return id == genreID, nil
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
			Data:      map[string]any{"genre_id": genreID},
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		genres := player.LastAction().State().ActivityFilter.Genres
		if !slices.Contains(genres, genreID) {
			t.Errorf("Genre %s was not added to player filter: %v", genreID, genres)
		}
	})
}
