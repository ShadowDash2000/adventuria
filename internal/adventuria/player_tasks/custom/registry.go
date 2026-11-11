package custom

import (
	"adventuria/internal/adventuria/player_tasks"
	"adventuria/internal/adventuria/player_tasks/custom/daily_reward_bonus"
)

func RegisterPlayerTasks() {
	player_tasks.Register(
		daily_reward_bonus.NewDef(),
	)
}
