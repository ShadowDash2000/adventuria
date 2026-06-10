package reroll

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

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.Action = (*Reroll)(nil)

const Type model.ActionType = "reroll"

type Reroll struct {
	actions.ActionBase
	cells   cells
	reviews reviews
	actions actionsService
}

func NewActionRerollDef(cells cells, reviews reviews, actionsService actionsService) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Reroll{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				reviews:    reviews,
				actions:    actionsService,
			}
		},
	)
}

func (r *Reroll) CanDo(ctx context.Context, events *model.Events, player *model.Player) bool {
	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if currentCell.Data().CantReroll() {
		return false
	}

	onBeforeRerollCheckEvent := &model.OnBeforeRerollCheckEvent{
		IsRerollBlocked: false,
	}
	err = events.OnBeforeRerollCheck().Trigger(onBeforeRerollCheckEvent)
	if err != nil {
		return false
	}

	if onBeforeRerollCheckEvent.IsRerollBlocked {
		return false
	}

	return player.LastAction().Type() == actions.ActionTypeRollWheel
}

type Request struct {
	Comment string `json:"comment"`
	Score   int    `json:"score"`
}

func (r *Reroll) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	review, err := model.NewReview(uuid.New(), req.Comment, req.Score)
	if err != nil {
		return nil, err
	}
	review, err = r.reviews.Save(ctx, review)
	if err != nil {
		return nil, err
	}

	currentCell, err := r.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return nil, errors.New("current cell is not refreshable")
	}

	lastAction := player.LastAction()
	lastAction.SetType(Type)
	lastAction.SetReview(review.ID())
	_, err = r.actions.Save(ctx, lastAction)
	if err != nil {
		return nil, err
	}

	newAction, err := model.NewAction(uuid.New(), model.ActionCreate{
		Player: player.ID(),
		Cell:   currentCell.Data().ID(),
		Type:   Type,
	})
	if err != nil {
		return nil, err
	}

	player.SetLastAction(newAction)

	err = cellRefreshable.RefreshItems(ctx, events, player)
	if err != nil {
		return nil, err
	}

	return nil, events.OnAfterReroll().Trigger(&model.OnAfterRerollEvent{})
}
