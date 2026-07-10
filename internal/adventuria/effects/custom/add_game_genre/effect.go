package add_game_genre

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
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

var _ model.Effect = (*AddGameGenre)(nil)

const Type model.EffectType = "add_game_genre"

type AddGameGenre struct {
	effects.EffectBase
	actions         actionsService
	cells           cells
	genres          genres
	activityFilters activityFilters
}

func NewDef(actions actionsService, cells cells, genres genres, activityFilters activityFilters) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &AddGameGenre{
				EffectBase:      effects.NewEffectBase(effect),
				actions:         actions,
				cells:           cells,
				genres:          genres,
				activityFilters: activityFilters,
			}
		},
	)
}

func (a *AddGameGenre) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !a.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	currentCell, err := a.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if !currentCell.InCategory("game") {
		return false
	}

	if filterId := currentCell.Data().Filter(); filterId != "" {
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

func (a *AddGameGenre) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			if effectCtx.InvItemID != e.InvItemId {
				return e.Next()
			}

			currentCell, err := a.cells.GetCurrentCellByProgress(ctx, player.Progress())
			if err != nil {
				return err
			}

			cellWheel, ok := currentCell.(model.Rollable)
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
