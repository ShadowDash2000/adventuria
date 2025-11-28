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
		if err := cf.Verify(e.Record.GetString("type"), e.Record.GetString("value")); err != nil {
			return err
		}
		return e.Next()
	})
	PocketBase.OnRecordUpdate(CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if err := cf.Verify(e.Record.GetString("type"), e.Record.GetString("value")); err != nil {
			return err
		}
		return e.Next()
	})
}

func (cf *CellVerifier) Verify(cellType, value string) error {
	cellCreator, ok := cellsList[CellType(cellType)]
	if !ok {
		return errors.New("unknown cell type")
	}

	cell := cellCreator()

	return cell.Verify(value)
}
