package usecases

type Usable interface {
	Use() error
}

const (
	ItemTypeNone           = "none"
	ItemTypeDiceMultiplier = "diceMultiplier"

	// ItemUseTypeInstant
	// В теории тип instant должен срабатывать при любом action.
	ItemUseTypeInstant  = "instant"
	ItemUseTypeOnDrop   = "useOnDrop"
	ItemUseTypeOnReroll = "useOnReroll"
	ItemUseTypeOnRoll   = "useOnRoll"
)

type Item struct {
}

func NewItem() *Item {
	return &Item{}
}
