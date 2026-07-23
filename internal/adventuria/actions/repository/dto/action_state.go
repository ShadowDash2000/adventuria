package dto

import (
	"adventuria/internal/adventuria/model"
	"errors"
	"fmt"
)

type ActionState struct {
	Activities ActivitiesState `json:"activities"`
	Items      ItemsState      `json:"items"`
	Shop       ShopState       `json:"shop"`
	Dealer     *DealerState    `json:"dealer,omitempty"`
}

type ActivitiesState struct {
	Ids []string `json:"ids"`
}

type ItemsState struct {
	Ids []string `json:"ids"`
}

type ShopState struct {
	Ids             []string `json:"ids"`
	PriceMultiplier float64  `json:"price_multiplier"`
}

type DealerState struct {
	Type         string            `json:"type"`
	Description  string            `json:"description"`
	CoinsForItem *DealCoinsForItem `json:"coins_for_item"`
}

type DealCoinsForItem struct {
	Coins  int    `json:"coins"`
	ItemId string `json:"item_id"`
}

func ActionStateToDTO(state model.ActionState) (ActionState, error) {
	dealerStateDTO, err := dealerStateToDTO(state.Dealer)
	if err != nil {
		return ActionState{}, err
	}

	return ActionState{
		Activities: ActivitiesState{
			Ids: state.Activities.Ids,
		},
		Items: ItemsState{
			Ids: state.Items.Ids,
		},
		Shop: ShopState{
			Ids:             state.Shop.Ids,
			PriceMultiplier: state.Shop.PriceMultiplier,
		},
		Dealer: dealerStateDTO,
	}, nil
}

func ActionStateFromDTO(dto ActionState) (model.ActionState, error) {
	dealerState, err := dealerStateFromDTO(dto.Dealer)
	if err != nil {
		return model.ActionState{}, err
	}

	return model.ActionState{
		Activities: model.ActionActivitiesState{
			Ids: dto.Activities.Ids,
		},
		Items: model.ActionItemsState{
			Ids: dto.Items.Ids,
		},
		Shop: model.ActionShopState{
			Ids:             dto.Shop.Ids,
			PriceMultiplier: dto.Shop.PriceMultiplier,
		},
		Dealer: dealerState,
	}, nil
}

func dealerStateToDTO(state *model.ActionDealerState) (*DealerState, error) {
	if state == nil {
		return nil, nil
	}

	dto := &DealerState{
		Type:        string(state.Type),
		Description: state.Description,
	}
	switch state.Type {
	case model.DealTypeCoinsForItem:
		if state.CoinsForItem == nil {
			return nil, errors.New("deal data is nil")
		}

		dto.CoinsForItem = &DealCoinsForItem{
			Coins:  state.CoinsForItem.Coins,
			ItemId: state.CoinsForItem.ItemId,
		}
	default:
		return nil, fmt.Errorf("unknown deal type: %s", state.Type)
	}

	return dto, nil
}

func dealerStateFromDTO(dto *DealerState) (*model.ActionDealerState, error) {
	if dto == nil {
		return nil, nil
	}

	state := &model.ActionDealerState{
		Type:        model.DealType(dto.Type),
		Description: dto.Description,
	}
	switch state.Type {
	case model.DealTypeCoinsForItem:
		if dto.CoinsForItem == nil {
			return nil, errors.New("deal data is nil")
		}

		state.CoinsForItem = &model.DealCoinsForItem{
			Coins:  dto.CoinsForItem.Coins,
			ItemId: dto.CoinsForItem.ItemId,
		}
	default:
		return nil, fmt.Errorf("unknown deal type: %s", state.Type)
	}

	return state, nil
}
