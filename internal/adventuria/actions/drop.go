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
		IsSafeDrop:    false,
		IsDropBlocked: false,
		PointsForDrop: adventuria.GameSettings.PointsForDrop(),
	}
	res, err := user.OnBeforeDrop().Trigger(onBeforeDropEvent)
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onBeforeDropEvent event",
		}, err
	}

	if onBeforeDropEvent.IsDropBlocked {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "drop is not allowed",
		}, nil
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
	action.SetGame("")
	action.SetDiceRoll(0)

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.IsSafeDrop() {
		user.SetPoints(user.Points() + onBeforeDropEvent.PointsForDrop)
		user.SetDropsInARow(user.DropsInARow() + 1)

		if user.IsSafeDrop() {
			action.SetCanMove(true)
		} else {
			_, err = user.MoveToClosestCellType("jail")
			if err != nil {
				return &adventuria.ActionResult{
					Success: false,
					Error:   "internal error",
				}, fmt.Errorf("drop.do(): %w", err)
			}

			user.SetIsInJail(true)

			res, err = user.OnAfterGoToJail().Trigger(&adventuria.OnAfterGoToJailEvent{})
			if res != nil && !res.Success {
				return &adventuria.ActionResult{
					Success: false,
					Error:   res.Error,
				}, err
			}
			if err != nil {
				return &adventuria.ActionResult{
					Success: false,
					Error:   "internal error: failed to trigger onAfterGoToJailEvent event",
				}, err
			}
		}
	}

	res, err = user.OnAfterDrop().Trigger(&adventuria.OnAfterDropEvent{})
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterDropEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}
