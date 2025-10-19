package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/games"
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

func (c *CellGame) Roll(user adventuria.User) (*adventuria.WheelRollResult, error) {
	res := &adventuria.WheelRollResult{}

	wheelItems, err := c.GetItems(user)
	if err != nil {
		return nil, err
	}

	res.WinnerId = helper.RandomItemFromSlice(wheelItems).Id

	return res, nil
}

func (c *CellGame) GetItems(user adventuria.User) ([]*adventuria.WheelItem, error) {
	var filter adventuria.GameFilterRecord

	if c.Filter() != "" {
		filterRecord, err := adventuria.PocketBase.FindRecordById(
			adventuria.GameCollections.Get(adventuria.CollectionGameFilters),
			c.Filter(),
		)
		if err != nil {
			return nil, err
		}

		filter = adventuria.NewGameFilterFromRecord(filterRecord)
	}

	q := adventuria.PocketBase.RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionGames)).
		Limit(20).
		OrderBy(fmt.Sprintf("sin(id + %d)", user.LastAction().Seed()))

	if filter != nil {
		q = c.setFilters(filter, q)
	}

	var records []*core.Record
	err := q.All(&records)
	if err != nil {
		return nil, err
	}

	gameRecords := make([]games.GameRecord, len(records))
	var res []*adventuria.WheelItem
	for i, record := range records {
		gameRecords[i] = games.NewGameFromRecord(record)
		wheelItem := &adventuria.WheelItem{
			Id:   gameRecords[i].ID(),
			Name: gameRecords[i].Name(),
			Icon: gameRecords[i].Cover(),
		}
		res = append(res, wheelItem)
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
