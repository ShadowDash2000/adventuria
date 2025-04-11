package adventuria

import (
	"errors"
	"github.com/pocketbase/pocketbase/core"
	"slices"
)

type InventoryItemBase struct {
	core.BaseRecordProxy
	gc   *GameComponents
	item Item
}

func NewInventoryItemFromRecord(record *core.Record, gc *GameComponents) (InventoryItem, error) {
	ii := &InventoryItemBase{
		gc: gc,
	}

	ii.SetProxyRecord(record)

	item, ok := gc.Items.GetById(ii.ItemId())
	if !ok {
		return nil, errors.New("item not found")
	}

	ii.item = item

	return ii, nil
}

func (ii *InventoryItemBase) ID() string {
	return ii.Id
}

func (ii *InventoryItemBase) UserId() string {
	return ii.GetString("user")
}

func (ii *InventoryItemBase) ItemId() string {
	return ii.GetString("item")
}

func (ii *InventoryItemBase) IsActive() bool {
	return ii.GetBool("isActive")
}

func (ii *InventoryItemBase) SetIsActive(isActive bool) {
	ii.Set("isActive", isActive)
}

func (ii *InventoryItemBase) IsUsingSlot() bool {
	return ii.item.IsUsingSlot()
}

func (ii *InventoryItemBase) AppliedEffects() []string {
	return ii.GetStringSlice("appliedEffects")
}

func (ii *InventoryItemBase) SetAppliedEffects(appliedEffects []string) {
	ii.Set("appliedEffects", appliedEffects)
}

func (ii *InventoryItemBase) CanDrop() bool {
	return ii.item.CanDrop()
}

func (ii *InventoryItemBase) Name() string {
	return ii.item.Name()
}

func (ii *InventoryItemBase) Order() int {
	return ii.item.Order()
}

func (ii *InventoryItemBase) EffectsByEvent(event EffectUse) []Effect {
	if !ii.IsActive() {
		return nil
	}

	appliedEffects := ii.AppliedEffectsMap()
	var effects []Effect
	for _, effect := range ii.item.EffectsByEvent(event) {
		if _, ok := appliedEffects[effect.ID()]; !ok {
			effects = append(effects, effect)
		}
	}

	return effects
}

func (ii *InventoryItemBase) EffectsByTypes(types []string) []Effect {
	if !ii.IsActive() {
		return nil
	}

	appliedEffects := ii.AppliedEffectsMap()
	var effects []Effect
	for _, effect := range ii.item.Effects() {
		if _, ok := appliedEffects[effect.ID()]; ok {
			continue
		}

		if !slices.Contains(types, effect.Type()) {
			continue
		}

		effects = append(effects, effect)
	}

	return effects
}

func (ii *InventoryItemBase) Effects() []Effect {
	if !ii.IsActive() {
		return nil
	}

	appliedEffects := ii.AppliedEffectsMap()
	var effects []Effect
	for _, effect := range ii.item.Effects() {
		if _, ok := appliedEffects[effect.ID()]; !ok {
			effects = append(effects, effect)
		}
	}

	return effects
}

func (ii *InventoryItemBase) ApplyEffects(effectsIds []string) error {
	ii.AppendAppliedEffects(effectsIds)

	if len(ii.AppliedEffects()) < ii.EffectsCount() {
		err := ii.gc.App.Save(ii)
		if err != nil {
			return err
		}
	} else {
		err := ii.gc.App.Delete(ii)
		if err != nil {
			return err
		}
	}

	return nil
}

func (ii *InventoryItemBase) AppliedEffectsMap() map[string]struct{} {
	var appliedEffects = make(map[string]struct{})
	for _, effect := range ii.AppliedEffects() {
		appliedEffects[effect] = struct{}{}
	}
	return appliedEffects
}

func (ii *InventoryItemBase) AppendAppliedEffects(effectsIds []string) {
	ii.SetAppliedEffects(append(ii.AppliedEffects(), effectsIds...))
}

func (ii *InventoryItemBase) EffectsCount() int {
	return ii.item.EffectsCount()
}

func (ii *InventoryItemBase) Use() error {
	if ii.IsActive() {
		return errors.New("item is active already")
	}

	ii.SetIsActive(true)
	err := ii.gc.App.Save(ii)
	if err != nil {
		return err
	}

	ii.gc.Log.Add(ii.UserId(), LogTypeItemUse, ii.Name())

	return nil
}

func (ii *InventoryItemBase) Drop() error {
	if !ii.CanDrop() {
		return errors.New("item isn't droppable")
	}

	err := ii.gc.App.Delete(ii)
	if err != nil {
		return err
	}

	ii.gc.Log.Add(ii.UserId(), LogTypeItemDrop, ii.Name())

	return nil
}
