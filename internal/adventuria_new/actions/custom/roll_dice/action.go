package roll_dice

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/model"
	"context"
	"strconv"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

type board interface {
	Move(ctx context.Context, events *model.Events, player *model.Player, steps int) ([]*model.MoveResult, error)
}

var _ model.Action = (*RollDice)(nil)

const Type model.ActionType = "roll_dice"

type RollDice struct {
	actions.ActionBase
	cells   cells
	actions actionsService
	board   board
}

func NewDef(cells cells, actionsService actionsService, board board) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RollDice{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				actions:    actionsService,
				board:      board,
			}
		},
	)
}

func (r *RollDice) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().CanMove()
}

type RollDiceResult struct {
	Roll      int                 `json:"roll"`
	DiceRolls []DiceRoll          `json:"dice_rolls"`
	Path      []*model.MoveResult `json:"path"`
	From      PositionSnapshot    `json:"from"`
	To        PositionSnapshot    `json:"to"`
	PathSteps []PathStep          `json:"path_steps"`
}

type PositionSnapshot struct {
	WorldId     string `json:"world_id"`
	CellsPassed int    `json:"cells_passed"`
}

type PathStep struct {
	WorldId        string `json:"world_id"`
	WorldSlug      string `json:"world_slug"`
	CellOrder      int    `json:"cell_order"`
	TotalSteps     int    `json:"total_steps"`
	PrevTotalSteps int    `json:"prev_total_steps"`
	Event          string `json:"event"`
}

type DiceRoll struct {
	Type string `json:"type"`
	Roll int    `json:"roll"`
}

func (r *RollDice) Do(ctx context.Context, events *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	onBeforeRollEvent := &model.OnBeforeRollEvent{
		Dices: []model.Dice{model.DiceD6(), model.DiceD6()},
	}
	err := events.OnBeforeRoll().Trigger(ctx, onBeforeRollEvent)
	if err != nil {
		return nil, err
	}

	onBeforeRollMoveEvent := &model.OnBeforeRollMoveEvent{
		N: 0,
	}
	diceRolls := make([]DiceRoll, len(onBeforeRollEvent.Dices))
	for i, dice := range onBeforeRollEvent.Dices {
		diceRolls[i] = DiceRoll{
			Type: "d" + strconv.Itoa(int(dice)),
			Roll: dice.Roll(),
		}
		onBeforeRollMoveEvent.N += diceRolls[i].Roll
	}
	err = events.OnBeforeRollMove().Trigger(ctx, onBeforeRollMoveEvent)
	if err != nil {
		return nil, err
	}

	// we need to save the latest action before Move(), because it creates a new one
	newAction, err := r.actions.Save(ctx, player.LastAction())
	if err != nil {
		return nil, err
	}
	player.SetLastAction(newAction)

	from := PositionSnapshot{
		WorldId:     player.Progress().CurrentWorld(),
		CellsPassed: player.Progress().CellsPassed(),
	}
	moves, err := r.board.Move(ctx, events, player, onBeforeRollMoveEvent.N)
	if err != nil {
		return nil, err
	}
	to := PositionSnapshot{
		WorldId:     player.Progress().CurrentWorld(),
		CellsPassed: player.Progress().CellsPassed(),
	}

	steps := make([]PathStep, 0, len(moves))
	for i, move := range moves {
		eventType := "move"
		if i > 0 && moves[i-1].CurrentWorld.ID() != move.CurrentWorld.ID() {
			eventType = "world_transition"
		}

		steps = append(steps, PathStep{
			WorldId:        move.CurrentWorld.ID(),
			WorldSlug:      move.CurrentWorld.Slug(),
			CellOrder:      move.CellLocalOrder,
			TotalSteps:     move.TotalSteps,
			PrevTotalSteps: move.PrevTotalSteps,
			Event:          eventType,
		})
	}

	player.LastAction().SetType(Type)

	err = events.OnAfterRoll().Trigger(ctx, &model.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	})
	if err != nil {
		return nil, err
	}

	return RollDiceResult{
		Roll:      onBeforeRollMoveEvent.N,
		DiceRolls: diceRolls,
		Path:      moves,
		From:      from,
		To:        to,
		PathSteps: steps,
	}, nil
}
