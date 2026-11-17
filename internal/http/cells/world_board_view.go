package cells

import (
	"adventuria/internal/adventuria/worlds/service/view"
)

type worldBoardView struct {
	Id   string `json:"id"`
	Sort int    `json:"sort"`
}

func worldBoardToView(world *view.WorldBoard) *worldBoardView {
	return &worldBoardView{
		Id:   world.Id,
		Sort: world.Sort,
	}
}

func worldsBoardToView(worlds []*view.WorldBoard) []*worldBoardView {
	res := make([]*worldBoardView, len(worlds))
	for i, world := range worlds {
		res[i] = worldBoardToView(world)
	}
	return res
}
