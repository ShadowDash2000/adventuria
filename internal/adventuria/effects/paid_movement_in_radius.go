package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"fmt"
	"strconv"
	"strings"
)

type PaidMovementInRadiusEffect struct {
	adventuria.EffectRecord
}

func (ef *PaidMovementInRadiusEffect) CanUse(appCtx adventuria.AppContext, ctx adventuria.EffectContext) bool {
	if adventuria.GameActions.HasActionsInCategories(appCtx, ctx.User, []string{"wheel_roll", "on_cell"}) {
		return false
	}

	currentCell, ok := ctx.User.CurrentCell()
	if !ok {
		return false
	}

	canDone := adventuria.GameActions.CanDo(appCtx, ctx.User, "done")
	canDrop := adventuria.GameActions.CanDo(appCtx, ctx.User, "drop")

	if canDone && !canDrop {
		if currentCell.Type() != "jail" {
			return false
		}
	}

	value, err := ef.DecodeValue(ef.GetString("value"))
	if err != nil {
		return false
	}

	return ctx.User.Balance() >= value.Price
}

func (ef *PaidMovementInRadiusEffect) Subscribe(
	ctx adventuria.EffectContext,
	callback adventuria.EffectCallback,
) ([]event.Unsubscribe, error) {
	return []event.Unsubscribe{
		ctx.User.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*event.Result, error) {
			if ctx.InvItemID == e.InvItemId {
				if cellId, ok := e.Data["cell_id"].(string); ok {
					currentCell, ok := ctx.User.CurrentCell()
					if !ok {
						return &event.Result{
							Success: false,
							Error:   "current cell not found",
						}, nil
					}

					distance, ok := adventuria.GameCells.GetDistanceByIds(currentCell.ID(), cellId)
					if !ok {
						return &event.Result{
							Success: false,
							Error:   "destination cell not found",
						}, nil
					}

					value, _ := ef.DecodeValue(ef.GetString("value"))
					if distance == 0 || distance > value.Radius {
						return &event.Result{
							Success: false,
							Error:   "destination cell is too far",
						}, nil
					}

					_, err := ctx.User.MoveToCellId(e.AppContext, cellId)
					if err != nil {
						return &event.Result{
							Success: false,
							Error:   "failed to move to cell",
						}, nil
					}

					ctx.User.SetBalance(ctx.User.Balance() - value.Price)

					callback(e.AppContext)
				} else {
					return &event.Result{
						Success: false,
						Error:   "cell_id not found in request",
					}, nil
				}
			}

			return e.Next()
		}),
	}, nil
}

func (ef *PaidMovementInRadiusEffect) Verify(_ adventuria.AppContext, value string) error {
	_, err := ef.DecodeValue(value)
	return err
}

func (ef *PaidMovementInRadiusEffect) GetVariants(_ adventuria.AppContext, ctx adventuria.EffectContext) any {
	value, _ := ef.DecodeValue(ef.GetString("value"))
	currentCellOrder := ctx.User.CurrentCellOrder()

	var variants []any

	startOrder := currentCellOrder - value.Radius
	if startOrder < 0 {
		startOrder = 0
	}

	endOrder := currentCellOrder + value.Radius
	if endOrder > adventuria.GameCells.Count()-1 {
		endOrder = adventuria.GameCells.Count() - 1
	}

	for i := startOrder; i <= endOrder; i++ {
		if i == currentCellOrder {
			continue
		}

		if cell, ok := adventuria.GameCells.GetByOrder(i); ok {
			variants = append(variants, cell)
		}
	}

	return variants
}

type PaidMovementInRadiusEffectValue struct {
	Radius int
	Price  int
}

func (ef *PaidMovementInRadiusEffect) DecodeValue(value string) (*PaidMovementInRadiusEffectValue, error) {
	vals := strings.Split(value, ";")
	if len(vals) != 2 {
		return nil, fmt.Errorf("paidMovementInRadius: invalid value, expected format 'radius;price': %s", value)
	}

	var (
		radius, price int
		err           error
	)
	if radius, err = strconv.Atoi(vals[0]); err != nil {
		return nil, fmt.Errorf("paidMovementInRadius: invalid radius value: %s", vals[0])
	} else if price, err = strconv.Atoi(vals[1]); err != nil {
		return nil, fmt.Errorf("paidMovementInRadius: invalid price value: %s", vals[1])
	}

	return &PaidMovementInRadiusEffectValue{
		Radius: radius,
		Price:  price,
	}, nil
}
