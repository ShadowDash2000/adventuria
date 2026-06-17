package nothing

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

var _ model.Effect = (*Nothing)(nil)

const Type model.EffectType = "nothing"

type Nothing struct {
	effects.EffectBase
	cells cellsService
}

func NewDef(cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &Nothing{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
			}
		},
	)
}

func (n *Nothing) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (n *Nothing) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	switch n.Value() {
	case useAfterItemAdd:
		return []event_new.Unsubscribe{
			events.OnAfterItemAdd().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemAddEvent) error {
				if e.Item.Inventory().ID() == effectCtx.InvItemID {
					callback(ctx)
				}
				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useAfterItemUse:
		return []event_new.Unsubscribe{
			events.OnAfterItemUse().BindFuncWithPriority(func(ctx context.Context, e *model.OnAfterItemUseEvent) error {
				if e.InvItemId == effectCtx.InvItemID {
					callback(ctx)
				}
				return e.Next()
			}, effectCtx.Priority),
		}, nil
	case useBeforeGameDone:
		return []event_new.Unsubscribe{
			events.OnBeforeDone().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeDoneEvent) error {
				currentCell, err := n.cells.GetCurrentCellByProgress(ctx, player.Progress())
				if err != nil {
					return err
				}

				if currentCell.InCategories([]string{"activity", "game"}) {
					callback(ctx)
				}

				return e.Next()
			}, effectCtx.Priority),
		}, nil
	}
	return nil, nil
}
