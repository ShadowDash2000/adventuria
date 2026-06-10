package buy

import (
	"adventuria/internal/adventuria_new/model"
)

func (b *Buy) calculatePrice(basePrice int, cellShopValue *cellShopValue) (int, error) {
	if cellShopValue.PriceMultiplier != 0 {
		basePrice = int(float32(basePrice) * cellShopValue.PriceMultiplier)
	}

	return basePrice, nil
}

func (b *Buy) triggerOnBeforeItemBuy(events *model.Events, item *model.Item, price int) (*model.OnBeforeItemBuyEvent, error) {
	onBeforeItemBuy := &model.OnBeforeItemBuyEvent{
		Item:  item,
		Price: price,
	}
	err := events.OnBeforeItemBuy().Trigger(onBeforeItemBuy)
	if err != nil {
		return nil, err
	}
	return onBeforeItemBuy, nil
}

func (b *Buy) triggerOnBuyGetView(events *model.Events, item *model.Item, price int) (*model.OnBuyGetViewEvent, error) {
	onBuyGetView := &model.OnBuyGetViewEvent{
		Item:  item,
		Price: price,
	}
	err := events.OnBuyGetView().Trigger(onBuyGetView)
	if err != nil {
		return nil, err
	}
	return onBuyGetView, nil
}
