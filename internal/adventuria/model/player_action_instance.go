package model

type PlayerActionInstance[T PlayerAction] struct {
	data   *PlayerActionInfo
	action T
}

func NewPlayerActionInstance[T PlayerAction](data *PlayerActionInfo, action T) *PlayerActionInstance[T] {
	return &PlayerActionInstance[T]{
		data:   data,
		action: action,
	}
}

func (p *PlayerActionInstance[T]) Data() *PlayerActionInfo {
	return p.data
}

func (p *PlayerActionInstance[T]) Action() T {
	return p.action
}
