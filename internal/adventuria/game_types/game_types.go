package game_types

import (
	"adventuria/internal/adventuria/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	GetOrCreate(ctx context.Context, id uuid.UUID, data model.GameTypeCreate) (*model.GameType, error)
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

func (t *GameTypes) GetOrCreate(ctx context.Context, id uuid.UUID, data model.GameTypeCreate) (*model.GameType, error) {
	return t.repository.GetOrCreate(ctx, id, data)
}

func (t *GameTypes) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return t.repository.GetChecksumsByIDs(ctx, ids)
}

func (t *GameTypes) Save(ctx context.Context, gameType *model.GameType) (*model.GameType, error) {
	return t.repository.Save(ctx, gameType)
}
