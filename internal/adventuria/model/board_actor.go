package model

import "context"

type ActorKind string

const (
	ActorKindPlayer ActorKind = "player"
	ActorKindNPC    ActorKind = "npc"
)

type ActorRef struct {
	Id   string
	Kind ActorKind
}

type BoardPosition struct {
	WorldId     string
	CellId      string
	CellsPassed int
}

type BoardActor interface {
	Ref() ActorRef
	Position() BoardPosition
	ChangeWorld(ctx context.Context, worldID string) error
	ApplyMovement(ctx context.Context, move *MoveResult) error
	LeaveCell(ctx context.Context, cell Cell) error
	ReachCell(ctx context.Context, cell Cell, reached *ReachedContext) error
	AfterMovement(ctx context.Context, move *MoveResult) error
	MovementFinished(ctx context.Context, moves []*MoveResult) error
}
