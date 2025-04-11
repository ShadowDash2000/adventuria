package adventuria

type CellStart struct {
	CellBase
}

func NewCellStart() CellCreator {
	return func() Cell {
		return &CellStart{}
	}
}
