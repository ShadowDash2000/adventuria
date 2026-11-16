package repository

import (
	"adventuria/internal/adventuria/board/actors"
	"adventuria/pkg/pbhelper"
	"time"

	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

type activeActorRow struct {
	Id                string           `db:"id"`
	Collection        string           `db:"collection"`
	Kind              actors.ActorKind `db:"kind"`
	Updated           types.DateTime   `db:"updated"`
	Name              string           `db:"name"`
	Avatar            string           `db:"avatar"`
	Color             string           `db:"color"`
	CurrentWorld      string           `db:"current_world"`
	CellsPassed       int              `db:"cells_passed"`
	NextShiftChangeAt *types.DateTime  `db:"next_shift_change_at"`
}

func (a *activeActorRow) toActiveActor() *actors.ActiveActor {
	var nextShiftChangeAt *time.Time
	if a.NextShiftChangeAt != nil {
		nextShiftChangeAt = new(a.NextShiftChangeAt.Time())
	}

	return &actors.ActiveActor{
		Id:                a.Id,
		Kind:              a.Kind,
		Updated:           a.Updated.Time(),
		Name:              a.Name,
		Avatar:            pbhelper.BuildFilePath(a.Collection, a.Id, a.Avatar),
		Color:             a.Color,
		CurrentWorld:      a.CurrentWorld,
		CellsPassed:       a.CellsPassed,
		NextShiftChangeAt: nextShiftChangeAt,
	}
}
