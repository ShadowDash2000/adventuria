package paid_movement_in_radius

import (
	"fmt"
	"strconv"
	"strings"
)

type effectValue struct {
	Radius int
	Price  int
}

func (p *PaidMovementInRadius) decodeValue(value string) (*effectValue, error) {
	values := strings.Split(value, ";")
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid values: %s", value)
	}

	var (
		radius, price int
		err           error
	)
	radius, err = strconv.Atoi(values[0])
	if err != nil {
		return nil, err
	}
	price, err = strconv.Atoi(values[1])
	if err != nil {
		return nil, err
	}

	return &effectValue{
		Radius: radius,
		Price:  price,
	}, nil
}
