package change_balance

import "context"

func (c *ChangeBalance) Verify(_ context.Context, value string) error {
	_, err := c.decodeValue(value)
	return err
}
