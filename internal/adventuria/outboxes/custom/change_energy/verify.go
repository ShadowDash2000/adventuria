package change_energy

import "context"

func (c *ChangeEnergy) Verify(_ context.Context, payload string) error {
	_, err := c.decodePayload(payload)
	return err
}
