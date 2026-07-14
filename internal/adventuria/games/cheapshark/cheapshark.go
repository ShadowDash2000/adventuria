package cheapshark

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	Save(ctx context.Context, cheapShark *model.CheapShark) (*model.CheapShark, error)
	ExistsByIdDb(ctx context.Context, idDb int) (bool, error)
	GetByAppID(ctx context.Context, id int) (*model.CheapShark, error)
}

type remoteRepository interface {
	FetchLatestRelease(ctx context.Context) ([]*CheapSharkResponse, error)
}

type CheapShark struct {
	repository       repository
	remoteRepository remoteRepository
}

func NewCheapShark(repository repository, remoteRepository remoteRepository) *CheapShark {
	return &CheapShark{
		repository:       repository,
		remoteRepository: remoteRepository,
	}
}

func (c *CheapShark) Parse(ctx context.Context) error {
	deals, err := c.remoteRepository.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, deal := range deals {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ok, err := c.repository.ExistsByIdDb(ctx, int(deal.SteamAppID))
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		cs, err := model.NewCheapShark(model.CheapSharkCreate{
			IdDb:  int(deal.SteamAppID),
			Name:  deal.Title,
			Price: deal.NormalPrice,
		})
		if err != nil {
			return err
		}

		_, err = c.repository.Save(ctx, cs)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *CheapShark) GetByAppID(ctx context.Context, id int) (*model.CheapShark, error) {
	return c.repository.GetByAppID(ctx, id)
}
