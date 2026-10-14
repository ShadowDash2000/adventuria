package repository

import (
	"adventuria/pkg/pbrealtime"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type Notify struct {
	pb core.App
}

func NewRepositoryNotify(pb core.App) *Notify {
	return &Notify{pb: pb}
}

func (n *Notify) BalanceChangeNotify(ctx context.Context, id string, balance int) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, fmt.Sprintf("balance_change#%s", id), struct {
			Balance int `json:"balance"`
		}{
			Balance: balance,
		})
	})

	return nil
}

func (n *Notify) EnergyChangeNotify(ctx context.Context, id string, energy int) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, fmt.Sprintf("energy_change#%s", id), struct {
			Energy int `json:"energy"`
		}{
			Energy: energy,
		})
	})

	return nil
}
