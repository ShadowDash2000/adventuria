package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"fmt"
	"slices"
)

type BuyAction struct {
	adventuria.ActionBase
}

func (a *BuyAction) CanDo() bool {
	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return false
	}

	if _, ok = currentCell.(*cells.CellShop); !ok {
		return false
	}

	return true
}

func (a *BuyAction) Do(req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	if _, ok := req["item_id"]; !ok {
		return nil, fmt.Errorf("buy.do(): item_id not specified")
	}

	itemId, ok := req["item_id"].(string)
	if !ok {
		return nil, fmt.Errorf("buy.do(): item_id is not string")
	}

	ids, err := a.User().LastAction().ItemsList()
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't get items list",
		}, fmt.Errorf("buy.do(): can't get items list: %w", err)
	}

	if !slices.Contains(ids, itemId) {
		return &adventuria.ActionResult{
			Success: false,
			Error:   fmt.Sprintf("item with id = %s not found", itemId),
		}, fmt.Errorf("buy.do(): item with id = %s not found", itemId)
	}

	itemRecord, err := adventuria.PocketBase.FindRecordById(
		adventuria.GameCollections.Get(adventuria.CollectionItems),
		itemId,
	)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't get item record",
		}, fmt.Errorf("buy.do(): can't get item record: %w", err)
	}

	item := adventuria.NewItemFromRecord(itemRecord)
	if a.User().Balance() < item.Price() {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "not enough money",
		}, nil
	}

	invItemId, err := a.User().Inventory().AddItemById(itemId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("buy.do(): can't add item to inventory: %w", err)
	}

	ids = slices.DeleteFunc(ids, func(s string) bool {
		return s == itemId
	})

	a.User().LastAction().SetItemsList(ids)
	a.User().SetBalance(a.User().Balance() - item.Price())

	return &adventuria.ActionResult{
		Success: true,
		Data:    invItemId,
	}, nil
}
