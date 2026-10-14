package custom

import (
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/outboxes"
	"adventuria/internal/adventuria/outboxes/custom/change_balance"
	"adventuria/internal/adventuria/outboxes/custom/change_energy"
	"adventuria/internal/adventuria/player_events"
	"adventuria/internal/adventuria/player_progress"
)

func RegisterOutboxes(
	progress *player_progress.PlayerProgress,
	actions *actions.Actions,
	playerEvents *player_events.PlayerEvents,
) {
	outboxes.Register(
		change_balance.NewDef(progress),
		change_energy.NewDef(progress, actions, playerEvents),
	)
}
