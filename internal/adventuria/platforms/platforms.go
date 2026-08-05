package platforms

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, platform *model.Platform) (*model.Platform, error)
	Update(ctx context.Context, platform *model.Platform) (*model.Platform, error)
	GetByIdDb(ctx context.Context, idDb string) (*model.Platform, error)
	GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error)
}

type Platforms struct {
	repository repository
}

func NewPlatforms(repository repository) *Platforms {
	return &Platforms{repository: repository}
}

func (p *Platforms) GetOrCreate(ctx context.Context, data model.PlatformCreate) (*model.Platform, error) {
	platform, err := p.repository.GetByIdDb(ctx, data.IdDb)
	if err != nil {
		if errors.Is(err, errs.ErrPlatformNotFound) {
			return model.NewPlatform(data)
		}
		return nil, err
	}

	return platform, nil
}

func (p *Platforms) GetChecksumsByIDs(ctx context.Context, ids []string) (map[string]string, error) {
	return p.repository.GetChecksumsByIDs(ctx, ids)
}

func (p *Platforms) Save(ctx context.Context, platform *model.Platform) (*model.Platform, error) {
	if platform.IsNew() {
		return p.repository.Create(ctx, platform)
	}

	return p.repository.Update(ctx, platform)
}
