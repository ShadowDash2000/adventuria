package view

import (
	"adventuria/internal/adventuria/cells/service/view"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"

	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

type cellDetailedRow struct {
	Id                  string         `db:"id"`
	Type                model.CellType `db:"type"`
	Name                string         `db:"name"`
	Icon                string         `db:"icon"`
	Color               string         `db:"color"`
	Description         string         `db:"description"`
	Filter              string         `db:"filter"`
	Points              int            `db:"points"`
	DoneEnergyConsume   int            `db:"done_energy_consume"`
	DropEnergyConsume   int            `db:"drop_energy_consume"`
	Coins               int            `db:"coins"`
	AverageCampaignTime float64        `db:"average_campaign_time"`
}

func (c *cellDetailedRow) toCellDetailed(filter *view.CellDetailedActivityFilter) *view.CellDetailed {
	if c == nil {
		return nil
	}

	return &view.CellDetailed{
		Id:                  c.Id,
		Type:                c.Type,
		Name:                c.Name,
		Icon:                pbhelper.BuildFilePath(schema.CollectionCells, c.Id, c.Icon),
		Color:               c.Color,
		Description:         c.Description,
		Points:              c.Points,
		DoneEnergyConsume:   c.DoneEnergyConsume,
		DropEnergyConsume:   c.DropEnergyConsume,
		Coins:               c.Coins,
		AverageCampaignTime: c.AverageCampaignTime,
		ActivityFilter:      filter,
	}
}

type cellDetailedActivityFilterRow struct {
	Id              string                  `db:"id"`
	Type            model.ActivityType      `db:"type"`
	Platforms       types.JSONArray[string] `db:"platforms"`
	Developers      types.JSONArray[string] `db:"developers"`
	Publishers      types.JSONArray[string] `db:"publishers"`
	Genres          types.JSONArray[string] `db:"genres"`
	Tags            types.JSONArray[string] `db:"tags"`
	Themes          types.JSONArray[string] `db:"themes"`
	MinPrice        int                     `db:"min_price"`
	MaxPrice        int                     `db:"max_price"`
	MinCampaignTime float64                 `db:"min_campaign_time"`
	MaxCampaignTime float64                 `db:"max_campaign_time"`
	Activities      types.JSONArray[string] `db:"activities"`
}

type cellDetailedPlatformRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func (c *cellDetailedPlatformRow) toDetailed() *view.CellDetailedPlatform {
	return &view.CellDetailedPlatform{
		Id:   c.Id,
		Name: c.Name,
	}
}

type cellDetailedCompanyRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func (c *cellDetailedCompanyRow) toDetailed() *view.CellDetailedCompany {
	return &view.CellDetailedCompany{
		Id:   c.Id,
		Name: c.Name,
	}
}

type cellDetailedGenreRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func (c *cellDetailedGenreRow) toDetailed() *view.CellDetailedGenre {
	return &view.CellDetailedGenre{
		Id:   c.Id,
		Name: c.Name,
	}
}

type cellDetailedTagRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func (c *cellDetailedTagRow) toDetailed() *view.CellDetailedTag {
	return &view.CellDetailedTag{
		Id:   c.Id,
		Name: c.Name,
	}
}

type cellDetailedThemeRow struct {
	Id   string `db:"id"`
	Name string `db:"name"`
}

func (c *cellDetailedThemeRow) toDetailed() *view.CellDetailedTheme {
	return &view.CellDetailedTheme{
		Id:   c.Id,
		Name: c.Name,
	}
}

type cellDetailedActivityRow struct {
	Id               string  `db:"id"`
	Name             string  `db:"name"`
	HLTBCampaignTime float64 `db:"hltb_campaign_time"`
	Cover            string  `db:"cover"`
	CoverAlt         string  `db:"cover_alt"`
}

func (c *cellDetailedActivityRow) toCellDetailedActivity() *view.CellDetailedActivity {
	return &view.CellDetailedActivity{
		Id:               c.Id,
		Name:             c.Name,
		HLTBCampaignTime: c.HLTBCampaignTime,
		Cover:            c.Cover,
		CoverAlt:         pbhelper.BuildFilePath(schema.CollectionActivities, c.Id, c.CoverAlt),
	}
}
