package repository

import (
	"adventuria/internal/adventuria/model"
	"errors"
	"fmt"
	"time"
)

type actionState struct {
	Activities     *activitiesState     `json:"activities,omitempty"`
	ActivityFilter *activityFilterState `json:"activity_filter,omitempty"`
	Items          *itemsState          `json:"items,omitempty"`
	Shop           *shopState           `json:"shop,omitempty"`
	ShopFilter     *shopFilterState     `json:"shop_filter,omitempty"`
	Dealer         *dealerState         `json:"dealer,omitempty"`
}

type activitiesState struct {
	Ids []string `json:"ids"`
}

type activityFilterState struct {
	Type            model.ActivityType `json:"type"`
	Platforms       []string           `json:"platforms"`
	PlatformsStrict bool               `json:"platforms_strict"`
	GameTypes       []string           `json:"game_types"`
	Developers      []string           `json:"developers"`
	Publishers      []string           `json:"publishers"`
	Genres          []string           `json:"genres"`
	Tags            []string           `json:"tags"`
	Themes          []string           `json:"themes"`
	MinPrice        int                `json:"min_price"`
	MaxPrice        int                `json:"max_price"`
	ReleaseDateFrom time.Time          `json:"release_date_from"`
	ReleaseDateTo   time.Time          `json:"release_date_to"`
	MinCampaignTime float64            `json:"min_campaign_time"`
	MaxCampaignTime float64            `json:"max_campaign_time"`
	Activities      []string           `json:"activities"`
}

type itemsState struct {
	Ids []string `json:"ids"`
}

type shopState struct {
	Type            model.ShopType `json:"type"`
	Ids             []string       `json:"ids"`
	PriceMultiplier float64        `json:"price_multiplier"`
	RefreshPrice    int            `json:"refresh_price"`
}

type shopFilterState struct {
	ItemType model.ItemType `json:"item_type"`
	Ids      []string       `json:"ids"`
}

type dealerState struct {
	Type         string            `json:"type"`
	Description  string            `json:"description"`
	CoinsForItem *dealCoinsForItem `json:"coins_for_item"`
}

type dealCoinsForItem struct {
	Coins  int    `json:"coins"`
	ItemId string `json:"item_id"`
}

func actionStateToDTO(state model.ActionState) (actionState, error) {
	dealerStateDTO, err := dealerStateToDTO(state.Dealer)
	if err != nil {
		return actionState{}, err
	}

	return actionState{
		Activities:     activitiesStateToDTO(state.Activities),
		ActivityFilter: activityFilterStateToDTO(state.ActivityFilter),
		Items:          itemsStateToDTO(state.Items),
		Shop:           shopStateToDTO(state.Shop),
		ShopFilter:     shopFilterStateToDTO(state.ShopFilter),
		Dealer:         dealerStateDTO,
	}, nil
}

func actionStateFromDTO(dto actionState) (model.ActionState, error) {
	dealerState, err := dealerStateFromDTO(dto.Dealer)
	if err != nil {
		return model.ActionState{}, err
	}

	return model.ActionState{
		Activities:     activitiesStateFromDTO(dto.Activities),
		ActivityFilter: activityFilterStateFromDTO(dto.ActivityFilter),
		Items:          itemsStateFromDTO(dto.Items),
		Shop:           shopStateFromDTO(dto.Shop),
		ShopFilter:     shopFilterStateFromDTO(dto.ShopFilter),
		Dealer:         dealerState,
	}, nil
}

func activitiesStateToDTO(state *model.ActionActivitiesState) *activitiesState {
	if state == nil {
		return nil
	}

	return &activitiesState{
		Ids: state.Ids,
	}
}

func activitiesStateFromDTO(dto *activitiesState) *model.ActionActivitiesState {
	if dto == nil {
		return nil
	}

	return &model.ActionActivitiesState{
		Ids: dto.Ids,
	}
}

