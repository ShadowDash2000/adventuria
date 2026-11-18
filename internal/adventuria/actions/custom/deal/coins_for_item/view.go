package coins_for_item

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/player_actions/custom/deals/coins_for_item"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"context"
)

var _ actions.View = (*CoinsForItem)(nil)

type itemView struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (c *CoinsForItem) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	dealEntry, ok := actor.Player().Actions().Get[*coins_for_item.CoinsForItem](coins_for_item.PlayerActionType)
	if !ok {
		return nil, errs.ErrNoActiveCoinsForItemDeal
	}

	dealAction := dealEntry.Action()

	item, err := c.items.GetByID(ctx, dealAction.ItemID())
	if err != nil {
		return nil, err
	}

	return struct {
		Item  itemView `json:"item"`
		Coins int      `json:"coins"`
	}{
		Item: itemView{
			Id:          item.ID(),
			Name:        item.Name(),
			Description: item.Description(),
			Icon:        pbhelper.BuildFilePath(schema.CollectionItems, item.ID(), item.Icon()),
		},
		Coins: dealAction.Coins(),
	}, nil
}
