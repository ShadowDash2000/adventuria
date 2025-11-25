package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type CellGame struct {
	adventuria.CellWheel
}

func NewCellGame() adventuria.CellCreator {
	return func() adventuria.Cell {
		return &CellGame{
			&adventuria.CellWheelBase{
				CellBase: adventuria.CellBase{},
			},
		}
	}
}

func (c *CellGame) Roll(user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
	items, err := user.LastAction().ItemsList()
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't unmarshal items list",
		}, fmt.Errorf("game.roll(): can't unmarshal items list: %w", err)
	}

	records, err := adventuria.PocketBase.FindRecordsByIds(
		adventuria.GameCollections.Get(adventuria.CollectionGames),
		items,
	)
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't fetch records",
		}, fmt.Errorf("game.roll(): can't fetch records: %w", err)
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

func (c *CellGame) OnCellReached(user adventuria.User) error {
	var filter adventuria.GameFilterRecord

	if c.Filter() != "" {
		filterRecord, err := adventuria.PocketBase.FindRecordById(
			adventuria.GameCollections.Get(adventuria.CollectionGameFilters),
			c.Filter(),
		)
		if err != nil {
			return err
		}

		filter = adventuria.NewGameFilterFromRecord(filterRecord)
	}

	res, err := FetchRecordsByFilter(filter)
	if err != nil {
		return err
	}

	user.LastAction().SetItemsList(res)

	return nil
}

func FetchRecordsByFilter(filter adventuria.GameFilterRecord) ([]string, error) {
	q := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		Limit(20).
		OrderBy("random()")

	if filter != nil {
		q = SetFilters(filter, q)
	}

	var records []*core.Record
	err := q.All(&records)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(records))
	for i, record := range records {
		res[i] = record.Id
	}

	return res, nil
}

func SetFilters(filter adventuria.GameFilterRecord, q *dbx.SelectQuery) *dbx.SelectQuery {
	if len(filter.Platforms()) > 0 {
		q = q.AndWhere(dbx.In("platforms", StringSliceToAny(filter.Platforms())...))
	}
	if len(filter.Developers()) > 0 {
		q = q.AndWhere(dbx.In("developers", StringSliceToAny(filter.Developers())...))
	}
	if len(filter.Publishers()) > 0 {
		q = q.AndWhere(dbx.In("publishers", StringSliceToAny(filter.Publishers())...))
	}
	if len(filter.Genres()) > 0 {
		q = q.AndWhere(dbx.In("genres", StringSliceToAny(filter.Genres())...))
	}
	if len(filter.Tags()) > 0 {
		q = q.AndWhere(dbx.In("tags", StringSliceToAny(filter.Tags())...))
	}
	if len(filter.Games()) > 0 {
		q = q.AndWhere(dbx.In("id", StringSliceToAny(filter.Games())...))
	}

	if filter.MinPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_app_price > {:price}", dbx.Params{"price": filter.MinPrice()}))
	}
	if filter.MaxPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_app_price < {:price}", dbx.Params{"price": filter.MaxPrice()}))
	}

	if !filter.ReleaseDateFrom().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date > {:date}", dbx.Params{"date": filter.ReleaseDateFrom()}))
	}
	if !filter.ReleaseDateTo().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date < {:date}", dbx.Params{"date": filter.ReleaseDateTo()}))
	}

	if filter.MinCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("campaign_time > {:time}", dbx.Params{"time": filter.MinCampaignTime()}))
	}
	if filter.MaxCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("campaign_time < {:time}", dbx.Params{"time": filter.MaxCampaignTime()}))
	}

	return q
}

func StringSliceToAny(slice []string) []any {
	res := make([]any, len(slice))
	for i, s := range slice {
		res[i] = s
	}
	return res
}
