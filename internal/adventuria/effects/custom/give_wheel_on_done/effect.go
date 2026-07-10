package give_wheel_on_done

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.EffectPersistent = (*GiveWheelOnDone)(nil)

const Type model.EffectType = "give_wheel_on_done"

type GiveWheelOnDone struct{}

func NewDef() effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&GiveWheelOnDone{},
	)
}

func (g *GiveWheelOnDone) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			err := player.Progress().ItemWheelsCountChange(1)
			if err != nil {
				return err
			}
			return e.Next()
		}),
	}, nil
}
