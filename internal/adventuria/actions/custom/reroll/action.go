package reroll

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"
	"errors"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
	GetByPlayerWrapped(ctx context.Context, player *model.Player) (model.Cell, error)
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.Action = (*Reroll)(nil)

const Type model.ActionType = "reroll"

type Reroll struct {
	actions.ActionBase
	cells   cells
	reviews reviewsService
	actions actionsService
}

func NewDef(cells cells, reviews reviewsService, actionsService actionsService) actions.ActionDef {
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
	if !r.actions.CanDo(ctx, events, player, actions.ActionTypeCompleteActivity) {
		return false
	}

	currentCell, err := r.cells.GetByPlayer(ctx, player)
	if err != nil {
		return false
	}

	if currentCell.CantReroll() {
		return false
	}

	onBeforeRerollCheckEvent := &model.OnBeforeRerollCheckEvent{
		IsRerollBlocked: false,
	}
	err = events.OnBeforeRerollCheck().Trigger(ctx, onBeforeRerollCheckEvent)
	if err != nil {
		return false
	}

	if onBeforeRerollCheckEvent.IsRerollBlocked {
		return false
	}

	return true
}

type Request struct {
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

func (r *Reroll) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	review, err := r.reviews.Create(ctx, reviews.CreateInput{
		Comment: req.Comment,
		Score:   req.Score,
	})
	if err != nil {
		return nil, err
	}

	currentCell, err := r.cells.GetByPlayerWrapped(ctx, player)
	if err != nil {
		return nil, err
	}

	cellRefreshable, ok := currentCell.(model.Refreshable)
	if !ok {
		return nil, errors.New("current cell is not refreshable")
	}

	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusReroll)
	lastAction.SetReview(review.ID())
	_, err = r.actions.Save(ctx, lastAction)
	if err != nil {
		return nil, err
	}

	newAction, err := model.NewAction(model.ActionCreate{
		Player: player.ID(),
		Cell:   currentCell.Data().ID(),
		Status: model.ActionStatusNeedToRollWheel,
	})
	if err != nil {
		return nil, err
	}

	newAction.SetState(lastAction.State())
	player.SetLastAction(newAction)

	err = cellRefreshable.RefreshItems(ctx, events, player)
	if err != nil {
		return nil, err
	}

	return nil, events.OnAfterReroll().Trigger(ctx, &model.OnAfterRerollEvent{})
}
