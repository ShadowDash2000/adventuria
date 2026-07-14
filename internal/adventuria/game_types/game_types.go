package game_types

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByIdDb(ctx context.Context, idDb string) (*model.GameType, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error)
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

func (t *GameTypes) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *GameTypes) Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	return t.repository.Save(ctx, gameType)
}
