package behaviors

import (
	"adventuria/internal/adventuria/model"
)

type NPCBehaviorCreator func(payloadData, data []byte) (Behavior, error)

type NPCBehaviorDef struct {
	t model.NPCBehaviorType
	c NPCBehaviorCreator
}

func (n *NPCBehaviorDef) New(behavior *model.NPCBehavior) (Behavior, error) {
	return n.c(behavior.Payload(), behavior.Data())
}

func (n *NPCBehaviorDef) Type() model.NPCBehaviorType {
	return n.t
}

var registry = &npcBehaviorRegistry{
	npcBehaviors: map[model.NPCBehaviorType]NPCBehaviorDef{},
}

type npcBehaviorRegistry struct {
	npcBehaviors map[model.NPCBehaviorType]NPCBehaviorDef
}

func (n *npcBehaviorRegistry) Register(npcBehaviors ...NPCBehaviorDef) {
	for _, npcBehavior := range npcBehaviors {
		n.npcBehaviors[npcBehavior.t] = npcBehavior
	}
}

func (n *npcBehaviorRegistry) Get(t model.NPCBehaviorType) (NPCBehaviorDef, bool) {
	def, ok := n.npcBehaviors[t]
	return def, ok
}

func NewNPCBehavior(t model.NPCBehaviorType, c NPCBehaviorCreator) NPCBehaviorDef {
	return NPCBehaviorDef{
		t: t,
		c: c,
	}
}

func Register(npcBehaviors ...NPCBehaviorDef) {
	registry.Register(npcBehaviors...)
}

func Get(t model.NPCBehaviorType) (NPCBehaviorDef, bool) {
	return registry.Get(t)
}
