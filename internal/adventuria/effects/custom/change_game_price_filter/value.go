package change_game_price_filter

import (
	"fmt"
	"strconv"
	"strings"
)

type effectValue struct {
	Price     int
	PriceType string
	UseType   string
}

const (
	priceTypeMin = "min"
	priceTypeMax = "max"

	useTypeUsable   = "usable"
	useTypeUnusable = "unusable"
)

var priceTypes = map[string]struct{}{
	priceTypeMin: {},
	priceTypeMax: {},
}

var useTypes = map[string]struct{}{
	useTypeUsable:   {},
	useTypeUnusable: {},
}

func (c *ChangeGamePriceFilter) decodeValue(value string) (*effectValue, error) {
	vals := strings.Split(value, ";")
	if len(vals) != 3 {
		return nil, fmt.Errorf("invalid value: %s", value)
	}

	var (
		res effectValue
		err error
	)
	res.Price, err = strconv.Atoi(vals[0])
	if err != nil {
		return nil, err
	}

	if _, ok := priceTypes[vals[1]]; !ok {
		return nil, fmt.Errorf("invalid price type: %s", vals[1])
	}
	res.PriceType = vals[1]

	if _, ok := useTypes[vals[2]]; !ok {
		return nil, fmt.Errorf("invalid use type: %s", vals[2])
	}
	res.UseType = vals[2]

	return &res, nil
}
