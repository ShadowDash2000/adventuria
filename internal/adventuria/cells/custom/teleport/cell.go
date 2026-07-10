package teleport

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
}

type board interface {
	MoveToCellId(ctx context.Context, events *model.Events, player *model.Player, cellId string) ([]*model.MoveResult, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

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

func (c *CellTeleport) OnCellReached(ctx context.Context, events *model.Events, player *model.Player, reachedCtx *model.ReachedContext) error {
	decodedValue, err := c.decodeValue(c.Value())
	if err != nil {
		return err
	}

	onBeforeTeleportOnCell := &model.OnBeforeTeleportOnCellEvent{
		CellId:       decodedValue.CellId,
		SkipTeleport: false,
	}

	err = events.OnBeforeTeleportOnCell().Trigger(ctx, onBeforeTeleportOnCell)
	if err != nil {
		return err
	}

	if onBeforeTeleportOnCell.SkipTeleport {
		player.LastAction().SetCanMove(true)
		return nil
	}

	player.LastAction().SetType(actions.ActionTypeTeleport)
	newAction, err := c.actions.Save(ctx, player.LastAction())
	if err != nil {
		return err
	}

	player.SetLastAction(newAction)

	moves, err := c.board.MoveToCellId(ctx, events, player, decodedValue.CellId)
	if err != nil {
		return err
	}

	reachedCtx.Moves = append(reachedCtx.Moves, moves...)

	return nil
}

func (c *CellTeleport) OnCellLeft(_ context.Context, _ *model.Events, _ *model.Player) error {
	return nil
}
