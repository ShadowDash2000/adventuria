package teleport_to_closest_cell_by_type

import (
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*TeleportToClosestCellByType)(nil)

func (t *TeleportToClosestCellByType) Verify(_ context.Context, value string) error {
	cellType := t.decodeValue(value)
	_, ok := cells.Get(cellType)
	if !ok {
		return errors.New("unknown cell type")
	}
	return nil
}
