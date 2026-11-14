package repository

import (
	"adventuria/internal/adventuria/npcs"
	"adventuria/internal/adventuria/schema"

	"ganigeorgiev.com/pocketbase/dbx"
)

type filterColumns struct {
	CurrentCell string
}

var defaultFilterColumns = filterColumns{
	CurrentCell: schema.NPCsSchema.CurrentCell,
}

func buildQueryFromFilter(filter npcs.Filter, q *dbx.SelectQuery, columns filterColumns) {
	if filter.CurrentCell != nil {
		q.AndWhere(dbx.HashExp{
			columns.CurrentCell: *filter.CurrentCell,
		})
	}
}
