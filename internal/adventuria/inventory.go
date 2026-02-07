package adventuria

type Inventory interface {
	Closable

	Refetch(ctx AppContext) error
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
	MustDropItem(ctx AppContext, invItemId string) error
	DropRandomItem(AppContext) error
	DropInventory(AppContext) error
	GetItemById(invItemId string) (Item, bool)
}
