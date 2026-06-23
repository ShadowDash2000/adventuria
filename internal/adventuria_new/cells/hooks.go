package cells

import (
	"adventuria/internal/adventuria/schema"
	repo "adventuria/internal/adventuria_new/cells/repository"
	"adventuria/internal/adventuria_new/errs"
	"adventuria/internal/adventuria_new/model"
	"context"
	"errors"
	"fmt"

	"github.com/pocketbase/pocketbase/core"
)

func BindHooks(pb core.App) {
	pb.OnRecordValidate(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		err := verify(e.Context, repo.RecordToCell(e.Record))
		if err != nil {
			return err
		}
		return e.Next()
	})
}

func verify(ctx context.Context, cellInfo *model.CellInfo) error {
	cellValue := cellInfo.Value()

	cellDef, ok := Get(cellInfo.Type())
	if !ok {
		return fmt.Errorf("%w: %s", errs.ErrUnknownCellType, cellInfo.Type())
	}

	cell := cellDef.new(*cellInfo)
	cellVerifiable, ok := cell.(model.Verifiable)
	if !ok {
		// cellValue is JSON value so we need to check those empty values
		if cellValue == "\"\"" || cellValue == "null" {
			return nil
		}
		return errors.New("cell type is not verifiable")
	}

	return cellVerifiable.Verify(ctx, cellValue)
}
