package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/result"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

type RollWheelAction struct {
	adventuria.ActionBase
}

func (a *RollWheelAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.Player.Progress().CurrentCell()
	if !ok {
		return false
	}

	if !currentCell.InCategory("activity") {
		return false
	}

	return !ctx.Player.LastAction().CanMove() && ctx.Player.LastAction().Type() != ActionTypeRollWheel
}

func (a *RollWheelAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*result.Result, error) {
	currentCell, ok := ctx.Player.Progress().CurrentCell()
	if !ok {
		return result.Err("internal error: current cell not found"),
			errors.New("roll_wheel.do(): current cell not found")
	}

	onBeforeWheelRollEvent := &adventuria.OnBeforeWheelRollEvent{
		AppContext:  ctx.AppContext,
		CurrentCell: currentCell.(adventuria.CellWheel),
	}
	eventRes, err := ctx.Player.OnBeforeWheelRoll().Trigger(onBeforeWheelRollEvent)
	if err != nil {
		return eventRes, err
	}
	if eventRes.Failed() {
		return eventRes, err
	}

	res, err := onBeforeWheelRollEvent.CurrentCell.Roll(ctx.AppContext, ctx.Player, adventuria.RollWheelRequest(req))
	if err != nil {
		return result.Err("internal error: failed to roll wheel"),
			fmt.Errorf("roll_wheel.do(): %w", err)
	}

	action := ctx.Player.LastAction()
	action.SetType(ActionTypeRollWheel)
	action.SetActivity(res.WinnerId)

	eventRes, err = ctx.Player.OnAfterWheelRoll().Trigger(&adventuria.OnAfterWheelRollEvent{
		AppContext: ctx.AppContext,
		ItemId:     res.WinnerId,
	})
	if err != nil {
		return eventRes, err
	}
	if eventRes.Failed() {
		return eventRes, err
	}

	return result.Ok().WithData(res), nil
}

func (a *RollWheelAction) GetVariants(ctx adventuria.ActionContext) any {
	ids, err := ctx.Player.LastAction().ItemsList()
	if err != nil {
		return nil
	}

	currentCell, ok := ctx.Player.Progress().CurrentCell()
	if !ok {
		return nil
	}

	records, err := ctx.App.FindRecordsByIds(schema.CollectionActivities, ids)
	if err != nil {
		return nil
	}

	errs := ctx.App.ExpandRecords(records, []string{
		schema.ActivitySchema.Platforms,
		schema.ActivitySchema.Developers,
		schema.ActivitySchema.Publishers,
		schema.ActivitySchema.Genres,
		schema.ActivitySchema.Tags,
		schema.ActivitySchema.Themes,
	}, nil)
	if len(errs) > 0 {
		return nil
	}

	return struct {
		Items         []*core.Record `json:"items"`
		AudioPresetId string         `json:"audio_preset_id,omitempty"`
	}{
		Items:         records,
		AudioPresetId: currentCell.AudioPreset(),
	}
}
