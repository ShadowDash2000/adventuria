package model

import "errors"

type PlayerEventType string

type PlayerEvent interface {
	Type() PlayerEventType
	Payload() string
}

type PlayerEventData struct {
	Id      string
	Player  string
	Season  string
	Type    PlayerEventType
	Action  string
	Payload string
}

type PlayerEventInfo struct {
	data  PlayerEventData
	isNew bool
}

type PlayerEventCreate struct {
	Player string
	Season string
	Action string
}

func NewPlayerEvent(data PlayerEventCreate, event PlayerEvent) (*PlayerEventInfo, error) {
	if data.Player == "" {
		return nil, errors.New("player_event: player is empty")
	}
	if data.Season == "" {
		return nil, errors.New("player_event: season is empty")
	}
	if event == nil {
		return nil, errors.New("player_event: event is empty")
	}
	eventType := event.Type()
	if eventType == "" {
		return nil, errors.New("player_event: type is empty")
	}
	payload := event.Payload()
	if payload == "" {
		return nil, errors.New("player_event: payload is empty")
	}

	return &PlayerEventInfo{
		data: PlayerEventData{
			Player:  data.Player,
			Season:  data.Season,
			Type:    eventType,
			Action:  data.Action,
			Payload: payload,
		},
		isNew: true,
	}, nil
}

func RestorePlayerEvent(data PlayerEventData) *PlayerEventInfo {
	return &PlayerEventInfo{
		data:  data,
		isNew: false,
	}
}

func (p *PlayerEventInfo) IsNew() bool {
	return p.isNew
}

func (p *PlayerEventInfo) ID() string {
	return p.data.Id
}

func (p *PlayerEventInfo) Player() string {
	return p.data.Player
}

func (p *PlayerEventInfo) Season() string {
	return p.data.Season
}

func (p *PlayerEventInfo) Type() PlayerEventType {
	return p.data.Type
}

func (p *PlayerEventInfo) Action() string {
	return p.data.Action
}

func (p *PlayerEventInfo) Payload() string {
	return p.data.Payload
}
