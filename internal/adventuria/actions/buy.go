package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
)

type BuyAction struct {
	adventuria.ActionBase
}

func (a *BuyAction) CanDo() bool {
	currentCell, ok := a.User().CurrentCell()
	if !ok {
		return false
	}

	if _, ok = currentCell.(*cells.CellShop); !ok {
		return false
	}

	return true
}

func (a *BuyAction) Do(req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	return &adventuria.ActionResult{
		Success: true,
	}, nil
}
