package custom

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/activity_filters"
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/cells/custom/activity"
	"adventuria/internal/adventuria/cells/custom/casino"
	"adventuria/internal/adventuria/cells/custom/jail"
	"adventuria/internal/adventuria/cells/custom/roll_item"
	"adventuria/internal/adventuria/cells/custom/shop"
	"adventuria/internal/adventuria/cells/custom/start"
	"adventuria/internal/adventuria/cells/custom/teleport"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/model"
)

func RegisterCells(
	activityFilters *activity_filters.ActivityFilters,
	items *items.Items,
	cellsService *cells.Cells,
	actions *actions.Actions,
	board *board.Board,
) {
	cells.Register(
		start.NewDef(),
		activity.NewDef(
			model.ActivityTypeGame, activityFilters,
			"wheel", "activity", "game",
		),
		activity.NewDef(
			model.ActivityTypeMovie, activityFilters,
			"wheel", "activity",
		),
		activity.NewDef(
			model.ActivityTypeGym, activityFilters,
			"wheel", "activity",
		),
		activity.NewDef(
			model.ActivityTypeKaraoke, activityFilters,
			"wheel", "activity",
		),
		jail.NewDef(
			activityFilters,
			"wheel", "activity", "game",
		),
		roll_item.NewDef(items, "wheel"),
		casino.NewDef(items, "shop"),
		shop.NewDef(
			cells.CellTypeShop, model.ItemTypeBuff, items,
			"shop",
		),
		teleport.NewDef(cellsService, board, actions),
	)
}
