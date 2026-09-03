package game_state

import "adventuria/internal/adventuria/model"

type view struct {
	Id              string `json:"id"`
	Disabled        bool   `json:"disabled"`
	Debug           bool   `json:"debug"`
	Season          string `json:"season"`
	CurrentWorld    string `json:"current_world"`
	Balance         int    `json:"balance"`
	Energy          int    `json:"energy"`
	DropsInARow     int    `json:"drops_in_a_row"`
	ItemWheelsCount int    `json:"item_wheels_count"`
}

func stateToView(currentSeason string, player *model.PlayerInfo, progress *model.PlayerProgress) *view {
	v := &view{
		Season: currentSeason,
	}

	if player != nil {
		v.Id = player.ID()
		v.Disabled = player.Disabled()
		v.Debug = player.Debug()
	}

	if progress != nil {
		v.CurrentWorld = progress.CurrentWorld()
		v.Balance = progress.Balance()
		v.Energy = progress.Energy()
		v.DropsInARow = progress.DropsInARow()
		v.ItemWheelsCount = progress.ItemWheelsCount()
	}

	return v
}
