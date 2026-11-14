package repository

import (
	"adventuria/internal/adventuria/model"
	"strings"
)

func stringToCellTypes(s string) []model.CellType {
	if s == "" {
		return nil
	}

	cellTypes := strings.Split(s, ";")
	res := make([]model.CellType, len(cellTypes))
	for i, cellType := range cellTypes {
		res[i] = model.CellType(cellType)
	}

	return res
}
