package buy

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"adventuria/internal/adventuria/schema"
	"context"
)

var _ actions.View = (*Buy)(nil)

type itemView struct {
	Id             string `json:"id"`
	CollectionName string `json:"collectionName"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
	Price          int    `json:"price"`
}

func (b *Buy) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	shopEntry, ok := actor.Player().Actions().Get[*shop.Shop](shop.PlayerActionType)
	if !ok {
		return nil, errs.ErrNoActiveShop
	}

	shopAction := shopEntry.Action()

	shopIds := shopAction.IDs()
	items, err := b.items.GetByIDs(ctx, shopIds)
	if err != nil {
		return nil, err
	}

	itemsViewMap := make(map[string]*itemView, len(items))
	for _, item := range items {
		buyResult, err := calculatePrice(ctx, actor.Events(), item, shopAction, model.EffectRunModePreview)
		if err != nil {
			return nil, err
		}

		itemView := itemToItemView(item)
		itemView.Price = buyResult.Price()
		itemsViewMap[item.ID()] = itemView
	}

	result := make([]*itemView, len(shopIds))
	for i, id := range shopIds {
		if itemView, ok := itemsViewMap[id]; ok {
			result[i] = itemView
		}
	}

	return struct {
		ShopType shop.Type   `json:"shop_type"`
		Items    []*itemView `json:"items"`
	}{
		ShopType: shopAction.ShopType(),
		Items:    result,
	}, nil
}

func itemToItemView(item *model.Item) *itemView {
	return &itemView{
		Id:             item.ID(),
		CollectionName: schema.CollectionItems,
		Name:           item.Name(),
		Description:    item.Description(),
		Icon:           item.Icon(),
		Price:          item.Price(),
	}
}
