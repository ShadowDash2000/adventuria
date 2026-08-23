package steam_spy

import (
	"adventuria/internal/adventuria/model"
	"context"
	"log/slog"
)

type repository interface {
	Create(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error)
	Update(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error)
	ExistsByIdDb(ctx context.Context, idDb int) (bool, error)
	GetByAppID(ctx context.Context, id int) (*model.SteamSpy, error)
}

type remoteRepository interface {
	FetchLatestRelease(ctx context.Context) ([]*SteamSpyResponse, error)
}

type SteamSpy struct {
	logger           *slog.Logger
	repository       repository
	remoteRepository remoteRepository
}

func NewSteamSpy(logger *slog.Logger, repository repository, remoteRepository remoteRepository) *SteamSpy {
	return &SteamSpy{
		logger:           logger,
		repository:       repository,
		remoteRepository: remoteRepository,
	}
}

func (s *SteamSpy) Parse(ctx context.Context) error {
	apps, err := s.remoteRepository.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, app := range apps {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ok, err := s.repository.ExistsByIdDb(ctx, app.AppId)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		steamSpy, err := model.NewSteamSpy(model.SteamSpyCreate{
			IdDb:  app.AppId,
			Name:  app.Name,
			Price: app.Price,
		})
		if err != nil {
			s.logger.Error("Failed to create SteamSpy", "error", err, "data", app)
			continue
		}

		_, err = s.Save(ctx, steamSpy)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *SteamSpy) Save(ctx context.Context, steamSpy *model.SteamSpy) (*model.SteamSpy, error) {
	if steamSpy.IsNew() {
		return s.repository.Create(ctx, steamSpy)
	}

	return s.repository.Update(ctx, steamSpy)
}

func (s *SteamSpy) GetByAppID(ctx context.Context, id int) (*model.SteamSpy, error) {
	return s.repository.GetByAppID(ctx, id)
}
