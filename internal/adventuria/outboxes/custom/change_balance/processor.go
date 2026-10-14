package change_balance

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/outboxes"
	"context"
)

type progress interface {
	ChangeBalance(ctx context.Context, id string, amount int) error
}

const Type model.OutboxType = "change_balance"

type ChangeBalanceProcessor struct {
	progress progress
}

func NewDef(progress progress) outboxes.OutboxDef {
	processor := ChangeBalanceProcessor{
		progress: progress,
	}

	return outboxes.NewOutbox(
		Type,
		func(ctx context.Context, payload string) error {
			task, err := NewFromJSON(payload)
			if err != nil {
				return err
			}

			return processor.Process(ctx, task)
		},
	)
}

func (c *ChangeBalanceProcessor) Process(ctx context.Context, task *ChangeBalance) error {
	return c.progress.ChangeBalance(ctx, task.payload.ProgressId, *task.payload.Amount)
}
