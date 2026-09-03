package move_to_cell_id

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

type board interface {
	MoveToCellId(ctx context.Context, events *model.Events, player *model.Player, cellId string) ([]*model.MoveResult, error)
}

var _ model.Action = (*MoveToCellId)(nil)

const Type model.ActionType = "move_to_cell_id"

type MoveToCellId struct {
	actions.ActionBase
	board board
}

func NewDef(board board) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &MoveToCellId{
				ActionBase: actions.NewActionBase(Type),
				board:      board,
			}
		},
	)
}

func (m *MoveToCellId) CanDo(_ context.Context, _ *model.Events, _ *model.Player) bool {
	return true
}

type Request struct {
	CellId string `json:"cell_id"`
}

func (m *MoveToCellId) Do(ctx context.Context, events *model.Events, player *model.Player, actionReq model.ActionRequest) (any, error) {
	req, ok := actionReq.(Request)
	if !ok {
		return nil, errors.New("invalid request")
	}

	_, err := m.board.MoveToCellId(ctx, events, player, req.CellId)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
