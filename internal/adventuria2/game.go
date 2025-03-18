package adventuria2

import (
	"github.com/ShadowDash2000/adventuria"
	"github.com/pocketbase/pocketbase/core"
)

type Game struct {
	adventuria.BaseGame
}

func NewGame(app core.App) adventuria.Game {
	g := adventuria.New(app)

	g = adventuria.WithBaseEvents(g)
	adventuria.WithBaseEffects()

	return g
}
