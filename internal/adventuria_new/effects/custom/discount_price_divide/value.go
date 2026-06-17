package discount_price_divide

import "strconv"

func (d *DiscountPriceDivide) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
