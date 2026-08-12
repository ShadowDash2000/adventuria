package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/players"
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

func (r *Repository) Exists(ctx context.Context, id string) (bool, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record struct {
		Id string `db:"id"`
	}
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Select(schema.PlayerSchema.Id).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*model.PlayerInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		All(&records)
	if err != nil {
		return nil, err
	}

	return RecordsToPlayerInfos(records), nil
}

func (r *Repository) NotifyChange(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var record core.Record
	err := pb.RecordQuery(schema.CollectionPlayers).
		WithContext(ctx).
		Where(dbx.HashExp{schema.PlayerSchema.Id: id}).
		Limit(1).
		One(&record)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrPlayerNotFound
		}
		return err
	}

	event := &core.ModelEvent{
		App:     pb,
		Context: ctx,
		Type:    core.ModelEventTypeUpdate,
	}
	event.Model = &record

	return pb.OnModelAfterUpdateSuccess().Trigger(event)
}

func (r *Repository) GetActivitiesByCellIDWithStatus(
	ctx context.Context,
	cellId string,
	statuses []model.ActionStatus,
	limit int,
) (map[string][]*players.CompletedActivity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	actionsTable := schema.CollectionActions
	playersTable := schema.CollectionPlayers
	actionActivity := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Activity)
	actionStatus := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Status)
	actionCell := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Cell)
	actionPlayer := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Player)
	actionCreated := pbhelper.DotExpand(actionsTable, "created")
	playerId := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Id)
	playerName := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Name)
	playerAvatar := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Avatar)
	playerColor := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Color)

	selectedActivitiesQuery :=
		pb.DB().
			Select(
				actionActivity,
				fmt.Sprintf(
					"MAX(%s) AS last_action_at",
					actionCreated,
				),
			).
			From(actionsTable).
			Where(dbx.And(
				dbx.HashExp{
					actionCell: cellId,
				},
				dbx.In(
					actionStatus,
					pbhelper.SliceToAny(statuses)...,
				),
				dbx.Not(dbx.HashExp{
					actionActivity: "",
				}),
			)).
			GroupBy(actionActivity).
			OrderBy(
				"last_action_at DESC",
				actionActivity,
			).
			Limit(int64(limit)).
			Build()

	var playerActivityRows []*playerActivityRow
	err := pb.DB().
		Select(
			fmt.Sprintf(
				"%s AS activity_id",
				actionActivity,
			),
			actionStatus,
			fmt.Sprintf(
				"%s AS player_id",
				playerId,
			),
			playerName,
			playerAvatar,
			playerColor,
		).
		From(fmt.Sprintf(
			"(%s) AS selected_activities",
			selectedActivitiesQuery.SQL(),
		)).
		InnerJoin(
			actionsTable,
			dbx.And(
				dbx.NewExp(pbhelper.Eq(
					actionActivity,
					fmt.Sprintf("selected_activities.%s", schema.ActionSchema.Activity),
				)),
				dbx.HashExp{
					actionCell: cellId,
				},
				dbx.In(
					actionStatus,
					pbhelper.SliceToAny(statuses)...,
				),
			),
		).
		InnerJoin(
			playersTable,
			dbx.NewExp(pbhelper.Eq(
				playerId,
				actionPlayer,
			)),
		).
		OrderBy(
			"selected_activities.last_action_at DESC",
			fmt.Sprintf("%s DESC", actionCreated),
		).
		Bind(selectedActivitiesQuery.Params()).
		WithContext(ctx).
		All(&playerActivityRows)
	if err != nil {
		return nil, err
	}

	res := make(map[string][]*players.CompletedActivity)
	for _, playerActivityRow := range playerActivityRows {
		res[playerActivityRow.ActivityId] = append(
			res[playerActivityRow.ActivityId],
			playerActivityRowToCompletedActivity(playerActivityRow),
		)
	}

	return res, nil
}
