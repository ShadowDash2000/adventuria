package npcs

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
)

type NPC struct {
	npc         *model.NPC
	behavior    *model.NPCBehavior
	actionEvent action_events.ActionEvent
	effects     []effects.Effect
}

func RestoreNPC(npc *model.NPC, behavior *model.NPCBehavior, actionEvent action_events.ActionEvent, effects []effects.Effect) *NPC {
	return &NPC{
		npc:         npc,
		behavior:    behavior,
		actionEvent: actionEvent,
		effects:     effects,
	}
}

func (n *NPC) NPC() *model.NPC {
	return n.npc
}

func (n *NPC) Behavior() *model.NPCBehavior {
	return n.behavior
}

func (n *NPC) ActionEvent() action_events.ActionEvent {
	return n.actionEvent
}

func (n *NPC) Effects() []effects.Effect {
	return n.effects
}
