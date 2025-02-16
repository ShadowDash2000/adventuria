package usecases

type Usable interface {
	Use() error
	SetUserId(userId string)
}

const (
	ItemTypeNone           = "none"
	ItemTypeDiceMultiplier = "diceMultiplier"

	// ItemUseTypeInstant
	// В теории тип instant должен срабатывать при любом action.
	// После применения предмета нужно сохранять эффект от его использования
	// в таблицу actionsEffects.
	ItemUseTypeInstant  = "instant"
	ItemUseTypeOnDrop   = "useOnDrop"
	ItemUseTypeOnReroll = "useOnReroll"
	ItemUseTypeOnRoll   = "useOnRoll"
)

var items = map[string]Usable{
	ItemTypeDiceMultiplier: &ItemDiceMultiplier{},
}

type Item struct {
	userId string
}

func (i *Item) SetUserId(userId string) {
	i.userId = userId
}

func NewItem(itemType string, userId string) Usable {
	item, ok := items[itemType]
	if !ok {
		return nil
	}

	item.SetUserId(userId)

	return item
}

type ItemDiceMultiplier struct {
	Item
}

func (idm *ItemDiceMultiplier) Use() error {
	return nil
}
