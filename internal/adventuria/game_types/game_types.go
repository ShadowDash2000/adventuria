package game_types

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, gameType *model.GameType) (*model.GameType, error)
	Update(ctx context.Context, gameType *model.GameType) (*model.GameType, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.GameType, error)
}

type GameTypes struct {
	repository repository
}

func NewGameTypes(repo repository) *GameTypes {
	return &GameTypes{
		repository: repo,
	}
}

func (t *GameTypes) GetOrCreate(ctx context.Context, data model.GameTypeCreate) (*model.GameType, error) {
	gameType, err := t.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrGameTypeNotFound) {
			return model.NewGameType(data)
		}
		return nil, err
	}

	return gameType, nil
}

func (t *GameTypes) Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	if gameType.IsNew() {
		return t.repository.Create(ctx, gameType)
	}

	return t.repository.Update(ctx, gameType)
}
