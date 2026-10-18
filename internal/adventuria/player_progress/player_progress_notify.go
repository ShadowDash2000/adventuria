package player_progress

import "context"

type repositoryNotify interface {
	BalanceChangeNotify(ctx context.Context, id string, amount int) error
	EnergyChangeNotify(ctx context.Context, id string, amount int) error
	ItemWheelsChangeNotify(ctx context.Context, id string, amount int) error
}

type PlayerProgressNotify struct {
	repository repositoryNotify
}

func NewPlayerProgressNotify(repository repositoryNotify) *PlayerProgressNotify {
	return &PlayerProgressNotify{
		repository: repository,
	}
}

func (p *PlayerProgressNotify) BalanceChangeNotify(ctx context.Context, id string, amount int) error {
	return p.repository.BalanceChangeNotify(ctx, id, amount)
}

func (p *PlayerProgressNotify) EnergyChangeNotify(ctx context.Context, id string, amount int) error {
	return p.repository.EnergyChangeNotify(ctx, id, amount)
}

func (p *PlayerProgressNotify) ItemWheelsChangeNotify(ctx context.Context, id string, amount int) error {
	return p.repository.ItemWheelsChangeNotify(ctx, id, amount)
}
