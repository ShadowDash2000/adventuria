package event_stats

import "github.com/pocketbase/pocketbase/core"

type EventStatsData struct {
	MostGamesCompleted    []EventStatEntry `json:"most_games_completed"`
	MostDrops             []EventStatEntry `json:"most_drops"`
	MostRerolls           []EventStatEntry `json:"most_rerolls"`
	MostGymsCompleted     []EventStatEntry `json:"most_gyms_completed"`
	MostMoviesWatched     []EventStatEntry `json:"most_movies_watched"`
	MostKaraokeCompleted  []EventStatEntry `json:"most_karaoke_completed"`
	MostWanted            []EventStatEntry `json:"most_wanted"`
	MostItemsUsed         []EventStatEntry `json:"most_items_used"`
	MostRobloxPlayed      []EventStatEntry `json:"most_roblox_played"`
	MostHappyWheelsPlayed []EventStatEntry `json:"most_happy_wheels_played"`
	MostVisitedCells      []EventStatEntry `json:"most_visited_cells"`
	LeastVisitedCells     []EventStatEntry `json:"least_visited_cells"`
	MostUsedItems         []EventStatEntry `json:"most_used_items"`
}

type EventStatEntry struct {
	Count  int          `json:"count"`
	Record *core.Record `json:"record"`
}
