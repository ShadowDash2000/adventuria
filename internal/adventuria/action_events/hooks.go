package action_events

import (
	repo "adventuria/internal/adventuria/action_events/repository"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
	"fmt"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
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

	actionEvent, err := actionEventDef.c([]byte(actionEventValue))
	if err != nil {
		return err
	}

	verifiable, ok := actionEvent.(model.Verifiable)
	if ok {
		return verifiable.Verify(ctx, actionEventValue)
	}

	return nil
}
