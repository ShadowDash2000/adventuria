package adventuria

import "github.com/pocketbase/pocketbase/core"

type InventoryItem interface {
	core.RecordProxy
	ID() string
	UserId() string
	IsActive() bool
	SetIsActive(bool)
	IsUsingSlot() bool
	AppliedEffects() []string
	SetAppliedEffects([]string)
	CanDrop() bool
	Name() string
	Order() int
	Effects() []Effect
	EffectsCount() int
	EffectsByEvent(EffectUse) []Effect
	EffectsByTypes([]string) []Effect
	ApplyEffects([]string) error
	AppliedEffectsMap() map[string]struct{}
	AppendAppliedEffects([]string)
	Use() error
	Drop() error
}
