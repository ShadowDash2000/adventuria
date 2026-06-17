package change_dices

import (
	"adventuria/internal/adventuria_new/model"
	"errors"
	"strings"
)

func (c *ChangeDices) decodeValue(value string) ([]model.Dice, error) {
	diceTypes := strings.Split(value, ";")
	var dices []model.Dice
	for _, diceType := range diceTypes {
		dice, ok := model.GetDice(model.DiceType(diceType))
		if !ok {
			return nil, errors.New("unknown dice type")
		}
		dices = append(dices, dice)
	}
	return dices, nil
}
