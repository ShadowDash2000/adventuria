package cells

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type CellTeleport struct {
	adventuria.CellRecord
}

func NewCellTeleport() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellTeleport{
			adventuria.CellRecord{},
		}
	}
}

func (c *CellTeleport) OnCellReached(ctx *adventuria.CellReachedContext) error {
	if err := adventuria.PocketBase.Save(ctx.User.LastAction().ProxyRecord()); err != nil {
		return err
	}
	res, err := ctx.User.MoveToCellName(c.Value())
	if err != nil {
		return err
	}

	ctx.Moves = append(ctx.Moves, res...)
	ctx.User.LastAction().SetType("teleport")

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
