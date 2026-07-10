package teleport_to_closest_cell_by_type

import (
	"adventuria/internal/adventuria/model"
)

func (t *TeleportToClosestCellByType) decodeValue(value string) model.CellType {
	return model.CellType(value)
}
