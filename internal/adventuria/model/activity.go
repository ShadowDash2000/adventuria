package model

import (
	"errors"
	"time"
)

type ActivityData struct {
	Id               string
	IdDb             string
	Type             ActivityType
	Name             string
	Slug             string
	ReleaseDate      time.Time
	Platforms        []string
	Developers       []string
	Publishers       []string
	Genres           []string
	Tags             []string
	Themes           []string
	GameType         string
	SteamAppId       uint
	SteamAppPrice    uint
	HltbId           uint
	HltbCampaignTime float64
	Cover            string
	CoverAlt         string
	Checksum         string
}

type Activity struct {
	data  ActivityData
	isNew bool
}

type ActivityCreate struct {
	IdDb     string
	Type     ActivityType
	Name     string
	Checksum string
}

func NewActivity(data ActivityCreate) (*Activity, error) {
	if data.IdDb == "" {
		return nil, errors.New("activity: id_db is empty")
	}
	if data.Type == "" {
		return nil, errors.New("activity: type is empty")
	}
	if data.Name == "" {
		return nil, errors.New("activity: name is empty")
	}
	if data.Checksum == "" {
		return nil, errors.New("activity: checksum is empty")
	}

	return &Activity{
		data: ActivityData{
			IdDb:     data.IdDb,
			Type:     data.Type,
			Name:     data.Name,
			Checksum: data.Checksum,
		},
		isNew: true,
	}, nil
}

func RestoreActivity(data ActivityData) *Activity {
	return &Activity{
		data:  data,
		isNew: false,
	}
}

func (a *Activity) IsNew() bool {
	return a.isNew
}

func (a *Activity) ID() string {
	return a.data.Id
}

func (a *Activity) IdDb() string {
	return a.data.IdDb
}

func (a *Activity) SetIdDb(id string) {
	a.data.IdDb = id
}

func (a *Activity) Type() ActivityType {
	return a.data.Type
}

func (a *Activity) Name() string {
	return a.data.Name
}

func (a *Activity) SetName(name string) {
	a.data.Name = name
}

func (a *Activity) Slug() string {
	return a.data.Slug
}

func (a *Activity) SetSlug(slug string) {
	a.data.Slug = slug
}

func (a *Activity) ReleaseDate() time.Time {
	return a.data.ReleaseDate
}

func (a *Activity) SetReleaseDate(releaseDate time.Time) {
	a.data.ReleaseDate = releaseDate
}

func (a *Activity) Platforms() []string {
	return a.data.Platforms
}

func (a *Activity) SetPlatforms(platforms []string) {
	a.data.Platforms = platforms
}

func (a *Activity) Developers() []string {
	return a.data.Developers
}

func (a *Activity) SetDevelopers(developers []string) {
	a.data.Developers = developers
}

func (a *Activity) Publishers() []string {
	return a.data.Publishers
}

func (a *Activity) SetPublishers(publishers []string) {
	a.data.Publishers = publishers
}

func (a *Activity) Genres() []string {
	return a.data.Genres
}

func (a *Activity) SetGenres(genres []string) {
	a.data.Genres = genres
}

func (a *Activity) Tags() []string {
	return a.data.Tags
}

func (a *Activity) SetTags(tags []string) {
	a.data.Tags = tags
}

func (a *Activity) Themes() []string {
	return a.data.Themes
}

func (a *Activity) SetThemes(themes []string) {
	a.data.Themes = themes
}

func (a *Activity) GameType() string {
	return a.data.GameType
}

func (a *Activity) SetGameType(gameType string) {
	a.data.GameType = gameType
}

func (a *Activity) SteamAppId() uint {
	return a.data.SteamAppId
}

func (a *Activity) SetSteamAppId(steamAppId uint) {
	a.data.SteamAppId = steamAppId
}

func (a *Activity) SteamAppPrice() uint {
	return a.data.SteamAppPrice
}

func (a *Activity) SetSteamAppPrice(steamAppPrice uint) {
	a.data.SteamAppPrice = steamAppPrice
}

func (a *Activity) HltbId() uint {
	return a.data.HltbId
}

func (a *Activity) SetHltbId(hltbId uint) {
	a.data.HltbId = hltbId
}

func (a *Activity) HltbCampaignTime() float64 {
	return a.data.HltbCampaignTime
}

func (a *Activity) SetHltbCampaignTime(hltbCampaignTime float64) {
	a.data.HltbCampaignTime = hltbCampaignTime
}

func (a *Activity) Cover() string {
	return a.data.Cover
}

func (a *Activity) SetCover(cover string) {
	a.data.Cover = cover
}

func (a *Activity) CoverAlt() string {
	return a.data.CoverAlt
}

func (a *Activity) SetCoverAlt(coverAlt string) {
	a.data.CoverAlt = coverAlt
}

func (a *Activity) Checksum() string {
	return a.data.Checksum
}

func (a *Activity) SetChecksum(checksum string) {
	a.data.Checksum = checksum
}
