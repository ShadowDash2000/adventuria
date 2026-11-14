package player_tasks

import (
	"adventuria/internal/adventuria/model"
	"bytes"
)

type PlayerTaskBase struct {
	t           model.PlayerTaskType
	trigger     model.PlayerTaskTrigger
	payloadJSON []byte
}

func NewPlayerTaskBase(t model.PlayerTaskType, trigger model.PlayerTaskTrigger, data []byte) PlayerTaskBase {
	return PlayerTaskBase{
		t:           t,
		trigger:     trigger,
		payloadJSON: data,
	}
}

func (p PlayerTaskBase) Type() model.PlayerTaskType {
	return p.t
}

func (p PlayerTaskBase) Trigger() model.PlayerTaskTrigger {
	return p.trigger
}

func (p PlayerTaskBase) Payload() []byte {
	return bytes.Clone(p.payloadJSON)
}
