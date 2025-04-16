package adventuria

import (
	"adventuria/pkg/helper"
	"errors"
)

type CellJail struct {
	CellBase
}

func NewCellJail() CellCreator {
	return func() Cell {
		return &CellJail{}
	}
}

func (c *CellJail) NextStep(user *User) string {
	nextStepType := ""

	if user.IsInJail() {
		switch user.LastAction.Type() {
		case ActionTypeRollDice,
			ActionTypeReroll:
			nextStepType = ActionTypeRollWheel
		case ActionTypeRollWheel:
			nextStepType = ActionTypeChooseResult
		case ActionTypeChooseResult:
			nextStepType = ActionTypeRollDice
		default:
			nextStepType = ActionTypeRollWheel
		}
	} else {
		nextStepType = ActionTypeRollDice
	}

	return nextStepType
}

func (c *CellJail) Roll(_ *User) (*WheelRollResult, error) {
	cellsCol, err := GameCollections.Get(TableCells)
	if err != nil {
		return nil, err
	}

	res := &WheelRollResult{
		Collection: cellsCol,
	}

	cells := GameCells.GetAllByType(CellTypeGame)

	if len(cells) == 0 {
		return nil, errors.New("game cells not found")
	}

	for _, cell := range cells {
		res.FillerItems = append(res.FillerItems, &WheelItem{
			Name: cell.Name(),
			Icon: cell.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(cells).ID()

	// TODO
	panic("implement me")

	return res, nil
}
