package give_wheel_on_done

import (
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

var _ model.EffectPersistent = (*GiveWheelOnDone)(nil)

const Type model.EffectType = "give_wheel_on_done"

type notify interface {
	ItemWheelsChangeNotify(ctx context.Context, id string, amount int) error
}

type GiveWheelOnDone struct {
	notify notify
}

func NewDef(notify notify) effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&GiveWheelOnDone{
			notify: notify,
		},
	)
}

func (g *GiveWheelOnDone) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			if e.CurrentCell.DontGiveItemWheel() {
				return e.Next()
			}

			err := player.Progress().ItemWheelsCountChange(1)
			if err != nil {
				return err
			}

			err = g.notify.ItemWheelsChangeNotify(ctx, player.Progress().ID(), 1)
			if err != nil {
				return err
			}

			return e.Next()
		}),
	}, nil
}
