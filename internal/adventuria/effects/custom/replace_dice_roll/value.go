package replace_dice_roll

import "strconv"

func (r *ReplaceDiceRoll) decodeValue(value string) (int, error) {
	return strconv.Atoi(value)
}
