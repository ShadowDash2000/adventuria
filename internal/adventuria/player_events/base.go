package player_events

import "adventuria/internal/adventuria/model"

type PlayerEventBase struct {
	t           model.PlayerEventType
	payloadJSON string
}

func NewPlayerEventBase(t model.PlayerEventType, payload string) PlayerEventBase {
	return PlayerEventBase{
		t:           t,
		payloadJSON: payload,
	}
}

func (p PlayerEventBase) Type() model.PlayerEventType {
	return p.t
}

func (p PlayerEventBase) Payload() string {
	return p.payloadJSON
}
