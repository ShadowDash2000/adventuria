package give_wheel_on_new_lap

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.EffectPersistent = (*GiveWheelOnNewLap)(nil)

const Type model.EffectType = "give_wheel_on_new_lap"

type GiveWheelOnNewLap struct{}

func NewDef() effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&GiveWheelOnNewLap{},
	)
}

func (g *GiveWheelOnNewLap) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnNewLap().BindFunc(func(ctx context.Context, e *model.OnNewLapEvent) error {
			err := player.Progress().ItemWheelsCountChange(e.Laps)
			if err != nil {
				return err
			}
			return e.Next()
		}),
	}, nil
}
