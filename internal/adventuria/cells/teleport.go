package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
)

var _ adventuria.CellVerifiable = (*CellTeleport)(nil)

type CellTeleport struct {
	adventuria.CellRecord
}

type cellTeleportValue struct {
	CellId string `json:"cell_id"`
}

func (c *CellTeleport) OnCellReached(ctx *adventuria.CellReachedContext) error {
	var decodedValue cellTeleportValue
	if err := c.UnmarshalJSONField(schema.CellSchema.Value, &decodedValue); err != nil {
		return fmt.Errorf("teleport.verify: invalid JSON: %w", err)
	}

	onBeforeTeleportOnCell := &adventuria.OnBeforeTeleportOnCell{
		AppContext:   ctx.AppContext,
		CellId:       decodedValue.CellId,
		SkipTeleport: false,
	}

	_, err := ctx.Player.OnBeforeTeleportOnCell().Trigger(onBeforeTeleportOnCell)
	if err != nil {
		return fmt.Errorf("teleport.onBeforeTeleportOnCell: %w", err)
	}

	if onBeforeTeleportOnCell.SkipTeleport {
		ctx.Player.LastAction().SetCanMove(true)
		return nil
	}

	ctx.Player.LastAction().SetType("teleport")
	if err := ctx.App.Save(ctx.Player.LastAction().ProxyRecord()); err != nil {
		return err
	}

	res, err := ctx.Player.MoveToCellId(ctx.AppContext, decodedValue.CellId)
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
		Where(dbx.HashExp{schema.CellSchema.Id: decodedValue.CellId}).
		Limit(1).
		Row(&exists)
	if err != nil {
		return fmt.Errorf("teleport.verify(): can't find cell: %w", err)
	}

	return nil
}
