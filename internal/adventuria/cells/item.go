package cells

import "adventuria/internal/adventuria"

type CellItem struct {
	adventuria.CellBase
}

func NewCellItem() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellItem{
			adventuria.CellBase{},
		}
	}
}

func (c *CellItem) OnCellReached(user adventuria.User) error {
	user.SetItemWheelsCount(user.ItemWheelsCount() + 1)
	user.LastAction().SetCanMove(true)
	return nil
}
