package model

import (
	"slices"
	"time"
)

type ActivityFilter struct {
	Type            ActivityType
	Platforms       []string
	PlatformsStrict bool
	GameTypes       []string
	Developers      []string
	Publishers      []string
	Genres          []string
	Tags            []string
	Themes          []string
	MinPrice        int
	MaxPrice        int
	ReleaseDateFrom time.Time
	ReleaseDateTo   time.Time
	MinCampaignTime float64
	MaxCampaignTime float64
	Activities      []string
}

func (a *ActivityFilter) Clone() ActivityFilter {
	if a == nil {
		return ActivityFilter{}
	}

	return ActivityFilter{
		Type:            a.Type,
		Platforms:       slices.Clone(a.Platforms),
		PlatformsStrict: a.PlatformsStrict,
		GameTypes:       slices.Clone(a.GameTypes),
		Developers:      slices.Clone(a.Developers),
		Publishers:      slices.Clone(a.Publishers),
		Genres:          slices.Clone(a.Genres),
		Tags:            slices.Clone(a.Tags),
		Themes:          slices.Clone(a.Themes),
		MinPrice:        a.MinPrice,
		MaxPrice:        a.MaxPrice,
		ReleaseDateFrom: a.ReleaseDateFrom,
		ReleaseDateTo:   a.ReleaseDateTo,
		MinCampaignTime: a.MinCampaignTime,
		MaxCampaignTime: a.MaxCampaignTime,
		Activities:      slices.Clone(a.Activities),
	}
}

func (a *ActivityFilter) CloneNil() *ActivityFilter {
	if a == nil {
		return nil
	}

	return new(a.Clone())
}

type ActivityFilterData struct {
	Id   string
	Name string
	ActivityFilter
}

type ActivityFilterInfo struct {
	data ActivityFilterData
}

func RestoreActivityFilter(data ActivityFilterData) *ActivityFilterInfo {
	return &ActivityFilterInfo{data: data}
}

func (a *ActivityFilterInfo) Id() string {
	return a.data.Id
}

func (a *ActivityFilterInfo) Filter() ActivityFilter {
	return a.data.ActivityFilter.Clone()
}

func (a *ActivityFilterInfo) Type() ActivityType {
	return a.data.Type
}

func (a *ActivityFilterInfo) SetType(t ActivityType) {
	a.data.Type = t
}

func (a *ActivityFilterInfo) Name() string {
	return a.data.Name
}

func (a *ActivityFilterInfo) Platforms() []string {
	return a.data.Platforms
}

func (a *ActivityFilterInfo) SetPlatforms(platforms []string) {
	a.data.Platforms = platforms
}

func (a *ActivityFilterInfo) PlatformsStrict() bool {
	return a.data.PlatformsStrict
}

func (a *ActivityFilterInfo) GameTypes() []string {
	return a.data.GameTypes
}

func (a *ActivityFilterInfo) Developers() []string {
	return a.data.Developers
}

func (a *ActivityFilterInfo) SetDevelopers(developers []string) {
	a.data.Developers = developers
}

func (a *ActivityFilterInfo) Publishers() []string {
	return a.data.Publishers
}

func (a *ActivityFilterInfo) SetPublishers(publishers []string) {
	a.data.Publishers = publishers
}

func (a *ActivityFilterInfo) Genres() []string {
	return a.data.Genres
}

func (a *ActivityFilterInfo) SetGenres(genres []string) {
	a.data.Genres = genres
}

func (a *ActivityFilterInfo) Tags() []string {
	return a.data.Tags
}

func (a *ActivityFilterInfo) SetTags(tags []string) {
	a.data.Tags = tags
}

func (a *ActivityFilterInfo) Themes() []string {
	return a.data.Themes
}

func (a *ActivityFilterInfo) SetThemes(themes []string) {
	a.data.Themes = themes
}

func (a *ActivityFilterInfo) MinPrice() int {
	return a.data.MinPrice
}

func (a *ActivityFilterInfo) SetMinPrice(minPrice int) {
	a.data.MinPrice = minPrice
}

func (a *ActivityFilterInfo) MaxPrice() int {
	return a.data.MaxPrice
}

func (a *ActivityFilterInfo) SetMaxPrice(maxPrice int) {
	a.data.MaxPrice = maxPrice
}

func (a *ActivityFilterInfo) ReleaseDateFrom() time.Time {
	return a.data.ReleaseDateFrom
}

func (a *ActivityFilterInfo) SetReleaseDateFrom(releaseDateFrom time.Time) {
	a.data.ReleaseDateFrom = releaseDateFrom
}

func (a *ActivityFilterInfo) ReleaseDateTo() time.Time {
	return a.data.ReleaseDateTo
}

func (a *ActivityFilterInfo) SetReleaseDateTo(releaseDateTo time.Time) {
	a.data.ReleaseDateTo = releaseDateTo
}

func (a *ActivityFilterInfo) MinCampaignTime() float64 {
	return a.data.MinCampaignTime
}

func (a *ActivityFilterInfo) SetMinCampaignTime(minCampaignTime float64) {
	a.data.MinCampaignTime = minCampaignTime
}

func (a *ActivityFilterInfo) MaxCampaignTime() float64 {
	return a.data.MaxCampaignTime
}

func (a *ActivityFilterInfo) SetMaxCampaignTime(maxCampaignTime float64) {
	a.data.MaxCampaignTime = maxCampaignTime
}

func (a *ActivityFilterInfo) Activities() []string {
	return a.data.Activities
}
