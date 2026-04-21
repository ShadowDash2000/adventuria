package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"iter"
	"time"
)

type Players struct {
	players *MemoryCacheWithClose[string, Player]
}

func NewPlayers(ctx AppContext) *Players {
	p := &Players{
		players: NewMemoryCacheWithClose[string, Player](ctx, time.Hour, false),
	}
	return p
}

func (p *Players) GetByID(ctx AppContext, playerId string) (Player, error) {
	player, ok := p.players.Get(playerId)
	if ok {
		return player, nil
	}

	player, err := NewPlayer(ctx, playerId)
	if err != nil {
		return nil, err
	}

	p.players.Set(playerId, player)
	return player, nil
}

func (p *Players) GetByName(ctx AppContext, playerName string) (Player, error) {
	for _, player := range p.players.GetAll() {
		if playerName == player.Name() {
			return player, nil
		}
	}

	player, err := NewPlayerFromName(ctx, playerName)
	if err != nil {
		return nil, err
	}

	p.players.Set(player.ID(), player)
	return player, nil
}

func (p *Players) GetAll(ctx AppContext) (iter.Seq2[string, Player], error) {
	var players []struct {
		Id string `db:"id"`
	}
	err := ctx.App.RecordQuery(schema.CollectionPlayers).
		Select(schema.PlayerSchema.Id).
		All(&players)
	if err != nil {
		return nil, err
	}

	for _, player := range players {
		_, ok := p.players.Get(player.Id)
		if ok {
			continue
		}

		player, err := NewPlayer(ctx, player.Id)
		if err != nil {
			return nil, err
		}

		p.players.Set(player.ID(), player)
	}

	return p.GetAllInMemory(), nil
}

func (p *Players) GetAllInMemory() iter.Seq2[string, Player] {
	return p.players.GetAll()
}

func (p *Players) RefetchAllInMemory(ctx AppContext) []error {
	var errs []error
	for _, player := range p.GetAllInMemory() {
		player.Lock()
		err := player.Refetch(ctx)
		if err != nil {
			errs = append(errs, err)
		}
		player.Unlock()
	}
	return errs
}
