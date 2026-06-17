package buy

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.WithView = (*Buy)(nil)

type itemView struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Price int    `json:"price"`
}

func (b *Buy) GetView(ctx context.Context, events *model.Events, player *model.Player) (any, error) {
	ids := player.LastAction().ItemsList()
	items, err := b.items.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	currentCell, err := b.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellShopValue, err := b.decodeValue(currentCell.Data().Value())
	if err != nil {
		return nil, err
	}

	itemsViewMap := make(map[string]*itemView, len(items))
	for _, item := range items {
		basePrice, err := b.calculatePrice(item.Price(), cellShopValue)
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

	result := make([]*itemView, len(ids))
	for i, id := range ids {
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
		Id:    item.ID(),
		Name:  item.Name(),
		Icon:  item.Icon(),
		Price: item.Price(),
	}
}
