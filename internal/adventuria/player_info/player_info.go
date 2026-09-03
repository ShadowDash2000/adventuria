package player_info

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.PlayerInfo, error)
	GetAll(ctx context.Context) ([]*model.PlayerInfo, error)
	IsDisabled(ctx context.Context, id string) (bool, error)
	IsDebugEnabled(ctx context.Context, id string) (bool, error)
	NotifyChange(ctx context.Context, id string) error
}

type PlayerInfo struct {
	repository repository
}

func NewPlayerInfo(repository repository) *PlayerInfo {
	return &PlayerInfo{
		repository: repository,
	}
}

func (p *PlayerInfo) GetByID(ctx context.Context, id string) (*model.PlayerInfo, error) {
	return p.repository.GetByID(ctx, id)
}

func (p *PlayerInfo) GetAll(ctx context.Context) ([]*model.PlayerInfo, error) {
	return p.repository.GetAll(ctx)
}

func (p *PlayerInfo) IsDisabled(ctx context.Context, id string) (bool, error) {
	return p.repository.IsDisabled(ctx, id)
}

func (p *PlayerInfo) IsDebugEnabled(ctx context.Context, id string) (bool, error) {
	return p.repository.IsDebugEnabled(ctx, id)
}

func (p *PlayerInfo) NotifyChange(ctx context.Context, id string) error {
	return p.repository.NotifyChange(ctx, id)
}
