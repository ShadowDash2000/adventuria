package nothing

import (
	"adventuria/internal/adventuria/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*Nothing)(nil)

func (n *Nothing) Verify(_ context.Context, value string) error {
	if _, ok := useEvents[value]; !ok {
		return errors.New("unknown event")
	}
	return nil
}
