package outboxes

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

type repository interface {
	Create(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
	Update(ctx context.Context, outbox *model.OutboxInfo) (*model.OutboxInfo, error)
	GetAndLockNextPending(ctx context.Context) (*model.OutboxInfo, error)
	UpdateStatus(ctx context.Context, id string, status model.OutboxStatus) error
}

type Outboxes struct {
	logger     *slog.Logger
	repository repository
}

func NewOutboxes(logger *slog.Logger, repository repository) *Outboxes {
	return &Outboxes{
		logger:     logger,
		repository: repository,
	}
}

func (o *Outboxes) Start(ctx context.Context) {
	go func() {
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
	}()
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

		err = o.process(ctx, outbox)
		if err != nil {
			o.logger.Error("Failed to process outbox", "outbox_id", outbox.ID(), "error", err)
		}
	}
}

func (o *Outboxes) process(ctx context.Context, outbox *model.OutboxInfo) error {
	outboxDef, ok := Get(outbox.Type())
	if !ok {
		return o.repository.UpdateStatus(ctx, outbox.ID(), model.OutboxStatusFailed)
	}

	err := outboxDef.New().Process(ctx, outbox)
	if err != nil {
		return o.repository.UpdateStatus(ctx, outbox.ID(), model.OutboxStatusFailed)
	}

	return o.repository.UpdateStatus(ctx, outbox.ID(), model.OutboxStatusCompleted)
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

	if outbox.IsNew() {
		return o.repository.Create(ctx, outbox)
	}

	return o.repository.Update(ctx, outbox)
}
