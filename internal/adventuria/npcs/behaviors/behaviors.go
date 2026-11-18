package behaviors

import (
	"adventuria/internal/adventuria/board/actors/npc"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"fmt"
)

type npcFactory interface {
	New(npc *model.NPC) *npc.Actor
}

type Behaviors struct {
	npcFactory npcFactory
}

func NewBehaviors(npcFactory npcFactory) *Behaviors {
	return &Behaviors{
		npcFactory: npcFactory,
	}
}

func (b *Behaviors) Move(ctx context.Context, npc *model.NPC, behavior *model.NPCBehavior) (model.ActorPosition, error) {
	behaviorDef, ok := Get(behavior.Type())
	if !ok {
		return model.ActorPosition{}, fmt.Errorf("%w: %s", errs.ErrUnknownNPCBehaviorType, behavior.Type())
	}

	bh, err := behaviorDef.New(behavior)
	if err != nil {
		return model.ActorPosition{}, err
	}

	actor := b.npcFactory.New(npc)
	position, err := bh.Move(ctx, actor)
	if err != nil {
		return model.ActorPosition{}, err
	}

	return position, nil
}
