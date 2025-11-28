package cells

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type CellSnake struct {
	adventuria.CellBase
}

func NewCellSnake() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellSnake{
			adventuria.CellBase{},
		}
	}
}

func (c *CellSnake) OnCellReached(user adventuria.User) error {
	if err := adventuria.PocketBase.Save(user.LastAction().ProxyRecord()); err != nil {
		return err
	}
	return user.MoveToCellName(c.GetString("value"))
}

func (c *CellSnake) Verify(value string) error {
	if _, err := adventuria.PocketBase.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionCells),
		fmt.Sprintf("name = '%s'", value),
	); err != nil {
		return fmt.Errorf("snake.verify(): can't find cell: %w", err)
	}

	return nil
}

func (c *CellSnake) DecodeValue(_ string) (any, error) {
	return nil, nil
}
