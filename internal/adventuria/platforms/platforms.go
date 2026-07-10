package platforms

import (
	"adventuria/internal/adventuria/model"
	"context"

	"github.com/google/uuid"
)

type repository interface {
	GetOrCreate(ctx context.Context, id uuid.UUID, data model.PlatformCreate) (*model.Platform, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
	Save(ctx context.Context, platform *model.Platform) (*model.Platform, error)
}

type Platforms struct {
	repository repository
}

func NewPlatforms(repository repository) *Platforms {
	return &Platforms{repository: repository}
}

func (p *Platforms) GetOrCreate(ctx context.Context, id uuid.UUID, data model.PlatformCreate) (*model.Platform, error) {
	return p.repository.GetOrCreate(ctx, id, data)
}

func (p *Platforms) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return p.repository.GetChecksumsByIDs(ctx, ids)
}

func (p *Platforms) Save(ctx context.Context, platform *model.Platform) (*model.Platform, error) {
	return p.repository.Save(ctx, platform)
}
