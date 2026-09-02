package coins_for_item

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
)

var _ model.WithView = (*CoinsForItem)(nil)

type itemView struct {
	Id             string `json:"id"`
	CollectionName string `json:"collectionName"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
}

func (c *CoinsForItem) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	deal, err := player.LastAction().State().Dealer.AsCoinsForItemDeal()
	if err != nil {
		return nil, err
	}

	item, err := c.items.GetByID(ctx, deal.ItemId)
	if err != nil {
		return nil, err
	}

	return struct {
		Item  itemView `json:"item"`
		Coins int      `json:"coins"`
	}{
		Item: itemView{
			Id:             item.ID(),
			CollectionName: schema.CollectionItems,
			Name:           item.Name(),
			Description:    item.Description(),
			Icon:           item.Icon(),
		},
		Coins: deal.Coins,
	}, nil
}
