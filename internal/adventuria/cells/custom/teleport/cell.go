package teleport

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
}

type board interface {
	MoveToCellId(ctx context.Context, actor model.BoardActor, cellId string) ([]*model.MoveResult, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.ActorCellReached = (*CellTeleport)(nil)
var _ player.CellReached = (*CellTeleport)(nil)

const Type model.CellType = "teleport"

type CellTeleport struct {
	cells.CellBase
	cells   cellsService
	board   board
	actions actionsService
}

func NewDef(
	cellsService cellsService,
	boardService board,
	actions actionsService,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) model.Cell {
			return &CellTeleport{
				CellBase: cells.NewCellBase(cell),
				cells:    cellsService,
				board:    boardService,
				actions:  actions,
			}
		},
		categories...,
	)
}

func (c *CellTeleport) OnPlayerReached(ctx context.Context, actor *player.Actor, reached *model.ReachedContext) error {
	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	onBeforeTeleportOnCell := &model.OnBeforeTeleportOnCellEvent{
		CellId:       decodedValue.CellId,
		SkipTeleport: false,
	}
	err = actor.Events().OnBeforeTeleportOnCell().Trigger(ctx, onBeforeTeleportOnCell)
	if err != nil {
		return err
	}

	if onBeforeTeleportOnCell.SkipTeleport {
		actor.Player().Progress().SetCanMove(true)
		return nil
	}

	actor.Player().LastAction().SetStatus(model.ActionStatusTeleport)
	newAction, err := c.actions.Save(ctx, actor.Player().LastAction())
	if err != nil {
		return err
	}

	actor.Player().SetLastAction(newAction)

	return c.OnActorReached(ctx, actor, reached)
}

func (c *CellTeleport) OnActorReached(ctx context.Context, actor model.BoardActor, reached *model.ReachedContext) error {
	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	moves, err := c.board.MoveToCellId(ctx, actor, decodedValue.CellId)
	if err != nil {
		return err
	}

	reached.Moves = append(reached.Moves, moves...)

	return nil
}
