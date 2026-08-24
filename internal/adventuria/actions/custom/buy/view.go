package buy

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
)

var _ model.WithView = (*Buy)(nil)

type itemView struct {
	Id             string `json:"id"`
	CollectionName string `json:"collectionName"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
	Price          int    `json:"price"`
}

func (b *Buy) GetView(ctx context.Context, events *model.Events, player *model.Player) (any, error) {
	shopState := player.LastAction().State().Shop
	if shopState == nil {
		return nil, errs.ErrNoActiveShop
	}

	items, err := b.items.GetByIDs(ctx, shopState.Ids)
	if err != nil {
		return nil, err
	}

	itemsViewMap := make(map[string]*itemView, len(items))
	for _, item := range items {
		buyResult, err := calculatePrice(ctx, events, item, shopState, model.EffectRunModePreview)
		if err != nil {
			return nil, err
		}

		itemView := itemToItemView(item)
		itemView.Price = buyResult.Price()
		itemsViewMap[item.ID()] = itemView
	}

	result := make([]*itemView, len(shopState.Ids))
	for i, id := range shopState.Ids {
		if itemView, ok := itemsViewMap[id]; ok {
			result[i] = itemView
		}
	}

	return struct {
		ShopType model.ShopType `json:"shop_type"`
		Items    []*itemView    `json:"items"`
	}{
		ShopType: shopState.Type,
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
