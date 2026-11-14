package walk

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/npc"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/npcs/behaviors"
	"bytes"
	"context"
	"math/rand/v2"
)

type board interface {
	Move(ctx context.Context, actor actors.Actor, steps int, moveType model.MoveType) ([]*model.MoveResult, error)
}

var _ behaviors.Behavior = (*Walk)(nil)

const Type model.NPCBehaviorType = "walk"

type Walk struct {
	behaviors.BehaviorBase
	board   board
	payload Payload
}

func NewDef(board board) behaviors.NPCBehaviorDef {
	return behaviors.NewNPCBehavior(
		Type,
		func(data, _ []byte) (behaviors.Behavior, error) {
			return NewFromJSON(board, data)
		},
	)
}

func NewFromJSON(board board, data []byte) (*Walk, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &Walk{
		BehaviorBase: behaviors.NewBehaviorBase(
			bytes.Clone(data),
			nil,
		),
		board:   board,
		payload: payload,
	}, nil
}

func (w *Walk) Move(ctx context.Context, actor *npc.Actor) (actors.Position, error) {
	steps := rand.IntN(w.payload.MaxDistance) + w.payload.MinDistance

	_, err := w.board.Move(ctx, actor, steps, model.MoveTypePath)
	if err != nil {
		return actors.Position{}, err
	}

	return actor.Position(), nil
}
