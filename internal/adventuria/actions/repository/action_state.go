package repository

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type actionState struct {
	ActivityCompletion *activityCompletionState `json:"activity_completion,omitzero"`
	Activities         *activitiesState         `json:"activities,omitempty"`
	UsedItems          usedItemsState           `json:"used_items,omitempty"`
	ActivityFilter     *activityFilterState     `json:"activity_filter,omitempty"`
	Items              *itemsState              `json:"items,omitempty"`
}

type activityCompletionState struct {
	Points                  int  `json:"points,omitzero"`
	DoneEnergyConsume       int  `json:"done_energy_consume,omitzero"`
	DropEnergyConsume       int  `json:"drop_energy_consume,omitzero"`
	Coins                   int  `json:"coins,omitzero"`
	CanDrop                 bool `json:"can_drop,omitzero"`
	CanReroll               bool `json:"can_reroll,omitzero"`
	IsSafeDrop              bool `json:"is_safe_drop,omitzero"`
	IsCustomFilterAllowed   bool `json:"is_custom_filter_allowed,omitzero"`
	IsChangeActivityAllowed bool `json:"is_change_activity_allowed,omitzero"`
	GiveItemWheel           bool `json:"give_item_wheel,omitzero"`
}

type activitiesState struct {
	Ids []string `json:"ids"`
}

type usedItemsState []usedItemState

type usedItemState struct {
	Id string `json:"id"`
}

type activityFilterState struct {
	Type            model.ActivityType `json:"type"`
	Platforms       []string           `json:"platforms,omitempty"`
	PlatformsStrict bool               `json:"platforms_strict"`
	GameTypes       []string           `json:"game_types,omitempty"`
	Developers      []string           `json:"developers,omitempty"`
	Publishers      []string           `json:"publishers,omitempty"`
	Genres          []string           `json:"genres,omitempty"`
	Tags            []string           `json:"tags,omitempty"`
	Themes          []string           `json:"themes,omitempty"`
	MinPrice        int                `json:"min_price"`
	MaxPrice        int                `json:"max_price"`
	ReleaseDateFrom time.Time          `json:"release_date_from,omitzero"`
	ReleaseDateTo   time.Time          `json:"release_date_to,omitzero"`
	MinCampaignTime float64            `json:"min_campaign_time"`
	MaxCampaignTime float64            `json:"max_campaign_time"`
	Activities      []string           `json:"activities,omitempty"`
}

type itemsState struct {
	Ids []string `json:"ids"`
}

func actionStateToDTO(state model.ActionState) (actionState, error) {
	return actionState{
		ActivityCompletion: activityCompletionStateToDTO(state.ActivityCompletion),
		Activities:         activitiesStateToDTO(state.Activities),
		UsedItems:          usedItemsStateToDTO(state.UsedItems),
		ActivityFilter:     activityFilterStateToDTO(state.ActivityFilter),
		Items:              itemsStateToDTO(state.Items),
	}, nil
}

func actionStateFromDTO(dto actionState) (model.ActionState, error) {
	return model.ActionState{
		ActivityCompletion: activityCompletionStateFromDTO(dto.ActivityCompletion),
		Activities:         activitiesStateFromDTO(dto.Activities),
		UsedItems:          usedItemsStateFromDTO(dto.UsedItems),
		ActivityFilter:     activityFilterStateFromDTO(dto.ActivityFilter),
		Items:              itemsStateFromDTO(dto.Items),
	}, nil
}

func activityCompletionStateToDTO(state *model.ActionActivityCompletionState) *activityCompletionState {
	if state == nil {
		return nil
	}

	return &activityCompletionState{
		Points:                  state.Points,
		DoneEnergyConsume:       state.DoneEnergyConsume,
		DropEnergyConsume:       state.DropEnergyConsume,
		Coins:                   state.Coins,
		CanDrop:                 state.CanDrop,
		CanReroll:               state.CanReroll,
		IsSafeDrop:              state.IsSafeDrop,
		IsCustomFilterAllowed:   state.IsCustomFilterAllowed,
		IsChangeActivityAllowed: state.IsChangeActivityAllowed,
		GiveItemWheel:           state.GiveItemWheel,
	}
}

func activityCompletionStateFromDTO(dto *activityCompletionState) *model.ActionActivityCompletionState {
	if dto == nil {
		return nil
	}

	return &model.ActionActivityCompletionState{
		Points:                  dto.Points,
		DoneEnergyConsume:       dto.DoneEnergyConsume,
		DropEnergyConsume:       dto.DropEnergyConsume,
		Coins:                   dto.Coins,
		CanDrop:                 dto.CanDrop,
		CanReroll:               dto.CanReroll,
		IsSafeDrop:              dto.IsSafeDrop,
		IsCustomFilterAllowed:   dto.IsCustomFilterAllowed,
		IsChangeActivityAllowed: dto.IsChangeActivityAllowed,
		GiveItemWheel:           dto.GiveItemWheel,
	}
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

func usedItemsStateToDTO(state model.ActionUsedItemsState) usedItemsState {
	if state == nil {
		return nil
	}

	res := make([]usedItemState, len(state))
	for i, item := range state {
		res[i] = usedItemStateToDTO(item)
	}

	return res
}

func usedItemStateToDTO(state model.ActionUsedItemState) usedItemState {
	return usedItemState{
		Id: state.Id,
	}
}

func usedItemsStateFromDTO(dto usedItemsState) model.ActionUsedItemsState {
	if dto == nil {
		return nil
	}

	res := make([]model.ActionUsedItemState, len(dto))
	for i, item := range dto {
		res[i] = usedItemStateFromDTO(item)
	}

	return res
}

func usedItemStateFromDTO(dto usedItemState) model.ActionUsedItemState {
	return model.ActionUsedItemState{
		Id: dto.Id,
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
