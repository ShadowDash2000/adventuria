package change_balance

import "context"

func (c *ChangeBalance) Verify(_ context.Context, payload string) error {
	_, err := c.decodePayload(payload)
	return err
}
