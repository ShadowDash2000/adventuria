package schema

const (
	CollectionUsers     = "users"
	CollectionActions   = "actions"
	CollectionCells     = "cells"
	CollectionItems     = "items"
	CollectionEffects   = "effects"
	CollectionInventory = "inventory"
	CollectionTimers    = "timers"
	CollectionSettings  = "settings"

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
)

var UserSchema = struct {
	Id   string
	Name string

	Points            string
	CellsPassed       string
	IsInJail          string
	DropsInARow       string
	MaxInventorySlots string
	ItemWheelsCount   string
	Balance           string
	Stats             string
	ClearStats        string

	Twitch           string
	YouTube          string
	YouTubeChannelId string
	IsStreamLive     string
}{
	Id:                "id",
	Name:              "name",
	Points:            "points",
	CellsPassed:       "cellsPassed",
	IsInJail:          "isInJail",
	DropsInARow:       "dropsInARow",
	MaxInventorySlots: "maxInventorySlots",
	ItemWheelsCount:   "itemWheelsCount",
	Balance:           "balance",
	Stats:             "stats",
	ClearStats:        "clear_stats",
	Twitch:            "twitch",
	YouTube:           "youtube",
	YouTubeChannelId:  "youtube_channel_id",
	IsStreamLive:      "is_stream_live",
}

var ActionSchema = struct {
	Id                   string
	User                 string
	Cell                 string
	Type                 string
	Activity             string
	Comment              string
	DiceRoll             string
	ItemsList            string
	UsedItems            string
	CanMove              string
	CustomActivityFilter string
}{
	Id:                   "id",
	User:                 "user",
	Cell:                 "cell",
	Type:                 "type",
	Activity:             "activity",
	Comment:              "comment",
	DiceRoll:             "diceRoll",
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

var InventorySchema = struct {
	Id             string
	Activated      string
	User           string
	Item           string
	IsActive       string
	AppliedEffects string
}{
	Id:             "id",
	Activated:      "activated",
	User:           "user",
	Item:           "item",
	IsActive:       "isActive",
	AppliedEffects: "appliedEffects",
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
	IsUsingSlot:       "isUsingSlot",
	IsActiveByDefault: "isActiveByDefault",
	CanDrop:           "canDrop",
	IsRollable:        "isRollable",
	Description:       "description",
	Type:              "type",
	Price:             "price",
}

var CellSchema = struct {
	Id                       string
	Sort                     string
	Type                     string
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
	Value                    string
}{
	Id:                       "id",
	Sort:                     "sort",
	Type:                     "type",
	Filter:                   "filter",
	AudioPreset:              "audio_preset",
	Icon:                     "icon",
	Name:                     "name",
	Points:                   "points",
	Coins:                    "coins",
	Description:              "description",
	Color:                    "color",
	CantDrop:                 "cantDrop",
	CantReroll:               "cantReroll",
	IsSafeDrop:               "isSafeDrop",
	IsCustomFilterNotAllowed: "is_custom_filter_not_allowed",
	Value:                    "value",
}

var SettingsSchema = struct {
	EventDateStart     string
	CurrentWeek        string
	TimerTimeLimit     string
	LimitExceedPenalty string
	BlockAllActions    string
	PointsForDrop      string
	DropsToJail        string

	IgdbGamesParsed         string
	DisableIgdbParser       string
	DisableSteamParser      string
	DisableCheapsharkParser string
	DisableHltbParser       string
	KillParser              string
	IgdbForceUpdateGames    string
}{
	EventDateStart:          "eventDateStart",
	CurrentWeek:             "currentWeek",
	TimerTimeLimit:          "timerTimeLimit",
	LimitExceedPenalty:      "limitExceedPenalty",
	BlockAllActions:         "blockAllActions",
	PointsForDrop:           "pointsForDrop",
	DropsToJail:             "dropsToJail",
	IgdbGamesParsed:         "igdb_games_parsed",
	DisableIgdbParser:       "disable_igdb_parser",
	DisableSteamParser:      "disable_steam_parser",
	DisableCheapsharkParser: "disable_cheapshark_parser",
	DisableHltbParser:       "disable_hltb_parser",
	KillParser:              "kill_parser",
	IgdbForceUpdateGames:    "igdb_force_update_games",
}

var TimerSchema = struct {
	Id         string
	User       string
	IsActive   string
	TimePassed string
	TimeLimit  string
	StartTime  string
}{
	Id:         "id",
	User:       "user",
	IsActive:   "isActive",
	TimePassed: "timePassed",
	TimeLimit:  "timeLimit",
	StartTime:  "startTime",
}
