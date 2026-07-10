package model

import "time"

type ActivityFilterData struct {
	Id              string
	Type            ActivityType
	Name            string
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

type ActivityFilter struct {
	data ActivityFilterData
}

func RestoreActivityFilter(data ActivityFilterData) *ActivityFilter {
	return &ActivityFilter{data: data}
}

func (a *ActivityFilter) Id() string {
	return a.data.Id
}

func (a *ActivityFilter) Type() ActivityType {
	return a.data.Type
}

func (a *ActivityFilter) SetType(t ActivityType) {
	a.data.Type = t
}

func (a *ActivityFilter) Name() string {
	return a.data.Name
}

func (a *ActivityFilter) Platforms() []string {
	return a.data.Platforms
}

func (a *ActivityFilter) SetPlatforms(platforms []string) {
	a.data.Platforms = platforms
}

func (a *ActivityFilter) PlatformsStrict() bool {
	return a.data.PlatformsStrict
}

func (a *ActivityFilter) GameTypes() []string {
	return a.data.GameTypes
}

func (a *ActivityFilter) Developers() []string {
	return a.data.Developers
}

func (a *ActivityFilter) SetDevelopers(developers []string) {
	a.data.Developers = developers
}

func (a *ActivityFilter) Publishers() []string {
	return a.data.Publishers
}

func (a *ActivityFilter) SetPublishers(publishers []string) {
	a.data.Publishers = publishers
}

func (a *ActivityFilter) Genres() []string {
	return a.data.Genres
}

func (a *ActivityFilter) SetGenres(genres []string) {
	a.data.Genres = genres
}

func (a *ActivityFilter) Tags() []string {
	return a.data.Tags
}

func (a *ActivityFilter) SetTags(tags []string) {
	a.data.Tags = tags
}

func (a *ActivityFilter) Themes() []string {
	return a.data.Themes
}

func (a *ActivityFilter) SetThemes(themes []string) {
	a.data.Themes = themes
}

func (a *ActivityFilter) MinPrice() int {
	return a.data.MinPrice
}

func (a *ActivityFilter) SetMinPrice(minPrice int) {
	a.data.MinPrice = minPrice
}

func (a *ActivityFilter) MaxPrice() int {
	return a.data.MaxPrice
}

func (a *ActivityFilter) SetMaxPrice(maxPrice int) {
	a.data.MaxPrice = maxPrice
}

func (a *ActivityFilter) ReleaseDateFrom() time.Time {
	return a.data.ReleaseDateFrom
}

func (a *ActivityFilter) SetReleaseDateFrom(releaseDateFrom time.Time) {
	a.data.ReleaseDateFrom = releaseDateFrom
}

func (a *ActivityFilter) ReleaseDateTo() time.Time {
	return a.data.ReleaseDateTo
}

func (a *ActivityFilter) SetReleaseDateTo(releaseDateTo time.Time) {
	a.data.ReleaseDateTo = releaseDateTo
}

func (a *ActivityFilter) MinCampaignTime() float64 {
	return a.data.MinCampaignTime
}

func (a *ActivityFilter) SetMinCampaignTime(minCampaignTime float64) {
	a.data.MinCampaignTime = minCampaignTime
}

func (a *ActivityFilter) MaxCampaignTime() float64 {
	return a.data.MaxCampaignTime
}

func (a *ActivityFilter) SetMaxCampaignTime(maxCampaignTime float64) {
	a.data.MaxCampaignTime = maxCampaignTime
}

func (a *ActivityFilter) Activities() []string {
	return a.data.Activities
}
