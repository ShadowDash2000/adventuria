package cells

import "adventuria/internal/adventuria"

type CellGame struct {
	adventuria.CellBase
}

func NewCellGame() adventuria.CellCreator {
	return func(_ adventuria.ServiceLocator) adventuria.Cell {
		return &CellGame{
			adventuria.CellBase{},
		}
	}
}

func (c *CellGame) NextStep(user adventuria.User) string {
	nextStepType := ""

	switch user.LastAction().Type() {
	case adventuria.ActionTypeRollDice,
		adventuria.ActionTypeReroll:
		nextStepType = adventuria.ActionTypeRollWheel
	case adventuria.ActionTypeRollWheel:
		nextStepType = adventuria.ActionTypeChooseResult
	case adventuria.ActionTypeChooseResult,
		adventuria.ActionTypeDrop:
		nextStepType = adventuria.ActionTypeRollDice
	default:
		nextStepType = adventuria.ActionTypeRollWheel
	}

	return nextStepType
}

func (c *CellGame) Roll(_ adventuria.User) (*adventuria.WheelRollResult, error) {
	res := &adventuria.WheelRollResult{}

	// TODO
	panic("implement me")

	return res, nil
}
