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

func (n *Notify) BalanceChangeNotify(ctx context.Context, id string, amount int) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, fmt.Sprintf("balance_change#%s", id), struct {
			Amount int `json:"amount"`
		}{
			Amount: amount,
		})
	})

	return nil
}

func (n *Notify) EnergyChangeNotify(ctx context.Context, id string, amount int) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, fmt.Sprintf("energy_change#%s", id), struct {
			Amount int `json:"amount"`
		}{
			Amount: amount,
		})
	})

	return nil
}

func (n *Notify) ItemWheelsChangeNotify(ctx context.Context, id string, amount int) error {
	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, fmt.Sprintf("item_wheels_change#%s", id), struct {
			Amount int `json:"amount"`
		}{
			Amount: amount,
		})
	})

	return nil
}
