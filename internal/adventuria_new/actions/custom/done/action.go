package done

import (
	"adventuria/internal/adventuria_new/actions"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"

	"github.com/google/uuid"
)

type cells interface {
	GetCurrentCellByProgress(ctx context.Context, progress *model.PlayerProgress) (model.Cell, error)
}

type reviews interface {
	Save(ctx context.Context, review *model.Review) (*model.Review, error)
}

var _ model.Action = (*Done)(nil)

const Type model.ActionType = "done"

type Done struct {
	actions.ActionBase
	cells   cells
	reviews reviews
}

func NewActionDoneDef(cells cells, reviews reviews) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Done{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				reviews:    reviews,
			}
		},
	)
}

func (d *Done) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Type() == actions.ActionTypeRollWheel
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

	review, err := model.NewReview(uuid.New(), req.Comment, req.Score)
	if err != nil {
		return nil, err
	}
	review, err = d.reviews.Save(ctx, review)
	if err != nil {
		return nil, err
	}

	currentCell, err := d.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	onBeforeDoneEvent := &model.OnBeforeDoneEvent{
		CellPoints: currentCell.Data().Points(),
		CellCoins:  currentCell.Data().Coins(),
	}
	err = events.OnBeforeDone().Trigger(onBeforeDoneEvent)
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetType(Type)
	lastAction.SetReview(review.ID())
	lastAction.SetCanMove(true)

	progress := player.Progress()
	progress.SetDropsInARow(0)
	progress.SetIsInJail(false)
	err = progress.PointsChange(onBeforeDoneEvent.CellPoints)
	if err != nil {
		return nil, err
	}
	err = progress.BalanceChange(onBeforeDoneEvent.CellCoins)
	if err != nil {
		return nil, err
	}

	return nil, events.OnAfterDone().Trigger(&model.OnAfterDoneEvent{})
}
