package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"fmt"
)

var _ adventuria.CellWheel = (*CellCinema)(nil)

type CellCinema struct {
	adventuria.CellRecord
}

func NewCellCinema() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellCinema{
			adventuria.CellRecord{},
		}
	}
}

func (c *CellCinema) Verify(_ string) error {
	return nil
}

func (c *CellCinema) DecodeValue(_ string) (any, error) {
	return nil, nil
}

func (c *CellCinema) Roll(user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
	items, err := user.LastAction().ItemsList()
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't unmarshal items list",
		}, fmt.Errorf("cinema.roll(): can't unmarshal items list: %w", err)
	}

	if len(items) == 0 {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: no items to roll",
		}, fmt.Errorf("cinema.roll(): no items to roll")
	}

	records, err := adventuria.PocketBase.FindRecordsByIds(
		adventuria.GameCollections.Get(adventuria.CollectionActivities),
		items,
	)
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't fetch records",
		}, fmt.Errorf("cinema.roll(): can't fetch records: %w", err)
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

func (c *CellCinema) RefreshItems(user adventuria.User) error {
	return c.refreshItems(user)
}

func (c *CellCinema) OnCellReached(ctx *adventuria.CellReachedContext) error {
	return c.refreshItems(ctx.User)
}

func (c *CellCinema) refreshItems(user adventuria.User) error {
	filter, err := newActivityFilterById(c.Filter())
	if err != nil {
		return err
	}
	filter.SetType(adventuria.ActivityTypeMovie)
	return updateActivitiesFromFilter(user, filter, true)
}
