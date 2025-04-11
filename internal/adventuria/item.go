package adventuria

type Item interface {
	ID() string
	Name() string
	Icon() string
	IsUsingSlot() bool
	IsActiveByDefault() bool
	CanDrop() bool
	IsRollable() bool
	Order() int
	Effects() []Effect
	EffectsCount() int
	EffectsByEvent(EffectUse) []Effect
}
