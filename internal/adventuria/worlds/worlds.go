package worlds

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/model"
	"context"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.World, error)
	GetAll(ctx context.Context) ([]*model.World, error)
	GetDefault(ctx context.Context) (*model.World, error)
}

type effectsService interface {
	GetByIDs(ctx context.Context, ids []string) ([]effects.Effect, error)
}

type Worlds struct {
	repository repository
	effects    effectsService
}

func NewWorlds(repo repository, effects effectsService) *Worlds {
	return &Worlds{
		repository: repo,
		effects:    effects,
	}
}

func (w *Worlds) SubscribeEffects(ctx context.Context, actor *player.Actor, worldId string) error {
	world, err := w.GetByID(ctx, worldId)
	if err != nil {
		return err
	}

	worldEffects, err := w.effects.GetByIDs(ctx, world.Effects())
	if err != nil {
		return err
	}

	unsubKeys := make([]string, len(worldEffects))
	for i, effect := range worldEffects {
		unsubs, err := effect.Subscribe(
			ctx,
			actor,
			effects.Context{
				InvItemID: "",
				Priority:  100,
			},
			func(_ context.Context) {},
		)
		if err != nil {
			return err
		}

		unsubKey := actor.Player().ID() + ":" + worldId + ":" + string(effect.Data().Type())
		actor.Events().AddUnsubs(unsubKey, unsubs...)
		unsubKeys[i] = unsubKey
	}

	actor.Events().OnWorldChanged().BindFuncOnce(func(ctx context.Context, e *model.OnWorldChangedEvent) error {
		for _, key := range unsubKeys {
			actor.Events().Unsubscribe(key)
		}

		err := w.SubscribeEffects(ctx, actor, e.NewWorldId)
		if err != nil {
			return err
		}

		return e.Next()
	})

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
