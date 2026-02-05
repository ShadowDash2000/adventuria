package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"fmt"
)

var _ adventuria.CellWheel = (*CellActivity)(nil)

type CellActivity struct {
	adventuria.CellRecord
	activityType adventuria.ActivityType
}

func NewCellActivity(activityType adventuria.ActivityType) adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellActivity{
			activityType: activityType,
		}
	}
}

func (c *CellActivity) Verify(_ adventuria.AppContext, _ string) error {
	return nil
}

func (c *CellActivity) Roll(ctx adventuria.AppContext, user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
	items, err := user.LastAction().ItemsList()
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't unmarshal items list",
		}, fmt.Errorf("%s.roll(): can't unmarshal items list: %w", c.activityType, err)
	}

	if len(items) == 0 {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: no items to roll",
		}, fmt.Errorf("%s.roll(): no items to roll", c.activityType)
	}

	records, err := ctx.App.FindRecordsByIds(
		adventuria.GameCollections.Get(adventuria.CollectionActivities),
		items,
	)
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't fetch records",
		}, fmt.Errorf("%s.roll(): can't fetch records: %w", c.activityType, err)
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

func (c *CellActivity) RefreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	return c.refreshItems(ctx, user)
}

func (c *CellActivity) OnCellReached(ctx *adventuria.CellReachedContext) error {
	return c.refreshItems(ctx.AppContext, ctx.User)
}

func (c *CellActivity) OnCellLeft(_ *adventuria.CellLeftContext) error {
	return nil
}

func (c *CellActivity) refreshItems(ctx adventuria.AppContext, user adventuria.User) error {
	filter, err := newActivityFilterById(ctx.App, c.Filter())
	if err != nil {
		return err
	}
	filter.SetType(c.activityType)

	if err = updateActivitiesFromFilter(ctx.App, user, filter, true); err != nil {
		return err
	}

	return ctx.App.Save(user.LastAction().ProxyRecord())
}
