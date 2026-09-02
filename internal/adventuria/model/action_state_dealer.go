package model

import (
	"adventuria/internal/adventuria/errs"
	"errors"
	"fmt"
)

type DealType string

const (
	DealTypeCoinsForItem DealType = "coins_for_item"
)

type ActionDealerState struct {
	Type         DealType
	CoinsForItem *DealCoinsForItem
}

func (s *ActionDealerState) Clone() *ActionDealerState {
	if s == nil {
		return nil
	}

	return &ActionDealerState{
		Type:         s.Type,
		CoinsForItem: s.CoinsForItem.Clone(),
	}
}

func (s *ActionDealerState) AsCoinsForItemDeal() (DealCoinsForItem, error) {
	if s == nil {
		return DealCoinsForItem{}, errs.ErrNoActiveDeals
	}
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

func (d *DealCoinsForItem) Clone() *DealCoinsForItem {
	if d == nil {
		return nil
	}

	return &DealCoinsForItem{
		Coins:  d.Coins,
		ItemId: d.ItemId,
	}
}

func NewCoinsForItemDeal(coins int, itemId string) *ActionDealerState {
	return &ActionDealerState{
		Type: DealTypeCoinsForItem,
		CoinsForItem: &DealCoinsForItem{
			Coins:  coins,
			ItemId: itemId,
		},
	}
}
