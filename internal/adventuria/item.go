package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type Effects struct {
	PointsIncrement  int    `json:"pointsIncrement" mapstructure:"pointsIncrement"`
	JailEscape       bool   `json:"jailEscape" mapstructure:"jailEscape"`
	DiceMultiplier   int    `json:"diceMultiplier" mapstructure:"diceMultiplier"`
	DiceIncrement    int    `json:"diceIncrement" mapstructure:"diceIncrement"`
	Dices            []Dice `json:"dices" mapstructure:"changeDices"`
	IsSafeDrop       bool   `json:"isSafeDrop" mapstructure:"isSafeDrop"`
	TimerIncrement   int    `json:"timerIncrement" mapstructure:"timerIncrement"`
	RollReverse      bool   `json:"rollReverse" mapstructure:"rollReverse"`
	DropInventory    bool   `json:"dropInventory" mapstructure:"dropInventory"`
	CellPointsDivide int    `json:"cellPointsDivide" mapstructure:"cellPointsDivide"`
}

type Item struct {
	gc      *GameComponents
	item    *core.Record
	effects []IEffect
}

func NewItem(record *core.Record, gc *GameComponents) (*Item, error) {
	errs := gc.app.ExpandRecord(record, []string{"effects"}, nil)
	if errs != nil {
		for _, err := range errs {
			return nil, err
		}
	}

	var effects []IEffect
	for _, effectRecord := range record.ExpandedAll("effects") {
		effect, err := NewEffect(effectRecord)
		if err != nil {
			return nil, err
		}

		effects = append(effects, effect)
	}

	item := &Item{
		gc:      gc,
		item:    record,
		effects: effects,
	}

	item.bindHooks()

	return item, nil
}

func (i *Item) bindHooks() {
	i.gc.app.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.item.Id {
			i.item = e.Record
		}
		return e.Next()
	})
}

func (i *Item) IsUsingSlot() bool {
	return i.item.GetBool("isUsingSlot")
}

func (i *Item) IsActiveByDefault() bool {
	return i.item.GetBool("isActiveByDefault")
}

func (i *Item) CanDrop() bool {
	return i.item.GetBool("canDrop")
}

func (i *Item) Order() int {
	return i.item.GetInt("order")
}

func (i *Item) EffectsCount() int {
	return len(i.effects)
}

func (i *Item) GetEffectsByEvent(event string) []IEffect {
	var effects []IEffect
	for _, e := range i.effects {
		if e.Event() == event {
			effects = append(effects, e)
		}
	}

	return effects
}

func (i *Item) Name() string {
	return i.item.GetString("name")
}
