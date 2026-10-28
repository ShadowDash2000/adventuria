package game_state

import (
	"adventuria/internal/adventuria/model"
	"time"
)

type view struct {
	Id               string    `json:"id"`
	Disabled         bool      `json:"disabled"`
	Debug            bool      `json:"debug"`
	Season           string    `json:"season"`
	NextDayAt        time.Time `json:"next_day_at"`
	PlayerProgressId string    `json:"player_progress_id"`
	CurrentWorld     string    `json:"current_world"`
	Balance          int       `json:"balance"`
	Energy           int       `json:"energy"`
	DropsInARow      int       `json:"drops_in_a_row"`
	ItemWheelsCount  int       `json:"item_wheels_count"`
}

func stateToView(
	currentSeason string,
	nextDayAt time.Time,
	player *model.PlayerInfo,
	progress *model.PlayerProgress,
) *view {
	v := &view{
		Season:    currentSeason,
		NextDayAt: nextDayAt,
	}

	if player != nil {
		v.Id = player.ID()
		v.Disabled = player.Disabled()
		v.Debug = player.Debug()
	}

	if progress != nil {
		v.PlayerProgressId = progress.ID()
		v.CurrentWorld = progress.CurrentWorld()
		v.Balance = progress.Balance()
		v.Energy = progress.Energy()
		v.DropsInARow = progress.DropsInARow()
		v.ItemWheelsCount = progress.ItemWheelsCount()
	}

	return v
}
