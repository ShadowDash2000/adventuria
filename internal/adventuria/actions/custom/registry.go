package custom

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/actions/custom/buy"
	"adventuria/internal/adventuria/actions/custom/done"
	"adventuria/internal/adventuria/actions/custom/drop"
	"adventuria/internal/adventuria/actions/custom/generate_wheel"
	"adventuria/internal/adventuria/actions/custom/refresh_shop"
	"adventuria/internal/adventuria/actions/custom/reroll"
	"adventuria/internal/adventuria/actions/custom/roll_dice"
	"adventuria/internal/adventuria/actions/custom/roll_item"
	"adventuria/internal/adventuria/actions/custom/roll_item_on_cell"
	"adventuria/internal/adventuria/actions/custom/roll_wheel"
	rollWheelRepo "adventuria/internal/adventuria/actions/custom/roll_wheel/repository"
	"adventuria/internal/adventuria/actions/custom/update_review"
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/inventories"
	"adventuria/internal/adventuria/items"
	"adventuria/internal/adventuria/players"
	"adventuria/internal/adventuria/reviews"
	"adventuria/internal/adventuria/settings"
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
		done.NewDef(cells, reviews),
		drop.NewDef(cells, reviews, players, settings, board),
		reroll.NewDef(cells, reviews, actionsService),
		buy.NewDef(cells, items, inventories),
		refresh_shop.NewActionRefreshShopDef(cells),
		roll_dice.NewDef(cells, actionsService, board),
		roll_item.NewDef(actionsService, inventories, items),
		roll_item_on_cell.NewDef(cells, inventories, items),
		roll_wheel.NewDef(cells, rollWheelRepo),
		update_review.NewDef(reviews),
		generate_wheel.NewDef(cells, actionsService),
	)
}
