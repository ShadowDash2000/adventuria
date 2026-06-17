package custom

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/activities"
	"adventuria/internal/adventuria_new/activity_filters"
	"adventuria/internal/adventuria_new/board"
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/cells/custom/activity"
	"adventuria/internal/adventuria_new/cells/custom/casino"
	"adventuria/internal/adventuria_new/cells/custom/jail"
	"adventuria/internal/adventuria_new/cells/custom/roll_item"
	"adventuria/internal/adventuria_new/cells/custom/shop"
	"adventuria/internal/adventuria_new/cells/custom/start"
	"adventuria/internal/adventuria_new/cells/custom/teleport"
	"adventuria/internal/adventuria_new/items"
	"adventuria/internal/adventuria_new/model"
)

func RegisterCells(
	activities *activities.Activities,
	activityFilters *activity_filters.ActivityFilters,
	items *items.Items,
	cellsService *cells.Cells,
	actions *actions.Actions,
	board *board.Board,
) {
	cells.Register(
		start.NewDef(),
		activity.NewDef(
			model.ActivityTypeGame, activities, activityFilters,
			"wheel", "activity", "game",
		),
		activity.NewDef(
			model.ActivityTypeMovie, activities, activityFilters,
			"wheel", "activity",
		),
		activity.NewDef(
			model.ActivityTypeGym, activities, activityFilters,
			"wheel", "activity",
		),
		activity.NewDef(
			model.ActivityTypeKaraoke, activities, activityFilters,
			"wheel", "activity",
		),
		jail.NewDef(
			activities, activityFilters,
			"wheel", "activity", "game",
		),
		roll_item.NewDef(items, "wheel"),
		casino.NewDef(items, "shop"),
		shop.NewDef(items, "shop"),
		teleport.NewDef(cellsService, board, actions),
	)
}
