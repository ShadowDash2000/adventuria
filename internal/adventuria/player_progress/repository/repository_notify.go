package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
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

type newActivePlayerNotification struct {
	Id       string                  `json:"id"`
	Name     string                  `json:"name"`
	Avatar   string                  `json:"avatar"`
	Progress newActivePlayerProgress `json:"progress"`
}

type newActivePlayerProgress struct {
	Id string `json:"id"`
}

func (n *Notify) NewActivePlayerNotify(ctx context.Context, playerInfo *model.PlayerInfo, progress *model.PlayerProgress) error {
	data := newActivePlayerNotification{
		Id:     playerInfo.ID(),
		Name:   playerInfo.Name(),
		Avatar: pbhelper.BuildFilePath(schema.CollectionPlayers, playerInfo.ID(), playerInfo.Avatar()),
		Progress: newActivePlayerProgress{
			Id: progress.ID(),
		},
	}

	pbtransaction.AfterCommit(ctx, n.pb, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		_ = pbrealtime.Notify(ctx, n.pb, "new_active_player", data)
	})

	return nil
}
