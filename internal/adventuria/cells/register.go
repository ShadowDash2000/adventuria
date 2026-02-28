package cells

import "adventuria/internal/adventuria"

func WithBaseCells() {
	adventuria.RegisterCells([]adventuria.CellDef{
		adventuria.NewCell(
			"game",
			NewCellActivity(adventuria.ActivityTypeGame),
			"wheel", "activity", "game",
		),
		adventuria.NewCell(
			"movie",
			NewCellActivity(adventuria.ActivityTypeMovie),
			"wheel", "activity",
		),
		adventuria.NewCell(
			"gym",
			NewCellActivity(adventuria.ActivityTypeGym),
			"wheel", "activity",
		),
		adventuria.NewCell(
			"karaoke",
			NewCellActivity(adventuria.ActivityTypeKaraoke),
			"wheel", "activity",
		),
		adventuria.NewCell(
			"start",
			func() adventuria.Cell { return &CellStart{} },
		),
		adventuria.NewCell(
			"jail",
			NewCellJail(),
			"wheel", "activity", "game",
		),
		adventuria.NewCell(
			"casino",
			func() adventuria.Cell { return &CellCasino{} },
			"shop",
		),
		adventuria.NewCell(
			"shop",
			func() adventuria.Cell { return &CellShop{} },
			"shop",
		),
		adventuria.NewCell(
			"teleport",
			func() adventuria.Cell { return &CellTeleport{} },
		),
		adventuria.NewCell(
			"rollItem",
			func() adventuria.Cell { return &CellRollItem{} },
			"wheel",
		),
	})
}
