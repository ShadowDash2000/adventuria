package model

import (
	"errors"
	"slices"
)

type ShopType string

const (
	ShopTypeBuffet ShopType = "buffet"
	ShopTypeCasino ShopType = "casino"
)

type ActionShopState struct {
	Type            ShopType
	Ids             []string
	PriceMultiplier float64
	RefreshPrice    int
}

func (a *ActionShopState) Clone() *ActionShopState {
	if a == nil {
		return nil
	}

	return &ActionShopState{
		Type:            a.Type,
		Ids:             slices.Clone(a.Ids),
		PriceMultiplier: a.PriceMultiplier,
		RefreshPrice:    a.RefreshPrice,
	}
}

type ActionShopStateCreate struct {
	Type            ShopType
	Ids             []string
	PriceMultiplier float64
	RefreshPrice    int
}

func NewShopState(data ActionShopStateCreate) (*ActionShopState, error) {
	if data.Type == "" {
		return nil, errors.New("shop: type is empty")
	}
	refreshPrice := data.RefreshPrice
	if refreshPrice == 0 {
		refreshPrice = 10
	}

	return &ActionShopState{
		Type:            data.Type,
		Ids:             data.Ids,
		PriceMultiplier: data.PriceMultiplier,
		RefreshPrice:    refreshPrice,
	}, nil
}
