package dto

import (
	"adventuria/internal/adventuria/model"
	"errors"
	"fmt"
	"time"
)

type ActionState struct {
	Activities     *ActivitiesState     `json:"activities,omitempty"`
	ActivityFilter *ActivityFilterState `json:"activity_filter,omitempty"`
	Items          *ItemsState          `json:"items,omitempty"`
	Shop           *ShopState           `json:"shop,omitempty"`
	Dealer         *DealerState         `json:"dealer,omitempty"`
}

type ActivitiesState struct {
	Ids []string `json:"ids"`
}

type ActivityFilterState struct {
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

type ItemsState struct {
	Ids []string `json:"ids"`
}

type ShopState struct {
	Type            model.ShopType `json:"type"`
	Ids             []string       `json:"ids"`
	PriceMultiplier float64        `json:"price_multiplier"`
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
		Activities:     activitiesStateToDTO(state.Activities),
		ActivityFilter: activityFilterStateToDTO(state.ActivityFilter),
		Items:          itemsStateToDTO(state.Items),
		Shop:           shopStateToDTO(state.Shop),
		Dealer:         dealerStateDTO,
	}, nil
}

func ActionStateFromDTO(dto ActionState) (model.ActionState, error) {
	dealerState, err := dealerStateFromDTO(dto.Dealer)
	if err != nil {
		return model.ActionState{}, err
	}

	return model.ActionState{
		Activities:     activitiesStateFromDTO(dto.Activities),
		ActivityFilter: activityFilterStateFromDTO(dto.ActivityFilter),
		Items:          itemsStateFromDTO(dto.Items),
		Shop:           shopStateFromDTO(dto.Shop),
		Dealer:         dealerState,
	}, nil
}

func activitiesStateToDTO(state *model.ActionActivitiesState) *ActivitiesState {
	if state == nil {
		return nil
	}

	return &ActivitiesState{
		Ids: state.Ids,
	}
}

func activitiesStateFromDTO(dto *ActivitiesState) *model.ActionActivitiesState {
	if dto == nil {
		return nil
	}

	return &model.ActionActivitiesState{
		Ids: dto.Ids,
	}
}

func activityFilterStateToDTO(state *model.ActivityFilter) *ActivityFilterState {
	if state == nil {
		return nil
	}

	return &ActivityFilterState{
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

func activityFilterStateFromDTO(dto *ActivityFilterState) *model.ActivityFilter {
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

func itemsStateToDTO(state *model.ActionItemsState) *ItemsState {
	if state == nil {
		return nil
	}

	return &ItemsState{
		Ids: state.Ids,
	}
}

func itemsStateFromDTO(dto *ItemsState) *model.ActionItemsState {
	if dto == nil {
		return nil
	}

	return &model.ActionItemsState{
		Ids: dto.Ids,
	}
}

func shopStateToDTO(state *model.ActionShopState) *ShopState {
	if state == nil {
		return nil
	}

	return &ShopState{
		Type:            state.Type,
		Ids:             state.Ids,
		PriceMultiplier: state.PriceMultiplier,
	}
}

func shopStateFromDTO(dto *ShopState) *model.ActionShopState {
	if dto == nil {
		return nil
	}

	return &model.ActionShopState{
		Type:            dto.Type,
		Ids:             dto.Ids,
		PriceMultiplier: dto.PriceMultiplier,
	}
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
