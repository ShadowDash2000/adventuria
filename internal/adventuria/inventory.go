package adventuria

type Inventory interface {
	Closable

	MaxSlots() int
	SetMaxSlots(int)
	AvailableSlots() int
	HasEmptySlots() bool
	HasItem(invItemId string) bool
	RegisterItem(item Item)
	AddItem(AppContext, ItemRecord) (string, error)
	AddItemById(AppContext, string) (string, error)
	MustAddItemById(AppContext, string) (string, error)
	CanUseItem(AppContext, string) bool
	UseItem(AppContext, string) (OnUseSuccess, OnUseFail, error)
	DropItem(AppContext, string) error
	DropRandomItem(AppContext) error
	DropInventory(AppContext) error
	GetItemById(invItemId string) (Item, bool)
}
