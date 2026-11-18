package buy

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_actions/custom/shop"
	"context"
)

func calculatePrice(
	ctx context.Context,
	events *model.Events,
	item *model.Item,
	shop *shop.Shop,
	mode model.EffectRunMode,
) (*model.BuyResult, error) {
	price := shop.ApplyPriceMultiplier(item.Price())

	onItemBuy := &model.OnItemBuyEvent{
		Mode:   mode,
		Item:   item,
		Result: model.NewBuyResult(price),
	}
	err := events.OnItemBuy().Trigger(ctx, onItemBuy)
	if err != nil {
		return nil, err
	}

	return onItemBuy.Result, nil
}
