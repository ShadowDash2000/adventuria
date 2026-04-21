package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/event"
	"adventuria/pkg/result"
)

type StatsEffect struct{}

func (ef *StatsEffect) Subscribe(player adventuria.Player) []event.Unsubscribe {
	return []event.Unsubscribe{
		player.OnAfterReroll().BindFunc(func(e *adventuria.OnAfterRerollEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.Rerolls++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterDrop().BindFunc(func(e *adventuria.OnAfterDropEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.Drops++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterGoToJail().BindFunc(func(e *adventuria.OnAfterGoToJailEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.WasInJail++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterDone().BindFunc(func(e *adventuria.OnAfterDoneEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.Finished++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterRoll().BindFunc(func(e *adventuria.OnAfterRollEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.DiceRolls++
			if e.N > stats.MaxDiceRoll {
				stats.MaxDiceRoll = e.N
			}
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterWheelRoll().BindFunc(func(e *adventuria.OnAfterWheelRollEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.WheelRolled++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
		player.OnAfterItemUse().BindFunc(func(e *adventuria.OnAfterItemUseEvent) (*result.Result, error) {
			stats, err := player.Progress().Stats()
			if err != nil {
				return result.Err("internal error: failed to get player stats"), err
			}
			stats.ItemsUsed++
			player.Progress().SetStats(*stats)
			return e.Next()
		}),
	}
}
