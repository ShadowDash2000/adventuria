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

func (a *BuyAction) Do(_ adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	// TODO: get id from ActionRequest
	const requestedItemId = "1"

	ids, err := a.User().LastAction().ItemsList()
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't get items list",
		}, fmt.Errorf("buy.do(): can't get items list: %w", err)
	}

	if !slices.Contains(ids, requestedItemId) {
		return &adventuria.ActionResult{
			Success: false,
			Error:   fmt.Sprintf("item with id = %s not found", requestedItemId),
		}, fmt.Errorf("buy.do(): item with id = %s not found", requestedItemId)
	}

	itemRecord, err := adventuria.PocketBase.FindRecordById(
		adventuria.GameCollections.Get(adventuria.CollectionItems),
		requestedItemId,
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

	invItemId, err := a.User().Inventory().AddItemById(requestedItemId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("buy.do(): can't add item to inventory: %w", err)
	}

	ids = slices.DeleteFunc(ids, func(s string) bool {
		return s == requestedItemId
	})

	a.User().LastAction().SetItemsList(ids)
	a.User().SetBalance(a.User().Balance() - item.Price())

	return &adventuria.ActionResult{
		Success: true,
		Data:    invItemId,
	}, nil
}
