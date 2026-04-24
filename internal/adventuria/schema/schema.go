package schema

const (
	CollectionPlayers         = "players"
	CollectionPlayersProgress = "players_progress"
	CollectionActions         = "actions"
	CollectionCells           = "cells"
	CollectionItems           = "items"
	CollectionEffects         = "effects"
	CollectionInventory       = "inventory"
	CollectionSettings        = "settings"
	CollectionSeasons         = "seasons"
	CollectionsWorlds         = "worlds"

	CollectionActivities     = "activities"
	CollectionCompanies      = "companies"
	CollectionPlatforms      = "platforms"
	CollectionGenres         = "genres"
	CollectionTags           = "tags"
	CollectionThemes         = "themes"
	CollectionGameTypes      = "game_types"
	CollectionActivityFilter = "activity_filter"
	CollectionHowLongToBeat  = "howlongtobeat"
	CollectionSteamSpy       = "steam_spy"
	CollectionCheapshark     = "cheapshark"

	CollectionActivitiesPlatforms  = "activities_platforms"
	CollectionActivitiesDevelopers = "activities_developers"
	CollectionActivitiesPublishers = "activities_publishers"
	CollectionActivitiesGenres     = "activities_genres"
	CollectionActivitiesTags       = "activities_tags"
	CollectionActivitiesThemes     = "activities_themes"
)

var PlayerSchema = struct {
	Id               string
	Name             string
	Avatar           string
	Color            string
	Twitch           string
	YouTube          string
	YouTubeChannelId string
	IsStreamLive     string
}{
	Id:               "id",
	Name:             "name",
	Avatar:           "avatar",
	Color:            "color",
	Twitch:           "twitch",
	YouTube:          "youtube",
	YouTubeChannelId: "youtube_channel_id",
	IsStreamLive:     "is_stream_live",
}

var PlayerProgressSchema = struct {
	Id                string
	Player            string
	Season            string
	CurrentWorld      string
	Points            string
	Balance           string
	CellsPassed       string
	IsInJail          string
	DropsInARow       string
	ItemWheelsCount   string
	MaxInventorySlots string
	Stats             string
	ClearStats        string
}{
	Id:                "id",
	Player:            "player",
	Season:            "season",
	CurrentWorld:      "current_world",
	Points:            "points",
	Balance:           "balance",
	CellsPassed:       "cells_passed",
	IsInJail:          "is_in_jail",
	DropsInARow:       "drops_in_a_row",
	ItemWheelsCount:   "item_wheels_count",
	MaxInventorySlots: "max_inventory_slots",
	Stats:             "stats",
	ClearStats:        "clear_stats",
}

var SeasonSchema = struct {
	Id              string
	Name            string
	Slug            string
	SeasonDateStart string
	SeasonDateEnd   string
}{
	Id:              "id",
	Name:            "name",
	Slug:            "slug",
	SeasonDateStart: "season_date_start",
	SeasonDateEnd:   "season_date_end",
}

var ActionSchema = struct {
	Id                   string
	Player               string
	Cell                 string
	Type                 string
	Activity             string
	Comment              string
	CellsPassed          string
	ItemsList            string
	UsedItems            string
	CanMove              string
	CustomActivityFilter string
}{
	Id:                   "id",
	Player:               "player",
	Cell:                 "cell",
	Type:                 "type",
	Activity:             "activity",
	Comment:              "comment",
	CellsPassed:          "cells_passed",
	ItemsList:            "items_list",
	UsedItems:            "used_items",
	CanMove:              "can_move",
	CustomActivityFilter: "custom_activity_filter",
}

var ActivitySchema = struct {
	Id               string
	IdDb             string
	Type             string
	Name             string
	Slug             string
	ReleaseDate      string
	Platforms        string
	Developers       string
	Publishers       string
	Genres           string
	Tags             string
	Themes           string
	GameType         string
	SteamAppId       string
	SteamAppPrice    string
	HltbId           string
	HltbCampaignTime string
	Cover            string
	CoverAlt         string
	Checksum         string
}{
	Id:               "id",
	IdDb:             "id_db",
	Type:             "type",
	Name:             "name",
	Slug:             "slug",
	ReleaseDate:      "release_date",
	Platforms:        "platforms",
	Developers:       "developers",
	Publishers:       "publishers",
	Genres:           "genres",
	Tags:             "tags",
	Themes:           "themes",
	GameType:         "game_type",
	SteamAppId:       "steam_app_id",
	SteamAppPrice:    "steam_app_price",
	HltbId:           "hltb_id",
	HltbCampaignTime: "hltb_campaign_time",
	Cover:            "cover",
	CoverAlt:         "cover_alt",
	Checksum:         "checksum",
}

var HowLongToBeatSchema = struct {
	Id       string
	IdDb     string
	Name     string
	Year     string
	Campaign string
}{
	Id:       "id",
	IdDb:     "id_db",
	Name:     "name",
	Year:     "year",
	Campaign: "campaign",
}

var InventorySchema = struct {
	Id             string
	Activated      string
	Player         string
	Item           string
	IsActive       string
	AppliedEffects string
}{
	Id:             "id",
	Activated:      "activated",
	Player:         "player",
	Item:           "item",
	IsActive:       "is_active",
	AppliedEffects: "applied_effects",
}

