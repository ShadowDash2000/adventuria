package adventuria

import (
	"slices"

	"github.com/pocketbase/pocketbase/core"
)

type ItemBase struct {
	ItemRecordBase
	locator       PocketBaseLocator
	invItemRecord core.BaseRecordProxy
	effects       []Effect
}

func NewItemFromInventoryRecord(locator PocketBaseLocator, user User, invItemRecord *core.Record) (Item, error) {
	itemRecord, err := GetRecordById(locator, TableItems, invItemRecord.GetString("itemId"), []string{"effects"})
	if err != nil {
		return nil, err
	}

	var effects []Effect
	for _, effectRecord := range itemRecord.ExpandedAll("effects") {
		effect, err := NewEffectFromRecord(user, effectRecord)
		if err != nil {
			return nil, err
		}

		effects = append(effects, effect)
	}

	item := &ItemBase{
		locator: locator,
		effects: effects,
	}

	item.invItemRecord.SetProxyRecord(invItemRecord)
	item.itemRecord.SetProxyRecord(itemRecord)
	item.bindHooks()
	item.Awake()

	return item, nil
}

func (i *ItemBase) Awake() {
	for _, effect := range i.effects {
		if slices.Contains(i.AppliedEffects(), effect.ID()) {
			continue
		}

		effect.Subscribe(func() {
			i.addAppliedEffect(effect)

			if i.AppliedEffectsCount() == i.EffectsCount() {
				i.locator.PocketBase().Delete(i.invItemRecord)
			}
		})
	}
}

func (i *ItemBase) Sleep() {
	for _, effect := range i.effects {
		effect.Unsubscribe()
	}
}

func (i *ItemBase) bindHooks() {
	i.locator.PocketBase().OnRecordAfterUpdateSuccess(TableItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.itemRecord.Id {
			i.itemRecord.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	i.locator.PocketBase().OnRecordAfterUpdateSuccess(TableInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.invItemRecord.Id {
			i.invItemRecord.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
}

func (i *ItemBase) IDInventory() string {
	return i.invItemRecord.Id
}

func (i *ItemBase) IsActive() bool {
	return i.invItemRecord.GetBool("isActive")
}

func (i *ItemBase) SetIsActive(b bool) {
	i.invItemRecord.Set("isActive", b)
}

func (i *ItemBase) EffectsCount() int {
	return len(i.effects)
}

func (i *ItemBase) AppliedEffectsCount() int {
	return len(i.invItemRecord.GetStringSlice("appliedEffects"))
}

func (i *ItemBase) AppliedEffects() []string {
	return i.invItemRecord.GetStringSlice("appliedEffects")
}

func (i *ItemBase) addAppliedEffect(effect Effect) {
	i.invItemRecord.Set(
		"appliedEffects",
		append(i.invItemRecord.GetStringSlice("appliedEffects"), effect.ID()),
	)
}

func (i *ItemBase) Use() error {
	if i.IsActive() {
		return nil
	}

	i.SetIsActive(true)
	return i.locator.PocketBase().Save(i.invItemRecord)
}

func (i *ItemBase) Drop() error {
	if !i.CanDrop() {
		return nil
	}

	return i.locator.PocketBase().Delete(i.invItemRecord)
}
