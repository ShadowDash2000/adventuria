package drop

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

type players interface {
	Save(ctx context.Context, player *model.Player) error
}

type settings interface {
	GetFirstOrDefault(ctx context.Context) (*model.Settings, error)
}

type board interface {
	MoveToClosestCellType(ctx context.Context, events *model.Events, player *model.Player, cellType model.CellType) ([]*model.MoveResult, error)
}

var _ model.Action = (*Drop)(nil)

const Type model.ActionType = "drop"

type Drop struct {
	actions.ActionBase
	cells    cells
	reviews  reviews
	players  players
	settings settings
	board    board
}

func NewDef(cells cells, reviews reviews, players players, settings settings, board board) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Drop{
				ActionBase: actions.NewActionBase(Type),
				cells:      cells,
				reviews:    reviews,
				players:    players,
				settings:   settings,
				board:      board,
			}
		},
	)
}

func (d *Drop) CanDo(ctx context.Context, events *model.Events, player *model.Player) bool {
	currentCell, err := d.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return false
	}

	if currentCell.Data().CantDrop() {
		return false
	}

	if player.Progress().IsInJail() {
		return false
	}

	onBeforeDropCheckEvent := &model.OnBeforeDropCheckEvent{
		IsDropBlocked: false,
	}
	err = events.OnBeforeDropCheck().Trigger(ctx, onBeforeDropCheckEvent)
	if err != nil {
		return false
	}

	if onBeforeDropCheckEvent.IsDropBlocked {
		return false
	}

	return player.LastAction().Type() == actions.ActionTypeRollWheel
}

type Request struct {
	Comment string `json:"comment"`
	Score   int    `json:"score"`
}

func (d *Drop) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
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

	settings, err := d.settings.GetFirstOrDefault(ctx)
	if err != nil {
		return nil, err
	}

	currentCell, err := d.cells.GetCurrentCellByProgress(ctx, player.Progress())
	if err != nil {
		return nil, err
	}

	onBeforeDropEvent := &model.OnBeforeDropEvent{
		IsSafeDrop:    false,
		IsDropBlocked: false,
		PointsForDrop: settings.PointsForDrop(),
	}
	err = events.OnBeforeDrop().Trigger(ctx, onBeforeDropEvent)
	if err != nil {
		return nil, err
	}

	if onBeforeDropEvent.IsDropBlocked {
		return nil, errors.New("drop is not allowed")
	}

	lastAction := player.LastAction()
	lastAction.SetType(Type)
	lastAction.SetReview(review.ID())
	lastAction.SetCanMove(true)
	err = d.players.Save(ctx, player)
	if err != nil {
		return nil, err
	}

	if !onBeforeDropEvent.IsSafeDrop && !currentCell.Data().IsSafeDrop() {
		progress := player.Progress()
		err = progress.PointsChange(onBeforeDropEvent.PointsForDrop)
		if err != nil {
			return nil, err
		}
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
