package player_tasks

import "adventuria/internal/adventuria/model"

type PlayerTaskBase struct {
	t           model.PlayerTaskType
	trigger     model.PlayerTaskTrigger
	payloadJSON string
}

func NewPlayerTaskBase(t model.PlayerTaskType, trigger model.PlayerTaskTrigger, payload string) PlayerTaskBase {
	return PlayerTaskBase{
		t:           t,
		trigger:     trigger,
		payloadJSON: payload,
	}
}

func (p PlayerTaskBase) Type() model.PlayerTaskType {
	return p.t
}

func (p PlayerTaskBase) Trigger() model.PlayerTaskTrigger {
	return p.trigger
}

func (p PlayerTaskBase) Payload() string {
	return p.payloadJSON
}
