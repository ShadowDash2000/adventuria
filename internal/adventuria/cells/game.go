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

func (c *CellGame) Verify(_ string) error {
	return nil
}

func (c *CellGame) DecodeValue(_ string) (any, error) {
	return nil, nil
}

func (c *CellGame) Roll(user adventuria.User, _ adventuria.RollWheelRequest) (*adventuria.WheelRollResult, error) {
	if err := c.checkCustomFilter(user); err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't apply custom filter",
		}, fmt.Errorf("game.roll(): can't apply custom filter: %w", err)
	}

	items, err := user.LastAction().ItemsList()
	if err != nil {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: can't unmarshal items list",
		}, fmt.Errorf("game.roll(): can't unmarshal items list: %w", err)
	}

	if len(items) == 0 {
		return &adventuria.WheelRollResult{
			Success: false,
			Error:   "internal error: no items to roll",
		}, fmt.Errorf("game.roll(): no items to roll")
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

func (c *CellGame) OnCellReached(ctx *adventuria.CellReachedContext) error {
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

	ctx.User.LastAction().SetItemsList(res)

	return nil
}

func (c *CellGame) checkCustomFilter(user adventuria.User) error {
	needToUpdate := false
	customFilter := user.LastAction().CustomGameFilter()
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
	} else {
		filter = adventuria.NewGameFilterFromRecord(
			core.NewRecord(
				adventuria.GameCollections.Get(adventuria.CollectionGameFilters),
			),
		)
	}

	if len(customFilter.Platforms) > 0 {
		filter.SetPlatforms(append(filter.Platforms(), customFilter.Platforms...))
		needToUpdate = true
	}
	if len(customFilter.Developers) > 0 {
		filter.SetDevelopers(append(filter.Developers(), customFilter.Developers...))
		needToUpdate = true
	}
	if len(customFilter.Publishers) > 0 {
		filter.SetPublishers(append(filter.Publishers(), customFilter.Publishers...))
		needToUpdate = true
	}
	if len(customFilter.Genres) > 0 {
		filter.SetGenres(append(filter.Genres(), customFilter.Genres...))
		needToUpdate = true
	}
	if len(customFilter.Tags) > 0 {
		filter.SetTags(append(filter.Tags(), customFilter.Tags...))
		needToUpdate = true
	}
	if customFilter.MinPrice != 0 {
		filter.SetMinPrice(customFilter.MinPrice)
		needToUpdate = true
	}
	if customFilter.MaxPrice != 0 {
		filter.SetMaxPrice(customFilter.MaxPrice)
		needToUpdate = true
	}
	if !customFilter.ReleaseDateFrom.IsZero() {
		filter.SetReleaseDateFrom(customFilter.ReleaseDateFrom)
		needToUpdate = true
	}
	if !customFilter.ReleaseDateTo.IsZero() {
		filter.SetReleaseDateTo(customFilter.ReleaseDateTo)
		needToUpdate = true
	}
	if customFilter.MinCampaignTime != 0 {
		filter.SetMinCampaignTime(customFilter.MinCampaignTime)
		needToUpdate = true
	}
	if customFilter.MaxCampaignTime != 0 {
		filter.SetMaxCampaignTime(customFilter.MaxCampaignTime)
		needToUpdate = true
	}

	if needToUpdate {
		res, err := FetchRecordsByFilter(filter)
		if err != nil {
			return err
		}

		user.LastAction().SetItemsList(res)
	}

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
		q = q.AndWhere(dbx.OrLike("platforms", filter.Platforms()...))
	}
	if len(filter.Developers()) > 0 {
		q = q.AndWhere(dbx.OrLike("developers", filter.Developers()...))
	}
	if len(filter.Publishers()) > 0 {
		q = q.AndWhere(dbx.OrLike("publishers", filter.Publishers()...))
	}
	if len(filter.Genres()) > 0 {
		q = q.AndWhere(dbx.OrLike("genres", filter.Genres()...))
	}
	if len(filter.Tags()) > 0 {
		q = q.AndWhere(dbx.OrLike("tags", filter.Tags()...))
	}
	if len(filter.Games()) > 0 {
		q = q.AndWhere(dbx.OrLike("id", filter.Games()...))
	}

	if filter.MinPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_spy.price > {:price}", dbx.Params{"price": filter.MinPrice()}))
	}
	if filter.MaxPrice() > 0 {
		q = q.AndWhere(dbx.NewExp("steam_spy.price < {:price}", dbx.Params{"price": filter.MaxPrice()}))
	}

	if !filter.ReleaseDateFrom().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date > {:date}", dbx.Params{"date": filter.ReleaseDateFrom()}))
	}
	if !filter.ReleaseDateTo().IsZero() {
		q = q.AndWhere(dbx.NewExp("release_date < {:date}", dbx.Params{"date": filter.ReleaseDateTo()}))
	}

	if filter.MinCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("hltb.campaign > {:time}", dbx.Params{"time": filter.MinCampaignTime()}))
	}
	if filter.MaxCampaignTime() > 0 {
		q = q.AndWhere(dbx.NewExp("hltb.campaign < {:time}", dbx.Params{"time": filter.MaxCampaignTime()}))
	}

	return q
}
