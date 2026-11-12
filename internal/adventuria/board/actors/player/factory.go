package player

import "adventuria/internal/adventuria/model"

type Factory struct {
	actions actionsService
	notify  notifyService
}

func NewFactory(actions actionsService, notify notifyService) *Factory {
	return &Factory{
		actions: actions,
		notify:  notify,
	}
}

func (f *Factory) New(events *model.Events, player *model.Player) *Actor {
	return NewPlayerActor(f.actions, f.notify, events, player)
}
