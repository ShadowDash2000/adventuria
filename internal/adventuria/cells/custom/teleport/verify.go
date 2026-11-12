package teleport

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*CellTeleport)(nil)

func (c *CellTeleport) Verify(ctx context.Context, value string) error {
	payload, err := NewPayloadFromJSON([]byte(value))
	if err != nil {
		return err
	}

	_, err = c.cells.GetByID(ctx, payload.CellId)
	if err != nil {
		return err
	}

	return nil
}
