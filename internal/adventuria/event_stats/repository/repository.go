package repository

import (
	"adventuria/internal/adventuria/event_stats"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_stats/repository"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"
	"slices"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsData, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	stats := &event_stats.EventStatsData{
		MostGamesCompleted:    []event_stats.EventStatEntry{},
		MostDrops:             []event_stats.EventStatEntry{},
		MostRerolls:           []event_stats.EventStatEntry{},
		MostGymsCompleted:     []event_stats.EventStatEntry{},
		MostMoviesWatched:     []event_stats.EventStatEntry{},
		MostKaraokeCompleted:  []event_stats.EventStatEntry{},
		MostWanted:            []event_stats.EventStatEntry{},
		MostItemsUsed:         []event_stats.EventStatEntry{},
		MostRobloxPlayed:      []event_stats.EventStatEntry{},
		MostHappyWheelsPlayed: []event_stats.EventStatEntry{},
		MostVisitedCells:      []event_stats.EventStatEntry{},
		LeastVisitedCells:     []event_stats.EventStatEntry{},
		MostUsedItems:         []event_stats.EventStatEntry{},
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

	err := pb.DB().NewQuery(query).WithContext(ctx).All(&rawStats)
	if err != nil {
		return nil, err
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

	var playersRecords []*core.Record
	var playerStatsRecords []*core.Record
	if len(playerIds) > 0 {
		err = pb.RecordQuery(schema.CollectionPlayers).
			WithContext(ctx).
			Where(dbx.In(
				schema.PlayerSchema.Id,
				pbhelper.SliceToAny(playerIds)...,
			)).
			All(&playersRecords)
		if err != nil {
			return nil, err
		}

		err = pb.RecordQuery(schema.CollectionPlayerStats).
			Where(
				dbx.In(
					schema.PlayerStatsSchema.Player,
					pbhelper.SliceToAny(playerIds)...,
				),
			).
			AndWhere(dbx.HashExp{
				schema.PlayerStatsSchema.Season: seasonId,
			}).
			WithContext(ctx).
			All(&playerStatsRecords)
		if err != nil {
			return nil, err
		}
	}

	playerStatsMap := make(map[string]*model.PlayerStats, len(playerStatsRecords))
	for _, p := range playerStatsRecords {
		playerStatsMap[p.Id], err = repository.RecordToPlayerStats(p)
		if err != nil {
			return nil, err
		}
	}

	playerMap := make(map[string]*core.Record, len(playersRecords))
	for _, u := range playersRecords {
		playerMap[u.Id] = u

		playerStats, ok := playerStatsMap[u.Id]
		if !ok {
			continue
		}

		stats.MostWanted = append(stats.MostWanted, event_stats.EventStatEntry{
			Count:  playerStats.WasInJail(),
			Record: u,
		})
		stats.MostItemsUsed = append(stats.MostItemsUsed, event_stats.EventStatEntry{
			Count:  playerStats.ItemsUsed(),
			Record: u,
		})
	}

	var cellRecords []*core.Record
	if len(cellIds) > 0 {
		err = pb.RecordQuery(schema.CollectionCells).
			Where(
				dbx.In(
					schema.CellSchema.Id,
					pbhelper.SliceToAny(cellIds)...,
				),
			).
			WithContext(ctx).
			All(&cellRecords)
		if err != nil {
			return nil, err
		}
	}

	cellMap := make(map[string]*core.Record, len(cellRecords))
	for _, c := range cellRecords {
		cellMap[c.Id] = c
	}

	var itemRecords []*core.Record
	if len(itemIds) > 0 {
		err = pb.RecordQuery(schema.CollectionItems).
			Where(
				dbx.In(
					schema.ItemSchema.Id,
					pbhelper.SliceToAny(itemIds)...,
				),
			).
			WithContext(ctx).
			All(&itemRecords)
		if err != nil {
			return nil, err
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
			stats.MostVisitedCells = append(stats.MostVisitedCells, event_stats.EventStatEntry{
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
			stats.MostUsedItems = append(stats.MostUsedItems, event_stats.EventStatEntry{
				Count:  rs.Count,
				Record: record,
			})
			continue
		}

		record, ok := playerMap[rs.RecordID]
		if !ok {
			continue
		}

		entry := event_stats.EventStatEntry{
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

	sortFn := func(a, b event_stats.EventStatEntry) int {
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

	return stats, nil
}
