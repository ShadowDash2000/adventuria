package actions

import (
	"adventuria/internal/adventuria"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type BuyAction struct {
	adventuria.ActionBase
}

type cellShopValue struct {
	PriceMultiplier float32 `json:"price_multiplier"`
}

func (a *BuyAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if !currentCell.InCategory("shop") {
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

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, fmt.Errorf("buy.do(): current cell not found")
	}

	var decodedValue cellShopValue
	if err = json.Unmarshal([]byte(currentCell.Value()), &decodedValue); err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't decode cell value",
		}, fmt.Errorf("buy.do(): can't decode cell value: %w", err)
	}

	item := adventuria.NewItemFromRecord(itemRecord)
	price := a.calculatePrice(item.Price(), decodedValue)

	if ctx.User.Balance() < price {
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
	ctx.User.SetBalance(ctx.User.Balance() - price)

	return &adventuria.ActionResult{
		Success: true,
		Data:    invItemId,
	}, nil
}

func (a *BuyAction) GetVariants(ctx adventuria.ActionContext) any {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return nil
	}

	var decodedValue cellShopValue
	if err := json.Unmarshal([]byte(currentCell.Value()), &decodedValue); err != nil {
		return nil
	}

	ids, err := ctx.User.LastAction().ItemsList()
	if err != nil {
		return nil
	}

	exp := make([]dbx.Expression, len(ids))
	for i, id := range ids {
		exp[i] = dbx.HashExp{"id": id}
	}

	var records []*core.Record
	err = adventuria.PocketBase.
		RecordQuery(adventuria.CollectionItems).
		Where(dbx.Or(exp...)).
		All(&records)
	if err != nil {
		return nil
	}

	recordsMaps := make(map[string]*core.Record, len(records))
	for _, record := range records {
		record.Set("price", a.calculatePrice(record.GetInt("price"), decodedValue))
		recordsMaps[record.Id] = record
	}

	result := make([]*core.Record, len(ids))
	for i, id := range ids {
		if record, ok := recordsMaps[id]; ok {
			result[i] = record
		}
	}

	return struct {
		Items []*core.Record `json:"items"`
	}{
		Items: result,
	}
}

func (a *BuyAction) calculatePrice(price int, cellShopValue cellShopValue) int {
	if cellShopValue.PriceMultiplier != 0 {
		price = int(float32(price) * cellShopValue.PriceMultiplier)
	}
	return price
}
