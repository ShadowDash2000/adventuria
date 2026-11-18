package custom

import (
	"adventuria/internal/adventuria/player_actions"
	"adventuria/internal/adventuria/player_actions/custom/deals/coins_for_item"
	"adventuria/internal/adventuria/player_actions/custom/shop"
)

func RegisterPlayerActions() {
	player_actions.Register(
		shop.NewDef(),
		coins_for_item.NewDef(),
	)
}
