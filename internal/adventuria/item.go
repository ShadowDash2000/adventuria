package adventuria

type ItemRecord interface {
	ID() string
	Name() string
	Icon() string
	IsUsingSlot() bool
	IsActiveByDefault() bool
	CanDrop() bool
	IsRollable() bool
	Order() int
	Type() string
	Price() int
}

type Item interface {
	ItemRecord

	EffectsCount() int
	AppliedEffectsCount() int
	Use() error
	Drop() error
}
