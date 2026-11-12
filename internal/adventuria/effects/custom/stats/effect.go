package stats

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"adventuria/pkg/mathhelper"
	"context"
)

type activityFilters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilterInfo, error)
}

var _ effects.EffectPersistent = (*Stats)(nil)

const Type model.EffectType = "stats"

type Stats struct {
	activityFilters activityFilters
}

func NewDef(activityFilters activityFilters) effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&Stats{
			activityFilters: activityFilters,
		},
	)
}

func (s *Stats) Subscribe(_ context.Context, actor *player.Actor) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			activityFilter, err := s.activityFilters.GetByID(ctx, e.CurrentCell.Filter())
			if err != nil {
				return err
			}

			switch activityFilter.Type() {
			case model.ActivityTypeGame:
				actor.Player().Stats().GamesCompletedChange(1)
			case model.ActivityTypeMovie:
				actor.Player().Stats().MoviesCompletedChange(1)
			case model.ActivityTypeGym:
				actor.Player().Stats().GymsCompletedChange(1)
			case model.ActivityTypeKaraoke:
				actor.Player().Stats().KaraokeCompletedChange(1)
			}

			return e.Next()
		}),
		actor.Events().OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			actor.Player().Stats().CellsPassedChange(mathhelper.Abs(e.Steps))
			return e.Next()
		}),
		actor.Events().OnAfterDrop().BindFunc(func(ctx context.Context, e *model.OnAfterDropEvent) error {
			actor.Player().Stats().DropsChange(1)
			return e.Next()
		}),
		actor.Events().OnAfterGoToJail().BindFunc(func(ctx context.Context, e *model.OnAfterGoToJailEvent) error {
			actor.Player().Stats().WasInJailChange(1)
			return e.Next()
		}),
		actor.Events().OnAfterItemUse().BindFunc(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			actor.Player().Stats().ItemsUsedChange(1)
			return e.Next()
		}),
		actor.Events().OnAfterRoll().BindFunc(func(ctx context.Context, e *model.OnAfterRollEvent) error {
			actor.Player().Stats().DiceRollsChange(1)
			return e.Next()
		}),
		actor.Events().OnAfterWheelRoll().BindFunc(func(ctx context.Context, e *model.OnAfterWheelRollEvent) error {
			actor.Player().Stats().WheelsRolledChange(1)
			return e.Next()
		}),
	}, nil
}
