package done

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type reviews interface {
	Save(ctx context.Context, review *model.Review) (*model.Review, error)
}

var _ model.Action = (*Done)(nil)

const Type model.ActionType = "done"

type Done struct {
	actions.ActionBase
	actions actionsService
	cells   cells
	reviews reviews
}

func NewDef(actionsService actionsService, cells cells, reviews reviews) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Done{
				ActionBase: actions.NewActionBase(Type),
				actions:    actionsService,
				cells:      cells,
				reviews:    reviews,
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

	return currentCell.EnergyConsume() <= player.Progress().Energy()
}

type Request struct {
	Comment string `json:"comment"`
	Score   int    `json:"score"`
}

func (d *Done) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	review, err := model.NewReview(req.Comment, req.Score)
	if err != nil {
		return nil, err
	}
	review, err = d.reviews.Save(ctx, review)
	if err != nil {
		return nil, err
	}

	currentCell, err := d.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return nil, err
	}

	if currentCell.Data().EnergyConsume() > player.Progress().Energy() {
		return nil, errs.ErrNotEnoughEnergy
	}

	onBeforeDoneEvent := &model.OnBeforeDoneEvent{
		CellPoints:        currentCell.Data().Points(),
		CellEnergyConsume: currentCell.Data().EnergyConsume(),
		CellCoins:         currentCell.Data().Coins(),
	}
	err = events.OnBeforeDone().Trigger(ctx, onBeforeDoneEvent)
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetType(Type)
	lastAction.SetReview(review.ID())

	progress := player.Progress()
	progress.SetCanMove(true)
	progress.SetDropsInARow(0)
	progress.SetIsInJail(false)
	err = progress.PointsChange(onBeforeDoneEvent.CellPoints)
	if err != nil {
		return nil, err
	}
	err = progress.EnergyChange(-onBeforeDoneEvent.CellEnergyConsume)
	if err != nil {
		return nil, err
	}
	err = progress.BalanceChange(onBeforeDoneEvent.CellCoins)
	if err != nil {
		return nil, err
	}

	return nil, events.OnAfterDone().Trigger(ctx, &model.OnAfterDoneEvent{
		CurrentCell: currentCell.Data(),
	})
}
