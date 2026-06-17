package give_wheel_on_done

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
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

func (g *GiveWheelOnDone) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			err := player.Progress().ItemWheelsCountChange(1)
			if err != nil {
				return err
			}
			return e.Next()
		}),
	}, nil
}
