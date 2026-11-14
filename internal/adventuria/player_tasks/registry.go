package player_tasks

import "adventuria/internal/adventuria/model"

type PlayerTaskCreator func(data []byte) (model.PlayerTask, error)

type PlayerTaskDef struct {
	t model.PlayerTaskType
	c PlayerTaskCreator
}

type SubscribablePlayerTaskDef struct {
	PlayerTaskDef
	PlayerTaskSubscriber
}

func (p *PlayerTaskDef) NewFromJSON(data []byte) (model.PlayerTask, error) {
	return p.c(data)
}

func (p *PlayerTaskDef) Type() model.PlayerTaskType {
	return p.t
}

var registry = &taskRegistry{
	playerTasks:             map[model.PlayerTaskType]PlayerTaskDef{},
	subscribablePlayerTasks: map[model.PlayerTaskType]SubscribablePlayerTaskDef{},
}

type taskRegistry struct {
	playerTasks             map[model.PlayerTaskType]PlayerTaskDef
	subscribablePlayerTasks map[model.PlayerTaskType]SubscribablePlayerTaskDef
}

func (r *taskRegistry) Register(playerTasks ...PlayerTaskDef) {
	for _, playerTask := range playerTasks {
		r.playerTasks[playerTask.t] = playerTask
	}
}

func (r *taskRegistry) RegisterSubscribable(playerTasks ...SubscribablePlayerTaskDef) {
	for _, playerTask := range playerTasks {
		r.playerTasks[playerTask.t] = playerTask.PlayerTaskDef
		r.subscribablePlayerTasks[playerTask.t] = playerTask
	}
}

func (r *taskRegistry) Get(t model.PlayerTaskType) (PlayerTaskDef, bool) {
	def, ok := r.playerTasks[t]
	return def, ok
}

func (r *taskRegistry) GetSubscribable(t model.PlayerTaskType) (SubscribablePlayerTaskDef, bool) {
	def, ok := r.subscribablePlayerTasks[t]
	return def, ok
}

func NewPlayerTask(t model.PlayerTaskType, c PlayerTaskCreator) PlayerTaskDef {
	return PlayerTaskDef{
		t: t,
		c: c,
	}
}

func NewSubscribablePlayerTask(t model.PlayerTaskType, c PlayerTaskCreator, subscriber PlayerTaskSubscriber) SubscribablePlayerTaskDef {
	return SubscribablePlayerTaskDef{
		PlayerTaskDef:        NewPlayerTask(t, c),
		PlayerTaskSubscriber: subscriber,
	}
}

func Register(playerTasks ...PlayerTaskDef) {
	registry.Register(playerTasks...)
}

func RegisterSubscribable(playerTasks ...SubscribablePlayerTaskDef) {
	registry.RegisterSubscribable(playerTasks...)
}
