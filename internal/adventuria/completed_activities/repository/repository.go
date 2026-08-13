package repository

import (
	"adventuria/internal/adventuria/completed_activities"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
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

func (r *Repository) GetActivitiesByCellIDWithStatus(
	ctx context.Context,
	cellId string,
	statuses []model.ActionStatus,
	limit int,
) ([]*completed_activities.CompletedActivity, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	actionsTable := schema.CollectionActions
	playersTable := schema.CollectionPlayers
	activitiesTable := schema.CollectionActivities
	actionActivity := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Activity)
	actionStatus := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Status)
	actionCell := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Cell)
	actionPlayer := pbhelper.DotExpand(actionsTable, schema.ActionSchema.Player)
	actionCreated := pbhelper.DotExpand(actionsTable, "created")
	playerId := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Id)
	playerName := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Name)
	playerAvatar := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Avatar)
	playerColor := pbhelper.DotExpand(playersTable, schema.PlayerSchema.Color)
	activityId := pbhelper.DotExpand(activitiesTable, schema.ActivitySchema.Id)
	activityName := pbhelper.DotExpand(activitiesTable, schema.ActivitySchema.Name)
	activityCover := pbhelper.DotExpand(activitiesTable, schema.ActivitySchema.Cover)
	activityCoverAlt := pbhelper.DotExpand(activitiesTable, schema.ActivitySchema.CoverAlt)

	selectedActivitiesQuery :=
		pb.DB().
			Select(
				actionActivity,
				fmt.Sprintf("MAX(%s) AS last_action_at", actionCreated),
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
			fmt.Sprintf("%s AS activity_id", actionActivity),
			fmt.Sprintf("%s AS activity_name", activityName),
			fmt.Sprintf("%s AS activity_cover", activityCover),
			fmt.Sprintf("%s AS activity_cover_alt", activityCoverAlt),
			actionStatus,
			fmt.Sprintf("%s AS player_id", playerId),
			fmt.Sprintf("%s AS player_name", playerName),
			fmt.Sprintf("%s AS player_avatar", playerAvatar),
			fmt.Sprintf("%s AS player_color", playerColor),
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
			dbx.NewExp(pbhelper.Eq(playerId, actionPlayer)),
		).
		InnerJoin(
			activitiesTable,
			dbx.NewExp(pbhelper.Eq(activityId, actionActivity)),
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

	res := make([]*completed_activities.CompletedActivity, 0)
	mapById := make(map[string]*completed_activities.CompletedActivity)
	for _, row := range playerActivityRows {
		activity, ok := mapById[row.ActivityId]
		if !ok {
			activity = playerActivityRowToCompletedActivity(row)
			mapById[row.ActivityId] = activity
			res = append(res, activity)
		}

		activity.Players = append(
			activity.Players,
			playerActivityRowToCompletedActivityPlayerStatus(row),
		)
	}

	return res, nil
}
