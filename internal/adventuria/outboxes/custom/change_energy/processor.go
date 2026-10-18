package change_energy

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"adventuria/internal/adventuria/player_events/custom/energy_changed"
	"context"
	"errors"
)

type progress interface {
	GetByID(ctx context.Context, id string) (*model.PlayerProgress, error)
	ChangeEnergy(ctx context.Context, id string, amount int) error
}

type actions interface {
	GetLast(ctx context.Context, playerId, seasonId string) (*model.ActionInfo, error)
}

type playerEvents interface {
	Save(ctx context.Context, playerEvent *model.PlayerEventInfo) (*model.PlayerEventInfo, error)
}

const Type model.OutboxType = "change_energy"

type ChangeEnergyProcessor struct {
	progress     progress
	actions      actions
	playerEvents playerEvents
}

func NewDef(progress progress, actions actions, playerEvents playerEvents) outboxes.OutboxDef {
	processor := ChangeEnergyProcessor{
		progress:     progress,
		actions:      actions,
		playerEvents: playerEvents,
	}

	return outboxes.NewOutbox(
		Type,
		func(ctx context.Context, payload string) error {
			task, err := NewFromJSON(payload)
			if err != nil {
				return err
			}

			return processor.Process(ctx, task)
		},
	)
}

func (c *ChangeEnergyProcessor) Process(ctx context.Context, task *ChangeEnergy) error {
	err := c.progress.ChangeEnergy(ctx, task.payload.ProgressId, *task.payload.Amount)
	if err != nil {
		return err
	}

	progress, err := c.progress.GetByID(ctx, task.payload.ProgressId)
	if err != nil {
		return err
	}

	lastAction, err := c.actions.GetLast(ctx, progress.Player(), progress.Season())
	if err != nil && !errors.Is(err, errs.ErrActionNotFound) {
		return err
	}

	actionId := ""
	if lastAction != nil {
		actionId = lastAction.ID()
	}

	energyChangedEvent, err := energy_changed.NewPlayerEvent(
		model.PlayerEventCreate{
			Player: progress.Player(),
			Season: progress.Season(),
			Action: actionId,
		},
		energy_changed.Payload{
			Amount: new(*task.payload.Amount),
			Source: task.payload.Source,
		},
	)
	if err != nil {
		return err
	}

	_, err = c.playerEvents.Save(ctx, energyChangedEvent)
	if err != nil {
		return err
	}

	return nil
}
