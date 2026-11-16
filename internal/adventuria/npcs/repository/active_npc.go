package repository

import (
	"adventuria/internal/adventuria/npcs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
)

type activeNPCRow struct {
	Id     string `db:"id"`
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
	Color  string `db:"color"`
}

func (a *activeNPCRow) ToActiveNPC() *npcs.ActiveNPC {
	if a == nil {
		return nil
	}

	return &npcs.ActiveNPC{
		Id:     a.Id,
		Name:   a.Name,
		Avatar: pbhelper.BuildFilePath(schema.CollectionNPCs, a.Id, a.Avatar),
		Color:  a.Color,
	}
}
