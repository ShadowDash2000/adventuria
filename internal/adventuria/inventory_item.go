package adventuria

import (
	"errors"
	"github.com/pocketbase/pocketbase/core"
)

type InventoryItem struct {
	app     core.App
	log     *Log
	invItem *core.Record
	item    *Item
}

func NewInventoryItem(record *core.Record, log *Log, app core.App) (*InventoryItem, error) {
	var err error
	ii := &InventoryItem{
		app:     app,
		log:     log,
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

func (ii *InventoryItem) Use() error {
	isActive := ii.invItem.GetBool("isActive")
	if isActive {
		return errors.New("item is already active")
	}

	ii.invItem.Set("isActive", true)
	err := ii.app.Save(ii.invItem)
	if err != nil {
		return err
	}

	return nil
}

func (ii *InventoryItem) CanDrop() bool {
	return ii.item.CanDrop()
}

func (ii *InventoryItem) GetName() string {
	return ii.item.GetName()
}
