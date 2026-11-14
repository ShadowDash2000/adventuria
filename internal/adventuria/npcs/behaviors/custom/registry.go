package custom

import (
	"adventuria/internal/adventuria/board"
	"adventuria/internal/adventuria/npcs/behaviors"
	"adventuria/internal/adventuria/npcs/behaviors/custom/walk"
)

func RegisterNPCBehaviors(
	board *board.Board,
) {
	behaviors.Register(
		walk.NewDef(board),
	)
}
