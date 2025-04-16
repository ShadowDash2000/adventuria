package adventuria

type CellItem struct {
	CellBase
}

func NewCellItem() CellCreator {
	return func() Cell {
		return &CellItem{}
	}
}

func (c *CellItem) OnCellReached(user *User) error {
	user.SetItemWheelsCount(user.ItemWheelsCount() + 1)
	return nil
}
