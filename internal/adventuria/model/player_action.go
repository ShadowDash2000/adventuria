package model

import "errors"

type PlayerActionType string

type PlayerAction interface {
	Type() PlayerActionType
	Payload() []byte
}

type PlayerActionData struct {
	Id          string
	Player      string
	Action      string
	NPC         string
	ActionEvent string
	Type        PlayerActionType
	Payload     []byte
}

type PlayerActionInfo struct {
	data  PlayerActionData
	isNew bool
}

type PlayerActionCreate struct {
	PlayerId      string
	ActionId      string
	NPCId         string
	ActionEventId string
	PlayerAction  PlayerAction
}

func NewPlayerAction(data PlayerActionCreate) (*PlayerActionInfo, error) {
	if data.PlayerId == "" {
		return nil, errors.New("player_action: player_id is empty")
	}
	if data.ActionId == "" {
		return nil, errors.New("player_action: action_id is empty")
	}
	if data.PlayerAction == nil {
		return nil, errors.New("player_action: player_action is empty")
	}

	return &PlayerActionInfo{
		data: PlayerActionData{
			Player:      data.PlayerId,
			Action:      data.ActionId,
			NPC:         data.NPCId,
			ActionEvent: data.ActionEventId,
			Type:        data.PlayerAction.Type(),
			Payload:     data.PlayerAction.Payload(),
		},
		isNew: true,
	}, nil
}

func RestorePlayerAction(data PlayerActionData) *PlayerActionInfo {
	return &PlayerActionInfo{
		data:  data,
		isNew: false,
	}
}

func (p *PlayerActionInfo) IsNew() bool {
	return p.isNew
}

func (p *PlayerActionInfo) ID() string {
	return p.data.Id
}

func (p *PlayerActionInfo) Player() string {
	return p.data.Player
}

func (p *PlayerActionInfo) Action() string {
	return p.data.Action
}

func (p *PlayerActionInfo) NPC() string {
	return p.data.NPC
}

func (p *PlayerActionInfo) ActionEvent() string {
	return p.data.ActionEvent
}

func (p *PlayerActionInfo) Type() PlayerActionType {
	return p.data.Type
}

func (p *PlayerActionInfo) Payload() []byte {
	return p.data.Payload
}
