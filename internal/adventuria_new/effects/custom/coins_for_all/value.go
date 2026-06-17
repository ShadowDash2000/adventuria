package coins_for_all

import (
	"fmt"
	"strconv"
	"strings"
)

type effectValue struct {
	CoinsForPlayer int
	CoinsForOther  int
}

func (c *CoinsForAll) decodeValue(value string) (*effectValue, error) {
	values := strings.Split(value, ";")
	if len(values) != 2 {
		return nil, fmt.Errorf("invalid values: %s", value)
	}

	var err error
	coins := make([]int, len(values))
	for i, value := range values {
		coins[i], err = strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid value: %s", value)
		}
	}

	return &effectValue{
		CoinsForPlayer: coins[0],
		CoinsForOther:  coins[1],
	}, nil
}
