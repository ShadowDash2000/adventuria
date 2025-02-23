package adventuria

import (
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"strings"
)

type Usable interface {
	GetEffects(string) any
	GetOrder() int
	SetItem(*core.Record)
}

const (
	ItemTypeDiceMultiplier = "diceMultiplier"
	ItemTypeDiceIncrement  = "diceIncrement"
	ItemTypeSafeDrop       = "safeDrop"
	ItemTypeChangeDices    = "changeDices"

	// ItemUseTypeInstant
	// В теории тип instant должен срабатывать при любом action.
	// После применения предмета нужно сохранять эффект от его использования
	// в таблицу actionsEffects.
	ItemUseTypeInstant  = "instant"
	ItemUseTypeOnDrop   = "useOnDrop"
	ItemUseTypeOnReroll = "useOnReroll"
	ItemUseTypeOnRoll   = "useOnRoll"
)

var items = map[string]Usable{
	ItemTypeDiceMultiplier: &ItemDiceMultiplier{},
	ItemTypeDiceIncrement:  &ItemDiceIncrement{},
	ItemTypeSafeDrop:       &ItemSafeDrop{},
}

type OnRollEffects struct {
	DiceMultiplier int    `json:"diceMultiplier"`
	DiceIncrement  int    `json:"diceIncrement"`
	Dices          []Dice `json:"dices"`
}

type OnDropEffects struct {
	IsSafeDrop bool `json:"isSafeDrop"`
}

type Item struct {
	item *core.Record
}

func NewItem(record *core.Record) (Usable, error) {
	item, ok := items[record.GetString("type")]
	if !ok {
		return nil, errors.New("item type not found")
	}

	item.SetItem(record)

	return item, nil
}

func (i *Item) parseString() []string {
	return strings.Split(i.item.GetString("value"), ", ")
}

func (i *Item) GetOrder() int {
	return i.item.GetInt("order")
}

func (i *Item) SetItem(item *core.Record) {
	i.item = item
}

type ItemDiceMultiplier struct {
	Item
}

func (i *ItemDiceMultiplier) GetEffects(event string) any {
	switch event {
	case ItemUseTypeOnRoll:
		return &OnRollEffects{
			DiceMultiplier: i.item.GetInt("value"),
		}
	}

	return nil
}

type ItemDiceIncrement struct {
	Item
}

func (i *ItemDiceIncrement) GetEffects(event string) any {
	switch event {
	case ItemUseTypeOnRoll:
		return &OnRollEffects{
			DiceIncrement: i.item.GetInt("value"),
		}
	}

	return nil
}

type ItemSafeDrop struct {
	Item
}

func (i *ItemSafeDrop) GetEffects(event string) any {
	switch event {
	case ItemUseTypeOnDrop:
		return OnDropEffects{IsSafeDrop: true}
	}

	return nil
}

type ItemChangeDices struct {
	Item
}

func (i *ItemChangeDices) GetEffects(event string) any {
	dicesTypes := i.parseString()

	dices := make([]Dice, len(dicesTypes))
	for _, diceType := range dicesTypes {
		if dice, ok := Dices[diceType]; ok {
			dices = append(dices, dice)
		}
	}

	switch event {
	case ItemUseTypeOnRoll:
		return &OnRollEffects{Dices: dices}
	}

	return nil
}
