package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/result"
	"fmt"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type EventStats struct {
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

func ComputeEventStats(app core.App) (*result.Result, error) {
	if !GameSettings.EventEnded() {
		return result.Err("Event has not ended"), nil
	}

	stats := &EventStats{
		MostGamesCompleted:    []EventStatEntry{},
		MostDrops:             []EventStatEntry{},
		MostRerolls:           []EventStatEntry{},
		MostGymsCompleted:     []EventStatEntry{},
		MostMoviesWatched:     []EventStatEntry{},
		MostKaraokeCompleted:  []EventStatEntry{},
		MostWanted:            []EventStatEntry{},
		MostItemsUsed:         []EventStatEntry{},
		MostRobloxPlayed:      []EventStatEntry{},
		MostHappyWheelsPlayed: []EventStatEntry{},
		MostVisitedCells:      []EventStatEntry{},
		LeastVisitedCells:     []EventStatEntry{},
		MostUsedItems:         []EventStatEntry{},
	}

	type rawStat struct {
		RecordID string `db:"record_id"`
		Count    int    `db:"count"`
		Type     string `db:"stat_type"`
	}

	var rawStats []rawStat

	query := fmt.Sprintf(`
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'games' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[3]s as cells ON cells.%[4]s = actions.%[5]s
		WHERE actions.%[6]s = 'done' AND (cells.%[7]s = 'game' OR cells.%[7]s = 'jail')
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'drops' as "stat_type"
		FROM %[2]s as actions
		WHERE actions.%[6]s = 'drop'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'rerolls' as "stat_type"
		FROM %[2]s as actions
		WHERE actions.%[6]s = 'reroll'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'gyms' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[3]s as cells ON cells.%[4]s = actions.%[5]s
		WHERE actions.%[6]s = 'done' AND cells.%[7]s = 'gym'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'movies' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[3]s as cells ON cells.%[4]s = actions.%[5]s
		WHERE actions.%[6]s = 'done' AND cells.%[7]s = 'movie'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'karaoke' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[3]s as cells ON cells.%[4]s = actions.%[5]s
		WHERE actions.%[6]s = 'done' AND cells.%[7]s = 'karaoke'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'roblox' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[8]s as activities ON activities.%[9]s = actions.%[11]s
		WHERE actions.%[6]s = 'done' AND activities.%[10]s = 'Roblox'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT %[1]s as "record_id", COUNT(*) as "count", 'happywheels' as "stat_type"
		FROM %[2]s as actions
		LEFT JOIN %[8]s as activities ON activities.%[9]s = actions.%[11]s
		WHERE actions.%[6]s = 'done' AND activities.%[10]s = 'Happy Wheels'
		GROUP BY actions.%[1]s
		UNION ALL
		SELECT actions.%[5]s as "record_id", COUNT(*) as "count", 'cell_visits' as "stat_type"
		FROM %[2]s as actions
		WHERE actions.%[6]s IN ('done', 'drop', 'rollDice', 'rollWheel', 'move', 'rollItemOnCell')
		GROUP BY actions.%[5]s
		UNION ALL
		SELECT items.value as "record_id", COUNT(*) as "count", 'items' as "stat_type"
		FROM %[2]s as actions
		CROSS JOIN json_each(actions.%[12]s) as items
		WHERE actions.%[12]s IS NOT NULL AND actions.%[12]s != '[]' AND actions.%[12]s != 'null'
		GROUP BY items.value
	`,
		schema.ActionSchema.Player,    // 1
		schema.CollectionActions,      // 2
		schema.CollectionCells,        // 3
		schema.CellSchema.Id,          // 4
		schema.ActionSchema.Cell,      // 5
		schema.ActionSchema.Type,      // 6
		schema.CellSchema.Type,        // 7
		schema.CollectionActivities,   // 8
		schema.ActivitySchema.Id,      // 9
		schema.ActivitySchema.Name,    // 10
		schema.ActionSchema.Activity,  // 11
		schema.ActionSchema.UsedItems, // 12
	)

	err := app.DB().NewQuery(query).All(&rawStats)
	if err != nil {
		return result.Err("Failed to fetch statistics"), err
	}

	playerIds := make([]string, 0, len(rawStats))
	cellIds := make([]string, 0, len(rawStats))
	itemIds := make([]string, 0, len(rawStats))
	for _, rs := range rawStats {
		if rs.RecordID == "" {
			continue
		}
		switch rs.Type {
		case "cell_visits":
			if !slices.Contains(cellIds, rs.RecordID) {
				cellIds = append(cellIds, rs.RecordID)
			}
		case "items":
			if !slices.Contains(itemIds, rs.RecordID) {
				itemIds = append(itemIds, rs.RecordID)
			}
		default:
			if !slices.Contains(playerIds, rs.RecordID) {
				playerIds = append(playerIds, rs.RecordID)
			}
		}
	}

	var players []*core.Record
	var playersProgress []*core.Record
	if len(playerIds) > 0 {
		players, err = app.FindRecordsByIds(schema.CollectionPlayers, playerIds)
		if err != nil {
			return result.Err("Failed to fetch players"), err
		}

		playerIdsAny := make([]any, len(playerIds))
		for i, id := range playerIds {
			playerIdsAny[i] = id
		}

		err = app.RecordQuery(schema.CollectionPlayersProgress).
			Where(
				dbx.In(schema.PlayerProgressSchema.Player, playerIdsAny...),
			).
			AndWhere(dbx.HashExp{
				schema.PlayerProgressSchema.Season: GameSettings.CurrentSeason(),
			}).
			All(&playersProgress)
		if err != nil {
			return result.Err("Failed to fetch players progress"), err
		}
	}

	playersProgressMap := make(map[string]*core.Record, len(playersProgress))
	for _, p := range playersProgress {
		playersProgressMap[p.Id] = p
	}

	playerMap := make(map[string]*core.Record, len(players))
	for _, u := range players {
		playerMap[u.Id] = u

		playerProgress, ok := playersProgressMap[u.Id]
		if !ok {
			continue
		}

		var playerStats Stats
		err = playerProgress.UnmarshalJSONField(schema.PlayerProgressSchema.Stats, &playerStats)
		if err != nil {
			return result.Err("Failed to unmarshal player stats"), err
		}

		stats.MostWanted = append(stats.MostWanted, EventStatEntry{
			Count:  playerStats.WasInJail,
			Record: u,
		})
		stats.MostItemsUsed = append(stats.MostItemsUsed, EventStatEntry{
			Count:  playerStats.ItemsUsed,
			Record: u,
		})
	}

	var cellRecords []*core.Record
	if len(cellIds) > 0 {
		cellRecords, err = app.FindRecordsByIds(schema.CollectionCells, cellIds)
		if err != nil {
			return result.Err("Failed to fetch cells"), err
		}
	}

	cellMap := make(map[string]*core.Record, len(cellRecords))
	for _, c := range cellRecords {
		cellMap[c.Id] = c
	}

	var itemRecords []*core.Record
	if len(itemIds) > 0 {
		itemRecords, err = app.FindRecordsByIds(schema.CollectionItems, itemIds)
		if err != nil {
			return result.Err("Failed to fetch items"), err
		}
	}

	itemMap := make(map[string]*core.Record, len(itemRecords))
	for _, i := range itemRecords {
		itemMap[i.Id] = i
	}

	for _, rs := range rawStats {
		if rs.Type == "cell_visits" {
			record, ok := cellMap[rs.RecordID]
			if !ok {
				continue
			}
			stats.MostVisitedCells = append(stats.MostVisitedCells, EventStatEntry{
				Count:  rs.Count,
				Record: record,
			})
			continue
		}

		if rs.Type == "items" {
			record, ok := itemMap[rs.RecordID]
			if !ok {
				continue
			}
			stats.MostUsedItems = append(stats.MostUsedItems, EventStatEntry{
				Count:  rs.Count,
				Record: record,
			})
			continue
		}

		record, ok := playerMap[rs.RecordID]
		if !ok {
			continue
		}

		entry := EventStatEntry{
			Count:  rs.Count,
			Record: record,
		}

		switch rs.Type {
		case "games":
			stats.MostGamesCompleted = append(stats.MostGamesCompleted, entry)
		case "drops":
			stats.MostDrops = append(stats.MostDrops, entry)
		case "rerolls":
			stats.MostRerolls = append(stats.MostRerolls, entry)
		case "gyms":
			stats.MostGymsCompleted = append(stats.MostGymsCompleted, entry)
		case "movies":
			stats.MostMoviesWatched = append(stats.MostMoviesWatched, entry)
		case "karaoke":
			stats.MostKaraokeCompleted = append(stats.MostKaraokeCompleted, entry)
		case "roblox":
			stats.MostRobloxPlayed = append(stats.MostRobloxPlayed, entry)
		case "happywheels":
			stats.MostHappyWheelsPlayed = append(stats.MostHappyWheelsPlayed, entry)
		}
	}

	sortFn := func(a, b EventStatEntry) int {
		if a.Count > b.Count {
			return -1
		}
		if a.Count < b.Count {
			return 1
		}
		return 0
	}
	slices.SortFunc(stats.MostGamesCompleted, sortFn)
	slices.SortFunc(stats.MostDrops, sortFn)
	slices.SortFunc(stats.MostRerolls, sortFn)
	slices.SortFunc(stats.MostGymsCompleted, sortFn)
	slices.SortFunc(stats.MostMoviesWatched, sortFn)
	slices.SortFunc(stats.MostKaraokeCompleted, sortFn)
	slices.SortFunc(stats.MostWanted, sortFn)
	slices.SortFunc(stats.MostItemsUsed, sortFn)
	slices.SortFunc(stats.MostRobloxPlayed, sortFn)
	slices.SortFunc(stats.MostHappyWheelsPlayed, sortFn)
	slices.SortFunc(stats.MostUsedItems, sortFn)

	slices.SortFunc(stats.MostVisitedCells, sortFn)
	mostVisitedLen := len(stats.MostVisitedCells)
	if mostVisitedLen > 0 {
		limit := 6
		if mostVisitedLen < limit {
			limit = mostVisitedLen
		}

		least := slices.Clone(stats.MostVisitedCells[mostVisitedLen-limit:])
		slices.Reverse(least)
		stats.LeastVisitedCells = least
	}

	if len(stats.MostVisitedCells) > 6 {
		stats.MostVisitedCells = stats.MostVisitedCells[:6]
	}

	return result.Ok().WithData(stats), nil
}
