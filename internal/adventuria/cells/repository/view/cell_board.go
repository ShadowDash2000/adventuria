package view

import (
	"adventuria/internal/adventuria/cells/service/view"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
)

type cellBoardRow struct {
	Id    string         `db:"id"`
	Type  model.CellType `db:"type"`
	Icon  string         `db:"icon"`
	World string         `db:"world"`
}

func (c *cellBoardRow) toCellBoard() *view.CellBoard {
	if c == nil {
		return nil
	}

	return &view.CellBoard{
		Id:   c.Id,
		Type: c.Type,
		Icon: pbhelper.BuildFilePath(schema.CollectionCells, c.Id, c.Icon, pbhelper.FilePathOpts{
			Thumb: &pbhelper.FileThumb{
				Width:  640,
				Height: 0,
			},
		}),
		World: c.World,
	}
}
