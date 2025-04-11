package adventuria

type Inventory interface {
	MaxSlots() int
	SetMaxSlots(int)
	AvailableSlots() int
	HasEmptySlots() bool
	AddItem(Item) error
	AddItemById(string) error
	MustAddItemById(string) error
	Effects(EffectUse) (*Effects, map[string][]string, error)
	ApplyEffects(map[string][]string) error
	ApplyEffectsByEvent(EffectUse) (*Effects, error)
	ApplyEffectsByTypes([]string) error
	UseItem(string) error
	DropItem(string) error
	DropRandomItem() error
	DropInventory() error
	Items() map[string]InventoryItem
}
