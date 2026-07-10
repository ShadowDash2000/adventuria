package how_long_to_beat

import (
	"adventuria/internal/adventuria/model"
	"context"
	"math"

	"github.com/google/uuid"
)

type repository interface {
	Save(ctx context.Context, hltb *model.HowLongToBeat) (*model.HowLongToBeat, error)
	ExistsByIdDb(ctx context.Context, idDb int) (bool, error)
	GetByNameAndYear(ctx context.Context, name string, year int) (*model.HowLongToBeat, error)
}

type remoteRepository interface {
	FetchLatestRelease(ctx context.Context) ([]*HowLongToBeatResponse, error)
}

type HowLongToBeat struct {
	repository       repository
	remoteRepository remoteRepository
}

func NewHowLongToBeat(repository repository, remoteRepository remoteRepository) *HowLongToBeat {
	return &HowLongToBeat{
		repository:       repository,
		remoteRepository: remoteRepository,
	}
}

func (h *HowLongToBeat) Parse(ctx context.Context) error {
	games, err := h.remoteRepository.FetchLatestRelease(ctx)
	if err != nil {
		return err
	}

	for _, game := range games {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		ok, err := h.repository.ExistsByIdDb(ctx, game.ID)
		if err != nil {
			return err
		}
		if ok {
			continue
		}

		hltb, err := model.NewHowLongToBeat(uuid.New(), model.HowLongToBeatCreate{
			IdDb:     game.ID,
			Name:     game.Name,
			Year:     game.ReleaseWorld,
			Campaign: math.Round(float64(game.CompMain) / 3600),
		})
		if err != nil {
			return err
		}

		_, err = h.repository.Save(ctx, hltb)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *HowLongToBeat) GetByNameAndYear(ctx context.Context, name string, year int) (*model.HowLongToBeat, error) {
	return h.repository.GetByNameAndYear(ctx, name, year)
}
