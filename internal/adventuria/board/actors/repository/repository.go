package repository

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) GetActiveActorsBySeasonID(ctx context.Context, seasonId string) ([]*actors.ActiveActor, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	playersTable := schema.CollectionPlayers
	progressTable := schema.CollectionPlayersProgress
	npcTable := schema.CollectionNPCs

	playerId := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Id)
	playerName := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Name)
	playerAvatar := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Avatar)
	playerColor := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Color)

	progressUpdated := pbhelper.DotExpand(progressTable, "updated")
	progressPlayer := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.Player)
	progressCurrentWorld := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.CurrentWorld)
	progressCellsPassed := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.CellsPassed)
	progressSeason := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.Season)

	npcId := pbhelper.DotExpand(npcTable, schema.NPCsSchema.Id)
	npcUpdated := pbhelper.DotExpand(npcTable, "updated")
	npcName := pbhelper.DotExpand(npcTable, schema.NPCsSchema.Name)
	npcAvatar := pbhelper.DotExpand(npcTable, schema.NPCsSchema.Avatar)
	npcColor := pbhelper.DotExpand(npcTable, schema.NPCsSchema.Color)
	npcCurrentWorld := pbhelper.DotExpand(npcTable, schema.NPCsSchema.CurrentWorld)
	npcCellsPassed := pbhelper.DotExpand(npcTable, schema.NPCsSchema.CellsPassed)
	npcNextShiftChangeAt := pbhelper.DotExpand(npcTable, schema.NPCsSchema.NextShiftChangeAt)

	npcQuery := "SELECT " +
		fmt.Sprintf("%s AS id, ", npcId) +
		fmt.Sprintf("'%s' as collection, ", npcTable) +
		fmt.Sprintf("'%s' as kind, ", actors.ActorKindNPC) +
		fmt.Sprintf("%s AS updated, ", npcUpdated) +
		fmt.Sprintf("%s AS name, ", npcName) +
		fmt.Sprintf("%s AS avatar, ", npcAvatar) +
		fmt.Sprintf("%s AS color, ", npcColor) +
		fmt.Sprintf("%s AS current_world, ", npcCurrentWorld) +
		fmt.Sprintf("%s AS cells_passed, ", npcCellsPassed) +
		fmt.Sprintf("%s AS next_shift_change_at ", npcNextShiftChangeAt) +
		fmt.Sprintf("FROM %s", npcTable)

	playersQuery := "SELECT " +
		fmt.Sprintf("%s AS id, ", playerId) +
		fmt.Sprintf("'%s' as collection, ", playersTable) +
		fmt.Sprintf("'%s' as kind, ", actors.ActorKindPlayer) +
		fmt.Sprintf("%s AS updated, ", progressUpdated) +
		fmt.Sprintf("%s AS name, ", playerName) +
		fmt.Sprintf("%s AS avatar, ", playerAvatar) +
		fmt.Sprintf("%s AS color, ", playerColor) +
		fmt.Sprintf("%s AS current_world, ", progressCurrentWorld) +
		fmt.Sprintf("%s AS cells_passed, ", progressCellsPassed) +
		"NULL AS next_shift_change_at " +
		fmt.Sprintf("FROM %s ", playersTable) +
		fmt.Sprintf("INNER JOIN %s ON %s = %s ", progressTable, playerId, progressPlayer) +
		fmt.Sprintf("WHERE %s = {:seasonId}", progressSeason)

	query := playersQuery + " UNION ALL " + npcQuery

	var activeActorsRows []*activeActorRow
	err := pb.DB().
		NewQuery(query).
		Bind(dbx.Params{
			"seasonId": seasonId,
		}).
		WithContext(ctx).
		All(&activeActorsRows)
	if err != nil {
		return nil, err
	}

	res := make([]*actors.ActiveActor, len(activeActorsRows))
	for i, row := range activeActorsRows {
		res[i] = row.toActiveActor()
	}

	return res, nil
}
