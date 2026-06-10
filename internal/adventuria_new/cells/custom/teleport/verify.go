package teleport

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*CellTeleport)(nil)

func (c *CellTeleport) Verify(ctx context.Context, value string) error {
	decodedValue, err := c.decodeValue(value)
	if err != nil {
		return err
	}

	_, err = c.cells.GetByID(ctx, decodedValue.CellId)
	if err != nil {
		return err
	}

	return nil
}
