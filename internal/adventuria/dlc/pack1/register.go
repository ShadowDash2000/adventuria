package pack1

import "adventuria/internal/adventuria"

func WithItemPack1(g adventuria.Game) adventuria.Game {
	WithBaseEffects()
	return WithBaseEvents(g)
}
