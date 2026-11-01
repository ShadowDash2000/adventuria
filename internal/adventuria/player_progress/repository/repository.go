package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_progress"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayersProgress)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerProgressToRecord(progress, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerProgress(record), nil
}

func (r *Repository) Update(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlayersProgress, progress.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProgressNotFound
		}
		return nil, err
	}

	PlayerProgressToRecord(progress, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerProgress(record), nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayersProgress).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.PlayerProgressSchema.Id: id,
		}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProgressNotFound
		}
		return nil, err
	}

	return RecordToPlayerProgress(&record), nil
}

func (r *Repository) GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayersProgress).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.PlayerProgressSchema.Player: playerId},
			dbx.HashExp{schema.PlayerProgressSchema.Season: seasonId},
		)).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrProgressNotFound
		}
		return nil, err
	}

	return RecordToPlayerProgress(&record), nil
}

func (r *Repository) GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.PlayerProgress, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionPlayersProgress).
		WithContext(ctx).
		Where(dbx.And(
			dbx.HashExp{schema.PlayerProgressSchema.Season: seasonId},
		)).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToPlayerProgresses(records), nil
}

func (r *Repository) ChangeBalance(ctx context.Context, id string, amount int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	_, err := pb.DB().
		NewQuery(fmt.Sprintf(
			"UPDATE %s SET %s = %s + {:amount} WHERE %s = {:id}",
			schema.CollectionPlayersProgress,
			schema.PlayerProgressSchema.Balance,
			schema.PlayerProgressSchema.Balance,
			schema.PlayerProgressSchema.Id,
		)).
		Bind(dbx.Params{
			"amount": amount,
			"id":     id,
		}).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) ChangeEnergy(ctx context.Context, id string, amount int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	_, err := pb.DB().
		NewQuery(fmt.Sprintf(
			"UPDATE %s SET %s = %s + {:amount} WHERE %s = {:id}",
			schema.CollectionPlayersProgress,
			schema.PlayerProgressSchema.Energy,
			schema.PlayerProgressSchema.Energy,
			schema.PlayerProgressSchema.Id,
		)).
		Bind(dbx.Params{
			"amount": amount,
			"id":     id,
		}).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetActivePlayersBySeasonID(ctx context.Context, seasonId string) ([]*player_progress.ActivePlayer, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	playersTable := schema.CollectionPlayers
	progressTable := schema.CollectionPlayersProgress
	playerId := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Id)
	playerName := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Name)
	playerAvatar := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Avatar)
	playerColor := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Color)
	playerIsStreamLive := pbhelper.DotExpand(playersTable, schema.PlayerSchema.IsStreamLive)
	progressId := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.Id)
	progressUpdated := pbhelper.DotExpand(progressTable, "updated")
	progressPlayer := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.Player)
	progressCurrentWorld := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.CurrentWorld)
	progressCellsPassed := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.CellsPassed)
	progressSeason := pbhelper.DotExpand(progressTable, schema.PlayerProgressSchema.Season)

	var activePlayersRows []*activePlayerRow
	err := pb.DB().
		Select(
			fmt.Sprintf("%s AS player_id", playerId),
			fmt.Sprintf("%s AS name", playerName),
			fmt.Sprintf("%s AS avatar", playerAvatar),
			fmt.Sprintf("%s AS color", playerColor),
			fmt.Sprintf("%s AS is_stream_live", playerIsStreamLive),
			fmt.Sprintf("%s AS progress_id", progressId),
			fmt.Sprintf("%s AS progress_updated", progressUpdated),
			fmt.Sprintf("%s AS current_world", progressCurrentWorld),
			fmt.Sprintf("%s AS cells_passed", progressCellsPassed),
		).
		From(progressTable).
		InnerJoin(
			playersTable,
			dbx.NewExp(pbhelper.Eq(playerId, progressPlayer)),
		).
		Where(dbx.HashExp{
			progressSeason: seasonId,
		}).
		WithContext(ctx).
		All(&activePlayersRows)
	if err != nil {
		return nil, err
	}

	res := make([]*player_progress.ActivePlayer, len(activePlayersRows))
	for i, row := range activePlayersRows {
		res[i] = row.toActivePlayer()
	}

	return res, nil
}
