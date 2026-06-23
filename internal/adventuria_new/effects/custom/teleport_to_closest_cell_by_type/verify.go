package teleport_to_closest_cell_by_type

import (
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"fmt"
)

var _ model.Verifiable = (*TeleportToClosestCellByType)(nil)

func (t *TeleportToClosestCellByType) Verify(_ context.Context, value string) error {
	cellType := t.decodeValue(value)
	_, ok := cells.Get(cellType)
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownCellType, cellType)
	}
	return nil
}
