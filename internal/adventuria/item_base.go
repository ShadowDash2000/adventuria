package adventuria

import (
	"adventuria/pkg/event"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

type ItemBase struct {
	ItemRecordBase
	isAwake           bool
	user              User
	invItemRecord     core.BaseRecordProxy
	effects           map[string]Effect
	effectsUnsubGroup map[string]event.UnsubGroup
	hookIds           []string
}

func NewItemFromInventoryRecord(user User, invItemRecord *core.Record) (Item, error) {
	itemRecord, err := GetRecordById(CollectionItems, invItemRecord.GetString("item"), []string{"effects"})
	if err != nil {
		return nil, err
	}

	effects := make(map[string]Effect)
	for _, effectRecord := range itemRecord.ExpandedAll("effects") {
		effect, err := NewEffectFromRecord(effectRecord)
		if err != nil {
			return nil, err
		}

		effects[effect.ID()] = effect
	}

	item := &ItemBase{
		user:    user,
		effects: effects,
	}

	item.invItemRecord.SetProxyRecord(invItemRecord)
	item.itemRecord.SetProxyRecord(itemRecord)
	item.bindHooks()

	if item.IsActive() {
		if err = item.awake(); err != nil {
			item.Close()
			return nil, err
		}
	}

	return item, nil
}

func (i *ItemBase) awake() error {
	appliedEffects := make(map[string]struct{}, i.AppliedEffectsCount())
	for _, appliedEffectId := range i.AppliedEffects() {
		appliedEffects[appliedEffectId] = struct{}{}
	}

	i.effectsUnsubGroup = make(map[string]event.UnsubGroup, len(i.effects)-len(appliedEffects))
	for _, effect := range i.effects {
		if _, ok := appliedEffects[effect.ID()]; ok {
			continue
		}

		unsubs, err := effect.Subscribe(
			EffectContext{
				User:      i.user,
				InvItemID: i.invItemRecord.Id,
			},
			func() {
				i.addAppliedEffect(effect)
				i.unsubEffectByID(effect.ID())

				if i.AppliedEffectsCount() == i.EffectsCount() {
					i.user.LastAction().UsedItemAppend(i.itemRecord.Id)
					i.sleep()
					PocketBase.Delete(i.invItemRecord.ProxyRecord())
				}
			})
		if len(unsubs) > 0 {
			i.effectsUnsubGroup[effect.ID()] = event.UnsubGroup{Fns: unsubs}
		}
		if err != nil {
			i.sleep()
			return err
		}
	}

	i.isAwake = true
	return nil
}

func (i *ItemBase) sleep() {
	for _, effect := range i.effects {
		i.unsubEffectByID(effect.ID())
	}

	i.isAwake = false
}

func (i *ItemBase) unsubEffectByID(id string) {
	if unsubGroup, ok := i.effectsUnsubGroup[id]; ok {
		unsubGroup.Unsubscribe()
		delete(i.effectsUnsubGroup, id)
	}
}

func (i *ItemBase) bindHooks() {
	i.hookIds = make([]string, 3)

	i.hookIds[0] = PocketBase.OnRecordAfterUpdateSuccess(CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.itemRecord.Id {
			i.itemRecord.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	i.hookIds[1] = PocketBase.OnRecordAfterUpdateSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.invItemRecord.Id {
			i.invItemRecord.SetProxyRecord(e.Record)

			if i.IsActive() && !i.isAwake {
				i.awake()
			} else if !i.IsActive() && i.isAwake {
				i.sleep()
			}
		}
		return e.Next()
	})
	i.hookIds[2] = PocketBase.OnRecordAfterDeleteSuccess(CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.invItemRecord.Id {
			i.sleep()
		}
		return e.Next()
	})
}

func (i *ItemBase) Close() {
	i.sleep()
	PocketBase.OnRecordAfterCreateSuccess(CollectionInventory).Unbind(i.hookIds[0])
	PocketBase.OnRecordAfterUpdateSuccess(CollectionInventory).Unbind(i.hookIds[1])
	PocketBase.OnRecordAfterDeleteSuccess(CollectionInventory).Unbind(i.hookIds[2])
}

func (i *ItemBase) IDInventory() string {
	return i.invItemRecord.Id
}

func (i *ItemBase) IsActive() bool {
	return i.invItemRecord.GetBool("isActive")
}

func (i *ItemBase) setIsActive(b bool) {
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

func (i *ItemBase) CanUse() bool {
	for _, effect := range i.effects {
		if !effect.CanUse(EffectContext{
			User:      i.user,
			InvItemID: i.invItemRecord.Id,
		}) {
			return false
		}
	}
	return true
}

func (i *ItemBase) Use() (OnUseSuccess, OnUseFail, error) {
	if i.IsActive() {
		return nil, nil, errors.New("item is already active")
	}

	if err := i.awake(); err != nil {
		return nil, nil, err
	}
	i.setIsActive(true)

	return func() error {
			// if an item is not awake, then it was removed from inventory
			if !i.isAwake {
				return nil
			}
			return PocketBase.Save(i.invItemRecord)
		}, func() {
			i.setIsActive(false)
			i.sleep()
		}, nil
}

func (i *ItemBase) Drop() error {
	if !i.CanDrop() {
		return nil
	}

	if i.IsActive() {
		i.sleep()
	}

	return PocketBase.Delete(i.invItemRecord.ProxyRecord())
}

func (i *ItemBase) GetEffectVariants(effectId string) (any, error) {
	effect, ok := i.effects[effectId]
	if !ok {
		return nil, errors.New("effect not found")
	}

	return effect.GetVariants(EffectContext{
		User:      i.user,
		InvItemID: i.invItemRecord.Id,
	}), nil
}
