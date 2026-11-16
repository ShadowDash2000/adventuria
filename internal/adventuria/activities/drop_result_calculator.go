package activities

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type DropResultCalculator struct{}

func NewDropResultCalculator() *DropResultCalculator {
	return &DropResultCalculator{}
}

func (d *DropResultCalculator) Calculate(ctx context.Context, actor *player.Actor, mode model.EffectRunMode) (*model.ActivityDropResult, error) {
	state := actor.Player().LastAction().State().ActivityCompletion
	if state == nil {
		return nil, errors.New("no activity completion state")
	}

	onActivityDrop := &model.OnActivityDropEvent{
		Mode: mode,
		Result: model.NewActivityDropResult(model.ActivityDropResultData{
			IsBlocked:     !state.CanDrop,
			IsSafe:        state.IsSafeDrop,
			EnergyConsume: state.DropEnergyConsume,
		}),
	}
	err := actor.Events().OnActivityDrop().Trigger(ctx, onActivityDrop)
	if err != nil {
		return nil, err
	}

	return onActivityDrop.Result, nil
}
