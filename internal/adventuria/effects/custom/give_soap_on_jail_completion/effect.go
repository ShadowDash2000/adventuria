package give_soap_on_jail_completion

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_events/custom/item_received"
	"adventuria/pkg/event"
	"context"
	"errors"
)

type inventories interface {
	AddItemBySlug(ctx context.Context, events *model.Events, player *model.Player, itemSlug string, source item_received.Source) (*model.InventoryItemWithEffects, error)
}

var _ model.EffectPersistent = (*GiveSoapOnJailCompletion)(nil)

const Type model.EffectType = "give_soap_on_jail_completion"

type GiveSoapOnJailCompletion struct {
	inventories inventories
}

func NewDef(inventories inventories) effects.EffectPersistentDef {
	return effects.NewEffectPersistentDef(
		Type,
		&GiveSoapOnJailCompletion{
			inventories: inventories,
		},
	)
}

func (g *GiveSoapOnJailCompletion) Subscribe(_ context.Context, events *model.Events, player *model.Player) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		events.OnAfterDone().BindFunc(func(ctx context.Context, e *model.OnAfterDoneEvent) error {
			if e.CurrentCell.Type() != cells.CellTypeJail {
				return e.Next()
			}

			_, err := g.inventories.AddItemBySlug(
				ctx,
				events,
				player,
				"soap",
				item_received.SourceCell,
			)
			if err != nil {
				if errors.Is(err, errs.ErrItemNotFound) {
					return e.Next()
				}
				return err
			}

			return e.Next()
		}),
	}, nil
}
