package actions

import (
	"adventuria/internal/adventuria"
	"errors"
	"fmt"
)

type DropAction struct {
	adventuria.ActionBase
}

func (a *DropAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if ok {
		if currentCell.CantDrop() {
			return false
		}
	}

	if user.IsInJail() {
		return false
	}

	return user.LastAction().Type() == ActionTypeRollWheel
}

func (a *DropAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	var comment string
	if c, ok := req["comment"]; ok {
		comment, ok = c.(string)
		if !ok {
			return &adventuria.ActionResult{
				Success: false,
				Error:   "request error: comment is not string",
			}, nil
		}
	}

	currentCell, ok := user.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, errors.New("drop.do(): current cell not found")
	}

	onBeforeDropEvent := &adventuria.OnBeforeDropEvent{
		IsSafeDrop: false,
	}
	err := user.OnBeforeDrop().Trigger(onBeforeDropEvent)
	if err != nil {
		adventuria.PocketBase.Logger().Error(
			"drop.do(): failed to trigger onBeforeDrop event",
			"error",
			err,
		)
	}

	action := user.LastAction()
	action.SetType(ActionTypeDrop)
	action.SetComment(comment)
	err = adventuria.PocketBase.Save(action.ProxyRecord())
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't save action record",
		}, fmt.Errorf("drop.do(): %w", err)
	}
	action.ProxyRecord().MarkAsNew()
	action.ProxyRecord().Set("id", "")
	action.SetComment("")

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.IsSafeDrop() {
		user.SetPoints(user.Points() + adventuria.GameSettings.PointsForDrop())
		user.SetDropsInARow(user.DropsInARow() + 1)

		if user.IsSafeDrop() {
			action.SetCanMove(true)
		} else {
			_, err = user.MoveToCellType("jail")
			if err != nil {
				return &adventuria.ActionResult{
					Success: false,
					Error:   "internal error",
				}, fmt.Errorf("drop.do(): %w", err)
			}

			user.SetIsInJail(true)

			err = user.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
			if err != nil {
				adventuria.PocketBase.Logger().Error(
					"drop.do(): failed to trigger onAfterGoToJail event",
					"error",
					err,
				)
			}
		}
	}

	err = user.OnAfterDrop().Trigger(&adventuria.OnAfterDropEvent{})
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("drop.do(): %w", err)
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}
