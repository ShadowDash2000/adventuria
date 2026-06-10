package worlds

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

type Repository interface {
	GetByID(ctx context.Context, id string) (*model.World, error)
	GetAll(ctx context.Context) ([]*model.World, error)
	GetDefault(ctx context.Context) (*model.World, error)
}

type Worlds struct {
	repository Repository
}

func NewWorlds(repo Repository) *Worlds {
	return &Worlds{repository: repo}
}

func (w *Worlds) GetByID(ctx context.Context, id string) (*model.World, error) {
	return w.repository.GetByID(ctx, id)
}

func (w *Worlds) GetAll(ctx context.Context) ([]*model.World, error) {
	return w.repository.GetAll(ctx)
}

func (w *Worlds) GetDefault(ctx context.Context) (*model.World, error) {
	return w.repository.GetDefault(ctx)
}
