package buy

import (
	"adventuria/internal/adventuria/model"
	"context"
)

func calculatePrice(
	ctx context.Context,
	events *model.Events,
	item *model.Item,
	shopState *model.ActionShopState,
	mode model.EffectRunMode,
) (*model.BuyResult, error) {
	basePrice := item.Price()
	if shopState.PriceMultiplier != 0 {
		basePrice = int(float64(basePrice) * shopState.PriceMultiplier)
	}

	onItemBuy := &model.OnItemBuyEvent{
		Mode:   mode,
		Item:   item,
		Result: model.NewBuyResult(basePrice),
	}
	err := events.OnItemBuy().Trigger(ctx, onItemBuy)
	if err != nil {
		return nil, err
	}

	return onItemBuy.Result, nil
}
