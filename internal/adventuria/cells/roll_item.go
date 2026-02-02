package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"fmt"

	"github.com/pocketbase/dbx"
)

var _ adventuria.CellWheel = (*CellRollItem)(nil)

type CellRollItem struct {
	adventuria.CellRecord
}

func NewCellRollItem() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellRollItem{
			adventuria.CellRecord{},
		}
	}
}

func (c *CellRollItem) Verify(val string) error {
	_, ok := adventuria.ItemTypes[adventuria.ItemType(val)]
	if !ok {
		return fmt.Errorf("roll_item: unknown item type %s", val)
	}

	return nil
}

func (c *CellRollItem) Roll(user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
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

	records, err := adventuria.PocketBase.FindRecordsByIds(
		adventuria.GameCollections.Get(adventuria.CollectionItems),
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

func (c *CellRollItem) RefreshItems(user adventuria.User) error {
	return c.refreshItems(user)
}

func (c *CellRollItem) OnCellReached(ctx *adventuria.CellReachedContext) error {
	return c.refreshItems(ctx.User)
}

func (c *CellRollItem) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellRollItem) refreshItems(user adventuria.User) error {
	var records []struct {
		Id string `db:"id"`
	}
	err := adventuria.PocketBase.
		RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionItems)).
		Where(dbx.And(
			dbx.HashExp{"type": c.Value()},
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
