package npc

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type notifyService interface {
	ActorMovedNotify(ctx context.Context, actorRef actors.ActorRef, moves []*model.MoveResult) error
}

var _ actors.Actor = (*Actor)(nil)

type Actor struct {
	notify notifyService
	npc    *model.NPC
}

func NewActor(notify notifyService, npc *model.NPC) *Actor {
	return &Actor{
		notify: notify,
		npc:    npc,
	}
}

func (a *Actor) Ref() actors.ActorRef {
	return actors.ActorRef{
		Id:   a.npc.ID(),
		Kind: actors.ActorKindNPC,
	}
}

func (a *Actor) Position() model.ActorPosition {
	return model.ActorPosition{
		WorldId:     a.npc.CurrentWorld(),
		CellId:      a.npc.CurrentCell(),
		CellsPassed: a.npc.CellsPassed(),
	}
}

func (a *Actor) ChangeWorld(_ context.Context, worldId string) error {
	a.npc.SetCurrentWorld(worldId)
	a.npc.SetCellsPassed(0)
	return nil
}

func (a *Actor) ApplyMovement(_ context.Context, move *model.MoveResult) error {
	a.npc.SetCurrentWorld(move.CurrentWorld.ID())
	a.npc.SetCurrentCell(move.CurrentCell.ID())
	a.npc.SetCellsPassed(move.CurrentCell.LocalOrder())
	return nil
}

func (a *Actor) LeaveCell(ctx context.Context, cell cells.Cell) error {
	if reaction, ok := cell.(actors.ActorCellLeft); ok {
		return reaction.OnActorLeft(ctx, a)
	}

	return nil
}

func (a *Actor) ReachCell(ctx context.Context, cell cells.Cell, reached *cells.ReachedContext) error {
	if reaction, ok := cell.(actors.ActorCellReached); ok {
		return reaction.OnActorReached(ctx, a, reached)
	}

	return nil
}

func (a *Actor) AfterMovement(_ context.Context, _ *model.MoveResult) error {
	return nil
}

func (a *Actor) MovementFinished(ctx context.Context, moves []*model.MoveResult) error {
	return a.notify.ActorMovedNotify(ctx, a.Ref(), moves)
}
