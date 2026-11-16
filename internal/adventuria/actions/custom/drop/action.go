package drop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"
	"errors"
)

type dropResultCalculator interface {
	Calculate(ctx context.Context, actor *player.Actor, mode model.EffectRunMode) (*model.ActivityDropResult, error)
}

type actionsService interface {
	CanDo(ctx context.Context, actor *player.Actor, t model.ActionType) bool
}

type reviewsService interface {
	Create(ctx context.Context, input reviews.CreateInput) (*model.Review, error)
}

type players interface {
	Save(ctx context.Context, player *model.Player) error
}

type settings interface {
	GetFirst(ctx context.Context) (*model.Settings, error)
}

type boardService interface {
	MoveToClosestCellType(ctx context.Context, actor actors.Actor, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ actions.Action = (*Drop)(nil)

const Type model.ActionType = "drop"

type Drop struct {
	actions.ActionBase
	dropResultCalculator dropResultCalculator
	actions              actionsService
	reviews              reviewsService
	players              players
	settings             settings
	board                boardService
}

func NewDef(
	dropResultCalculator dropResultCalculator,
	actionsService actionsService,
	reviews reviewsService,
	players players,
	settings settings,
	board boardService,
) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &Drop{
				ActionBase:           actions.NewActionBase(Type),
				dropResultCalculator: dropResultCalculator,
				actions:              actionsService,
				reviews:              reviews,
				players:              players,
				settings:             settings,
				board:                board,
			}
		},
	)
}

func (d *Drop) CanDo(ctx context.Context, actor *player.Actor) bool {
	if !d.actions.CanDo(ctx, actor, actions.ActionTypeCompleteActivity) {
		return false
	}

	dropResult, err := d.dropResultCalculator.Calculate(ctx, actor, model.EffectRunModePreview)
	if err != nil {
		return false
	}

	if dropResult.IsBlocked() {
		return false
	}

	return dropResult.EnergyConsume() <= actor.Player().Progress().Energy()
}

func (d *Drop) Do(ctx context.Context, actor *player.Actor, actionReq actions.Request) (any, error) {
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

	settings, err := d.settings.GetFirst(ctx)
	if err != nil {
		return nil, err
	}

	dropResult, err := d.dropResultCalculator.Calculate(ctx, actor, model.EffectRunModeApply)
	if err != nil {
		return nil, err
	}

	lastAction := actor.Player().LastAction()
	lastAction.SetStatus(model.ActionStatusDrop)
	lastAction.SetReview(review.ID())
	err = d.players.Save(ctx, actor.Player())
	if err != nil {
		return nil, err
	}

	progress := actor.Player().Progress()
	progress.SetCanMove(true)

	if !dropResult.IsSafe() {
		err = progress.DropsInARowChange(1)
		if err != nil {
			return nil, err
		}

		if progress.DropsInARow() >= settings.DropsToJail() {
			progress.SetIsInJail(true)

			_, err = d.board.MoveToClosestCellType(ctx, actor, cells.CellTypeJail)
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, actor.Events().OnAfterDrop().Trigger(ctx, &model.OnAfterDropEvent{})
}
