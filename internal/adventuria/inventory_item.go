package adventuria

import (
	"errors"
	"github.com/pocketbase/pocketbase/core"
)

type InventoryItem struct {
	core.BaseRecordProxy
	gc   *GameComponents
	item *Item
}

func NewInventoryItem(record *core.Record, gc *GameComponents) (*InventoryItem, error) {
	var err error
	ii := &InventoryItem{
		gc: gc,
	}

	ii.SetProxyRecord(record)

	errs := gc.app.ExpandRecord(ii.Record, []string{"item"}, nil)
	if errs != nil {
		for _, err = range errs {
			return nil, err
		}
	}

	ii.item, err = NewItem(ii.ExpandedOne("item"), gc)
	if err != nil {
		return nil, err
	}

	return ii, nil
}

func (ii *InventoryItem) GetEffectsByEvent(event string) []Effect {
	if !ii.IsActive() {
		return nil
	}

	appliedEffects := ii.AppliedEffectsMap()
	var effects []Effect
	for _, effect := range ii.item.GetEffectsByEvent(event) {
		if _, ok := appliedEffects[effect.GetId()]; !ok {
			effects = append(effects, effect)
		}
	}

	return effects
}

func (ii *InventoryItem) IsActive() bool {
	return ii.GetBool("isActive")
}

func (ii *InventoryItem) SetIsActive(isActive bool) {
	ii.Set("isActive", isActive)
}

func (ii *InventoryItem) IsUsingSlot() bool {
	return ii.item.IsUsingSlot()
}

func (ii *InventoryItem) AppliedEffects() []string {
	return ii.GetStringSlice("appliedEffects")
}

func (ii *InventoryItem) AppliedEffectsMap() map[string]struct{} {
	var appliedEffects = make(map[string]struct{})
	for _, effect := range ii.GetStringSlice("appliedEffects") {
		appliedEffects[effect] = struct{}{}
	}
	return appliedEffects
}

func (ii *InventoryItem) SetAppliedEffects(appliedEffects []string) {
	ii.Set("appliedEffects", appliedEffects)
}

func (ii *InventoryItem) EffectsCount() int {
	return ii.item.EffectsCount()
}

func (ii *InventoryItem) Use() error {
	if ii.IsActive() {
		return errors.New("item is already active")
	}

	ii.SetIsActive(true)
	err := ii.gc.app.Save(ii)
	if err != nil {
		return err
	}

	return nil
}

func (ii *InventoryItem) CanDrop() bool {
	return ii.item.CanDrop()
}

func (ii *InventoryItem) GetName() string {
	return ii.item.Name()
}
