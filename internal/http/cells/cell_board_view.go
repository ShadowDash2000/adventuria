package cells

import (
	"adventuria/internal/adventuria/cells/service/view"
	"adventuria/internal/adventuria/model"
)

type cellBoardView struct {
	Id    string         `json:"id"`
	Type  model.CellType `json:"type"`
	Icon  string         `json:"icon"`
	World string         `json:"world"`
}

func cellBoardToView(cell *view.CellBoard) *cellBoardView {
	return &cellBoardView{
		Id:    cell.Id,
		Type:  cell.Type,
		Icon:  cell.Icon,
		World: cell.World,
	}
}

func cellsBoardToView(cells []*view.CellBoard) []*cellBoardView {
	res := make([]*cellBoardView, len(cells))
	for i, cell := range cells {
		res[i] = cellBoardToView(cell)
	}
	return res
}
