package cells

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type CellTeleport struct {
	adventuria.CellBase
}

func NewCellTeleport() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellTeleport{
			adventuria.CellBase{},
		}
	}
}

func (c *CellTeleport) OnCellReached(user adventuria.User) error {
	if err := adventuria.PocketBase.Save(user.LastAction().ProxyRecord()); err != nil {
		return err
	}
	if err := user.MoveToCellName(c.GetString("value")); err != nil {
		return err
	}
	user.LastAction().SetType("teleport")
	return nil
}

func (c *CellTeleport) Verify(value string) error {
	if _, err := adventuria.PocketBase.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionCells),
		fmt.Sprintf("name = '%s'", value),
	); err != nil {
		return fmt.Errorf("teleport.verify(): can't find cell: %w", err)
	}

	return nil
}

func (c *CellTeleport) DecodeValue(_ string) (any, error) {
	return nil, nil
}
