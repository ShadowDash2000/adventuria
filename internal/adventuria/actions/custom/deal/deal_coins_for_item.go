package deal

import (
	"adventuria/internal/adventuria/model"
	"context"
)

func (d *Deal) doCoinsForItemDeal(ctx context.Context, events *model.Events, player *model.Player, dealerState *model.ActionDealerState) error {
	coinsForItemDeal, err := dealerState.AsCoinsForItemDeal()
	if err != nil {
		return err
	}

	_, err = d.inventories.AddItemByID(ctx, events, player, coinsForItemDeal.ItemId)
	if err != nil {
		return err
	}

	err = player.Progress().BalanceChange(-coinsForItemDeal.Coins)
	if err != nil {
		return err
	}

	return nil
}