var ItemSchema = struct {
	Id                string
	Name              string
	Icon              string
	Effects           string
	Order             string
	IsUsingSlot       string
	IsActiveByDefault string
	CanDrop           string
	IsRollable        string
	Description       string
	Type              string
	Price             string
}{
	Id:                "id",
	Name:              "name",
	Icon:              "icon",
	Effects:           "effects",
	Order:             "order",
	IsUsingSlot:       "is_using_slot",
	IsActiveByDefault: "is_active_by_default",
	CanDrop:           "can_drop",
	IsRollable:        "is_rollable",
	Description:       "description",
	Type:              "type",
	Price:             "price",
}

var EffectSchema = struct {
	Id    string
	Name  string
	Type  string
	Value string
}{
	Id:    "id",
	Name:  "name",
	Type:  "type",
	Value: "value",
}

var CellSchema = struct {
	Id                       string
	Disabled                 string
	Sort                     string
	Type                     string
	World                    string
	Filter                   string
	AudioPreset              string
	Icon                     string
	Name                     string
	Points                   string
	Coins                    string
	Description              string
	Color                    string
	CantDrop                 string
	CantReroll               string
	IsSafeDrop               string
	IsCustomFilterNotAllowed string
	IsChangeGameNotAllowed   string
	Value                    string
}{
	Id:                       "id",
	Disabled:                 "disabled",
	Sort:                     "sort",
	Type:                     "type",
	World:                    "world",
	Filter:                   "filter",
	AudioPreset:              "audio_preset",
	Icon:                     "icon",
	Name:                     "name",
	Points:                   "points",
	Coins:                    "coins",
	Description:              "description",
	Color:                    "color",
	CantDrop:                 "cant_drop",
	CantReroll:               "cant_reroll",
	IsSafeDrop:               "is_safe_drop",
	IsCustomFilterNotAllowed: "is_custom_filter_not_allowed",
	IsChangeGameNotAllowed:   "is_change_game_not_allowed",
	Value:                    "value",
}

var SettingsSchema = struct {
	EventEnded        string
	CurrentSeason     string
	CurrentWeek       string
	BlockAllActions   string
	MaxInventorySlots string
	PointsForDrop     string
	DropsToJail       string

	IgdbGamesParsed         string
	DisableIgdbParser       string
	DisableSteamParser      string
	DisableCheapsharkParser string
	DisableHltbParser       string
	DisableRefreshHltbTime  string
	KillParser              string
	IgdbForceUpdateGames    string
}{
	EventEnded:              "event_ended",
	CurrentSeason:           "current_season",
	CurrentWeek:             "current_week",
	BlockAllActions:         "block_all_actions",
	MaxInventorySlots:       "max_inventory_slots",
	PointsForDrop:           "points_for_drop",
	DropsToJail:             "drops_to_jail",
	IgdbGamesParsed:         "igdb_games_parsed",
	DisableIgdbParser:       "disable_igdb_parser",
	DisableSteamParser:      "disable_steam_parser",
	DisableCheapsharkParser: "disable_cheapshark_parser",
	DisableHltbParser:       "disable_hltb_parser",
	DisableRefreshHltbTime:  "disable_refresh_hltb_time",
	KillParser:              "kill_parser",
	IgdbForceUpdateGames:    "igdb_force_update_games",
}

var ActivitiesPlatformsSchema = struct {
	Id       string
	Activity string
	Platform string
}{
	Id:       "id",
	Activity: "activity",
	Platform: "platform",
}

var ActivitiesDevelopersSchema = struct {
	Id        string
	Activity  string
	Developer string
}{
	Id:        "id",
	Activity:  "activity",
	Developer: "developer",
}

var ActivitiesPublishersSchema = struct {
	Id        string
	Activity  string
	Publisher string
}{
	Id:        "id",
	Activity:  "activity",
	Publisher: "publisher",
}

var ActivitiesGenresSchema = struct {
	Id       string
	Activity string
	Genre    string
}{
	Id:       "id",
	Activity: "activity",
	Genre:    "genre",
}

var ActivitiesTagsSchema = struct {
	Id       string
	Activity string
	Tag      string
}{
	Id:       "id",
	Activity: "activity",
	Tag:      "tag",
}

var ActivitiesThemesSchema = struct {
	Id       string
	Activity string
	Theme    string
}{
	Id:       "id",
	Activity: "activity",
	Theme:    "theme",
}

var TagSchema = struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}{
	Id:       "id",
	IdDb:     "id_db",
	Name:     "name",
	Checksum: "checksum",
}

var GenreSchema = struct {
	Id       string
	IdDb     string
	Name     string
	Checksum string
}{
	Id:       "id",
	IdDb:     "id_db",
	Name:     "name",
	Checksum: "checksum",
}

var WorldsSchema = struct {
	Id                string
	Name              string
	Slug              string
	Sort              string
	IsLoop            string
	IsDefaultWorld    string
	TransitionToWorld string
	Effects           string
}{
	Id:                "id",
	Name:              "name",
	Slug:              "slug",
	Sort:              "sort",
	IsLoop:            "is_loop",
	IsDefaultWorld:    "is_default_world",
	TransitionToWorld: "transition_to_world",
	Effects:           "effects",
}
