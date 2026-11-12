package actors

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/model"
	"context"
)

type ActorKind string

const (
	ActorKindPlayer ActorKind = "player"
	ActorKindNPC    ActorKind = "npc"
)

type ActorRef struct {
	Id   string
	Kind ActorKind
}

type Position struct {
	WorldId     string
	CellId      string
	CellsPassed int
}

type Actor interface {
	Ref() ActorRef
	Position() Position
	ChangeWorld(ctx context.Context, worldID string) error
	ApplyMovement(ctx context.Context, move *model.MoveResult) error
	LeaveCell(ctx context.Context, cell cells.Cell) error
	ReachCell(ctx context.Context, cell cells.Cell, reached *cells.ReachedContext) error
	AfterMovement(ctx context.Context, move *model.MoveResult) error
	MovementFinished(ctx context.Context, moves []*model.MoveResult) error
}

type ActorCellReached interface {
	OnActorReached(ctx context.Context, actor Actor, reached *cells.ReachedContext) error
}

type ActorCellLeft interface {
	OnActorLeft(ctx context.Context, actor Actor) error
}
