package replace_dice_roll

import (
	"adventuria/internal/adventuria_new/model"
	"context"
)

var _ model.Verifiable = (*ReplaceDiceRoll)(nil)

func (r *ReplaceDiceRoll) Verify(_ context.Context, value string) error {
	_, err := r.decodeValue(value)
	return err
}
