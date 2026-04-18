package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/result"
	"fmt"
	"slices"

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
		schema.ActionSchema.User,      // 1
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

	userIds := make([]string, 0, len(rawStats))
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
			if !slices.Contains(userIds, rs.RecordID) {
				userIds = append(userIds, rs.RecordID)
			}
		}
	}

	var users []*core.Record
	if len(userIds) > 0 {
		users, err = app.FindRecordsByIds(schema.CollectionUsers, userIds)
		if err != nil {
			return result.Err("Failed to fetch users"), err
		}
	}

	userMap := make(map[string]*core.Record, len(users))
	for _, u := range users {
		userMap[u.Id] = u

		var userStats Stats
		err = u.UnmarshalJSONField(schema.UserSchema.Stats, &userStats)
		if err != nil {
			return result.Err("Failed to unmarshal user stats"), err
		}

		stats.MostWanted = append(stats.MostWanted, EventStatEntry{
			Count:  userStats.WasInJail,
			Record: u,
		})
		stats.MostItemsUsed = append(stats.MostItemsUsed, EventStatEntry{
			Count:  userStats.ItemsUsed,
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

		record, ok := userMap[rs.RecordID]
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
