package item_received

import "context"

func (i *ItemReceived) Verify(_ context.Context, payload string) error {
	_, err := i.decodePayload(payload)
	return err
}
