package change_game_price_filter

import (
	"adventuria/internal/adventuria/actions"
	actionsMocks "adventuria/internal/adventuria/actions/mocks"
	filtersMocks "adventuria/internal/adventuria/activity_filters/mocks"
	"adventuria/internal/adventuria/cells"
	cellsMocks "adventuria/internal/adventuria/cells/mocks"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
	"testing"
)

func TestChangeGamePriceFilter_CanUse(t *testing.T) {
	ctx := t.Context()

	setup := func() (*ChangeGamePriceFilter, *actionsMocks.Actions, *cellsMocks.Cells, *filtersMocks.ActivityFilters) {
		mActions := &actionsMocks.Actions{}
		mCells := &cellsMocks.Cells{}
		mFilters := &filtersMocks.ActivityFilters{}

		eff := &ChangeGamePriceFilter{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:   "eff1",
					Type: Type,
				}),
			),
			actions:         mActions,
			cells:           mCells,
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

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				Type: cells.CellTypeGame,
			}), nil
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

	t.Run("wrong cell type", func(t *testing.T) {
		eff, mActions, mCells, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				Type: cells.CellTypeStart,
			}), nil
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})

	t.Run("custom filter not allowed", func(t *testing.T) {
		eff, mActions, mCells, _ := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				Type:                     cells.CellTypeGame,
				IsCustomFilterNotAllowed: true,
			}), nil
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})

	t.Run("filter has activities", func(t *testing.T) {
		eff, mActions, mCells, mFilters := setup()
		player := model.RestorePlayer(model.PlayerData{}, &model.PlayerProgress{}, nil, nil)

		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				Type:   cells.CellTypeGame,
				Filter: "filter1",
			}), nil
		}

		mFilters.GetByIDFunc = func(ctx context.Context, id string) (*model.ActivityFilter, error) {
			return model.RestoreActivityFilter(model.ActivityFilterData{
				Activities: []string{"act1"},
			}), nil
		}

		if eff.CanUse(ctx, nil, player) {
			t.Error("CanUse should return false")
		}
	})
}

func TestChangeGamePriceFilter_Subscribe(t *testing.T) {
	ctx := t.Context()

	setup := func(value string) (
		*ChangeGamePriceFilter,
		*model.Events,
		*model.Player,
		*bool,
		model.EffectCallback,
		*cellsMocks.Cells,
	) {
		mActions := &actionsMocks.Actions{}
		mCells := &cellsMocks.Cells{}
		mFilters := &filtersMocks.ActivityFilters{}

		eff := &ChangeGamePriceFilter{
			EffectBase: effects.NewEffectBase(
				*model.RestoreEffectInfo(model.EffectData{
					Id:    "eff1",
					Type:  Type,
					Value: value,
				}),
			),
			actions:         mActions,
			cells:           mCells,
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

		return eff, events, player, &callbackCalled, callback, mCells
	}

	effectCtx := model.EffectContext{
		InvItemID: "item1",
		Priority:  10,
	}

	t.Run("usable min price", func(t *testing.T) {
		eff, events, player, called, callback, mCells := setup("100;min;usable")

		mCell := &cellsMocks.RollableCell{
			Cell: &cellsMocks.Cell{
				CellInfo: model.RestoreCellInfo(model.CellData{Type: cells.CellTypeGame}),
			},
		}
		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return mCell, nil
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterItemUse().Trigger(ctx, &model.OnAfterItemUseEvent{
			InvItemId: "item1",
		})
		if err != nil {
			t.Errorf("Trigger failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called")
		}

		filter := player.LastAction().CustomActivityFilter()
		if filter.MinPrice != 100 || filter.MaxPrice != -1 {
			t.Errorf("Invalid filter prices: %+v", filter)
		}

		if !mCell.RefreshCalled {
			t.Error("RefreshItems was not called")
		}
	})

	t.Run("unusable max price", func(t *testing.T) {
		eff, events, player, called, callback, mCells := setup("500;max;unusable")

		mActions := eff.actions.(*actionsMocks.Actions)
		mActions.CanDoFunc = func(ctx context.Context, events *model.Events, p *model.Player, t model.ActionType) bool {
			return true
		}

		mCell := &cellsMocks.RollableCell{
			Cell: &cellsMocks.Cell{
				CellInfo: model.RestoreCellInfo(model.CellData{Type: cells.CellTypeGame}),
			},
		}
		mCells.GetByPlayerWrappedFunc = func(ctx context.Context, player *model.Player) (model.Cell, error) {
			return mCell, nil
		}
		mCells.GetByPlayerFunc = func(ctx context.Context, player *model.Player) (*model.CellInfo, error) {
			return model.RestoreCellInfo(model.CellData{
				Type: cells.CellTypeGame,
			}), nil
		}

		_, err := eff.Subscribe(ctx, events, player, effectCtx, callback)
		if err != nil {
			t.Fatalf("Subscribe failed: %v", err)
		}

		err = events.OnAfterMove().Trigger(ctx, &model.OnAfterMoveEvent{})
		if err != nil {
			t.Errorf("Trigger OnAfterMove failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called on OnAfterMove")
		}

		filter := player.LastAction().CustomActivityFilter()
		if filter.MaxPrice != 500 || filter.MinPrice != -1 {
			t.Errorf("Invalid filter prices: %+v", filter)
		}

		*called = false
		player.LastAction().SetCustomActivityFilter(model.CustomActivityFilter{})

		item := model.RestoreInventoryItem(
			model.RestoreInventory(model.InventoryData{Id: "item1"}),
			nil,
		)

		err = events.OnAfterItemAdd().Trigger(ctx, &model.OnAfterItemAddEvent{
			Item: item,
		})
		if err != nil {
			t.Errorf("Trigger OnAfterItemAdd failed: %v", err)
		}

		if !*called {
			t.Error("Callback was not called on OnAfterItemAdd")
		}
	})
}

func TestChangeGamePriceFilter_Verify(t *testing.T) {
	ctx := t.Context()
	eff := &ChangeGamePriceFilter{}

	t.Run("success min", func(t *testing.T) {
		err := eff.Verify(ctx, "100;min;usable")
		if err != nil {
			t.Errorf("Verify failed: %v", err)
		}
	})

	t.Run("success max unusable", func(t *testing.T) {
		err := eff.Verify(ctx, "500;max;unusable")
		if err != nil {
			t.Errorf("Verify failed: %v", err)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		err := eff.Verify(ctx, "100;min")
		if err == nil {
			t.Error("Verify should fail for invalid format")
		}
	})

	t.Run("invalid price", func(t *testing.T) {
		err := eff.Verify(ctx, "abc;min;usable")
		if err == nil {
			t.Error("Verify should fail for invalid price")
		}
	})

	t.Run("invalid price type", func(t *testing.T) {
		err := eff.Verify(ctx, "100;middle;usable")
		if err == nil {
			t.Error("Verify should fail for invalid price type")
		}
	})

	t.Run("invalid use type", func(t *testing.T) {
		err := eff.Verify(ctx, "100;min;always")
		if err == nil {
			t.Error("Verify should fail for invalid use type")
		}
	})
}
