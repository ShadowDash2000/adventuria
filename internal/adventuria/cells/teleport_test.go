package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_TeleportCell(t *testing.T) {
	WithBaseCells()

	_, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	_, err = createTeleportCell()
	if err != nil {
		t.Fatal(err)
	}

	firstCell, ok := adventuria.GameCells.GetByOrder(0)
	if !ok {
		t.Fatal("Test_CellTeleport(): Could not find cell 0")
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(ctx, 4)
	if err != nil {
		t.Fatal(err)
	}

	currentCell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_CellTeleport(): Current cell not found")
	}

	if currentCell.ID() != firstCell.ID() {
		t.Fatalf(
			"Test_CellTeleport(): Expected cell = %s (%s), got = %s (%s)",
			firstCell.Name(),
			firstCell.ID(),
			currentCell.Name(),
			currentCell.ID(),
		)
	}
}

func createTeleportCell() (*core.Record, error) {
	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionCells))
	record.Set("name", "Cell Teleport")
	record.Set("icon", icon)
	record.Set("sort", 500)
	record.Set("type", "teleport")
	record.Set("value", "{\"cell_name\": \"Cell 1 (start)\"}")
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
