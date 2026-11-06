package drop

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"context"
	"errors"
)

type dropResultCalculator interface {
	Calculate(ctx context.Context, events *model.Events, cell *model.CellInfo, mode model.EffectRunMode) (*model.ActivityDropResult, error)
}

type actionsService interface {
	CanDo(ctx context.Context, events *model.Events, player *model.Player, t model.ActionType) bool
}

type cells interface {
	GetByPlayer(ctx context.Context, player *model.Player) (*model.CellInfo, error)
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

type board interface {
	MoveToClosestCellType(ctx context.Context, events *model.Events, player *model.Player, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ model.Action = (*Drop)(nil)

const Type model.ActionType = "drop"

type Drop struct {
	actions.ActionBase
	dropResultCalculator dropResultCalculator
	actions              actionsService
	cells                cells
	reviews              reviewsService
	players              players
	settings             settings
	board                board
}

func NewDef(
	dropResultCalculator dropResultCalculator,
	actionsService actionsService,
	cells cells,
	reviews reviewsService,
	players players,
	settings settings,
	board board,
) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Drop{
				ActionBase:           actions.NewActionBase(Type),
				dropResultCalculator: dropResultCalculator,
				actions:              actionsService,
				cells:                cells,
				reviews:              reviews,
				players:              players,
				settings:             settings,
				board:                board,
			}
		},
	)
}

func (d *Drop) CanDo(ctx context.Context, events *model.Events, player *model.Player) bool {
	if !d.actions.CanDo(ctx, events, player, actions.ActionTypeCompleteActivity) {
		return false
	}

	currentCell, err := d.cells.GetByPlayer(ctx, player)
	if err != nil {
		return false
	}

	dropResult, err := d.dropResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return false
	}

	if dropResult.IsBlocked() {
		return false
	}

	return true
}

type Request struct {
	Comment string  `json:"comment"`
	Score   float64 `json:"score"`
}

func (d *Drop) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
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

	currentCell, err := d.cells.GetByPlayer(ctx, player)
	if err != nil {
		return nil, err
	}

	dropResult, err := d.dropResultCalculator.Calculate(ctx, events, currentCell, model.EffectRunModePreview)
	if err != nil {
		return nil, err
	}

	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusDrop)
	lastAction.SetReview(review.ID())
	err = d.players.Save(ctx, player)
	if err != nil {
		return nil, err
	}

	progress := player.Progress()
	progress.SetCanMove(true)

	if !dropResult.IsSafe() {
		err = progress.DropsInARowChange(1)
		if err != nil {
			return nil, err
		}

		if progress.DropsInARow() >= settings.DropsToJail() {
			progress.SetIsInJail(true)

			_, err = d.board.MoveToClosestCellType(ctx, events, player, "jail")
			if err != nil {
				return nil, err
			}
		}
	}

	return nil, events.OnAfterDrop().Trigger(ctx, &model.OnAfterDropEvent{})
}
