package player_actions

import "adventuria/internal/adventuria/model"

type PlayerActionCreator func(info *model.PlayerActionInfo) (model.PlayerActionEntry, error)

type PlayerActionDef struct {
	t model.PlayerActionType
	c PlayerActionCreator
}

func (p *PlayerActionDef) New(info *model.PlayerActionInfo) (model.PlayerActionEntry, error) {
	return p.c(info)
}

func (p *PlayerActionDef) Type() model.PlayerActionType {
	return p.t
}

var registry = &actionsRegistry{
	playerActions: map[model.PlayerActionType]PlayerActionDef{},
}

type actionsRegistry struct {
	playerActions map[model.PlayerActionType]PlayerActionDef
}

func (r *actionsRegistry) Register(playerActions ...PlayerActionDef) {
	for _, playerAction := range playerActions {
		r.playerActions[playerAction.t] = playerAction
	}
}
func (r *actionsRegistry) Get(t model.PlayerActionType) (PlayerActionDef, bool) {
	def, ok := r.playerActions[t]
	return def, ok
}
func NewPlayerAction(t model.PlayerActionType, c PlayerActionCreator) PlayerActionDef {
	return PlayerActionDef{
		t: t,
		c: c,
	}
}

func Register(playerActions ...PlayerActionDef) {
	registry.Register(playerActions...)
}

func Get(t model.PlayerActionType) (PlayerActionDef, bool) {
	return registry.Get(t)
}
