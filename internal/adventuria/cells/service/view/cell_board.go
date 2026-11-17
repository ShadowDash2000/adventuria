package view

import "adventuria/internal/adventuria/model"

type CellBoard struct {
	Id    string
	Type  model.CellType
	Icon  string
	World string
}
