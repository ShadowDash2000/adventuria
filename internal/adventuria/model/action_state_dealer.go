package model

import (
	"errors"
	"fmt"
)

type DealType string

const (
	DealTypeCoinsForItem DealType = "coins_for_item"
)

type ActionDealerState struct {
	Type         DealType
	Description  string
	CoinsForItem *DealCoinsForItem
}

func (s ActionDealerState) AsCoinsForItemDeal() (DealCoinsForItem, error) {
	if s.Type != DealTypeCoinsForItem {
		return DealCoinsForItem{}, fmt.Errorf("deal type is %q, expected %q", s.Type, DealTypeCoinsForItem)
	}
	if s.CoinsForItem == nil {
		return DealCoinsForItem{}, errors.New("deal data is nil")
	}

	return *s.CoinsForItem, nil
}

type DealCoinsForItem struct {
	Coins  int
	ItemId string
}

func NewCoinsForItemDeal(description string, coins int, itemId string) ActionDealerState {
	return ActionDealerState{
		Type:        DealTypeCoinsForItem,
		Description: description,
		CoinsForItem: &DealCoinsForItem{
			Coins:  coins,
			ItemId: itemId,
		},
	}
}
