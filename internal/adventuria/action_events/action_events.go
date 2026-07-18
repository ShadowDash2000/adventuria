package action_events

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error)
}

type ActionEvents struct {
	repository repository
}

func NewActionEvents(repository repository) *ActionEvents {
	return &ActionEvents{repository: repository}
}

func (a *ActionEvents) ListenToActionEvents(events *model.Events, player *model.Player) error {
	events.OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
		cellEventInfo, err := a.repository.GetByActiveCellID(ctx, e.CurrentCell.ID())
		if err != nil {
			if errors.Is(err, errs.ErrActionEventNotFound) {
				return e.Next()
			}
			return err
		}

		cellEvent, err := Create(*cellEventInfo)
		if err != nil {
			return err
		}

		err = cellEvent.Init(ctx, player)

		return e.Next()
	})

	return nil
}

func (a *ActionEvents) GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error) {
	return a.repository.GetByActiveCellID(ctx, activeCellId)
}
