package board

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/mathhelper"
	"context"
	"errors"
)

type pathMovement struct {
	worlds worlds
	cells  cellsService

	steps int
}

func newPathMovement(worlds worlds, cells cellsService, steps int) *pathMovement {
	return &pathMovement{
		worlds: worlds,
		cells:  cells,
		steps:  steps,
	}
}

func (p *pathMovement) Move(ctx context.Context, actor actors.Actor) (*moveResult, error) {
	return p.move(ctx, actor, p.steps, 0)
}

func (p *pathMovement) move(ctx context.Context, actor actors.Actor, steps, laps int) (*moveResult, error) {
	position := actor.Position()
	currentWorldId := position.WorldId

	world, err := p.worlds.GetByID(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	cellsCount, err := p.cells.CountLocal(ctx, currentWorldId)
	if err != nil {
		return nil, err
	}

	if cellsCount == 0 {
		return nil, errors.New("world has no cells")
	}

	currentCellLocalOrder := position.CellLocalOrder
	totalSteps := currentCellLocalOrder + steps

	if !world.IsLoop() && totalSteps >= cellsCount {
		stepsToLastCell := cellsCount - 1 - currentCellLocalOrder

		var moves []*model.MoveResult

		if stepsToLastCell > 0 {
			lastCell, err := p.cells.GetByLocalOrder(ctx, currentWorldId, cellsCount-1)
			if err != nil {
				return nil, err
			}

			moves = append(moves, &model.MoveResult{
				Type:           model.MoveTypePath,
				TotalSteps:     cellsCount - 1,
				PrevTotalSteps: currentCellLocalOrder,
				World:          world,
				Cell:           lastCell,
			})
		}

		var nextWorld *model.World
		if world.TransitionToWorld() != "" {
			nextWorld, err = p.worlds.GetByID(ctx, world.TransitionToWorld())
		} else {
			nextWorld, err = p.worlds.GetDefault(ctx)
		}
		if err != nil {
			return nil, err
		}

		firstCell, err := p.cells.GetByLocalOrder(ctx, nextWorld.ID(), 0)
		if err != nil {
			return nil, err
		}

		transitionMove := &model.MoveResult{
			Type:  model.MoveTypeWorldTransition,
			World: nextWorld,
			Cell:  firstCell,
		}

		moves = append(moves, transitionMove)

		err = actor.ChangeWorld(ctx, nextWorld.ID())
		if err != nil {
			return nil, err
		}

		remainingSteps := steps - stepsToLastCell - 1
		if remainingSteps == 0 {
			return &moveResult{
				ActorMovement: model.ActorMovement{
					Position: model.ActorPosition{
						WorldId:        transitionMove.World.ID(),
						CellId:         transitionMove.Cell.ID(),
						CellLocalOrder: transitionMove.Cell.LocalOrder(),
					},
					Steps: p.steps,
					Laps:  laps,
				},
				Moves: moves,
			}, nil
		}

		res, err := p.move(ctx, actor, remainingSteps, laps)
		if err != nil {
			return nil, err
		}

		res.Moves = append(moves, res.Moves...)

		return res, nil
	}

	currentCellOrder := mathhelper.Mod(totalSteps, cellsCount)
	lapsPassed := mathhelper.FloorDiv(totalSteps, cellsCount) - mathhelper.FloorDiv(currentCellLocalOrder, cellsCount)

	currentCell, err := p.cells.GetByLocalOrder(ctx, currentWorldId, currentCellOrder)
	if err != nil {
		return nil, err
	}

	return &moveResult{
		ActorMovement: model.ActorMovement{
			Position: model.ActorPosition{
				WorldId:        world.ID(),
				CellId:         currentCell.ID(),
				CellLocalOrder: currentCell.LocalOrder(),
			},
			Steps: p.steps,
			Laps:  laps + lapsPassed,
		},
		Moves: []*model.MoveResult{
			{
				Type:           model.MoveTypePath,
				TotalSteps:     totalSteps,
				PrevTotalSteps: currentCellLocalOrder,
				World:          world,
				Cell:           currentCell,
				Laps:           lapsPassed,
			},
		},
	}, nil
}
