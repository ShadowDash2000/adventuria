package players

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"time"
)

type repository interface {
	Exists(ctx context.Context, id string) (bool, error)
	GetAllIDs(ctx context.Context) ([]string, error)
}

type actions interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
	GetLastOrDefault(ctx context.Context, playerId string, timeFrom time.Time) (*model.ActionInfo, error)
}

type playerProgress interface {
	Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
	GetFirstOrDefault(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error)
}

type seasons interface {
	GetByID(ctx context.Context, id string) (*model.Season, error)
}

type Players struct {
	repository repository
	actions    actions
	progress   playerProgress
	seasons    seasons
}

func NewPlayers(repository repository, actions actions, progress playerProgress, seasons seasons) *Players {
	return &Players{
		repository: repository,
		actions:    actions,
		progress:   progress,
		seasons:    seasons,
	}
}

func (p *Players) GetByID(ctx context.Context, playerId, seasonId string) (*model.Player, error) {
	ok, err := p.repository.Exists(ctx, playerId)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errs.ErrPlayerNotFound
	}

	progress, err := p.progress.GetFirstOrDefault(ctx, playerId, seasonId)
	if err != nil {
		return nil, err
	}

	season, err := p.seasons.GetByID(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	action, err := p.actions.GetLastOrDefault(ctx, playerId, season.SeasonDateStart())
	if err != nil {
		return nil, err
	}

	player := model.RestorePlayer(
		model.PlayerData{Id: playerId},
		progress,
		action,
	)

	return player, nil
}

func (p *Players) GetAll(ctx context.Context, seasonId string) ([]*model.Player, error) {
	playerIds, err := p.repository.GetAllIDs(ctx)
	if err != nil {
		return nil, err
	}

	season, err := p.seasons.GetByID(ctx, seasonId)
	if err != nil {
		return nil, err
	}

	players := make([]*model.Player, len(playerIds))
	for i, playerId := range playerIds {
		progress, err := p.progress.GetFirstOrDefault(ctx, playerId, seasonId)
		if err != nil {
			return nil, err
		}

		action, err := p.actions.GetLastOrDefault(ctx, playerId, season.SeasonDateStart())
		if err != nil {
			return nil, err
		}

		players[i] = model.RestorePlayer(
			model.PlayerData{Id: playerId},
			progress,
			action,
		)
	}

	return players, nil
}

func (p *Players) Save(ctx context.Context, player *model.Player) error {
	action, err := p.actions.Save(ctx, player.LastAction())
	if err != nil {
		return err
	}

	progress, err := p.progress.Save(ctx, player.Progress())
	if err != nil {
		return err
	}

	player.SetLastAction(action)
	player.SetProgress(progress)

	return nil
}
