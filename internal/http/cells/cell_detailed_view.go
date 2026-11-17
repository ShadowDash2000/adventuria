package cells

import (
	"adventuria/internal/adventuria/cells/service/view"
	"adventuria/internal/adventuria/model"
)

type cellDetailedView struct {
	Id                  string         `json:"id"`
	Type                model.CellType `json:"type"`
	Name                string         `json:"name"`
	Icon                string         `json:"icon"`
	Color               string         `json:"color"`
	Description         string         `json:"description"`
	Points              int            `json:"points"`
	DoneEnergyConsume   int            `json:"done_energy_consume"`
	DropEnergyConsume   int            `json:"drop_energy_consume"`
	Coins               int            `json:"coins"`
	AverageCampaignTime float64        `json:"average_campaign_time"`

	ActivityFilter *cellDetailedActivityFilterView `json:"activity_filter"`
}

func cellDetailedToView(cell *view.CellDetailed) *cellDetailedView {
	return &cellDetailedView{
		Id:                  cell.Id,
		Type:                cell.Type,
		Name:                cell.Name,
		Icon:                cell.Icon,
		Color:               cell.Color,
		Description:         cell.Description,
		Points:              cell.Points,
		DoneEnergyConsume:   cell.DoneEnergyConsume,
		DropEnergyConsume:   cell.DropEnergyConsume,
		Coins:               cell.Coins,
		AverageCampaignTime: cell.AverageCampaignTime,
		ActivityFilter:      cellDetailedActivityFilterToView(cell.ActivityFilter),
	}
}

type cellDetailedActivityFilterView struct {
	Type            model.ActivityType          `json:"type"`
	Platforms       []*cellDetailedPlatformView `json:"platforms"`
	Developers      []*cellDetailedCompanyView  `json:"developers"`
	Publishers      []*cellDetailedCompanyView  `json:"publishers"`
	Genres          []*cellDetailedGenreView    `json:"genres"`
	Tags            []*cellDetailedTagsView     `json:"tags"`
	Themes          []*cellDetailedThemeView    `json:"themes"`
	MinPrice        int                         `json:"min_price"`
	MaxPrice        int                         `json:"max_price"`
	MinCampaignTime float64                     `json:"min_campaign_time"`
	MaxCampaignTime float64                     `json:"max_campaign_time"`
	Activities      []*cellDetailedActivityView `json:"activities"`
}

func cellDetailedActivityFilterToView(filter *view.CellDetailedActivityFilter) *cellDetailedActivityFilterView {
	if filter == nil {
		return nil
	}

	return &cellDetailedActivityFilterView{
		Type:            filter.Type,
		Platforms:       cellDetailedPlatformsToView(filter.Platforms),
		Developers:      cellDetailedCompaniesToView(filter.Developers),
		Publishers:      cellDetailedCompaniesToView(filter.Publishers),
		Genres:          cellDetailedGenresToView(filter.Genres),
		Tags:            cellDetailedTagsToView(filter.Tags),
		Themes:          cellDetailedThemesToView(filter.Themes),
		MinPrice:        filter.MinPrice,
		MaxPrice:        filter.MaxPrice,
		MinCampaignTime: filter.MinCampaignTime,
		MaxCampaignTime: filter.MaxCampaignTime,
		Activities:      cellDetailedActivitiesToView(filter.Activities),
	}
}

type cellDetailedPlatformView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func cellDetailedPlatformToView(platform *view.CellDetailedPlatform) *cellDetailedPlatformView {
	return &cellDetailedPlatformView{
		Id:   platform.Id,
		Name: platform.Name,
	}
}

func cellDetailedPlatformsToView(platforms []*view.CellDetailedPlatform) []*cellDetailedPlatformView {
	res := make([]*cellDetailedPlatformView, len(platforms))
	for i, platform := range platforms {
		res[i] = cellDetailedPlatformToView(platform)
	}
	return res
}

type cellDetailedCompanyView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func cellDetailedCompanyToView(company *view.CellDetailedCompany) *cellDetailedCompanyView {
	return &cellDetailedCompanyView{
		Id:   company.Id,
		Name: company.Name,
	}
}

func cellDetailedCompaniesToView(companies []*view.CellDetailedCompany) []*cellDetailedCompanyView {
	res := make([]*cellDetailedCompanyView, len(companies))
	for i, company := range companies {
		res[i] = cellDetailedCompanyToView(company)
	}
	return res
}

type cellDetailedGenreView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func cellDetailedGenreToView(genre *view.CellDetailedGenre) *cellDetailedGenreView {
	return &cellDetailedGenreView{
		Id:   genre.Id,
		Name: genre.Name,
	}
}

func cellDetailedGenresToView(genres []*view.CellDetailedGenre) []*cellDetailedGenreView {
	res := make([]*cellDetailedGenreView, len(genres))
	for i, genre := range genres {
		res[i] = cellDetailedGenreToView(genre)
	}
	return res
}

type cellDetailedTagsView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func cellDetailedTagToView(tag *view.CellDetailedTag) *cellDetailedTagsView {
	return &cellDetailedTagsView{
		Id:   tag.Id,
		Name: tag.Name,
	}
}

func cellDetailedTagsToView(tags []*view.CellDetailedTag) []*cellDetailedTagsView {
	res := make([]*cellDetailedTagsView, len(tags))
	for i, tag := range tags {
		res[i] = cellDetailedTagToView(tag)
	}
	return res
}

type cellDetailedThemeView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func cellDetailedThemeToView(theme *view.CellDetailedTheme) *cellDetailedThemeView {
	return &cellDetailedThemeView{
		Id:   theme.Id,
		Name: theme.Name,
	}
}

func cellDetailedThemesToView(themes []*view.CellDetailedTheme) []*cellDetailedThemeView {
	res := make([]*cellDetailedThemeView, len(themes))
	for i, theme := range themes {
		res[i] = cellDetailedThemeToView(theme)
	}
	return res
}

type cellDetailedActivityView struct {
	Id               string  `json:"id"`
	Name             string  `json:"name"`
	HLTBCampaignTime float64 `json:"hltb_campaign_time"`
	Cover            string  `json:"cover"`
	CoverAlt         string  `json:"cover_alt"`
}

func cellDetailedActivityToView(activity *view.CellDetailedActivity) *cellDetailedActivityView {
	return &cellDetailedActivityView{
		Id:               activity.Id,
		Name:             activity.Name,
		HLTBCampaignTime: activity.HLTBCampaignTime,
		Cover:            activity.Cover,
		CoverAlt:         activity.CoverAlt,
	}
}

func cellDetailedActivitiesToView(activities []*view.CellDetailedActivity) []*cellDetailedActivityView {
	res := make([]*cellDetailedActivityView, len(activities))
	for i, activity := range activities {
		res[i] = cellDetailedActivityToView(activity)
	}
	return res
}
