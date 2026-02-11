package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/result"
	"errors"
	"fmt"
)

type RefreshShopAction struct {
	adventuria.ActionBase
}

type cellShopRefreshValue struct {
	RefreshPrice int `json:"refresh_price"`
}

func (a *RefreshShopAction) CanDo(ctx adventuria.ActionContext) bool {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	if !currentCell.InCategory("shop") {
		return false
	}

	if _, ok = currentCell.(adventuria.CellRefreshable); !ok {
		return false
	}

	return true
}

func (a *RefreshShopAction) Do(ctx adventuria.ActionContext, _ adventuria.ActionRequest) (*result.Result, error) {
	value, err := decodeCellShopRefreshValue(ctx)
	if err != nil {
		return result.Err("internal error: can't decode cell shop value"),
			fmt.Errorf("refreshShop.do(): %w", err)
	}

	if ctx.User.Balance() < value.RefreshPrice {
		return result.Err("not enough money"), nil
	}

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return result.Err("internal error: current cell not found"),
			errors.New("refreshShop.do(): current cell not found")
	}

	cellRefreshable, ok := currentCell.(adventuria.CellRefreshable)
	if !ok {
		return result.Err("internal error: current cell is not refreshable"),
			errors.New("refreshShop.do(): current cell is not refreshable")
	}

	err = cellRefreshable.RefreshItems(ctx.AppContext, ctx.User)
	if err != nil {
		return result.Err("internal error: can't refresh items on cell"),
			fmt.Errorf("refreshShop.do(): %w", err)
	}

	ctx.User.AddBalance(-value.RefreshPrice)

	return result.Ok(), nil
}

func (a *RefreshShopAction) GetVariants(ctx adventuria.ActionContext) any {
	value, err := decodeCellShopRefreshValue(ctx)
	if err != nil {
		return nil
	}

	return struct {
		RefreshPrice int `json:"refresh_price"`
	}{
		RefreshPrice: value.RefreshPrice,
	}
}

const defaultCellShopRefreshPrice = 10

func decodeCellShopRefreshValue(ctx adventuria.ActionContext) (*cellShopRefreshValue, error) {
	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return nil, errors.New("refreshShop.decodeCellShopRefreshValue(): current cell not found")
	}

	var decodedValue *cellShopRefreshValue
	if err := currentCell.UnmarshalValue(&decodedValue); err != nil {
		return decodedValue, err
	}

	if decodedValue == nil {
		return &cellShopRefreshValue{
			RefreshPrice: defaultCellShopRefreshPrice,
		}, nil
	}

	if decodedValue.RefreshPrice == 0 {
		decodedValue.RefreshPrice = defaultCellShopRefreshPrice
	}

	return decodedValue, nil
}
