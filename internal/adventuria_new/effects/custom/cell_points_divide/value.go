package cell_points_divide

import (
	"strconv"
)

func (c *CellPointsDivide) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
