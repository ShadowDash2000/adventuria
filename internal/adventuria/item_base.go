package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/event"
	"errors"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
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

func NewItemFromInventoryRecord(ctx AppContext, user User, invItemRecord *core.Record) (Item, error) {
	itemRecord, err := GetRecordById(
		schema.CollectionItems,
		invItemRecord.GetString(schema.InventorySchema.Item),
		[]string{schema.ItemSchema.Effects},
	)
	if err != nil {
		return nil, err
	}

	effects := make(map[string]Effect)
	for _, effectRecord := range itemRecord.ExpandedAll(schema.ItemSchema.Effects) {
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
	item.bindHooks(ctx)

	if item.IsActive() {
		if err = item.awake(); err != nil {
			item.Close(ctx)
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
			func(ctx AppContext) {
				i.addAppliedEffect(effect)
				i.unsubEffectByID(effect.ID())

				if i.AppliedEffectsCount() == i.EffectsCount() {
					i.user.LastAction().UsedItemAppend(i.itemRecord.Id)
					i.sleep()
					if err := ctx.App.Delete(i.invItemRecord.ProxyRecord()); err != nil {
						ctx.App.Logger().Error(
							"Failed to delete item after all effects applied",
							"error", err,
						)
					}
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

func (i *ItemBase) bindHooks(ctx AppContext) {
	i.hookIds = make([]string, 3)

	i.hookIds[0] = ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionItems).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.itemRecord.Id {
			i.itemRecord.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	i.hookIds[1] = ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
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
	i.hookIds[2] = ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionInventory).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == i.invItemRecord.Id {
			i.sleep()
		}
		return e.Next()
	})
}

func (i *ItemBase) Close(ctx AppContext) {
	i.sleep()
	ctx.App.OnRecordAfterCreateSuccess(schema.CollectionInventory).Unbind(i.hookIds[0])
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionInventory).Unbind(i.hookIds[1])
	ctx.App.OnRecordAfterDeleteSuccess(schema.CollectionInventory).Unbind(i.hookIds[2])
}

func (i *ItemBase) IDInventory() string {
	return i.invItemRecord.Id
}

func (i *ItemBase) IsActive() bool {
	return i.invItemRecord.GetBool(schema.InventorySchema.IsActive)
}

func (i *ItemBase) setIsActive(b bool) {
	i.invItemRecord.Set(schema.InventorySchema.IsActive, b)
}

func (i *ItemBase) setActivated(date types.DateTime) {
	i.invItemRecord.Set(schema.InventorySchema.Activated, date)
}

func (i *ItemBase) EffectsCount() int {
	return len(i.effects)
}

func (i *ItemBase) AppliedEffectsCount() int {
	return len(i.AppliedEffects())
}

func (i *ItemBase) AppliedEffects() []string {
	return i.invItemRecord.GetStringSlice(schema.InventorySchema.AppliedEffects)
}

func (i *ItemBase) addAppliedEffect(effect Effect) {
	i.invItemRecord.Set(
		schema.InventorySchema.AppliedEffects,
		append(i.AppliedEffects(), effect.ID()),
	)
}

func (i *ItemBase) CanUse(ctx AppContext) bool {
	for _, effect := range i.effects {
		if !effect.CanUse(ctx, EffectContext{
			User:      i.user,
			InvItemID: i.invItemRecord.Id,
		}) {
			return false
		}
	}
	return true
}

func (i *ItemBase) Use(ctx AppContext) (OnUseSuccess, OnUseFail, error) {
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
			i.setActivated(types.NowDateTime())
			return ctx.App.Save(i.invItemRecord)
		}, func() {
			i.setIsActive(false)
			i.sleep()
		}, nil
}

func (i *ItemBase) Drop(ctx AppContext) error {
	if !i.CanDrop() {
		return nil
	}

	if i.IsActive() {
		i.sleep()
	}

	return ctx.App.Delete(i.invItemRecord.ProxyRecord())
}

func (i *ItemBase) GetEffectVariants(ctx AppContext, effectId string) (any, error) {
	effect, ok := i.effects[effectId]
	if !ok {
		return nil, errors.New("effect not found")
	}

	return effect.GetVariants(ctx, EffectContext{
		User:      i.user,
		InvItemID: i.invItemRecord.Id,
	}), nil
}
