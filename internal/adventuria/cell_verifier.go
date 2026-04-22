package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

// CellVerifier
// Binds hooks on cell's collection for record creation and update
// that verifies that a cell type really exists and calls Verify()
// method of cell that should try to parse record's value
type CellVerifier struct{}

func NewCellVerifier(ctx AppContext) *CellVerifier {
	ef := &CellVerifier{}
	ef.bindHooks(ctx)
	return ef
}

func (cf *CellVerifier) bindHooks(ctx AppContext) {
	ctx.App.OnRecordValidate(schema.CollectionCells).BindFunc(func(e *core.RecordEvent) error {
		if err := cf.Verify(AppContext{App: e.App}, e.Record); err != nil {
			return err
		}
		return e.Next()
	})
}

func (cf *CellVerifier) Verify(ctx AppContext, record *core.Record) error {
	cellType := record.GetString(schema.CellSchema.Type)
	cellValue := record.GetString(schema.CellSchema.Value)

	cellCreator, ok := cellsList[CellType(cellType)]
	if !ok {
		return errors.New("unknown cell type")
	}

	cell := cellCreator.New(record)
	cellVerifiable, ok := cell.(CellVerifiable)

	if !ok {
		// cellValue is JSON value so we need to check those empty values
		if cellValue == "\"\"" || cellValue == "null" {
			return nil
		}
		return errors.New("cell type is not verifiable")
	}

	return cellVerifiable.Verify(ctx, cellValue)
}
