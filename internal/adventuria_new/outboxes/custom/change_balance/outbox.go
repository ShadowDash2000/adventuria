package change_balance

import (
	"adventuria/internal/adventuria_new/model"
	"adventuria/internal/adventuria_new/outboxes"
	"context"
)

type progress interface {
	ChangeBalance(ctx context.Context, id string, amount int) error
	NotifyChange(ctx context.Context, id string) error
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
	outboxValue, err := c.decodeValue(outbox.Payload())
	if err != nil {
		return err
	}

	err = c.progress.ChangeBalance(ctx, outboxValue.ProgressId, outboxValue.Amount)
	if err != nil {
		return err
	}

	return c.progress.NotifyChange(ctx, outboxValue.ProgressId)
}
