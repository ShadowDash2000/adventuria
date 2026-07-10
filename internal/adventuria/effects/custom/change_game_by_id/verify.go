package change_game_by_id

import (
	"adventuria/internal/adventuria/model"
	"context"
)

var _ model.Verifiable = (*ChangeGameById)(nil)

func (c *ChangeGameById) Verify(ctx context.Context, value string) error {
	_, err := c.activities.GetByID(ctx, value)
	if err != nil {
		return err
	}
	return nil
}
