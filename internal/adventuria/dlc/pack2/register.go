package pack2

import "adventuria/internal/adventuria"

func WithItemPack2(g adventuria.Game) adventuria.Game {
	WithBaseEffects()
	return WithBaseEvents(g)
}
