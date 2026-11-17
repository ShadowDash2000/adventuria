package view

import "context"

type repository interface {
	GetAllBoard(ctx context.Context) ([]*WorldBoard, error)
}

type Worlds struct {
	repository repository
}

func NewWorlds(repository repository) *Worlds {
	return &Worlds{
		repository: repository,
	}
}

func (w *Worlds) GetAllBoard(ctx context.Context) ([]*WorldBoard, error) {
	return w.repository.GetAllBoard(ctx)
}
