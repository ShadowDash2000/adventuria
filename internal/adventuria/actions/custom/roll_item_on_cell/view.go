package roll_item_on_cell

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
)

var _ actions.View = (*RollItemOnCell)(nil)

type itemView struct {
	Id             string         `json:"id"`
	CollectionName string         `json:"collectionName"`
	Name           string         `json:"name"`
	Icon           string         `json:"icon"`
	Description    string         `json:"description"`
	Type           model.ItemType `json:"type"`
}

func (r *RollItemOnCell) GetView(ctx context.Context, actor *player.Actor) (any, error) {
	currentCell, err := r.cells.GetByPlayer(ctx, actor.Player())
	if err != nil {
		return nil, err
	}

	itemsState := actor.Player().LastAction().State().Items
	if itemsState == nil {
		return nil, errs.ErrNoActiveItems
	}

	items, err := r.items.GetByIDs(ctx, itemsState.Ids)
	if err != nil {
		return nil, err
	}

	return struct {
		Items         []*itemView `json:"items"`
		AudioPresetId string      `json:"audio_preset_id,omitempty"`
	}{
		Items:         itemsToItemViews(items),
		AudioPresetId: currentCell.AudioPreset(),
	}, nil
}

func itemToItemView(item *model.Item) *itemView {
	return &itemView{
		Id:             item.ID(),
		CollectionName: schema.CollectionItems,
		Name:           item.Name(),
		Icon:           item.Icon(),
		Description:    item.Description(),
		Type:           item.Type(),
	}
}

func itemsToItemViews(items []*model.Item) []*itemView {
	views := make([]*itemView, len(items))
	for i, item := range items {
		views[i] = itemToItemView(item)
	}
	return views
}
