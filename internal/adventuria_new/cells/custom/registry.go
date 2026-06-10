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
		start.NewCellStartDef(),
		activity.NewCellActivityDef(
			model.ActivityTypeGame,
			activities,
			activityFilters,
			"wheel", "activity", "game",
		),
		activity.NewCellActivityDef(
			model.ActivityTypeMovie,
			activities,
			activityFilters,
			"wheel", "activity",
		),
		activity.NewCellActivityDef(
			model.ActivityTypeGym,
			activities,
			activityFilters,
			"wheel", "activity",
		),
		activity.NewCellActivityDef(
			model.ActivityTypeKaraoke,
			activities,
			activityFilters,
			"wheel", "activity",
		),
		jail.NewCellJailDef(
			activities,
			activityFilters,
			"wheel", "activity", "game",
		),
		roll_item.NewCellRollItemDef(
			items,
			"wheel",
		),
		casino.NewCellCasinoDef(
			items,
			"shop",
		),
		shop.NewCellShopDef(
			items,
			"shop",
		),
		teleport.NewCellTeleportDef(
			cellsService,
			board,
			actions,
		),
	)
}
