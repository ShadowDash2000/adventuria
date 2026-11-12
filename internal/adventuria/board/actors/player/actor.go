package player

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type actionsService interface {
	Save(ctx context.Context, action *model.ActionInfo) (*model.ActionInfo, error)
}

type notifyService interface {
	PlayerMovedNotify(ctx context.Context, player *model.Player, moves []*model.MoveResult) error
}

var _ actors.Actor = (*Actor)(nil)

type Actor struct {
	actions actionsService
	notify  notifyService
	player  *model.Player
	events  *model.Events
}

func NewPlayerActor(actions actionsService, notify notifyService, events *model.Events, player *model.Player) *Actor {
	return &Actor{
		actions: actions,
		notify:  notify,
		player:  player,
		events:  events,
	}
}

func (a *Actor) Player() *model.Player {
	return a.player
}

func (a *Actor) Events() *model.Events {
	return a.events
}

func (a *Actor) Ref() actors.ActorRef {
	return actors.ActorRef{
		Id:   a.player.ID(),
		Kind: actors.ActorKindPlayer,
	}
}

func (a *Actor) Position() actors.Position {
	return actors.Position{
		WorldId:     a.player.Progress().CurrentWorld(),
		CellId:      a.player.LastAction().Cell(),
		CellsPassed: a.player.Progress().CellsPassed(),
	}
}

func (a *Actor) ChangeWorld(ctx context.Context, worldID string) error {
	progress := a.player.Progress()
	oldWorldID := progress.CurrentWorld()

	progress.SetCurrentWorld(worldID)
	progress.SetCellsPassed(0)

	return a.events.OnWorldChanged().Trigger(ctx, &model.OnWorldChangedEvent{
		OldWorldId: oldWorldID,
		NewWorldId: worldID,
	})
}

func (a *Actor) ApplyMovement(ctx context.Context, move *model.MoveResult) error {
	progress := a.player.Progress()
	progress.SetCellsPassed(move.TotalSteps)

	newAction, err := model.NewAction(model.ActionCreate{
		Player: a.player.ID(),
		Cell:   move.CurrentCell.ID(),
		Status: model.ActionStatusMove,
	})
	if err != nil {
		return err
	}

	newAction.SetCellsPassed(move.Steps)

	newAction, err = a.actions.Save(ctx, newAction)
	if err != nil {
		return err
	}

	a.player.SetLastAction(newAction)
	progress.SetCanMove(false)

	return a.events.OnAfterMove().Trigger(ctx, &model.OnAfterMoveEvent{
		Steps:          move.Steps,
		TotalSteps:     move.TotalSteps,
		PrevTotalSteps: move.PrevTotalSteps,
		CurrentCell:    move.CurrentCell,
		CurrentWorld:   move.CurrentWorld,
		Laps:           move.Laps,
	})
}

func (a *Actor) LeaveCell(ctx context.Context, cell cells.Cell) error {
	if reaction, ok := cell.(CellLeft); ok {
		return reaction.OnPlayerLeft(ctx, a)
	}

	if reaction, ok := cell.(actors.ActorCellLeft); ok {
		return reaction.OnActorLeft(ctx, a)
	}

	return nil
}

func (a *Actor) ReachCell(ctx context.Context, cell cells.Cell, reached *cells.ReachedContext) error {
	if reaction, ok := cell.(CellReached); ok {
		return reaction.OnPlayerReached(ctx, a, reached)
	}

	if reaction, ok := cell.(actors.ActorCellReached); ok {
		return reaction.OnActorReached(ctx, a, reached)
	}

	return nil
}

func (a *Actor) AfterMovement(ctx context.Context, move *model.MoveResult) error {
	if move.Steps > 0 && move.Laps > 0 {
		return a.events.OnNewLap().Trigger(ctx, &model.OnNewLapEvent{
			Laps: move.Laps,
		})
	}

	return nil
}

func (a *Actor) MovementFinished(ctx context.Context, moves []*model.MoveResult) error {
	return a.notify.PlayerMovedNotify(ctx, a.player, moves)
}
