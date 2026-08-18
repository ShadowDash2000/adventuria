package reroll

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"
	"errors"
)

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

type activities interface {
	GetRandomIDsByFilter(ctx context.Context, filter model.ActivityFilter) ([]string, error)
}

var _ model.Action = (*Reroll)(nil)

const Type model.ActionType = "reroll"

type Reroll struct {
	actions.ActionBase
	cells      cells
	reviews    reviewsService
	actions    actionsService
	activities activities
}

func NewDef(cells cells, reviews reviewsService, actionsService actionsService, activities activities) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Reroll{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				reviews:    reviews,
				actions:    actionsService,
				activities: activities,
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

	actionState := player.LastAction().State()
	if actionState.ActivityFilter == nil {
		return nil, errs.ErrNoActiveActivityFilter
	}

	review, err := r.reviews.Create(ctx, reviews.CreateInput{
		Comment: req.Comment,
		Score:   req.Score,
	})
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusReroll)
	lastAction.SetReview(review.ID())
	lastAction, err = r.actions.Save(ctx, lastAction)
	if err != nil {
		return nil, err
	}

	newAction, err := model.NewAction(model.ActionCreate{
		Player: player.ID(),
		Cell:   lastAction.Cell(),
		Status: model.ActionStatusNeedToRollWheel,
	})
	if err != nil {
		return nil, err
	}

	ids, err := r.activities.GetRandomIDsByFilter(ctx, *actionState.ActivityFilter)
	if err != nil {
		return nil, err
	}

	actionState.Activities = &model.ActionActivitiesState{
		Ids: ids,
	}
	newAction.SetState(actionState)

	rootActionId := lastAction.RootAction()
	if rootActionId == "" {
		rootActionId = lastAction.ID()
	}
	newAction.SetRootAction(rootActionId)

	player.SetLastAction(newAction)

	return nil, events.OnAfterReroll().Trigger(ctx, &model.OnAfterRerollEvent{})
}
