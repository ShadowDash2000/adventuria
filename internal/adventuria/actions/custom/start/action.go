package start

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Action = (*Start)(nil)

const Type model.ActionType = "start"

type Start struct {
	actions.ActionBase
}

func NewDef() actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Start{
				ActionBase: actions.NewActionBase(Type),
			}
		},
	)
}

func (s *Start) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Status() == model.ActionStatusNone
}

func (s *Start) Do(_ context.Context, _ *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	player.LastAction().SetStatus(model.ActionStatusStart)
	return nil, nil
}
