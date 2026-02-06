package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
)

type CellTeleport struct {
	adventuria.CellRecord
}

type cellTeleportValue struct {
	CellName string `json:"cell_name"`
}

func (c *CellTeleport) OnCellReached(ctx *adventuria.CellReachedContext) error {
	var decodedValue cellTeleportValue
	if err := c.UnmarshalJSONField(schema.CellSchema.Value, &decodedValue); err != nil {
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

	var exists bool
	err := ctx.App.
		RecordQuery(schema.CollectionCells).
		Select("count(*)").
		Where(dbx.HashExp{"name": decodedValue.CellName}).
		Limit(1).
		Row(&exists)
	if err != nil {
		return fmt.Errorf("teleport.verify(): can't find cell: %w", err)
	}

	return nil
}
