package buy

import (
	"adventuria/internal/adventuria/model"
	"context"
)

func (b *Buy) calculatePrice(basePrice int, shopState model.ActionShopState) (int, error) {
	if shopState.PriceMultiplier != 0 {
		basePrice = int(float64(basePrice) * shopState.PriceMultiplier)
	}

	return basePrice, nil
}

func (b *Buy) triggerOnBeforeItemBuy(ctx context.Context, events *model.Events, item *model.Item, price int) (*model.OnBeforeItemBuyEvent, error) {
	onBeforeItemBuy := &model.OnBeforeItemBuyEvent{
		Item:  item,
		Price: price,
	}
	err := events.OnBeforeItemBuy().Trigger(ctx, onBeforeItemBuy)
	if err != nil {
		return nil, err
	}
	return onBeforeItemBuy, nil
}

func (b *Buy) triggerOnBuyGetView(ctx context.Context, events *model.Events, item *model.Item, price int) (*model.OnBuyGetViewEvent, error) {
	onBuyGetView := &model.OnBuyGetViewEvent{
		Item:  item,
		Price: price,
	}
	err := events.OnBuyGetView().Trigger(ctx, onBuyGetView)
	if err != nil {
		return nil, err
	}
	return onBuyGetView, nil
}
