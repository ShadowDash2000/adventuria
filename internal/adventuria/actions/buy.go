package actions

import (
	"adventuria/internal/adventuria"
	"fmt"
	"slices"
)

type BuyAction struct {
	adventuria.ActionBase
}

func (a *BuyAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if currentCell.Type() != "shop" {
		return false
	}

	return true
}

func (a *BuyAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	if _, ok := req["item_id"]; !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: item_id not specified",
		}, nil
	}

	itemId, ok := req["item_id"].(string)
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: item_id is not string",
		}, nil
	}

	ids, err := ctx.User.LastAction().ItemsList()
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
	if ctx.User.Balance() < item.Price() {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "not enough money",
		}, nil
	}

	invItemId, err := ctx.User.Inventory().AddItemById(itemId)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't add item to inventory",
		}, fmt.Errorf("buy.do(): can't add item to inventory: %w", err)
	}

	if index := slices.Index(ids, itemId); index != -1 {
		ids = slices.Delete(ids, index, index+1)
	}

	ctx.User.LastAction().SetItemsList(ids)
	ctx.User.SetBalance(ctx.User.Balance() - item.Price())

	return &adventuria.ActionResult{
		Success: true,
		Data:    invItemId,
	}, nil
}
