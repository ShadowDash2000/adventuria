package change_energy

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"context"
)

type progress interface {
	ChangeEnergy(ctx context.Context, id string, amount int) error
}

var _ model.Outbox = (*ChangeEnergy)(nil)

const Type model.OutboxType = "change_energy"

type ChangeEnergy struct {
	outboxes.OutboxBase
	progress progress
}

func NewDef(progress progress) outboxes.OutboxDef {
	return outboxes.NewOutbox(
		Type,
		func() model.Outbox {
			return &ChangeEnergy{
				OutboxBase: outboxes.NewOutboxBase(Type),
				progress:   progress,
			}
		},
	)
}

func (c *ChangeEnergy) Process(ctx context.Context, outbox *model.OutboxInfo) error {
	payload, err := c.decodePayload(outbox.Payload())
	if err != nil {
		return err
	}

	return c.progress.ChangeEnergy(ctx, payload.ProgressId, payload.Amount)
}
