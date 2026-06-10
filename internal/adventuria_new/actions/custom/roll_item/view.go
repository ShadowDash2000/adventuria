package roll_item

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.WithView = (*RollItem)(nil)

type itemView struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	Type        model.ItemType `json:"type"`
}

func (r *RollItem) GetView(ctx context.Context, _ *model.Events, _ *model.Player) (any, error) {
	items, err := r.items.GetAllRollable(ctx)
	if err != nil {
		return nil, err
	}
	return struct {
		Items []*itemView `json:"items"`
	}{
		Items: itemsToItemViews(items),
	}, nil
}

func itemToItemView(item *model.Item) *itemView {
	return &itemView{
		Id:          item.ID(),
		Name:        item.Name(),
		Icon:        item.Icon(),
		Description: item.Description(),
		Type:        item.Type(),
	}
}

func itemsToItemViews(items []*model.Item) []*itemView {
	views := make([]*itemView, len(items))
	for i, item := range items {
		views[i] = itemToItemView(item)
	}
	return views
}
