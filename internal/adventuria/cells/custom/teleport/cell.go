package teleport

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
}

type boardService interface {
	MoveToCellId(ctx context.Context, actor actors.Actor, cellId string) ([]*model.MoveResult, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ actors.ActorCellReached = (*CellTeleport)(nil)
var _ player.CellReached = (*CellTeleport)(nil)

const Type model.CellType = "teleport"

type CellTeleport struct {
	cells.CellBase
	payload Payload
	cells   cellsService
	board   boardService
	actions actionsService
}

func NewDef(
	cellsService cellsService,
	boardService boardService,
	actions actionsService,
	categories ...string,
) cells.CellDef {
	return cells.NewCell(
		Type,
		func(cell model.CellInfo) (cells.Cell, error) {
			return New(cellsService, boardService, actions, cell)
		},
		categories...,
	)
}

func New(cellsService cellsService, boardService boardService, actionsService actionsService, cell model.CellInfo) (*CellTeleport, error) {
	payload, err := NewPayloadFromJSON([]byte(cell.Value()))
	if err != nil {
		return nil, err
	}

	return &CellTeleport{
		CellBase: cells.NewCellBase(cell),
		payload:  payload,
		cells:    cellsService,
		board:    boardService,
		actions:  actionsService,
	}, nil
}

func (c *CellTeleport) OnPlayerReached(ctx context.Context, actor *player.Actor, reached *cells.ReachedContext) error {
	onBeforeTeleportOnCell := &model.OnBeforeTeleportOnCellEvent{
		CellId:       c.payload.CellId,
		SkipTeleport: false,
	}
	err := actor.Events().OnBeforeTeleportOnCell().Trigger(ctx, onBeforeTeleportOnCell)
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

func (c *CellTeleport) OnActorReached(ctx context.Context, actor actors.Actor, reached *cells.ReachedContext) error {
	moves, err := c.board.MoveToCellId(ctx, actor, c.payload.CellId)
	if err != nil {
		return err
	}

	reached.Moves = append(reached.Moves, moves...)

	return nil
}
