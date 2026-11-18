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
	cellsService "adventuria/internal/adventuria/cells/service"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions"
)

func RegisterCells(
	activityFilters *activity_filters.ActivityFilters,
	items *items.Items,
	cellsService *cellsService.Cells,
	actions *actions.Actions,
	board *board.Board,
	playerActions *player_actions.PlayerActions,
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
		casino.NewDef(
			items, playerActions,
			"shop",
		),
		shop.NewDef(
			cells.CellTypeShop, model.ItemTypeBuff,
			items, playerActions,
			"shop",
		),
		teleport.NewDef(cellsService, board, actions),
	)
}
