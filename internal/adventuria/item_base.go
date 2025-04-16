package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type ItemBase struct {
	core.BaseRecordProxy
	effects []Effect
}

func NewBaseItemFromRecord(record *core.Record) Item {
	item := &ItemBase{}
	item.SetProxyRecord(record)
	return item
}

func NewItemFromRecord(record *core.Record) (Item, error) {
	errs := GameApp.ExpandRecord(record, []string{"effects"}, nil)
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

	item := &ItemBase{
		effects: effects,
	}

	item.SetProxyRecord(record)
	item.bindHooks()

	return item, nil
}

func (i *ItemBase) bindHooks() {
	GameApp.OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.Id {
			i.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
}

func (i *ItemBase) ID() string {
	return i.Id
}

func (i *ItemBase) Name() string {
	return i.GetString("name")
}

func (i *ItemBase) Icon() string {
	return i.GetString("icon")
}

func (i *ItemBase) IsUsingSlot() bool {
	return i.GetBool("isUsingSlot")
}

func (i *ItemBase) IsActiveByDefault() bool {
	return i.GetBool("isActiveByDefault")
}

func (i *ItemBase) CanDrop() bool {
	return i.GetBool("canDrop")
}

func (i *ItemBase) IsRollable() bool {
	return i.GetBool("isRollable")
}

func (i *ItemBase) Order() int {
	return i.GetInt("order")
}

func (i *ItemBase) EffectsCount() int {
	return len(i.effects)
}

func (i *ItemBase) EffectsByEvent(event EffectUse) []Effect {
	var effects []Effect
	for _, e := range i.effects {
		if e.Event() == event {
			effects = append(effects, e)
		}
	}

	return effects
}

func (i *ItemBase) Effects() []Effect {
	return i.effects
}
