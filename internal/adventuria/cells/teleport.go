package cells

import (
	"adventuria/internal/adventuria"
	"encoding/json"
	"fmt"
)

type CellTeleport struct {
	adventuria.CellRecord
}

type cellTeleportValue struct {
	CellName string `json:"cell_name"`
}

func (c *CellTeleport) OnCellReached(ctx *adventuria.CellReachedContext) error {
	var decodedValue cellTeleportValue
	if err := c.UnmarshalJSONField("value", &decodedValue); err != nil {
		return fmt.Errorf("teleport.verify: invalid JSON: %w", err)
	}

	if err := adventuria.PocketBase.Save(ctx.User.LastAction().ProxyRecord()); err != nil {
		return err
	}
	res, err := ctx.User.MoveToCellName(decodedValue.CellName)
	if err != nil {
		return err
	}

	ctx.Moves = append(ctx.Moves, res...)
	ctx.User.LastAction().SetType("teleport")

	return nil
}

func (c *CellTeleport) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellTeleport) Verify(value string) error {
	var decodedValue cellTeleportValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return fmt.Errorf("teleport.verify: invalid JSON: %w", err)
	}

	if _, err := adventuria.PocketBase.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionCells),
		fmt.Sprintf("name = '%s'", decodedValue.CellName),
	); err != nil {
		return fmt.Errorf("teleport.verify(): can't find cell: %w", err)
	}

	return nil
}
