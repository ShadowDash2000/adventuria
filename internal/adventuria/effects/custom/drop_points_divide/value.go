package drop_points_divide

import (
	"strconv"
)

func (d *DropPointsDivide) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
