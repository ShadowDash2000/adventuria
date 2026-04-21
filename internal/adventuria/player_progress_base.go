package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"database/sql"
	"errors"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type PlayerProgressBase struct {
	core.BaseRecordProxy

	hookIds []string
}

func NewPlayerProgress(ctx AppContext, playerId, seasonId string) (PlayerProgress, error) {
	p := &PlayerProgressBase{}
	err := p.init(ctx, playerId, seasonId)
	if err != nil {
		return nil, err
	}
	p.bindHooks(ctx)
	return p, nil
}

func (p *PlayerProgressBase) init(ctx AppContext, playerId, seasonId string) error {
	record, err := fetchPlayerProgress(ctx, playerId, seasonId)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if record == nil {
		record = defaultPlayerProgress(playerId, seasonId)
		err = ctx.App.Save(record)
		if err != nil {
			return err
		}
	}

	p.SetProxyRecord(record)

	return nil
}

func defaultPlayerProgress(playerId, seasonId string) *core.Record {
	record := core.NewRecord(GameCollections.Get(schema.CollectionPlayersProgress))
	record.Set(schema.PlayerProgressSchema.Player, playerId)
	record.Set(schema.PlayerProgressSchema.Season, seasonId)
	record.Set(schema.PlayerProgressSchema.MaxInventorySlots, GameSettings.MaxInventorySlots())
	return record
}

func fetchPlayerProgress(ctx AppContext, playerId, seasonId string) (*core.Record, error) {
	var record core.Record
	err := ctx.App.RecordQuery(schema.CollectionPlayersProgress).
		Where(
			dbx.HashExp{
				schema.PlayerProgressSchema.Player: playerId,
				schema.PlayerProgressSchema.Season: seasonId,
			},
		).
		One(&record)
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (p *PlayerProgressBase) bindHooks(ctx AppContext) {
	p.hookIds = make([]string, 1)

	p.hookIds[0] = ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionPlayersProgress).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == p.Id {
			p.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
}

func (p *PlayerProgressBase) Refetch(ctx AppContext) error {
	record, err := ctx.App.FindRecordById(schema.CollectionPlayersProgress, p.Id)
	if err != nil {
		return err
	}
	p.SetProxyRecord(record)
	return nil
}

func (p *PlayerProgressBase) Close(ctx AppContext) {
	ctx.App.OnRecordAfterUpdateSuccess(schema.CollectionPlayersProgress).Unbind(p.hookIds[0])
}

func (p *PlayerProgressBase) ID() string {
	return p.Id
}

func (p *PlayerProgressBase) Player() string {
	return p.GetString(schema.PlayerProgressSchema.Player)
}

func (p *PlayerProgressBase) SetPlayer(playerId string) {
	p.Set(schema.PlayerProgressSchema.Player, playerId)
}

func (p *PlayerProgressBase) Season() string {
	return p.GetString(schema.PlayerProgressSchema.Season)
}

func (p *PlayerProgressBase) SetSeason(seasonId string) {
	p.Set(schema.PlayerProgressSchema.Season, seasonId)
}

func (p *PlayerProgressBase) Points() int {
	return p.GetInt(schema.PlayerProgressSchema.Points)
}

func (p *PlayerProgressBase) AddPoints(amount int) {
	if p.Points()+amount < 0 {
		p.Set(schema.PlayerProgressSchema.Points, 0)
	} else {
		p.Set(schema.PlayerProgressSchema.Points+"+", amount)
	}
}

func (p *PlayerProgressBase) Balance() int {
	return p.GetInt(schema.PlayerProgressSchema.Balance)
}

func (p *PlayerProgressBase) AddBalance(amount int) {
	if p.Balance()+amount < 0 {
		p.Set(schema.PlayerProgressSchema.Balance, 0)
	} else {
		p.Set(schema.PlayerProgressSchema.Balance+"+", amount)
	}
}

func (p *PlayerProgressBase) DropsInARow() int {
	return p.GetInt(schema.PlayerProgressSchema.DropsInARow)
}

func (p *PlayerProgressBase) SetDropsInARow(amount int) {
	p.Set(schema.PlayerProgressSchema.DropsInARow, amount)
}

func (p *PlayerProgressBase) CellsPassed() int {
	return p.GetInt(schema.PlayerProgressSchema.CellsPassed)
}

func (p *PlayerProgressBase) addCellsPassed(amount int) {
	p.Set(schema.PlayerProgressSchema.CellsPassed+"+", amount)
}

func (p *PlayerProgressBase) IsInJail() bool {
	return p.GetBool(schema.PlayerProgressSchema.IsInJail)
}

func (p *PlayerProgressBase) SetIsInJail(b bool) {
	p.Set(schema.PlayerProgressSchema.IsInJail, b)
}

func (p *PlayerProgressBase) ItemWheelsCount() int {
	return p.GetInt(schema.PlayerProgressSchema.ItemWheelsCount)
}

func (p *PlayerProgressBase) AddItemWheelsCount(amount int) {
	if p.ItemWheelsCount()+amount < 0 {
		p.Set(schema.PlayerProgressSchema.ItemWheelsCount, 0)
	} else {
		p.Set(schema.PlayerProgressSchema.ItemWheelsCount+"+", amount)
	}
}

func (p *PlayerProgressBase) MaxInventorySlots() int {
	return p.GetInt(schema.PlayerProgressSchema.MaxInventorySlots)
}

func (p *PlayerProgressBase) Stats() (*Stats, error) {
	var stats Stats
	err := p.UnmarshalJSONField(schema.PlayerProgressSchema.Stats, &stats)
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (p *PlayerProgressBase) SetStats(stats Stats) {
	p.Set(schema.PlayerProgressSchema.Stats, stats)
}

func (p *PlayerProgressBase) IsSafeDrop() bool {
	return p.DropsInARow() < GameSettings.DropsToJail()
}

func (p *PlayerProgressBase) CurrentCell() (Cell, bool) {
	currentCellNum := p.CellsPassed() % GameCells.Count()
	cell, ok := GameCells.GetByOrder(currentCellNum)
	return cell, ok
}

func (p *PlayerProgressBase) CurrentCellOrder() int {
	return mod(p.CellsPassed(), GameCells.Count())
}
