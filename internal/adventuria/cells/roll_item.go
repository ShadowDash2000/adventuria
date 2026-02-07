package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/helper"
	"encoding/json"
	"fmt"

	"github.com/pocketbase/dbx"
)

var _ adventuria.CellWheel = (*CellRollItem)(nil)

type CellRollItem struct {
	adventuria.CellRecord
}

type cellRollItemValue struct {
	ItemsType adventuria.ItemType `json:"items_type"`
}

func (c *CellRollItem) Verify(_ adventuria.AppContext, value string) error {
	var decodedValue cellRollItemValue
	if err := json.Unmarshal([]byte(value), &decodedValue); err != nil {
		return fmt.Errorf("roll_item.refreshItems: invalid JSON: %w", err)
	}

	_, ok := adventuria.ItemTypes[decodedValue.ItemsType]
	if !ok {
		return fmt.Errorf("roll_item: unknown item type %s", decodedValue.ItemsType)
	}

	return nil
}

func (c *CellRollItem) Roll(ctx adventuria.AppContext, user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
	items, err := user.LastAction().ItemsList()
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't unmarshal items list",
		}, fmt.Errorf("roll_item.roll(): can't unmarshal items list: %w", err)
	}

	if len(items) == 0 {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: no items to roll",
		}, fmt.Errorf("roll_item.roll(): no items to roll")
	}

	records, err := ctx.App.FindRecordsByIds(
		adventuria.GameCollections.Get(schema.CollectionItems),
		items,
	)
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't fetch records",
		}, fmt.Errorf("roll_item.roll(): can't fetch records: %w", err)
	}

	var fillerItems []adventuria.WheelItem
	for _, record := range records {
		fillerItems = append(fillerItems, adventuria.WheelItem{
			Id:   record.Id,
			Name: record.GetString("name"),
			Icon: record.GetString("icon"),
		})
	}

	return &adventuria.WheelRollResult{
		FillerItems: fillerItems,
		WinnerId:    helper.RandomItemFromSlice(items),
		Success:     true,
	}, nil
}

func (c *CellRollItem) RefreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	return c.refreshItems(ctx, user)
}

func (c *CellRollItem) OnCellReached(ctx *adventuria.CellReachedContext) error {
	return c.refreshItems(ctx.AppContext, ctx.User)
}

func (c *CellRollItem) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellRollItem) refreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	var decodedValue cellRollItemValue
	if err := c.UnmarshalJSONField("value", &decodedValue); err != nil {
		return fmt.Errorf("roll_item.refreshItems: invalid JSON: %w", err)
	}

	var records []struct {
		Id string `db:"id"`
	}
	err := ctx.App.
		RecordQuery(adventuria.GameCollections.Get(schema.CollectionItems)).
		Where(dbx.And(
			dbx.HashExp{"type": decodedValue.ItemsType},
			dbx.NewExp("isRollable = true"),
		)).
		Select("id").
		All(&records)
	if err != nil {
		return fmt.Errorf("roll_item.refreshItems(): %w", err)
	}

	ids := make([]string, len(records))
	for i, record := range records {
		ids[i] = record.Id
	}

	user.LastAction().SetItemsList(ids)

	return nil
}
