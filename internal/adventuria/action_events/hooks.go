package action_events

import (
	repo "adventuria/internal/adventuria/action_events/repository"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func BindHooks(pb core.App) {
	pb.OnRecordValidate(schema.CollectionActionEvents).BindFunc(func(e *core.RecordEvent) error {
		err := verify(e.Context, repo.RecordToActionEvent(e.Record))
		if err != nil {
			return err
		}
		return e.Next()
	})
}

func verify(ctx context.Context, actionEventInfo *model.ActionEventInfo) error {
	actionEventValue := actionEventInfo.Value()

	actionEventDef, ok := Get(actionEventInfo.Type())
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownActionEventType, actionEventInfo.Type())
	}

	actionEvent := actionEventDef.new(*actionEventInfo)
	verifiable, ok := actionEvent.(model.Verifiable)
	if !ok {
		// actionEventValue is JSON value so we need to check those empty values
		if actionEventValue == "\"\"" || actionEventValue == "null" {
			return nil
		}
		return errors.New("action event type is not verifiable")
	}

	err := verifiable.Verify(ctx, actionEventValue)
	if err != nil {
		return err
	}

	actionDef, ok := actions.Get(actionEventInfo.ActionType())
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownAction, actionEventInfo.ActionType())
	}

	_, ok = actionDef.New().(model.ActionEventCompatible)
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrActionIsNotEventCompatible, actionEventInfo.ActionType())
	}

	return nil
}
