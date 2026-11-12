package board

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/mathhelper"
	"context"
	"fmt"
)

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetByIDWrapped(ctx context.Context, id string) (cells.Cell, error)
	GetByLocalOrderWrapped(ctx context.Context, worldId string, order int) (cells.Cell, error)
	GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error)
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

func (b *Board) Move(ctx context.Context, actor actors.Actor, steps int, moveType model.MoveType) ([]*model.MoveResult, error) {
	if isMovementStarted(ctx) {
		return b.move(ctx, actor, steps, moveType)
	}

	ctx = withMovementStarted(ctx)
	moves, err := b.move(ctx, actor, steps, moveType)
	if err != nil {
		return nil, err
	}

	err = actor.MovementFinished(ctx, moves)
	if err != nil {
		return nil, err
	}

	return moves, nil
}

func (b *Board) move(ctx context.Context, actor actors.Actor, steps int, moveType model.MoveType) ([]*model.MoveResult, error) {
	prevCell, err := b.cells.GetByIDWrapped(ctx, actor.Position().CellId)
	if err != nil {
		return nil, err
	}

	err = actor.LeaveCell(ctx, prevCell)
	if err != nil {
		return nil, err
	}

	currentWorldId := actor.Position().WorldId
	world, err := b.worlds.GetByID(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	cellsPassed := actor.Position().CellsPassed
	cellsCount, err := b.cells.CountLocal(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	totalSteps := cellsPassed + steps

	if !world.IsLoop() && totalSteps >= cellsCount {
		var nextWorld *model.World
		if world.TransitionToWorld() != "" {
			nextWorld, err = b.worlds.GetByID(ctx, world.TransitionToWorld())
			if err != nil {
				return nil, err
			}
		} else {
			nextWorld, err = b.worlds.GetDefault(ctx)
			if err != nil {
				return nil, err
			}
		}

		err = actor.ChangeWorld(ctx, nextWorld.ID())
		if err != nil {
			return nil, err
		}

		return b.Move(ctx, actor, 0, model.MoveTypeWorldTransition)
	}

	currentCellNum := mathhelper.Mod(totalSteps, cellsCount)
	lapsPassed := mathhelper.FloorDiv(totalSteps, cellsCount) - mathhelper.FloorDiv(cellsPassed, cellsCount)

	currentCell, err := b.cells.GetByLocalOrderWrapped(ctx, currentWorldId, currentCellNum)
	if err != nil {
		return nil, err
	}

	move := &model.MoveResult{
		Type:           moveType,
		Steps:          steps,
		TotalSteps:     totalSteps,
		PrevTotalSteps: cellsPassed,
		CurrentCell:    currentCell.Data(),
		CurrentWorld:   world,
		Laps:           lapsPassed,
	}
	err = actor.ApplyMovement(ctx, move)
	if err != nil {
		return nil, err
	}

	cellReachedCtx := cells.ReachedContext{
		Moves: []*model.MoveResult{move},
	}
	err = actor.ReachCell(ctx, currentCell, &cellReachedCtx)
	if err != nil {
		return nil, err
	}

	err = actor.AfterMovement(ctx, move)
	if err != nil {
		return nil, err
	}

	return cellReachedCtx.Moves, nil
}

func (b *Board) MoveToCell(ctx context.Context, actor actors.Actor, cell *model.CellInfo) ([]*model.MoveResult, error) {
	if cell.Disabled() {
		return nil, errs.ErrCellIsDisabled
	}

	cellWorldId := cell.World()
	currentWorld := actor.Position().WorldId
	if cellWorldId != currentWorld {
		err := actor.ChangeWorld(ctx, cellWorldId)
		if err != nil {
			return nil, err
		}
	}

	cellsCount, err := b.cells.CountLocal(ctx, actor.Position().WorldId)
	if err != nil {
		return nil, err
	}

	currentCellOrder := mathhelper.Mod(actor.Position().CellsPassed, cellsCount)

	return b.Move(ctx, actor, cell.LocalOrder()-currentCellOrder, model.MoveTypeTeleport)
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
