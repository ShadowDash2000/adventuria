package custom

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/actions/custom/buy"
	"adventuria/internal/adventuria_new/actions/custom/done"
	"adventuria/internal/adventuria_new/actions/custom/drop"
	"adventuria/internal/adventuria_new/actions/custom/reroll"
	"adventuria/internal/adventuria_new/actions/custom/roll_dice"
	"adventuria/internal/adventuria_new/actions/custom/roll_item"
	"adventuria/internal/adventuria_new/actions/custom/roll_item_on_cell"
	"adventuria/internal/adventuria_new/actions/custom/roll_wheel"
	rollWheelRepo "adventuria/internal/adventuria_new/actions/custom/roll_wheel/repository"
	"adventuria/internal/adventuria_new/board"
	"adventuria/internal/adventuria_new/cells"
	"adventuria/internal/adventuria_new/inventories"
	"adventuria/internal/adventuria_new/items"
	"adventuria/internal/adventuria_new/players"
	"adventuria/internal/adventuria_new/reviews"
	"adventuria/internal/adventuria_new/settings"
)

func RegisterActions(
	cells *cells.Cells,
	reviews *reviews.Reviews,
	players *players.Players,
	settings *settings.Settings,
	board *board.Board,
	actionsService *actions.Actions,
	items *items.Items,
	inventories *inventories.Inventories,
	rollWheelRepo *rollWheelRepo.Repository,
) {
	actions.Register(
		done.NewActionDoneDef(cells, reviews),
		drop.NewActionDropDef(cells, reviews, players, settings, board),
		reroll.NewActionRerollDef(cells, reviews, actionsService),
		buy.NewActionBuyDef(cells, items, inventories),
		roll_dice.NewActionRollDiceDef(cells, actionsService, board),
		roll_item.NewActionRollItemDef(actionsService, inventories, items),
		roll_item_on_cell.NewActionRollItemOnCellDef(cells, inventories, items),
		roll_wheel.NewActionRollWheelDef(cells, rollWheelRepo),
	)
}
