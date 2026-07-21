package repository

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type actionDataListDTO struct {
	Activities activitiesDataDTO `json:"activities"`
	Items      itemsDataDTO      `json:"items"`
}

type activitiesDataDTO struct {
	Ids []string `json:"ids"`
}

type itemsDataDTO struct {
	Ids             []string `json:"ids"`
	PriceMultiplier float64  `json:"price_multiplier"`
}

func actionDataListToDTO(dataList model.ActionDataList) actionDataListDTO {
	return actionDataListDTO{
		Activities: activitiesDataDTO{
			Ids: dataList.Activities.Ids,
		},
		Items: itemsDataDTO{
			Ids:             dataList.Items.Ids,
			PriceMultiplier: dataList.Items.PriceMultiplier,
		},
	}
}

func dtoToActionDataList(dto actionDataListDTO) model.ActionDataList {
	return model.ActionDataList{
		Activities: model.ActivitiesData{
			Ids: dto.Activities.Ids,
		},
		Items: model.ItemsData{
			Ids:             dto.Items.Ids,
			PriceMultiplier: dto.Items.PriceMultiplier,
		},
	}
}

type customActivityFilterDTO struct {
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

func customActivityFilterToDTO(filter model.CustomActivityFilter) customActivityFilterDTO {
	return customActivityFilterDTO{
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

func dtoToCustomActivityFilter(dto customActivityFilterDTO) model.CustomActivityFilter {
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
