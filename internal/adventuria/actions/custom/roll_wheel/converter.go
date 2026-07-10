package roll_wheel

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type activityViewDetailed struct {
	Activity   activityView    `json:"activity"`
	Platforms  []platformView  `json:"platforms"`
	Developers []developerView `json:"developers"`
	Publishers []publisherView `json:"publishers"`
	Genres     []genreView     `json:"genres"`
	Tags       []tagView       `json:"tags"`
	Themes     []themeView     `json:"themes"`
}

func toActivityViewDetailed(activityDetailed *model.ActivityViewDetailed) activityViewDetailed {
	return activityViewDetailed{
		Activity:   toActivityView(activityDetailed.Activity()),
		Platforms:  toPlatformViews(activityDetailed.Platforms()),
		Developers: toDeveloperViews(activityDetailed.Developers()),
		Publishers: toPublisherViews(activityDetailed.Publishers()),
		Genres:     toGenreViews(activityDetailed.Genres()),
		Tags:       toTagViews(activityDetailed.Tags()),
		Themes:     toThemeViews(activityDetailed.Themes()),
	}
}

func toActivityViewDetailedList(activitiesDetailed []*model.ActivityViewDetailed) []activityViewDetailed {
	res := make([]activityViewDetailed, len(activitiesDetailed))
	for i, activityDetailed := range activitiesDetailed {
		res[i] = toActivityViewDetailed(activityDetailed)
	}
	return res
}

type activityView struct {
	Id               string             `json:"id"`
	Type             model.ActivityType `json:"type"`
	Name             string             `json:"name"`
	Slug             string             `json:"slug"`
	ReleaseDate      time.Time          `json:"release_date"`
	Platforms        []string           `json:"platforms"`
	Developers       []string           `json:"developers"`
	Publishers       []string           `json:"publishers"`
	Genres           []string           `json:"genres"`
	Tags             []string           `json:"tags"`
	Themes           []string           `json:"themes"`
	GameType         string             `json:"game_type"`
	SteamAppId       uint               `json:"steam_app_id"`
	SteamAppPrice    uint               `json:"steam_app_price"`
	HltbId           uint               `json:"hltb_id"`
	HltbCampaignTime float64            `json:"hltb_campaign_time"`
	Cover            string             `json:"cover"`
	CoverAlt         string             `json:"cover_alt"`
}

func toActivityView(activity *model.Activity) activityView {
	return activityView{
		Id:               activity.ID(),
		Type:             activity.Type(),
		Name:             activity.Name(),
		Slug:             activity.Slug(),
		ReleaseDate:      activity.ReleaseDate(),
		Platforms:        activity.Platforms(),
		Developers:       activity.Developers(),
		Publishers:       activity.Publishers(),
		Genres:           activity.Genres(),
		Tags:             activity.Tags(),
		Themes:           activity.Themes(),
		GameType:         activity.GameType(),
		SteamAppId:       activity.SteamAppId(),
		SteamAppPrice:    activity.SteamAppPrice(),
		HltbId:           activity.HltbId(),
		HltbCampaignTime: activity.HltbCampaignTime(),
		Cover:            activity.Cover(),
		CoverAlt:         activity.CoverAlt(),
	}
}

type platformView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toPlatformView(platform *model.Platform) platformView {
	return platformView{
		Id:   platform.ID(),
		Name: platform.Name(),
	}
}

func toPlatformViews(platforms []*model.Platform) []platformView {
	res := make([]platformView, len(platforms))
	for i, platform := range platforms {
		res[i] = toPlatformView(platform)
	}
	return res
}

type developerView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toDeveloperView(developer *model.Developer) developerView {
	return developerView{
		Id:   developer.ID(),
		Name: developer.Name(),
	}
}

func toDeveloperViews(developers []*model.Developer) []developerView {
	res := make([]developerView, len(developers))
	for i, developer := range developers {
		res[i] = toDeveloperView(developer)
	}
	return res
}

type publisherView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toPublisherView(publisher *model.Publisher) publisherView {
	return publisherView{
		Id:   publisher.ID(),
		Name: publisher.Name(),
	}
}

func toPublisherViews(publishers []*model.Publisher) []publisherView {
	res := make([]publisherView, len(publishers))
	for i, publisher := range publishers {
		res[i] = toPublisherView(publisher)
	}
	return res
}

type genreView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toGenreView(genre *model.Genre) genreView {
	return genreView{
		Id:   genre.ID(),
		Name: genre.Name(),
	}
}

func toGenreViews(genres []*model.Genre) []genreView {
	res := make([]genreView, len(genres))
	for i, genre := range genres {
		res[i] = toGenreView(genre)
	}
	return res
}

type tagView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toTagView(tag *model.Tag) tagView {
	return tagView{
		Id:   tag.ID(),
		Name: tag.Name(),
	}
}

func toTagViews(tags []*model.Tag) []tagView {
	res := make([]tagView, len(tags))
	for i, tag := range tags {
		res[i] = toTagView(tag)
	}
	return res
}

type themeView struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func toThemeView(theme *model.Theme) themeView {
	return themeView{
		Id:   theme.ID(),
		Name: theme.Name(),
	}
}

func toThemeViews(themes []*model.Theme) []themeView {
	res := make([]themeView, len(themes))
	for i, theme := range themes {
		res[i] = toThemeView(theme)
	}
	return res
}
