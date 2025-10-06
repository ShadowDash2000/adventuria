package adventuria

type Inventory interface {
	MaxSlots() int
	SetMaxSlots(int)
	AvailableSlots() int
	HasEmptySlots() bool
	AddItem(ItemRecord) error
	AddItemById(string) error
	MustAddItemById(string) error
	UseItem(string) error
	DropItem(string) error
	DropRandomItem() error
	DropInventory() error
}
