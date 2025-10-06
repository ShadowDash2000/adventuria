package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"
)

type CellJail struct {
	adventuria.CellBase
	locator adventuria.ServiceLocator
}

func NewCellJail() adventuria.CellCreator {
	return func(locator adventuria.ServiceLocator) adventuria.Cell {
		return &CellJail{
			CellBase: adventuria.CellBase{},
			locator:  locator,
		}
	}
}

func (c *CellJail) NextStep(user adventuria.User) string {
	nextStepType := ""

	if user.IsInJail() {
		switch user.LastAction().Type() {
		case adventuria.ActionTypeRollDice,
			adventuria.ActionTypeReroll:
			nextStepType = adventuria.ActionTypeRollWheel
		case adventuria.ActionTypeRollWheel:
			nextStepType = adventuria.ActionTypeChooseResult
		case adventuria.ActionTypeChooseResult:
			nextStepType = adventuria.ActionTypeRollDice
		default:
			nextStepType = adventuria.ActionTypeRollWheel
		}
	} else {
		nextStepType = adventuria.ActionTypeRollDice
	}

	return nextStepType
}

func (c *CellJail) Roll(_ adventuria.User) (*adventuria.WheelRollResult, error) {
	cellsCol, err := c.locator.Collections().Get(adventuria.TableCells)
	if err != nil {
		return nil, err
	}

	res := &adventuria.WheelRollResult{
		Collection: cellsCol,
	}

	cells := c.locator.Cells().GetAllByType(CellTypeGame)

	if len(cells) == 0 {
		return nil, errors.New("game cells not found")
	}

	for _, cell := range cells {
		res.FillerItems = append(res.FillerItems, &adventuria.WheelItem{
			Name: cell.Name(),
			Icon: cell.Icon(),
		})
	}

	res.WinnerId = helper.RandomItemFromSlice(cells).ID()

	// TODO
	panic("implement me")

	return res, nil
}

func (c *CellJail) OnCellReached(user adventuria.User) error {
	if user.LastAction().Type() == adventuria.ActionTypeDrop &&
		user.DropsInARow() >= c.locator.Settings().DropsToJail() {
		user.SetIsInJail(true)
	}
	return nil
}
