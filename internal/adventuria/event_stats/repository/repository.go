package repository

import (
	"adventuria/internal/adventuria/event_stats"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type seasons interface {
	GetByID(ctx context.Context, id string) (*model.Season, error)
}

type playerStats interface {
	GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.PlayerStats, error)
}

type Repository struct {
	pb          core.App
	seasons     seasons
	playerStats playerStats
}

func NewRepository(pb core.App, seasons seasons, playerStats playerStats) *Repository {
	return &Repository{
		pb:          pb,
		seasons:     seasons,
		playerStats: playerStats,
	}
}

func (r *Repository) ComputeStats(ctx context.Context, seasonId string) (*event_stats.EventStatsEntries, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	season, err := r.seasons.GetByID(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	playersStats, err := r.playerStats.GetAllBySeasonID(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	playerIds := make([]any, len(playersStats))
	for i, playerStat := range playersStats {
		playerIds[i] = playerStat.Player()
	}

	var playerRecords []*core.Record
	err = pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Where(dbx.In(schema.PlayerSchema.Id, playerIds...)).
		All(&playerRecords)
	if err != nil {
		return nil, err
	}

	playerRecordsMap := make(map[string]*core.Record, len(playerRecords))
	for _, playerRecord := range playerRecords {
		playerRecordsMap[playerRecord.Id] = playerRecord
	}

	stats := &event_stats.EventStatsEntries{
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

	for _, playerStat := range playersStats {
		player, ok := playerRecordsMap[playerStat.Player()]
		if !ok {
			continue
		}

		stats.MostGamesCompleted = append(stats.MostGamesCompleted, event_stats.EventStatEntry{
			Count:  playerStat.ActivitiesStats().GamesCompleted,
			Record: player,
		})
		stats.MostGymsCompleted = append(stats.MostGymsCompleted, event_stats.EventStatEntry{
			Count:  playerStat.ActivitiesStats().GymsCompleted,
			Record: player,
		})
		stats.MostMoviesWatched = append(stats.MostMoviesWatched, event_stats.EventStatEntry{
			Count:  playerStat.ActivitiesStats().MoviesCompleted,
			Record: player,
		})
		stats.MostKaraokeCompleted = append(stats.MostKaraokeCompleted, event_stats.EventStatEntry{
			Count:  playerStat.ActivitiesStats().KaraokeCompleted,
			Record: player,
		})

		stats.MostDrops = append(stats.MostDrops, event_stats.EventStatEntry{
			Count:  playerStat.Drops(),
			Record: player,
		})
		stats.MostRerolls = append(stats.MostRerolls, event_stats.EventStatEntry{
			Count:  playerStat.Rerolls(),
			Record: player,
		})
		stats.MostWanted = append(stats.MostWanted, event_stats.EventStatEntry{
			Count:  playerStat.WasInJail(),
			Record: player,
		})
		stats.MostItemsUsed = append(stats.MostItemsUsed, event_stats.EventStatEntry{
			Count:  playerStat.ItemsUsed(),
			Record: player,
		})
	}

	cellsVisitsStats, err := r.getCellsVisitsStats(ctx, season.SeasonDateStart())
	if err != nil {
		return nil, err
	}

	stats.MostVisitedCells = cellsVisitsStats

	usedItemsStats, err := r.getUsedItemsStats(ctx, season.SeasonDateStart())
	if err != nil {
		return nil, err
	}

	stats.MostUsedItems = usedItemsStats

	sortStatsDesc(stats.MostGamesCompleted)
	sortStatsDesc(stats.MostDrops)
	sortStatsDesc(stats.MostRerolls)
	sortStatsDesc(stats.MostGymsCompleted)
	sortStatsDesc(stats.MostMoviesWatched)
	sortStatsDesc(stats.MostKaraokeCompleted)
	sortStatsDesc(stats.MostWanted)
	sortStatsDesc(stats.MostItemsUsed)
	sortStatsDesc(stats.MostVisitedCells)
	sortStatsDesc(stats.MostUsedItems)

	cellsLimit := 6
	stats.LeastVisitedCells = getReversedStatsTail(stats.MostVisitedCells, cellsLimit)

	if len(stats.MostVisitedCells) > cellsLimit {
		stats.MostVisitedCells = stats.MostVisitedCells[:cellsLimit]
	}

	return stats, nil
}

func sortStatsDesc(stats []event_stats.EventStatEntry) {
	sortFn := func(a, b event_stats.EventStatEntry) int {
		if a.Count > b.Count {
			return -1
		}
		if a.Count < b.Count {
			return 1
		}
		return 0
	}

	slices.SortFunc(stats, sortFn)
}

func getReversedStatsTail(stats []event_stats.EventStatEntry, limit int) []event_stats.EventStatEntry {
	statsLen := len(stats)
	if statsLen == 0 {
		return nil
	}

	limit = min(statsLen, limit)

	res := slices.Clone(stats[statsLen-limit:])
	slices.Reverse(res)

	return res
}

func (r *Repository) getUsedItemsStats(ctx context.Context, timeFrom time.Time) ([]event_stats.EventStatEntry, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	dateFrom, err := types.ParseDateTime(timeFrom)
	if err != nil {
		return nil, err
	}

	type itemStatEntry struct {
		ItemId string `db:"item_id"`
		Count  int    `db:"count"`
	}

	var itemsStats []itemStatEntry
	err = pb.DB().NewQuery(fmt.Sprintf(`
		SELECT items.value as "item_id", COUNT(*) as "count"
		FROM %[1]s as actions
		CROSS JOIN json_each(actions.%[2]s) as items
		WHERE actions.%[2]s IS NOT NULL
			AND actions.%[2]s != '[]'
			AND actions.%[2]s != 'null'
			AND actions.created > {:date}
		GROUP BY items.value
		`,
		schema.CollectionActions,
		schema.ActionSchema.UsedItems,
	)).Bind(dbx.Params{
		"date": dateFrom,
	}).WithContext(ctx).All(&itemsStats)
	if err != nil {
		return nil, err
	}

	var itemIds []any
	itemsMap := make(map[string]itemStatEntry, len(itemsStats))
	for _, itemStat := range itemsStats {
		itemsMap[itemStat.ItemId] = itemStat
		itemIds = append(itemIds, itemStat.ItemId)
	}

	var records []*core.Record
	err = pb.RecordQuery(schema.CollectionItems).
		Where(dbx.Not(dbx.HashExp{
			schema.ItemSchema.Type: model.ItemTypeDev,
		})).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	res := make([]event_stats.EventStatEntry, len(records))
	for i, record := range records {
		res[i] = event_stats.EventStatEntry{
			Count:  itemsMap[record.Id].Count,
			Record: record,
		}
	}

	return res, nil
}

func (r *Repository) getCellsVisitsStats(ctx context.Context, timeFrom time.Time) ([]event_stats.EventStatEntry, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	dateFrom, err := types.ParseDateTime(timeFrom)
	if err != nil {
		return nil, err
	}

	type cellVisitStatEntry struct {
		CellId string `db:"cell_id"`
		Count  int    `db:"count"`
	}

	var cellVisitsStats []cellVisitStatEntry
	err = pb.DB().NewQuery(fmt.Sprintf(`
		SELECT actions.%[2]s as "cell_id", COUNT(*) as "count"
		FROM %[1]s as actions
		WHERE actions.%[3]s IN (%[4]s)
			AND actions.created > {:date}
		GROUP BY actions.%[2]s
		`,
		schema.CollectionActions,
		schema.ActionSchema.Cell,
		schema.ActionSchema.Status,
		pbhelper.SliceToSqlString([]model.ActionStatus{
			model.ActionStatusDone,
			model.ActionStatusDrop,
			model.ActionStatusRollDice,
			model.ActionStatusRollWheel,
			model.ActionStatusMove,
			model.ActionStatusRollItemOnCell,
		}),
	)).Bind(dbx.Params{
		"date": dateFrom,
	}).WithContext(ctx).All(&cellVisitsStats)
	if err != nil {
		return nil, err
	}

	var cellIds []any
	cellsMap := make(map[string]cellVisitStatEntry, len(cellVisitsStats))
	for _, cellVisitStat := range cellVisitsStats {
		cellsMap[cellVisitStat.CellId] = cellVisitStat
		cellIds = append(cellIds, cellVisitStat.CellId)
	}

	var records []*core.Record
	err = pb.RecordQuery(schema.CollectionCells).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	res := make([]event_stats.EventStatEntry, len(records))
	for i, record := range records {
		res[i] = event_stats.EventStatEntry{
			Count:  cellsMap[record.Id].Count,
			Record: record,
		}
	}

	return res, nil
}
