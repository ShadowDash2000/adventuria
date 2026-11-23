package board

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/mathhelper"
	"context"
	"errors"
	"fmt"
)

type cellsService interface {
	Wrap(cell *model.CellInfo) (cells.Cell, error)
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByIDWrapped(ctx context.Context, id string) (cells.Cell, error)
	GetByLocalOrder(ctx context.Context, worldId string, order int) (*model.CellInfo, error)
	GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
}

type worlds interface {
	GetByID(ctx context.Context, id string) (*model.World, error)
	GetDefault(ctx context.Context) (*model.World, error)
}

type Board struct {
	cells  cellsService
	worlds worlds
}

func NewBoard(cells cellsService, worlds worlds) *Board {
	return &Board{
		cells:  cells,
		worlds: worlds,
	}
}

type movementStartedCtxKey struct{}

func withMovementStarted(ctx context.Context) context.Context {
	return context.WithValue(ctx, movementStartedCtxKey{}, true)
}

func isMovementStarted(ctx context.Context) bool {
	_, ok := ctx.Value(movementStartedCtxKey{}).(bool)
	return ok
}

type moveResult struct {
	ActorMovement model.ActorMovement
	Moves         []*model.MoveResult
}

type movementHandler interface {
	Move(ctx context.Context, actor actors.Actor) (*moveResult, error)
}

func (b *Board) Move(ctx context.Context, actor actors.Actor, handler movementHandler) ([]*model.MoveResult, error) {
	if isMovementStarted(ctx) {
		return b.move(ctx, actor, handler)
	}

	ctx = withMovementStarted(ctx)
	moves, err := b.move(ctx, actor, handler)
	if err != nil {
		return nil, err
	}

	err = actor.MovementFinished(ctx, moves)
	if err != nil {
		return nil, err
	}

	return moves, nil
}

func (b *Board) move(ctx context.Context, actor actors.Actor, handler movementHandler) ([]*model.MoveResult, error) {
	prevCell, err := b.cells.GetByIDWrapped(ctx, actor.Position().CellId)
	if err != nil {
		return nil, err
	}

	err = actor.LeaveCell(ctx, prevCell)
	if err != nil {
		return nil, err
	}

	res, err := handler.Move(ctx, actor)
	if err != nil {
		return nil, err
	}

	err = actor.ApplyMovement(ctx, res.ActorMovement)
	if err != nil {
		return nil, err
	}

	if len(res.Moves) == 0 {
		return nil, errors.New("no moves")
	}

	finalMove := res.Moves[len(res.Moves)-1]

	finalCell, err := b.cells.Wrap(finalMove.Cell)
	if err != nil {
		return nil, err
	}

	reachedCtx := cells.ReachedContext{
		Moves: res.Moves,
	}
	err = actor.ReachCell(ctx, finalCell, &reachedCtx)
	if err != nil {
		return nil, err
	}

	err = actor.AfterMovement(ctx, res.ActorMovement)
	if err != nil {
		return nil, err
	}

	return reachedCtx.Moves, nil
}

func (b *Board) MovePath(ctx context.Context, actor actors.Actor, steps int) ([]*model.MoveResult, error) {
	return b.Move(ctx, actor, newPathMovement(
		b.worlds,
		b.cells,
		steps,
	))
}

func (b *Board) MoveToCell(ctx context.Context, actor actors.Actor, cell *model.CellInfo) ([]*model.MoveResult, error) {
	return b.Move(ctx, actor, newTeleportMovement(
		b.worlds,
		cell,
	))
}

func (b *Board) MoveToCellId(ctx context.Context, actor actors.Actor, cellId string) ([]*model.MoveResult, error) {
	cell, err := b.cells.GetByID(ctx, cellId)
	if err != nil {
		return nil, err
	}

	return b.MoveToCell(ctx, actor, cell)
}

func (b *Board) MoveToClosestCellType(ctx context.Context, actor actors.Actor, cellType model.CellType) ([]*model.MoveResult, error) {
	currentCell, err := b.cells.GetByIDWrapped(ctx, actor.Position().CellId)
	if err != nil {
		return nil, err
	}

	cellsOfType, err := b.cells.GetAllGlobalByType(ctx, cellType)
	if err != nil {
		return nil, err
	}

	var (
		closest     int
		minDistance int
		closestCell *model.CellInfo
	)
	for _, cellOfType := range cellsOfType {
		globalOrder := cellOfType.GlobalOrder()
		distance := mathhelper.Abs(globalOrder - currentCell.Data().GlobalOrder())
		if closestCell == nil ||
			distance < minDistance ||
			(distance == minDistance && globalOrder > closest) {
			closest = globalOrder
			minDistance = distance
			closestCell = cellOfType
		}
	}

	if closestCell == nil {
		return nil, fmt.Errorf("cell of type %s not found", cellType)
	}

	return b.MoveToCell(ctx, actor, closestCell)
}

func (b *Board) GetLocalDistanceBetweenCells(cell1, cell2 *model.CellInfo) (int, error) {
	if cell1.World() != cell2.World() {
		return 0, fmt.Errorf("cells are in different worlds")
	}
	return mathhelper.Abs(cell1.LocalOrder() - cell2.LocalOrder()), nil
}
