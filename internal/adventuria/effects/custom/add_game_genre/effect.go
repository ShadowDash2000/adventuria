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

type genres interface {
	Exists(ctx context.Context, id string) (bool, error)
}

type activities interface {
	GetByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ model.Effect = (*AddGameGenre)(nil)

const Type model.EffectType = "add_game_genre"

type AddGameGenre struct {
	effects.EffectBase
	actions    actionsService
	genres     genres
	activities activities
}

func NewDef(actions actionsService, genres genres, activities activities) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &AddGameGenre{
				EffectBase: effects.NewEffectBase(effect),
				actions:    actions,
				genres:     genres,
				activities: activities,
			}
		},
	)
}

func (a *AddGameGenre) CanUse(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !a.actions.CanDo(ctx, events, player, actions.ActionTypeRollWheel) {
		return false
	}

	activityFilter := player.LastAction().State().ActivityFilter
	if activityFilter == nil {
		return false
	}

	if activityFilter.Type != model.ActivityTypeGame {
		return false
	}

	if len(activityFilter.Developers) > 0 {
		return false
	}
	if len(activityFilter.Publishers) > 0 {
		return false
	}
	if len(activityFilter.Activities) > 0 {
		return false
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

			genreId, ok := e.Data["genre_id"].(string)
			if !ok {
				return errors.New("genre_id not specified")
			}

			ok, err := a.genres.Exists(ctx, genreId)
			if err != nil {
				return err
			}
			if !ok {
				return errs.ErrGenreNotFound
			}

			actionState := player.LastAction().State()
			if actionState.ActivityFilter == nil {
				return errs.ErrNoActiveActivityFilter
			}

			if index := slices.Index(actionState.ActivityFilter.Genres, genreId); index != -1 {
				return errors.New("genre already exists")
			}

			actionState.ActivityFilter.Genres = append(actionState.ActivityFilter.Genres, genreId)

			ids, err := a.activities.GetByFilter(ctx, *actionState.ActivityFilter)
			if err != nil {
				return err
			}

			actionState.Activities = &model.ActionActivitiesState{
				Ids: ids,
			}
			player.LastAction().SetState(actionState)

			callback(ctx)
			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
