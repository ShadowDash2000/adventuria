package actions

import (
	"adventuria/internal/adventuria_new/model"
)

const (
	ActionTypeRollDice       model.ActionType = "roll_dice"
	ActionTypeDone           model.ActionType = "done"
	ActionTypeReroll         model.ActionType = "reroll"
	ActionTypeDrop           model.ActionType = "drop"
	ActionTypeRollWheel      model.ActionType = "roll_wheel"
	ActionTypeRollItem       model.ActionType = "roll_item"
	ActionTypeRollItemOnCell model.ActionType = "roll_item_on_cell"
	ActionTypeTeleport       model.ActionType = "teleport"
	ActionTypeBuy            model.ActionType = "buy"
	ActionTypeRefreshShop    model.ActionType = "refresh_shop"
	ActionTypeUpdateReview   model.ActionType = "update_review"
)
