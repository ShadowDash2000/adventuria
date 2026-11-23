package board

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
)

type teleportMovement struct {
	worlds worlds

	cell *model.CellInfo
}

func newTeleportMovement(worlds worlds, cell *model.CellInfo) *teleportMovement {
	return &teleportMovement{
		worlds: worlds,
		cell:   cell,
	}
}

func (t *teleportMovement) Move(ctx context.Context, actor actors.Actor) (*moveResult, error) {
	if t.cell.Disabled() {
		return nil, errs.ErrCellIsDisabled
	}

	cellWorldId := t.cell.World()
	currentWorld := actor.Position().WorldId
	if cellWorldId != currentWorld {
		err := actor.ChangeWorld(ctx, cellWorldId)
		if err != nil {
			return nil, err
		}
	}

	world, err := t.worlds.GetByID(ctx, cellWorldId)
	if err != nil {
		return nil, err
	}

	return &moveResult{
		ActorMovement: model.ActorMovement{
			Position: model.ActorPosition{
				WorldId:        cellWorldId,
				CellId:         t.cell.ID(),
				CellLocalOrder: t.cell.LocalOrder(),
			},
			Steps: 0,
			Laps:  0,
		},
		Moves: []*model.MoveResult{
			{
				Type:           model.MoveTypeTeleport,
				TotalSteps:     0,
				PrevTotalSteps: 0,
				World:          world,
				Cell:           t.cell,
				Laps:           0,
			},
		},
	}, nil
}
