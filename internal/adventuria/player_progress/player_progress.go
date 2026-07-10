package player_progress

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"

	"github.com/google/uuid"
)

type repository interface {
	Create(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
	Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
	GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error)
	ChangeBalance(ctx context.Context, id string, amount int) error
	NotifyChange(ctx context.Context, id string) error
}

type worlds interface {
	GetDefault(ctx context.Context) (*model.World, error)
}

type PlayerProgress struct {
	repository       repository
	worldsRepository worlds
}

func NewPlayerProgress(repository repository, worlds worlds) *PlayerProgress {
	return &PlayerProgress{
		repository:       repository,
		worldsRepository: worlds,
	}
}

func (p *PlayerProgress) GetFirstOrDefault(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error) {
	progress, err := p.repository.GetByPlayerId(ctx, playerId, seasonId)
	if err == nil {
		return progress, nil
	} else if !errors.Is(err, errs.ErrProgressNotFound) {
		return nil, err
	}

	world, err := p.worldsRepository.GetDefault(ctx)
	if err != nil {
		return nil, err
	}

	progress, err = model.NewPlayerProgress(uuid.New(), model.PlayerProgressCreate{
		Player:            playerId,
		Season:            seasonId,
		CurrentWorld:      world.ID(),
		MaxInventorySlots: 6,
	})
	if err != nil {
		return nil, err
	}

	progress, err = p.repository.Create(ctx, progress)
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (p *PlayerProgress) Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	return p.repository.Save(ctx, progress)
}

func (p *PlayerProgress) ChangeBalance(ctx context.Context, id string, amount int) error {
	return p.repository.ChangeBalance(ctx, id, amount)
}

func (p *PlayerProgress) NotifyChange(ctx context.Context, id string) error {
	return p.repository.NotifyChange(ctx, id)
}
