package actions

import (
	"adventuria/internal/adventuria/model"
)

const (
	ActionTypeMove            model.ActionType = "move"
	ActionTypeRollDice        model.ActionType = "roll_dice"
	ActionTypeDone            model.ActionType = "done"
	ActionTypeReroll          model.ActionType = "reroll"
	ActionTypeDrop            model.ActionType = "drop"
	ActionTypeGenerateWheel   model.ActionType = "generate_wheel"
	ActionTypeRollWheel       model.ActionType = "roll_wheel"
	ActionTypeRollItem        model.ActionType = "roll_item"
	ActionTypeRollItemOnCell  model.ActionType = "roll_item_on_cell"
	ActionTypeTeleport        model.ActionType = "teleport"
	ActionTypeBuy             model.ActionType = "buy"
	ActionTypeRefreshShop     model.ActionType = "refresh_shop"
	ActionTypeUpdateReview    model.ActionType = "update_review"
	ActionTypeNeedToRollWheel model.ActionType = "need_to_roll_wheel"
)
