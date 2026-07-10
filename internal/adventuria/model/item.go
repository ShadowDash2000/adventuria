package model

type ItemData struct {
	Id                string
	Name              string
	Icon              string
	Effects           []string
	IsUsingSlot       bool
	IsActiveByDefault bool
	CanDrop           bool
	IsRollable        bool
	Description       string
	Type              ItemType
	Price             int
}

type Item struct {
	data ItemData
}

func RestoreItem(data ItemData) *Item {
	return &Item{data: data}
}

func (i *Item) ID() string {
	return i.data.Id
}

func (i *Item) Name() string {
	return i.data.Name
}

func (i *Item) Icon() string {
	return i.data.Icon
}

func (i *Item) Effects() []string {
	return i.data.Effects
}

func (i *Item) IsUsingSlot() bool {
	return i.data.IsUsingSlot
}

func (i *Item) IsActiveByDefault() bool {
	return i.data.IsActiveByDefault
}

func (i *Item) CanDrop() bool {
	return i.data.CanDrop
}

func (i *Item) IsRollable() bool {
	return i.data.IsRollable
}

func (i *Item) Description() string {
	return i.data.Description
}

func (i *Item) Type() ItemType {
	return i.data.Type
}

func (i *Item) Price() int {
	return i.data.Price
}

func (i *Item) EffectsCount() int {
	return len(i.Effects())
}
