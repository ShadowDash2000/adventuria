package board

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/mathhelper"
	"context"
	"fmt"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

type playerProgress interface {
	Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
}

type cellsService interface {
	GetByID(ctx context.Context, id string) (*model.CellInfo, error)
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
	GetByLocalOrderWrapped(ctx context.Context, worldId string, order int) (model.Cell, error)
	GetByGlobalOrder(ctx context.Context, order int) (*model.CellInfo, error)
	GetAllGlobalByType(ctx context.Context, t model.CellType) ([]*model.CellInfo, error)
	CountLocal(ctx context.Context, worldId string) (int, error)
}

type worlds interface {
	GetByID(ctx context.Context, id string) (*model.World, error)
	GetDefault(ctx context.Context) (*model.World, error)
}

type Board struct {
	actions  actionsService
	progress playerProgress
	cells    cellsService
	worlds   worlds
}

func NewBoard(actions actionsService, progress playerProgress, cells cellsService, worlds worlds) *Board {
	return &Board{
		actions:  actions,
		progress: progress,
		cells:    cells,
		worlds:   worlds,
	}
}

func (b *Board) Move(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	steps int,
	moveType model.MoveType,
) ([]*model.MoveResult, error) {
	prevCell, err := b.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	err = prevCell.OnCellLeft(ctx, events, player)
	if err != nil {
		return nil, err
	}

	currentWorldId := player.Progress().CurrentWorld()
	world, err := b.worlds.GetByID(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	cellsPassed := player.Progress().CellsPassed()
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

		err = b.changeWorld(ctx, events, player, currentWorldId, nextWorld.ID())
		if err != nil {
			return nil, err
		}

		return b.Move(ctx, events, player, 0, model.MoveTypeWorldTransition)
	}

	currentCellNum := mathhelper.Mod(totalSteps, cellsCount)
	lapsPassed := mathhelper.FloorDiv(totalSteps, cellsCount) - mathhelper.FloorDiv(cellsPassed, cellsCount)

	currentCell, err := b.cells.GetByLocalOrderWrapped(ctx, currentWorldId, currentCellNum)
	if err != nil {
		return nil, err
	}

	err = player.Progress().CellsPassedChange(steps)
	if err != nil {
		return nil, err
	}

	newAction, err := model.NewAction(model.ActionCreate{
		Player: player.ID(),
		Cell:   currentCell.Data().ID(),
		Type:   actions.ActionTypeMove,
	})
	if err != nil {
		return nil, err
	}

	newAction.SetCellsPassed(steps)

	progress := player.Progress()
	progress.SetCanMove(false)

	newAction, err = b.actions.Save(ctx, newAction)
	if err != nil {
		return nil, err
	}
	newProgress, err := b.progress.Save(ctx, progress)
	if err != nil {
		return nil, err
	}

	player.SetLastAction(newAction)
	player.SetProgress(newProgress)

	onAfterMoveEvent := model.OnAfterMoveEvent{
		Steps:          steps,
		TotalSteps:     totalSteps,
		PrevTotalSteps: cellsPassed,
		CurrentCell:    currentCell.Data(),
		CurrentWorld:   world,
		Laps:           lapsPassed,
	}
	err = events.OnAfterMove().Trigger(ctx, &onAfterMoveEvent)
	if err != nil {
		return nil, err
	}

	cellReachedCtx := model.ReachedContext{
		Moves: []*model.MoveResult{
			{
				Type:           moveType,
				Steps:          onAfterMoveEvent.Steps,
				TotalSteps:     onAfterMoveEvent.TotalSteps,
				PrevTotalSteps: onAfterMoveEvent.PrevTotalSteps,
				CurrentCell:    onAfterMoveEvent.CurrentCell,
				CurrentWorld:   onAfterMoveEvent.CurrentWorld,
				Laps:           onAfterMoveEvent.Laps,
			},
		},
	}
	err = currentCell.OnCellReached(ctx, events, player, &cellReachedCtx)
	if err != nil {
		return nil, err
	}

	// Check if we're not moving backwards and passed new lap(-s)
	if steps > 0 && lapsPassed > 0 {
		err = events.OnNewLap().Trigger(ctx, &model.OnNewLapEvent{
			Laps: lapsPassed,
		})
		if err != nil {
			return nil, err
		}
	}

	return cellReachedCtx.Moves, nil
}

func (b *Board) MoveToCellId(ctx context.Context, events *model.Events, player *model.Player, cellId string) ([]*model.MoveResult, error) {
	cell, err := b.cells.GetByID(ctx, cellId)
	if err != nil {
		return nil, err
	}

	cellWorldId := cell.World()
	currentWorld := player.Progress().CurrentWorld()
	if cellWorldId != currentWorld {
		err = b.changeWorld(ctx, events, player, currentWorld, cellWorldId)
		if err != nil {
			return nil, err
		}
	}

	cellsCount, err := b.cells.CountLocal(ctx, currentWorld)
	if err != nil {
		return nil, err
	}

	currentCellOrder := mathhelper.Mod(player.Progress().CellsPassed(), cellsCount)

	return b.Move(ctx, events, player, cell.LocalOrder()-currentCellOrder, model.MoveTypeTeleport)
}

func (b *Board) changeWorld(ctx context.Context, events *model.Events, player *model.Player, oldWorldId, newWorldId string) error {
	player.Progress().SetCurrentWorld(newWorldId)
	player.Progress().SetCellsPassed(0)

	return events.OnWorldChanged().Trigger(ctx, &model.OnWorldChangedEvent{
		OldWorldId: oldWorldId,
		NewWorldId: newWorldId,
	})
}

func (b *Board) MoveToClosestCellType(
	ctx context.Context,
	events *model.Events,
	player *model.Player,
	cellType model.CellType,
) ([]*model.MoveResult, error) {
	currentCell, err := b.cells.GetCurrentCellByProgress(ctx, player.Progress())
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

	currentWorld := player.Progress().CurrentWorld()
	if closestCell.World() != currentWorld {
		err = b.changeWorld(ctx, events, player, currentWorld, closestCell.World())
		if err != nil {
			return nil, err
		}
	}

	return b.Move(ctx, events, player, closestCell.LocalOrder(), model.MoveTypeTeleport)
}

func (b *Board) GetLocalDistanceBetweenCells(cell1, cell2 *model.CellInfo) (int, error) {
	if cell1.World() != cell2.World() {
		return 0, fmt.Errorf("cells are in different worlds")
	}
	return mathhelper.Abs(cell1.LocalOrder() - cell2.LocalOrder()), nil
}
