package repository

import (
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"context"
	"fmt"

	"ganigeorgiev.com/pocketbase/pocketbase/core"
)

func BindHooks(pb core.App) {
	pb.OnRecordValidate(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := verify(e.Context, RecordToCell(e.Record))
		if err != nil {
			return err
		}
		return e.Next()
	})
}

func verify(ctx context.Context, cellInfo *model.CellInfo) error {
	cellDef, ok := cells.Get(cellInfo.Type())
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownCellType, cellInfo.Type())
	}

	cell, err := cellDef.New(*cellInfo)
	if err != nil {
		return err
	}

	cellVerifiable, ok := cell.(model.Verifiable)
	if ok {
		return cellVerifiable.Verify(ctx, cellInfo.Value())
	}

	return nil
}
