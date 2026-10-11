package player_progress

import "context"

type repositoryNotify interface {
	BalanceChangeNotify(ctx context.Context, id string, balance int) error
	EnergyChangeNotify(ctx context.Context, id string, energy int) error
}

type PlayerProgressNotify struct {
	repository repositoryNotify
}

func NewPlayerProgressNotify(repository repositoryNotify) *PlayerProgressNotify {
	return &PlayerProgressNotify{
		repository: repository,
	}
}

func (p *PlayerProgressNotify) BalanceChangeNotify(ctx context.Context, id string, balance int) error {
	return p.repository.BalanceChangeNotify(ctx, id, balance)
}

func (p *PlayerProgressNotify) EnergyChangeNotify(ctx context.Context, id string, energy int) error {
	return p.repository.EnergyChangeNotify(ctx, id, energy)
}
