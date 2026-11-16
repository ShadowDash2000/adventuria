package activities

import (
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type DoneResultCalculator struct{}

func NewDoneResultCalculator() *DoneResultCalculator {
	return &DoneResultCalculator{}
}

func (d *DoneResultCalculator) Calculate(ctx context.Context, actor *player.Actor, mode model.EffectRunMode) (*model.ActivityDoneResult, error) {
	state := actor.Player().LastAction().State().ActivityCompletion
	if state == nil {
		return nil, errors.New("no activity completion state")
	}

	onActivityComplete := &model.OnActivityDoneEvent{
		Mode: mode,
		Result: model.NewActivityDoneResult(model.ActivityDoneResultData{
			Points:        state.Points,
			EnergyConsume: state.DoneEnergyConsume,
			Coins:         state.Coins,
		}),
	}
	err := actor.Events().OnActivityDone().Trigger(ctx, onActivityComplete)
	if err != nil {
		return nil, err
	}

	return onActivityComplete.Result, nil
}
