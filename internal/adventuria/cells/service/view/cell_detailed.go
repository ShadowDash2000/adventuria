package view

import "adventuria/internal/adventuria/model"

type CellDetailed struct {
	Id                  string
	Type                model.CellType
	Name                string
	Icon                string
	Color               string
	Description         string
	Points              int
	DoneEnergyConsume   int
	DropEnergyConsume   int
	Coins               int
	AverageCampaignTime float64

	ActivityFilter *CellDetailedActivityFilter
}

type CellDetailedActivityFilter struct {
	Type            model.ActivityType
	Platforms       []*CellDetailedPlatform
	Developers      []*CellDetailedCompany
	Publishers      []*CellDetailedCompany
	Genres          []*CellDetailedGenre
	Tags            []*CellDetailedTag
	Themes          []*CellDetailedTheme
	MinPrice        int
	MaxPrice        int
	MinCampaignTime float64
	MaxCampaignTime float64
	Activities      []*CellDetailedActivity
}

type CellDetailedPlatform struct {
	Id   string
	Name string
}

type CellDetailedCompany struct {
	Id   string
	Name string
}

type CellDetailedGenre struct {
	Id   string
	Name string
}

type CellDetailedTag struct {
	Id   string
	Name string
}

type CellDetailedTheme struct {
	Id   string
	Name string
}

type CellDetailedActivity struct {
	Id               string
	Name             string
	HLTBCampaignTime float64
	Cover            string
	CoverAlt         string
}
