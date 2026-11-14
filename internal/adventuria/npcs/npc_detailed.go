package npcs

import (
	"adventuria/internal/adventuria/model"
)

type NPCDetailed struct {
	NPC         *model.NPC
	Behavior    *model.NPCBehavior
	ActionEvent *model.ActionEventInfo
	Effects     []*model.EffectInfo
}
