package buy

import (
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
	itemsData := player.LastAction().DataList().Items
	items, err := b.items.GetByIDs(ctx, itemsData.Ids)
	if err != nil {
		return nil, err
	}

	itemsViewMap := make(map[string]*itemView, len(items))
	for _, item := range items {
		basePrice, err := b.calculatePrice(item.Price(), itemsData)
		if err != nil {
			return nil, err
		}

		onBuyGetVariants, err := b.triggerOnBuyGetView(ctx, events, item, basePrice)
		if err != nil {
			return nil, err
		}

		itemView := itemToItemView(item)
		itemView.Price = onBuyGetVariants.Price
		itemsViewMap[item.ID()] = itemView
	}

	result := make([]*itemView, len(itemsData.Ids))
	for i, id := range itemsData.Ids {
		if itemView, ok := itemsViewMap[id]; ok {
			result[i] = itemView
		}
	}

	return struct {
		Items []*itemView `json:"items"`
	}{
		Items: result,
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
