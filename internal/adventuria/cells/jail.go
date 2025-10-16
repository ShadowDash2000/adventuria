package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"errors"
)

type CellJail struct {
	adventuria.CellBase
}

func NewCellJail() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellJail{
			CellBase: adventuria.CellBase{},
		}
	}
}

func (c *CellJail) Roll(_ adventuria.User) (*adventuria.WheelRollResult, error) {
	res := &adventuria.WheelRollResult{
		Collection: adventuria.GameCollections.Get(adventuria.CollectionCells),
	}

	cells := adventuria.GameCells.GetAllByType(CellTypeGame)

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

func (c *CellJail) OnCellReached(_ adventuria.User) error {
	return nil
}
