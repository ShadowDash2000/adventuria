package add_game_genre

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	filtersMocks "adventuria/internal/adventuria/activity_filters/mocks"
	cellsMocks "adventuria/internal/adventuria/cells/mocks"
	"adventuria/internal/adventuria/effects"
	genresMocks "adventuria/internal/adventuria/genres/mocks"
	"adventuria/internal/adventuria/model"
	"context"
	"slices"
	"testing"
)

func TestAddGameGenre_CanUse(t *testing.T) {
	ctx := t.Context()

	setup := func() (*AddGameGenre, *actionsMocks.Actions, *cellsMocks.Cells, *filtersMocks.ActivityFilters) {
		mActions := &actionsMocks.Actions{}
		mCells := &cellsMocks.Cells{}
		mGenres := &genresMocks.Genres{}
		mFilters := &filtersMocks.ActivityFilters{}

		eff := &AddGameGenre{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:         mActions,
			cells:           mCells,
			genres:          mGenres,
			activityFilters: mFilters,
		}

		return eff, mActions, mCells, mFilters
	}

	t.Run("success", func(t *testing.T) {
		eff, mActions, mCells, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return t == actions.ActionTypeRollWheel
		}

		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return &cellsMocks.Cell{
				CellInfo:        model.RestoreCellInfo(model.CellData{}),
				CategoriesValue: []string{"game"},
			}, nil
		}

		if !eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return true")
		}
	})

	t.Run("cannot roll wheel", func(t *testing.T) {
		eff, mActions, _, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return false
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})

	t.Run("wrong cell category", func(t *testing.T) {
		eff, mActions, mCells, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return &cellsMocks.Cell{
				CategoriesValue: []string{"other"},
			}, nil
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false when cell is not in 'game' category")
		}
	})

	t.Run("has developers filter", func(t *testing.T) {
		eff, mActions, mCells, mFilters := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return &cellsMocks.Cell{
				CellInfo: model.RestoreCellInfo(model.CellData{
					Filter: "filter1",
				}),
				CategoriesValue: []string{"game"},
			}, nil
		}

		mFilters.GetByIDFunc = func(ctx context.Context, id string) (*model.ActivityFilter, error) {
			return model.RestoreActivityFilter(model.ActivityFilterData{
				Developers: []string{"dev1"},
			}), nil
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
		*cellsMocks.Cells,
		*genresMocks.Genres,
	) {
		mActions := &actionsMocks.Actions{}
		mCells := &cellsMocks.Cells{}
		mGenres := &genresMocks.Genres{}
		mFilters := &filtersMocks.ActivityFilters{}

		eff := &AddGameGenre{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:         mActions,
			cells:           mCells,
			genres:          mGenres,
			activityFilters: mFilters,
		}

		events := model.NewEvents()
		action := model.RestoreAction(model.ActionData{})
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

		return eff, events, player, &callbackCalled, callback, mCells, mGenres
	}

	effectCtx := model.EffectContext{
		InvItemID: "item1",
		Priority:  10,
	}

	t.Run("successful activation", func(t *testing.T) {
		eff, events, player, called, callback, mCells, mGenres := setup()
		genreID := "rpg"

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		mCell := &cellsMocks.RollableCell{
			Cell: &cellsMocks.Cell{CategoriesValue: []string{"game"}},
		}

		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return mCell, nil
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

		if !mCell.RefreshCalled {
			t.Error("RefreshItems was not called on cell")
		}

		genres := player.LastAction().CustomActivityFilter().Genres
		if !slices.Contains(genres, genreID) {
			t.Errorf("Genre %s was not added to player filter: %v", genreID, genres)
		}
	})
}
