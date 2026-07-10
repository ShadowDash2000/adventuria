package balance_change

import "strconv"

func (b *BalanceChange) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
