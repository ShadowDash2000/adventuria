package player_tasks

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	Create(ctx context.Context, task *model.PlayerTaskInfo) (*model.PlayerTaskInfo, error)
	Update(ctx context.Context, task *model.PlayerTaskInfo) (*model.PlayerTaskInfo, error)
}

type PlayerTasks struct {
	repository repository
}

func NewPlayerTasks(repository repository) *PlayerTasks {
	return &PlayerTasks{
		repository: repository,
	}
}

func (p *PlayerTasks) Save(ctx context.Context, task *model.PlayerTaskInfo) (*model.PlayerTaskInfo, error) {
	if task.IsNew() {
		return p.repository.Create(ctx, task)
	}

	return p.repository.Update(ctx, task)
}
