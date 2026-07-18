package roll_dice

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"strconv"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

type board interface {
	Move(ctx context.Context, events *model.Events, player *model.Player, steps int, moveType model.MoveType) ([]*model.MoveResult, error)
}

var _ model.Action = (*RollDice)(nil)

const Type model.ActionType = "roll_dice"

type RollDice struct {
	actions.ActionBase
	actions actionsService
	board   board
}

func NewDef(actionsService actionsService, board board) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &RollDice{
				ActionBase: actions.NewActionBase(Type),
				actions:    actionsService,
				board:      board,
			}
		},
	)
}

func (r *RollDice) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.Progress().CanMove()
}

type rollDiceResult struct {
	Roll      int        `json:"roll"`
	DiceRolls []diceRoll `json:"dice_rolls"`
	Moves     []move     `json:"moves"`
}

type move struct {
	Type            model.MoveType `json:"type"`
	WorldId         string         `json:"world_id"`
	WorldSlug       string         `json:"world_slug"`
	CellLocalOrder  int            `json:"cell_local_order"`
	CellGlobalOrder int            `json:"cell_global_order"`
	TotalSteps      int            `json:"total_steps"`
	PrevTotalSteps  int            `json:"prev_total_steps"`
}

type diceRoll struct {
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
	diceRolls := make([]diceRoll, len(onBeforeRollEvent.Dices))
	for i, dice := range onBeforeRollEvent.Dices {
		diceRolls[i] = diceRoll{
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

	moves, err := r.board.Move(ctx, events, player, onBeforeRollMoveEvent.N, model.MoveTypePath)
	if err != nil {
		return nil, err
	}

	movesRes := make([]move, len(moves))
	for i, m := range moves {
		movesRes[i] = move{
			Type:            m.Type,
			WorldId:         m.CurrentWorld.ID(),
			WorldSlug:       m.CurrentWorld.Slug(),
			CellLocalOrder:  m.CurrentCell.LocalOrder(),
			CellGlobalOrder: m.CurrentCell.GlobalOrder(),
			TotalSteps:      m.TotalSteps,
			PrevTotalSteps:  m.PrevTotalSteps,
		}
	}

	player.LastAction().SetType(Type)

	err = events.OnAfterRoll().Trigger(ctx, &model.OnAfterRollEvent{
		Dices: onBeforeRollEvent.Dices,
		N:     onBeforeRollMoveEvent.N,
	})
	if err != nil {
		return nil, err
	}

	return rollDiceResult{
		Roll:      onBeforeRollMoveEvent.N,
		DiceRolls: diceRolls,
		Moves:     movesRes,
	}, nil
}
