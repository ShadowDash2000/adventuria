package start

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/model"
	"context"
)

type actionsService interface {
	CreateAndSetLastAction(ctx context.Context, player *model.Player, action *model.ActionInfo) (*model.ActionInfo, error)
}

var _ model.Action = (*Start)(nil)

const Type model.ActionType = "start"

type Start struct {
	actions.ActionBase
	actions actionsService
}

func NewDef(actionsService actionsService) actions.ActionDef {
	return actions.NewAction(
		Type,
		func() model.Action {
			return &Start{
				ActionBase: actions.NewActionBase(Type),
				actions:    actionsService,
			}
		},
	)
}

func (s *Start) CanDo(_ context.Context, _ *model.Events, player *model.Player) bool {
	return player.LastAction().Status() == model.ActionStatusNone
}

func (s *Start) Do(ctx context.Context, _ *model.Events, player *model.Player, _ model.ActionRequest) (any, error) {
	lastAction := player.LastAction()
	lastAction.SetStatus(model.ActionStatusStart)
	_, err := s.actions.CreateAndSetLastAction(ctx, player, lastAction)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
