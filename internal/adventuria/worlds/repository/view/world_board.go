package view

import (
	"adventuria/internal/adventuria/worlds/service/view"
)

type worldBoardRow struct {
	Id   string `json:"id"`
	Sort int    `json:"sort"`
}

func (w *worldBoardRow) toWorldBoard() *view.WorldBoard {
	if w == nil {
		return nil
	}

	return &view.WorldBoard{
		Id:   w.Id,
		Sort: w.Sort,
	}
}
