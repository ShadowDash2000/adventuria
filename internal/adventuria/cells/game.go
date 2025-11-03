package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
	"fmt"

	"github.com/mitchellh/mapstructure"
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

type GameFilterRequest struct {
	Platforms  []string `json:"platforms"`
	Developers []string `json:"developers"`
	Publishers []string `json:"publishers"`
	Genres     []string `json:"genres"`
	Tags       []string `json:"tags"`
	MinPrice   int      `json:"min_price"`
	MaxPrice   int      `json:"max_price"`
	//ReleaseDateFrom string   `json:"release_date_from"`
	//ReleaseDateTo   string   `json:"release_date_to"`
	MinCampaignTime float64  `json:"min_campaign_time"`
	MaxCampaignTime float64  `json:"max_campaign_time"`
	Games           []string `json:"games"`
}

func (c *CellGame) Roll(req adventuria.RollWheelRequest, user adventuria.User) (*adventuria.WheelRollResult, error) {
	var filter *GameFilterRequest
	if f, ok := req["filter"]; ok {
		filter = &GameFilterRequest{}
		err := mapstructure.Decode(f, &filter)
		if err != nil {
			return &adventuria.WheelRollResult{
				Success: false,
				Error:   "internal error: can't decode filter",
			}, fmt.Errorf("game.roll(): can't decode filter: %w", err)
		}
	}

	if filter != nil {
		filterRecord := adventuria.NewGameFilterFromRecord(
			core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionGameFilters)),
		)
		filterRecord.SetPlatforms(filter.Platforms)
		filterRecord.SetDevelopers(filter.Developers)
		filterRecord.SetPublishers(filter.Publishers)
		filterRecord.SetGenres(filter.Genres)
		filterRecord.SetTags(filter.Tags)
		filterRecord.SetMinPrice(filter.MinPrice)
		filterRecord.SetMaxPrice(filter.MaxPrice)
		filterRecord.SetMinCampaignTime(filter.MinCampaignTime)
		filterRecord.SetMaxCampaignTime(filter.MaxCampaignTime)
		filterRecord.SetGames(filter.Games)

		res, err := c.fetchRecordsByFilter(filterRecord)
		if err != nil {
			return &adventuria.WheelRollResult{
				Success: false,
				Error:   "internal error: can't fetch records by filter",
			}, fmt.Errorf("game.roll(): can't fetch records by filter: %w", err)
		}

		user.LastAction().SetItemsList(res)
	}

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

	res, err := c.fetchRecordsByFilter(filter)
	if err != nil {
		return err
	}

	user.LastAction().SetItemsList(res)

	return nil
}

func (c *CellGame) fetchRecordsByFilter(filter adventuria.GameFilterRecord) ([]string, error) {
	q := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		Limit(20).
		OrderBy("random()")

	if filter != nil {
		q = c.setFilters(filter, q)
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

func (c *CellGame) setFilters(filter adventuria.GameFilterRecord, q *dbx.SelectQuery) *dbx.SelectQuery {
	if len(filter.Platforms()) > 0 {
		q = q.AndWhere(dbx.In("platforms", c.stringSliceToAny(filter.Platforms())...))
	}
	if len(filter.Developers()) > 0 {
		q = q.AndWhere(dbx.In("developers", c.stringSliceToAny(filter.Developers())...))
	}
	if len(filter.Publishers()) > 0 {
		q = q.AndWhere(dbx.In("publishers", c.stringSliceToAny(filter.Publishers())...))
	}
	if len(filter.Genres()) > 0 {
		q = q.AndWhere(dbx.In("genres", c.stringSliceToAny(filter.Genres())...))
	}
	if len(filter.Tags()) > 0 {
		q = q.AndWhere(dbx.In("tags", c.stringSliceToAny(filter.Tags())...))
	}
	if len(filter.Games()) > 0 {
		q = q.AndWhere(dbx.In("id", c.stringSliceToAny(filter.Games())...))
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

func (c *CellGame) stringSliceToAny(slice []string) []any {
	res := make([]any, len(slice))
	for i, s := range slice {
		res[i] = s
	}
	return res
}
