package points_change

import "strconv"

func (p *PointsChange) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
