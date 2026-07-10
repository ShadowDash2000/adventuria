package balance_change

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*BalanceChange)(nil)

func (b *BalanceChange) Verify(_ context.Context, value string) error {
	_, err := b.decodeValue(value)
	return err
}
