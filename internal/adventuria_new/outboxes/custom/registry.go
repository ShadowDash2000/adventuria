package custom

import (
	"adventuria/internal/adventuria_new/outboxes"
	"adventuria/internal/adventuria_new/outboxes/custom/change_balance"
	"adventuria/internal/adventuria_new/player_progress"
)

func RegisterOutboxes(
	progress *player_progress.PlayerProgress,
) {
	outboxes.Register(
		change_balance.NewDef(progress),
	)
}
