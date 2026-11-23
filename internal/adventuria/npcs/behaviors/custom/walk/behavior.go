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
	MovePath(ctx context.Context, actor actors.Actor, steps int) ([]*model.MoveResult, error)
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

func (w *Walk) Move(ctx context.Context, actor *npc.Actor) (model.ActorPosition, error) {
	steps := rand.IntN(w.payload.MaxDistance-w.payload.MinDistance+1) + w.payload.MinDistance

	_, err := w.board.MovePath(ctx, actor, steps)
	if err != nil {
		return model.ActorPosition{}, err
	}

	return actor.Position(), nil
}
