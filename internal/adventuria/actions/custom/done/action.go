package done

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type activityResultCalculator interface {
	Calculate(ctx context.Context, events *model.Events, cell *model.CellInfo, mode model.EffectRunMode) (*model.ActivityCompletionResult, error)
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

var _ model.Action = (*Done)(nil)

const Type model.ActionType = "done"

type Done struct {
	actions.ActionBase
	actions                  actionsService
	cells                    cells
	activityResultCalculator activityResultCalculator
	reviews                  reviewsService
}

func NewDef(
	actionsService actionsService,
	cells cells,
	activityResultCalculator activityResultCalculator,
	reviews reviewsService,
) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Done{
				ActionBase:               actions.NewActionBase(Type),
				actions:                  actionsService,
				cells:                    cells,
				activityResultCalculator: activityResultCalculator,
				reviews:                  reviews,
			}
		},
	)
}

func (d *Done) CanDo(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !d.actions.CanDo(ctx, events, player, actions.ActionTypeCompleteActivity) {
		return false
	}

	currentCell, err := d.cells.GetByPlayer(ctx, player)
	if err != nil {
		return false
	}

	activityResult, err := d.activityResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return false
	}

	return activityResult.EnergyConsume() <= player.Progress().Energy()
}

type Request struct {
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

func (d *Done) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	review, err := d.reviews.Create(ctx, reviews.CreateInput{
		Comment: req.Comment,
		Score:   req.Score,
	})
	if err != nil {
		return nil, err
	}

	currentCell, err := d.cells.GetByPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	activityResult, err := d.activityResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModeApply)
	if err != nil {
		return nil, err
	}

	if activityResult.EnergyConsume() > player.Progress().Energy() {
		return nil, errs.ErrNotEnoughEnergy
	}

	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusDone)
	lastAction.SetReview(review.ID())

	progress := player.Progress()
	progress.SetCanMove(true)
	progress.SetDropsInARow(0)
	progress.SetIsInJail(false)
	err = progress.PointsChange(activityResult.Points())
	if err != nil {
		return nil, err
	}
	err = progress.EnergyChange(-activityResult.EnergyConsume())
	if err != nil {
		return nil, err
	}
	err = progress.BalanceChange(activityResult.Coins())
	if err != nil {
		return nil, err
	}

	return nil, events.OnAfterDone().Trigger(ctx, &model.OnAfterDoneEvent{
		CurrentCell: currentCell,
	})
}
