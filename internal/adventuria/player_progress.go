package adventuria

import "github.com/pocketbase/pocketbase/core"

type PlayerProgress interface {
	core.RecordProxy
	Closable

	Refetch(ctx AppContext) error

	ID() string
	Player() string
	SetPlayer(playerId string)
	Season() string
	SetSeason(seasonId string)
	Points() int
	AddPoints(amount int)
	Balance() int
	AddBalance(amount int)
	DropsInARow() int
	SetDropsInARow(amount int)
	CellsPassed() int
	addCellsPassed(amount int)
	IsInJail() bool
	SetIsInJail(b bool)
	ItemWheelsCount() int
	AddItemWheelsCount(amount int)
	MaxInventorySlots() int
	Stats() (*Stats, error)
	SetStats(stats Stats)

	IsSafeDrop() bool
	CurrentCell() (Cell, bool)
	CurrentCellOrder() int
}
