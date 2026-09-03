package actions

import (
	"adventuria/internal/adventuria/model"
)

const (
	ActionTypeStart            model.ActionType = "start"
	ActionTypeRollDice         model.ActionType = "roll_dice"
	ActionTypeCompleteActivity model.ActionType = "complete_activity"
	ActionTypeDone             model.ActionType = "done"
	ActionTypeReroll           model.ActionType = "reroll"
	ActionTypeDrop             model.ActionType = "drop"
	ActionTypeGenerateWheel    model.ActionType = "generate_wheel"
	ActionTypeRollWheel        model.ActionType = "roll_wheel"
	ActionTypeRollItem         model.ActionType = "roll_item"
	ActionTypeRollItemOnCell   model.ActionType = "roll_item_on_cell"
	ActionTypeBuy              model.ActionType = "buy"
	ActionTypeRefreshShop      model.ActionType = "refresh_shop"
	ActionTypeCoinsForItem     model.ActionType = "coins_for_item"

	ActionTypeMoveToCellId model.ActionType = "move_to_cell_id"
)
