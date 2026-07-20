package stats

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"adventuria/pkg/mathhelper"
	"context"
)

type activityFilters interface {
	GetByID(ctx context.Context, id string) (*model.ActivityFilter, error)
}

var _ model.EffectPersistent = (*Stats)(nil)

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

func (s *Stats) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			activityFilter, err := s.activityFilters.GetByID(ctx, e.CurrentCell.Filter())
			if err != nil {
				return err
			}

			switch activityFilter.Type() {
			case model.ActivityTypeGame:
				player.Stats().GamesCompletedChange(1)
			case model.ActivityTypeMovie:
				player.Stats().MoviesCompletedChange(1)
			case model.ActivityTypeGym:
				player.Stats().GymsCompletedChange(1)
			case model.ActivityTypeKaraoke:
				player.Stats().KaraokeCompletedChange(1)
			}

			return e.Next()
		}),
		events.OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
			player.Stats().CellsPassedChange(mathhelper.Abs(e.Steps))
			return e.Next()
		}),
		events.OnAfterDrop().BindFunc(func(ctx context.Context, e *model.OnAfterDropEvent) error {
			player.Stats().DropsChange(1)
			return e.Next()
		}),
		events.OnAfterGoToJail().BindFunc(func(ctx context.Context, e *model.OnAfterGoToJailEvent) error {
			player.Stats().WasInJailChange(1)
			return e.Next()
		}),
		events.OnAfterItemUse().BindFunc(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
			player.Stats().ItemsUsedChange(1)
			return e.Next()
		}),
		events.OnAfterRoll().BindFunc(func(ctx context.Context, e *model.OnAfterRollEvent) error {
			player.Stats().DiceRollsChange(1)
			return e.Next()
		}),
		events.OnAfterWheelRoll().BindFunc(func(ctx context.Context, e *model.OnAfterWheelRollEvent) error {
			player.Stats().WheelsRolledChange(1)
			return e.Next()
		}),
	}, nil
}
