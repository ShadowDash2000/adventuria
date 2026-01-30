package actions

import (
	"adventuria/internal/adventuria"
	"fmt"
)

type RerollAction struct {
	adventuria.ActionBase
}

func (a *RerollAction) CanDo(user adventuria.User) bool {
	currentCell, ok := user.CurrentCell()
	if ok {
		if currentCell.CantReroll() {
			return false
		}
	}

	return user.LastAction().Type() == ActionTypeRollWheel
}

func (a *RerollAction) Do(user adventuria.User, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
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

	cell, ok := user.CurrentCell()
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell not found",
		}, fmt.Errorf("reroll.do(): current cell not found")
	}

	cellWheel, ok := cell.(adventuria.CellWheel)
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: current cell isn't wheel cell",
		}, fmt.Errorf("reroll.do(): current cell isn't wheel cell")
	}

	action := user.LastAction()
	action.SetType(ActionTypeReroll)
	action.SetComment(comment)
	err := adventuria.PocketBase.Save(action.ProxyRecord())
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't save action record",
		}, fmt.Errorf("reroll.do(): %w", err)
	}
	action.MarkAsNew()

	err = cellWheel.RefreshItems(user)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: can't refresh items on cell",
		}, fmt.Errorf("reroll.do(): %w", err)
	}

	res, err := user.OnAfterReroll().Trigger(&adventuria.OnAfterRerollEvent{})
	if res != nil && !res.Success {
		return &adventuria.ActionResult{
			Success: false,
			Error:   res.Error,
		}, err
	}
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error: failed to trigger onAfterRerollEvent event",
		}, err
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}
