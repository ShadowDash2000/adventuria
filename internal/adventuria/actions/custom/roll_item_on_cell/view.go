package roll_item_on_cell

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.WithView = (*RollItemOnCell)(nil)

type itemView struct {
	Id          string         `json:"id"`
	Name        string         `json:"name"`
	Icon        string         `json:"icon"`
	Description string         `json:"description"`
	Type        model.ItemType `json:"type"`
}

func (r *RollItemOnCell) GetView(ctx context.Context, _ *model.Events, player *model.Player) (any, error) {
	items, err := r.items.GetByIDs(ctx, player.LastAction().DataList().Items.Ids)
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
