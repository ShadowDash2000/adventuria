package save_from_hole

import (
	"adventuria/internal/adventuria_new/effects"
	"adventuria/internal/adventuria_new/model"
	"adventuria/pkg/event_new"
	"context"
)

type cellsService interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
}

var _ model.Effect = (*SaveFromHole)(nil)

const Type model.EffectType = "save_from_hole"

type SaveFromHole struct {
	effects.EffectBase
	cells cellsService
}

func NewDef(cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) model.Effect {
			return &SaveFromHole{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
			}
		},
	)
}

func (s *SaveFromHole) CanUse(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

func (s *SaveFromHole) Subscribe(
	_ context.Context,
	events *model.Events,
	player *model.Player,
	effectCtx model.EffectContext,
	callback model.EffectCallback,
) ([]event_new.Unsubscribe, error) {
	return []event_new.Unsubscribe{
		events.OnBeforeTeleportOnCell().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeTeleportOnCellEvent) error {
			if e.SkipTeleport {
				return e.Next()
			}

			currentCell, err := s.cells.GetCurrentCellByProgress(ctx, player.Progress())
			if err != nil {
				return err
			}

			destinationCell, err := s.cells.GetByID(ctx, e.CellId)
			if err != nil {
				return err
			}

			if destinationCell.GlobalOrder() < currentCell.Data().GlobalOrder() {
				e.SkipTeleport = true
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
