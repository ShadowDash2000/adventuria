package add_game_genre

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
	"errors"
	"slices"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type genres interface {
	Exists(ctx context.Context, id string) (bool, error)
}

type activityFilters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilter, error)
}

var _ model.Effect = (*AddGameGenreEffect)(nil)

const Type model.EffectType = "add_game_genre"

type AddGameGenreEffect struct {
	effects.EffectBase
	actions         actionsService
	cells           cells
	genres          genres
	activityFilters activityFilters
}

func NewAddGameGenreEffectDef(actions actionsService, cells cells, genres genres, activityFilters activityFilters) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &AddGameGenreEffect{
				EffectBase:      effects.NewEffectBase(effect),
				actions:         actions,
				cells:           cells,
				genres:          genres,
				activityFilters: activityFilters,
			}
		},
	)
}

func (a *AddGameGenreEffect) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !a.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	cell, err := a.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if !cell.InCategory("game") {
		return false
	}

	if filterId := cell.Data().Filter(); filterId != "" {
		filter, err := a.activityFilters.GetByID(ctx, filterId)
		if err != nil {
			return false
		}

		if len(filter.Developers()) > 0 {
			return false
		}
		if len(filter.Publishers()) > 0 {
			return false
		}
		if len(filter.Activities()) > 0 {
			return false
		}
	}

	return true
}

func (a *AddGameGenreEffect) Subscribe(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(e *model.OnAfterItemUseEvent) error {
			if effectCtx.InvItemID != e.InvItemId {
				return e.Next()
			}

			cell, err := a.cells.GetCurrentCellByProgress(ctx, player.Progress())
			if err != nil {
				return err
			}

			cellWheel, ok := cell.(model.Rollable)
			if !ok {
				return errors.New("current cell is not refreshable")
			}

			genreId, ok := e.Data["genre_id"].(string)
			if !ok {
				return errors.New("genre_id not specified")
			}

			ok, err = a.genres.Exists(ctx, genreId)
			if err != nil {
				return err
			}
			if !ok {
				return errs.ErrGenreNotFound
			}

			filter := player.LastAction().CustomActivityFilter()
			if index := slices.Index(filter.Genres, genreId); index != -1 {
				return errors.New("genre already exists")
			}

			filter.Genres = append(filter.Genres, genreId)
			player.LastAction().SetCustomActivityFilter(filter)

			err = cellWheel.RefreshItems(ctx, events, player)
			if err != nil {
				return err
			}

			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
