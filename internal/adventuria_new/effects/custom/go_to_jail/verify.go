package go_to_jail

import (
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
)

var _ model.Verifiable = (*GoToJail)(nil)

func (g *GoToJail) Verify(_ context.Context, value string) error {
	if _, ok := useEvents[value]; !ok {
		return errors.New("unknown event")
	}
	return nil
}
