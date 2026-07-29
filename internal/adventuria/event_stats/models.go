package event_stats

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type EventStatsData struct {
	NextUpdateAt time.Time          `json:"next_update_at"`
	Stats        *EventStatsEntries `json:"stats"`
}

type EventStatsEntries struct {
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

func (e *EventStatsEntries) Clone() *EventStatsEntries {
	if e == nil {
		return nil
	}

	return &EventStatsEntries{
		MostGamesCompleted:    cloneEntries(e.MostGamesCompleted),
		MostDrops:             cloneEntries(e.MostDrops),
		MostRerolls:           cloneEntries(e.MostRerolls),
		MostGymsCompleted:     cloneEntries(e.MostGymsCompleted),
		MostMoviesWatched:     cloneEntries(e.MostMoviesWatched),
		MostKaraokeCompleted:  cloneEntries(e.MostKaraokeCompleted),
		MostWanted:            cloneEntries(e.MostWanted),
		MostItemsUsed:         cloneEntries(e.MostItemsUsed),
		MostRobloxPlayed:      cloneEntries(e.MostRobloxPlayed),
		MostHappyWheelsPlayed: cloneEntries(e.MostHappyWheelsPlayed),
		MostVisitedCells:      cloneEntries(e.MostVisitedCells),
		LeastVisitedCells:     cloneEntries(e.LeastVisitedCells),
		MostUsedItems:         cloneEntries(e.MostUsedItems),
	}
}

func cloneEntries(entries []EventStatEntry) []EventStatEntry {
	if entries == nil {
		return nil
	}
	res := make([]EventStatEntry, len(entries))
	for i, entry := range entries {
		res[i] = entry.Clone()
	}
	return res
}

type EventStatEntry struct {
	Count  int          `json:"count"`
	Record *core.Record `json:"record"`
}

func (e EventStatEntry) Clone() EventStatEntry {
	var record *core.Record
	if e.Record != nil {
		record = e.Record.Clone()
	}
	return EventStatEntry{
		Count:  e.Count,
		Record: record,
	}
}
