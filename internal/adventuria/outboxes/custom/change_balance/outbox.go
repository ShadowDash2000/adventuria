package change_balance

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"context"
)

type progress interface {
	ChangeBalance(ctx context.Context, id string, amount int) error
}

var _ model.Outbox = (*ChangeBalance)(nil)

const Type model.OutboxType = "change_balance"

type ChangeBalance struct {
	outboxes.OutboxBase
	progress progress
}

func NewDef(progress progress) outboxes.OutboxDef {
	return outboxes.NewOutbox(
		Type,
		func() model.Outbox {
			return &ChangeBalance{
				OutboxBase: outboxes.NewOutboxBase(Type),
				progress:   progress,
			}
		},
	)
}

func (c *ChangeBalance) Process(ctx context.Context, outbox *model.OutboxInfo) error {
	payload, err := c.decodePayload(outbox.Payload())
	if err != nil {
		return err
	}

	return c.progress.ChangeBalance(ctx, payload.ProgressId, payload.Amount)
}
