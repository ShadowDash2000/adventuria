package save_from_hole

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"context"
)

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
}

var _ effects.Effect = (*SaveFromHole)(nil)

const Type model.EffectType = "save_from_hole"

type SaveFromHole struct {
	effects.EffectBase
	cells cellsService
}

func NewDef(cells cellsService) effects.EffectDef {
	return effects.NewEffectDef(
		Type,
		func(effect model.EffectInfo) effects.Effect {
			return &SaveFromHole{
				EffectBase: effects.NewEffectBase(effect),
				cells:      cells,
			}
		},
	)
}

func (s *SaveFromHole) CanUse(_ context.Context, _ *player.Actor) bool {
	return true
}

func (s *SaveFromHole) Subscribe(
	_ context.Context,
	actor *player.Actor,
	effectCtx effects.Context,
	callback effects.Callback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		actor.Events().OnBeforeTeleportOnCell().BindFuncWithPriority(func(ctx context.Context, e *model.OnBeforeTeleportOnCellEvent) error {
			if e.SkipTeleport {
				return e.Next()
			}

			currentCell, err := s.cells.GetByID(ctx, actor.Player().LastAction().Cell())
			if err != nil {
				return err
			}

			destinationCell, err := s.cells.GetByID(ctx, e.CellId)
			if err != nil {
				return err
			}

			if destinationCell.GlobalOrder() < currentCell.GlobalOrder() {
				e.SkipTeleport = true
				callback(ctx)
			}

			return e.Next()
		}, effectCtx.Priority),
	}, nil
}
