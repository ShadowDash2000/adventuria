package action_events

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type repository interface {
	GetByID(ctx context.Context, id string) (*model.ActionEventInfo, error)
	GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error)
}

type ActionEvents struct {
	repository repository
}

func NewActionEvents(repository repository) *ActionEvents {
	return &ActionEvents{repository: repository}
}

func toActionEvent(actionEventInfo *model.ActionEventInfo) (model.ActionEvent, error) {
	return Create(*actionEventInfo)
}

func (a *ActionEvents) ListenToActionEvents(events *model.Events, player *model.Player) error {
	events.OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
		actionEventInfo, err := a.repository.GetByActiveCellID(ctx, e.CurrentCell.ID())
		if err != nil {
			if errors.Is(err, errs.ErrActionEventNotFound) {
				return e.Next()
			}
			return err
		}

		actionEvent, err := Create(*actionEventInfo)
		if err != nil {
			return err
		}

		err = actionEvent.Init(ctx, player)

		return e.Next()
	})

	return nil
}

func (a *ActionEvents) GetByID(ctx context.Context, id string) (*model.ActionEventInfo, error) {
	return a.repository.GetByID(ctx, id)
}

func (a *ActionEvents) GetByIDWrapped(ctx context.Context, id string) (model.ActionEvent, error) {
	actionEventInfo, err := a.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toActionEvent(actionEventInfo)
}

func (a *ActionEvents) GetByActiveCellID(ctx context.Context, activeCellId string) (*model.ActionEventInfo, error) {
	return a.repository.GetByActiveCellID(ctx, activeCellId)
}
