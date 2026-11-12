package start

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/model"
	"context"
)

var _ actions.Action = (*Start)(nil)

const Type model.ActionType = "start"

type Start struct {
	actions.ActionBase
}

func NewDef() actions.ActionDef {
	return actions.NewAction(
		Type,
		func() actions.Action {
			return &Start{
				ActionBase: actions.NewActionBase(Type),
			}
		},
	)
}

func (s *Start) CanDo(_ context.Context, actor *player.Actor) bool {
	return actor.Player().LastAction().Status() == model.ActionStatusNone
}

func (s *Start) Do(_ context.Context, actor *player.Actor, _ actions.Request) (any, error) {
	actor.Player().LastAction().SetStatus(model.ActionStatusStart)
	actor.Player().Progress().SetCanMove(true)

	return nil, nil
}
