package cells

import (
	"adventuria/internal/adventuria"
)

type CellBigWin struct {
	CellPreset
}

func NewCellBigWin() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellBigWin{
			CellPreset: CellPreset{
				CellBase: adventuria.CellBase{},
			},
		}
	}
}

func (c *CellBigWin) NextStep(user adventuria.User) string {
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
