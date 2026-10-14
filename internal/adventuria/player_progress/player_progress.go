package player_progress

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes/custom/change_energy"
	"context"
	"errors"
)

type repository interface {
	Create(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
	Update(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error)
	GetByID(ctx context.Context, id string) (*model.PlayerProgress, error)
	GetByPlayerId(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error)
	GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.PlayerProgress, error)
	ChangeBalance(ctx context.Context, id string, amount int) error
	ChangeEnergy(ctx context.Context, id string, amount int) error
}

type notify interface {
	BalanceChangeNotify(ctx context.Context, id string, balance int) error
	EnergyChangeNotify(ctx context.Context, id string, energy int) error
}

type worlds interface {
	GetDefault(ctx context.Context) (*model.World, error)
}

type settings interface {
	EnergyDefault(ctx context.Context) (int, error)
}

type outboxes interface {
	Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
}

type PlayerProgress struct {
	repository repository
	notify     notify
	worlds     worlds
	settings   settings
	outboxes   outboxes
}

func NewPlayerProgress(
	repository repository,
	notify notify,
	worlds worlds,
	settings settings,
	outboxes outboxes,
) *PlayerProgress {
	return &PlayerProgress{
		repository: repository,
		notify:     notify,
		worlds:     worlds,
		settings:   settings,
		outboxes:   outboxes,
	}
}

func (p *PlayerProgress) GetFirstOrDefault(ctx context.Context, playerId, seasonId string) (*model.PlayerProgress, error) {
	progress, err := p.repository.GetByPlayerId(ctx, playerId, seasonId)
	if err == nil {
		return progress, nil
	} else if !errors.Is(err, errs.ErrProgressNotFound) {
		return nil, err
	}

	world, err := p.worlds.GetDefault(ctx)
	if err != nil {
		return nil, err
	}

	energyDefault, err := p.settings.EnergyDefault(ctx)
	if err != nil {
		return nil, err
	}

	progress, err = model.NewPlayerProgress(model.PlayerProgressCreate{
		Player:            playerId,
		Season:            seasonId,
		CurrentWorld:      world.ID(),
		Energy:            energyDefault,
		MaxInventorySlots: 6,
	})
	if err != nil {
		return nil, err
	}

	progress.SetCanMove(true)

	return progress, nil
}

func (p *PlayerProgress) GetByID(ctx context.Context, id string) (*model.PlayerProgress, error) {
	return p.repository.GetByID(ctx, id)
}

func (p *PlayerProgress) GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.PlayerProgress, error) {
	return p.repository.GetAllBySeasonID(ctx, seasonId)
}

func (p *PlayerProgress) Save(ctx context.Context, progress *model.PlayerProgress) (*model.PlayerProgress, error) {
	if progress.IsNew() {
		return p.repository.Create(ctx, progress)
	}

	return p.repository.Update(ctx, progress)
}

func (p *PlayerProgress) ChangeBalance(ctx context.Context, id string, amount int) error {
	err := p.repository.ChangeBalance(ctx, id, amount)
	if err != nil {
		return err
	}

	return p.notify.BalanceChangeNotify(ctx, id, amount)
}

func (p *PlayerProgress) ChangeEnergy(ctx context.Context, id string, amount int) error {
	err := p.repository.ChangeEnergy(ctx, id, amount)
	if err != nil {
		return err
	}

	return p.notify.EnergyChangeNotify(ctx, id, amount)
}

func (p *PlayerProgress) ChangeEnergyOutbox(ctx context.Context, id string, amount int) (*model.OutboxInfo, error) {
	task, err := change_energy.New(change_energy.Payload{
		ProgressId: id,
		Amount:     new(amount),
	})
	if err != nil {
		return nil, err
	}

	outbox, err := model.NewOutbox(model.OutboxCreate{
		Task: task,
	})
	if err != nil {
		return nil, err
	}

	return p.outboxes.Save(ctx, outbox)
}
