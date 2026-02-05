package adventuria

import (
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

// CellVerifier
// Binds hooks on cell's collection for record creation and update
// that verifies that a cell type really exists and calls Verify()
// method of cell that should try to parse record's value
type CellVerifier struct{}

func NewCellVerifier() *CellVerifier {
	ef := &CellVerifier{}
	ef.bindHooks()
	return ef
}

func (cf *CellVerifier) bindHooks() {
	PocketBase.OnRecordCreate(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if err := cf.Verify(AppContext{App: e.App}, e.Record); err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordUpdate(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if err := cf.Verify(AppContext{App: e.App}, e.Record); err != nil {
			return err
		}
		return e.Next()
	})
}

func (cf *CellVerifier) Verify(ctx AppContext, record *core.Record) error {
	cellType := record.GetString("type")
	cellValue := record.GetString("value")

	cellCreator, ok := cellsList[CellType(cellType)]
	if !ok {
		return errors.New("unknown cell type")
	}

	cell := cellCreator.New(record)

	return cell.Verify(ctx, cellValue)
}