func activityFilterStateToDTO(state *model.ActivityFilter) *activityFilterState {
	if state == nil {
		return nil
	}

	return &activityFilterState{
		Type:            state.Type,
		Platforms:       state.Platforms,
		PlatformsStrict: state.PlatformsStrict,
		GameTypes:       state.GameTypes,
		Developers:      state.Developers,
		Publishers:      state.Publishers,
		Genres:          state.Genres,
		Tags:            state.Tags,
		Themes:          state.Themes,
		MinPrice:        state.MinPrice,
		MaxPrice:        state.MaxPrice,
		ReleaseDateFrom: state.ReleaseDateFrom,
		ReleaseDateTo:   state.ReleaseDateTo,
		MinCampaignTime: state.MinCampaignTime,
		MaxCampaignTime: state.MaxCampaignTime,
		Activities:      state.Activities,
	}
}

func activityFilterStateFromDTO(dto *activityFilterState) *model.ActivityFilter {
	if dto == nil {
		return nil
	}

	return &model.ActivityFilter{
		Type:            dto.Type,
		Platforms:       dto.Platforms,
		PlatformsStrict: dto.PlatformsStrict,
		GameTypes:       dto.GameTypes,
		Developers:      dto.Developers,
		Publishers:      dto.Publishers,
		Genres:          dto.Genres,
		Tags:            dto.Tags,
		Themes:          dto.Themes,
		MinPrice:        dto.MinPrice,
		MaxPrice:        dto.MaxPrice,
		ReleaseDateFrom: dto.ReleaseDateFrom,
		ReleaseDateTo:   dto.ReleaseDateTo,
		MinCampaignTime: dto.MinCampaignTime,
		MaxCampaignTime: dto.MaxCampaignTime,
		Activities:      dto.Activities,
	}
}

func itemsStateToDTO(state *model.ActionItemsState) *itemsState {
	if state == nil {
		return nil
	}

	return &itemsState{
		Ids: state.Ids,
	}
}

func itemsStateFromDTO(dto *itemsState) *model.ActionItemsState {
	if dto == nil {
		return nil
	}

	return &model.ActionItemsState{
		Ids: dto.Ids,
	}
}

func shopStateToDTO(state *model.ActionShopState) *shopState {
	if state == nil {
		return nil
	}

	return &shopState{
		Type:            state.Type,
		Ids:             state.Ids,
		PriceMultiplier: state.PriceMultiplier,
		RefreshPrice:    state.RefreshPrice,
	}
}

func shopStateFromDTO(dto *shopState) *model.ActionShopState {
	if dto == nil {
		return nil
	}

	return &model.ActionShopState{
		Type:            dto.Type,
		Ids:             dto.Ids,
		PriceMultiplier: dto.PriceMultiplier,
		RefreshPrice:    dto.RefreshPrice,
	}
}

func shopFilterStateToDTO(state *model.ActionShopFilterState) *shopFilterState {
	if state == nil {
		return nil
	}

	return &shopFilterState{
		ItemType: state.ItemType,
		Ids:      state.Ids,
	}
}

func shopFilterStateFromDTO(dto *shopFilterState) *model.ActionShopFilterState {
	if dto == nil {
		return nil
	}

	return &model.ActionShopFilterState{
		ItemType: dto.ItemType,
		Ids:      dto.Ids,
	}
}

func dealerStateToDTO(state *model.ActionDealerState) (*dealerState, error) {
	if state == nil {
		return nil, nil
	}

	dto := &dealerState{
		Type:        string(state.Type),
		Description: state.Description,
	}
	switch state.Type {
	case model.DealTypeCoinsForItem:
		if state.CoinsForItem == nil {
			return nil, errors.New("deal data is nil")
		}

		dto.CoinsForItem = &dealCoinsForItem{
			Coins:  state.CoinsForItem.Coins,
			ItemId: state.CoinsForItem.ItemId,
		}
	default:
		return nil, fmt.Errorf("unknown deal type: %s", state.Type)
	}

	return dto, nil
}

func dealerStateFromDTO(dto *dealerState) (*model.ActionDealerState, error) {
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
