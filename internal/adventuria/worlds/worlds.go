package worlds

import (
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.World, error)
	GetAll(ctx context.Context) ([]*model.World, error)
	GetDefault(ctx context.Context) (*model.World, error)
}

type effects interface {
	GetByIDs(ctx context.Context, ids []string) ([]model.Effect, error)
}

type Worlds struct {
	repository repository
	effects    effects
}

func NewWorlds(repo repository, effects effects) *Worlds {
	return &Worlds{
		repository: repo,
		effects:    effects,
	}
}

func (w *Worlds) SubscribeEffects(ctx context.Context, events *model.Events, player *model.Player, worldId string) error {
	world, err := w.GetByID(ctx, worldId)
	if err != nil {
		return err
	}

	effects, err := w.effects.GetByIDs(ctx, world.Effects())
	if err != nil {
		return err
	}

	for _, effect := range effects {
		unsubs, err := effect.Subscribe(
			ctx,
			events,
			player,
			model.EffectContext{
				InvItemID: "",
				Priority:  100,
			},
			func(_ context.Context) {},
		)
		if err != nil {
			return err
		}

		events.AddUnsubs(player.ID()+":"+worldId+":"+string(effect.Data().Type()), unsubs...)
	}

	return nil
}

func (w *Worlds) GetByID(ctx context.Context, id string) (*model.World, error) {
	return w.repository.GetByID(ctx, id)
}

func (w *Worlds) GetAll(ctx context.Context) ([]*model.World, error) {
	return w.repository.GetAll(ctx)
}

func (w *Worlds) GetDefault(ctx context.Context) (*model.World, error) {
	return w.repository.GetDefault(ctx)
}
