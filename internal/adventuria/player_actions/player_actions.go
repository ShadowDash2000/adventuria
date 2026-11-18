package player_actions

import (
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
	"time"
)

type repository interface {
	Create(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
	Update(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error)
	UpdatePayloadByID(ctx context.Context, id string, payload []byte) error
	DeleteByID(ctx context.Context, id string) error
	DeleteByNPC(ctx context.Context, playerId, actionId, npcId, actionEventId string) error
	DeleteOlderThanExceptActionIDs(ctx context.Context, createdBefore time.Time, actionIds []string) error
	GetAllByPlayerAndActionID(ctx context.Context, playerId, actionId string) ([]*model.PlayerActionInfo, error)
}

type settingsService interface {
	CurrentSeason(ctx context.Context) (string, error)
}

type actionsService interface {
	GetAllLast(ctx context.Context, seasonId string) ([]*model.ActionInfo, error)
}

type PlayerActions struct {
	repository repository
	settings   settingsService
	actions    actionsService
}

func NewPlayerActions(repository repository, settings settingsService, actions actionsService) *PlayerActions {
	return &PlayerActions{
		repository: repository,
		settings:   settings,
		actions:    actions,
	}
}

func (p *PlayerActions) toPlayerActionEntry(actionInfo *model.PlayerActionInfo) (model.PlayerActionEntry, error) {
	actionDef, ok := Get(actionInfo.Type())
	if !ok {
		return nil, fmt.Errorf("player action type %s not found", actionInfo.Type())
	}

	return actionDef.New(actionInfo)
}

func (p *PlayerActions) Save(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error) {
	if action.IsNew() {
		return p.repository.Create(ctx, action)
	}

	return p.repository.Update(ctx, action)
}

func (p *PlayerActions) UpdatePayloadByID(ctx context.Context, id string, payload []byte) error {
	return p.repository.UpdatePayloadByID(ctx, id, payload)
}

func (p *PlayerActions) DeleteByID(ctx context.Context, id string) error {
	return p.repository.DeleteByID(ctx, id)
}

func (p *PlayerActions) DeleteByNPC(ctx context.Context, npc *model.NPC, player *model.Player) error {
	return p.repository.DeleteByNPC(ctx, player.ID(), player.LastAction().ActionChainID(), npc.ID(), npc.ActionEvent())
}

func (p *PlayerActions) GetPlayerActions(ctx context.Context, playerId, actionId string) (*model.PlayerActions, error) {
	actions, err := p.GetAllInstances(ctx, playerId, actionId)
	if err != nil {
		return nil, err
	}

	return model.NewPlayerActions(actions...), nil
}

func (p *PlayerActions) GetAllInstances(ctx context.Context, playerId, actionId string) ([]model.PlayerActionEntry, error) {
	actions, err := p.repository.GetAllByPlayerAndActionID(ctx, playerId, actionId)
	if err != nil {
		return nil, err
	}

	res := make([]model.PlayerActionEntry, len(actions))
	for i, action := range actions {
		res[i], err = p.toPlayerActionEntry(action)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (p *PlayerActions) Cleanup(ctx context.Context, createdBefore time.Time) error {
	currentSeason, err := p.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	actions, err := p.actions.GetAllLast(ctx, currentSeason)
	if err != nil {
		return err
	}

	ids := make([]string, len(actions))
	for i, action := range actions {
		ids[i] = action.ActionChainID()
	}

	return p.repository.DeleteOlderThanExceptActionIDs(ctx, createdBefore, ids)
}
