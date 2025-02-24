package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
)

type InventoryItem struct {
	app     core.App
	invItem *core.Record
	item    *Item
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

	ii.item, err = NewItem(ii.invItem.ExpandedOne("item"), app)
	if err != nil {
		return nil, err
	}

	return ii, nil
}

func (ii *InventoryItem) GetEffects(event string) []*Effect {
	if !ii.invItem.GetBool("isActive") {
		return nil
	}

	if ii.item.GetEvent() != event {
		return nil
	}

	return ii.item.GetEffects()
}

func (ii *InventoryItem) IsUsingSlot() bool {
	return ii.item.IsUsingSlot()
}
