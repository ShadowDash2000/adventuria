package repository

import (
	"adventuria/pkg/pbrealtime"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

type Notify struct {
	pb core.App
}

func NewRepositoryNotify(pb core.App) *Notify {
	return &Notify{pb: pb}
}

func (n *Notify) BalanceChangeNotify(ctx context.Context, id string, balance int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, n.pb)

	return pbrealtime.Notify(ctx, pb, fmt.Sprintf("balance_change#%s", id), struct {
		Balance int `json:"balance"`
	}{
		Balance: balance,
	})
}

func (n *Notify) EnergyChangeNotify(ctx context.Context, id string, energy int) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, n.pb)

	return pbrealtime.Notify(ctx, pb, fmt.Sprintf("energy_change#%s", id), struct {
		Energy int `json:"energy"`
	}{
		Energy: energy,
	})
}
