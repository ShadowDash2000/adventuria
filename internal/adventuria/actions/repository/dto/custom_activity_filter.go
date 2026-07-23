package dto

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type CustomActivityFilter struct {
	Platforms       []string  `json:"platforms"`
	Developers      []string  `json:"developers"`
	Publishers      []string  `json:"publishers"`
	Genres          []string  `json:"genres"`
	Tags            []string  `json:"tags"`
	Themes          []string  `json:"themes"`
	MinPrice        int       `json:"min_price"`
	MaxPrice        int       `json:"max_price"`
	ReleaseDateFrom time.Time `json:"release_date_from"`
	ReleaseDateTo   time.Time `json:"release_date_to"`
	MinCampaignTime float64   `json:"min_campaign_time"`
	MaxCampaignTime float64   `json:"max_campaign_time"`
}

func CustomActivityFilterToDTO(filter model.CustomActivityFilter) CustomActivityFilter {
	return CustomActivityFilter{
		Platforms:       filter.Platforms,
		Developers:      filter.Developers,
		Publishers:      filter.Publishers,
		Genres:          filter.Genres,
		Tags:            filter.Tags,
		Themes:          filter.Themes,
		MinPrice:        filter.MinPrice,
		MaxPrice:        filter.MaxPrice,
		ReleaseDateFrom: filter.ReleaseDateFrom,
		ReleaseDateTo:   filter.ReleaseDateTo,
		MinCampaignTime: filter.MinCampaignTime,
		MaxCampaignTime: filter.MaxCampaignTime,
	}
}

func CustomActivityFilterFromDTO(dto CustomActivityFilter) model.CustomActivityFilter {
	return model.CustomActivityFilter{
		Platforms:       dto.Platforms,
		Developers:      dto.Developers,
		Publishers:      dto.Publishers,
		Genres:          dto.Genres,
		Tags:            dto.Tags,
		Themes:          dto.Themes,
		MinPrice:        dto.MinPrice,
		MaxPrice:        dto.MaxPrice,
		ReleaseDateFrom: dto.ReleaseDateFrom,
		ReleaseDateTo:   dto.ReleaseDateTo,
		MinCampaignTime: dto.MinCampaignTime,
		MaxCampaignTime: dto.MaxCampaignTime,
	}
}
