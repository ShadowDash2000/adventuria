package custom

import (
	"adventuria/internal/adventuria/outboxes"
	"adventuria/internal/adventuria/outboxes/custom/change_balance"
	"adventuria/internal/adventuria/player_progress"
)

func RegisterOutboxes(
	progress *player_progress.PlayerProgress,
) {
	outboxes.Register(
		change_balance.NewDef(progress),
	)
}
