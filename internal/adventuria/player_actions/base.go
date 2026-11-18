package player_actions

import (
	"adventuria/internal/adventuria/model"
	"bytes"
)

type PlayerActionBase struct {
	t           model.PlayerActionType
	payloadJSON []byte
}

func NewPlayerActionBase(t model.PlayerActionType, data []byte) *PlayerActionBase {
	return &PlayerActionBase{
		t:           t,
		payloadJSON: data,
	}
}

func (p *PlayerActionBase) Type() model.PlayerActionType {
	return p.t
}

func (p *PlayerActionBase) Payload() []byte {
	return bytes.Clone(p.payloadJSON)
}

func (p *PlayerActionBase) SetPayload(data []byte) {
	p.payloadJSON = data
}
