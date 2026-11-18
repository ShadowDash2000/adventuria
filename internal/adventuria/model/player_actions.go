package model

type PlayerActionEntry interface {
	Data() *PlayerActionInfo
}

type PlayerActionContainer[T PlayerAction] interface {
	PlayerActionEntry
	Action() T
}

type PlayerActions struct {
	actions map[PlayerActionType]PlayerActionEntry
}

func NewPlayerActions(actions ...PlayerActionEntry) *PlayerActions {
	actionsMap := make(map[PlayerActionType]PlayerActionEntry, len(actions))

	for _, action := range actions {
		actionsMap[action.Data().Type()] = action
	}

	return &PlayerActions{
		actions: actionsMap,
	}
}

func (p *PlayerActions) Get[T PlayerAction](t PlayerActionType) (PlayerActionContainer[T], bool) {
	instance, ok := p.actions[t]
	if !ok {
		return nil, false
	}

	typed, ok := instance.(PlayerActionContainer[T])

	return typed, ok
}

func (p *PlayerActions) Has[T PlayerAction](t PlayerActionType) bool {
	_, ok := p.Get[T](t)

	return ok
}
