package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type Item struct {
	core.BaseRecordProxy
	gc      *GameComponents
	effects []Effect
}

func NewBaseItem(record *core.Record) *Item {
	item := &Item{}
	item.SetProxyRecord(record)
	return item
}

func NewItem(record *core.Record, gc *GameComponents) (*Item, error) {
	errs := gc.App.ExpandRecord(record, []string{"effects"}, nil)
	if errs != nil {
		for _, err := range errs {
			return nil, err
		}
	}

	var effects []Effect
	for _, effectRecord := range record.ExpandedAll("effects") {
		effect, err := NewEffectRecord(effectRecord)
		if err != nil {
			return nil, err
		}

		effects = append(effects, effect)
	}

	item := &Item{
		gc:      gc,
		effects: effects,
	}

	item.SetProxyRecord(record)
	item.bindHooks()

	return item, nil
}

func (i *Item) bindHooks() {
	i.gc.App.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.Id {
			i.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
}

func (i *Item) Name() string {
	return i.GetString("name")
}

func (i *Item) IsUsingSlot() bool {
	return i.GetBool("isUsingSlot")
}

func (i *Item) IsActiveByDefault() bool {
	return i.GetBool("isActiveByDefault")
}

func (i *Item) CanDrop() bool {
	return i.GetBool("canDrop")
}

func (i *Item) Order() int {
	return i.GetInt("order")
}

func (i *Item) EffectsCount() int {
	return len(i.effects)
}

func (i *Item) GetEffectsByEvent(event string) []Effect {
	var effects []Effect
	for _, e := range i.effects {
		if e.Event() == event {
			effects = append(effects, e)
		}
	}

	return effects
}
