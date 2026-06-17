package outboxes

import (
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
	"fmt"
	"time"
)

type repository interface {
	Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
	GetAndLockNextPending(ctx context.Context) (*model.OutboxInfo, error)
}

type Outboxes struct {
	repository repository
}

func NewOutboxes(repository repository) *Outboxes {
	return &Outboxes{repository: repository}
}

func (o *Outboxes) Start(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			o.processAllPending(ctx)
		}
	}
}

func (o *Outboxes) processAllPending(ctx context.Context) {
	for {
		outbox, err := o.repository.GetAndLockNextPending(ctx)
		if err != nil {
			if errors.Is(err, errs.ErrNoPendingOutbox) {
				break
			}
			break
		}

		o.process(ctx, outbox)
	}
}

func (o *Outboxes) process(ctx context.Context, outbox *model.OutboxInfo) {
	outboxDef, ok := Get(outbox.Type())
	if !ok {
		outbox.SetStatus(model.OutboxStatusFailed)
		_, _ = o.repository.Save(ctx, outbox)
		return
	}

	err := outboxDef.New().Process(ctx, outbox)
	if err != nil {
		outbox.SetStatus(model.OutboxStatusFailed)
	} else {
		outbox.SetStatus(model.OutboxStatusCompleted)
	}

	_, _ = o.repository.Save(ctx, outbox)
}

func (o *Outboxes) Save(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error) {
	outboxDef, ok := Get(outbox.Type())
	if !ok {
		return nil, fmt.Errorf("unknown outbox type: %s", outbox.Type())
	}

	err := outboxDef.New().Verify(ctx, outbox.Payload())
	if err != nil {
		return nil, err
	}

	return o.repository.Save(ctx, outbox)
}
