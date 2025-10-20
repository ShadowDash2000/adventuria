package cells

import (
	"adventuria/internal/adventuria"
)

type CellShop struct {
	adventuria.CellBase
}

func NewCellShop() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellShop{
			adventuria.CellBase{},
		}
	}
}

func (c *CellShop) Roll(user adventuria.User) (*adventuria.WheelRollResult, error) {
	return nil, nil
}

func (c *CellShop) OnCellReached(user adventuria.User) error {
	user.LastAction().SetCanMove(true)

	return nil
}
