package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"fmt"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type RerollFilterAction struct {
	adventuria.ActionBase
}

func (a *RerollFilterAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if !ok {
		return false
	}

	if _, ok = currentCell.(*cells.CellGame); !ok {
		return false
	}

	return true
}

type GameFilterRequest struct {
	Platforms       []string `json:"platforms"`
	Developers      []string `json:"developers"`
	Publishers      []string `json:"publishers"`
	Genres          []string `json:"genres"`
	Tags            []string `json:"tags"`
	MinPrice        int      `json:"min_price"`
	MaxPrice        int      `json:"max_price"`
	ReleaseYearFrom int      `json:"release_year_from"`
	ReleaseYearTo   int      `json:"release_year_to"`
	MinCampaignTime float64  `json:"min_campaign_time"`
	MaxCampaignTime float64  `json:"max_campaign_time"`
	Games           []string `json:"games"`
}

func (a *RerollFilterAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	filter := &GameFilterRequest{}
	if f, ok := req["filter"]; ok {
		err := mapstructure.Decode(f, &filter)
		if err != nil {
			return &adventuria.ActionResult{
				Success: false,
				Error:   "internal error: can't decode filter",
			}, fmt.Errorf("roll_filter.do(): can't decode filter: %w", err)
		}
	} else {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: filter is not set",
		}, nil
	}

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
	if filter.ReleaseYearFrom > 0 {
		yearFrom, err := types.ParseDateTime(
			time.Date(filter.ReleaseYearFrom, 1, 1, 0, 0, 0, 0, time.UTC),
		)
		if err != nil {
			return &adventuria.ActionResult{
				Success: false,
				Error:   "request error: failed to parse release year from",
			}, nil
		}
		filterRecord.SetReleaseDateFrom(yearFrom)
	}
	if filter.ReleaseYearTo > 0 {
		yearTo, err := types.ParseDateTime(
			time.Date(filter.ReleaseYearTo, 1, 1, 0, 0, 0, 0, time.UTC).
				AddDate(1, 0, 0),
		)
		if err != nil {
			return &adventuria.ActionResult{
				Success: false,
				Error:   "request error: failed to parse release year to",
			}, nil
		}
		filterRecord.SetReleaseDateTo(yearTo)
	}
	filterRecord.SetMinCampaignTime(filter.MinCampaignTime)
	filterRecord.SetMaxCampaignTime(filter.MaxCampaignTime)
	filterRecord.SetGames(filter.Games)

	res, err := cells.FetchRecordsByFilter(filterRecord)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't fetch records by filter",
		}, fmt.Errorf("roll_filter.do(): can't fetch records by filter: %w", err)
	}

	user.LastAction().SetItemsList(res)

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}
