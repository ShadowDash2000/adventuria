package give_wheel_on_new_lap

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ effects.EffectPersistent = (*GiveWheelOnNewLap)(nil)

const Type model.EffectType = "give_wheel_on_new_lap"

type notify interface {
	ItemWheelsChangeNotify(ctx context.Context, id string, amount int) error
}

type GiveWheelOnNewLap struct {
	notify notify
}

func NewDef(notify notify) effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&GiveWheelOnNewLap{
			notify: notify,
		},
	)
}

func (g *GiveWheelOnNewLap) Subscribe(_ context.Context, actor *player.Actor) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnNewLap().BindFunc(func(ctx context.Context, e *model.OnNewLapEvent) error {
			err := actor.Player().Progress().ItemWheelsCountChange(e.Laps)
			if err != nil {
				return err
			}

			err = g.notify.ItemWheelsChangeNotify(ctx, actor.Player().Progress().ID(), e.Laps)
			if err != nil {
				return err
			}

			return e.Next()
		}),
	}, nil
}
