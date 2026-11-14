package npc

import "adventuria/internal/adventuria/model"

type Factory struct {
	notify notifyService
}

func NewFactory(notify notifyService) *Factory {
	return &Factory{
		notify: notify,
	}
}

func (f *Factory) New(npc *model.NPC) *Actor {
	return NewActor(f.notify, npc)
}
