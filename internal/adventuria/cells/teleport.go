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

	ctx.User.LastAction().SetType("teleport")
	if err := ctx.App.Save(ctx.User.LastAction().ProxyRecord()); err != nil {
		return err
	}

	res, err := ctx.User.MoveToCellName(ctx.AppContext, decodedValue.CellName)
	if err != nil {
		return err
	}

	ctx.Moves = append(ctx.Moves, res...)

	return nil
}

func (c *CellTeleport) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellTeleport) Verify(ctx adventuria.AppContext, value string) error {
	var decodedValue cellTeleportValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return fmt.Errorf("teleport.verify: invalid JSON: %w", err)
	}

	if _, err := ctx.App.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionCells),
		fmt.Sprintf("name = '%s'", decodedValue.CellName),
	); err != nil {
		return fmt.Errorf("teleport.verify(): can't find cell: %w", err)
	}

	return nil
}
