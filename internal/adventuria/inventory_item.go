package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type InventoryItem struct {
	app     core.App
	invItem *core.Record
	item    Usable
}

func NewInventoryItem(record *core.Record, app core.App) (*InventoryItem, error) {
	var err error
	ii := &InventoryItem{
		app:     app,
		invItem: record,
	}

	errs := app.ExpandRecord(ii.invItem, []string{"item"}, nil)
	if errs != nil {
		for _, err = range errs {
			return nil, err
		}
	}

	ii.item, err = NewItem(ii.invItem.ExpandedOne("item"))
	if err != nil {
		return nil, err
	}

	return ii, nil
}

func (ii *InventoryItem) GetEffects(event string) any {
	if !ii.invItem.GetBool("isActive") {
		return nil
	}

	return ii.item.GetEffects(event)
}

func (ii *InventoryItem) IsUsingSlot() bool {
	return ii.item.IsUsingSlot()
}
