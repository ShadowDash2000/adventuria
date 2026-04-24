package cells

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria/tests"
	"errors"
	"fmt"
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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = createTeleportCell(player.Progress().CurrentWorld())
	if err != nil {
		t.Fatal(err)
	}

	firstCell, ok := adventuria.GameCells.GetByOrder(player.Progress().CurrentWorld(), 0)
	if !ok {
		t.Fatal("Test_CellTeleport(): Could not find cell 0")
	}

	_, err = player.Move(ctx, 4)
	if err != nil {
		t.Fatal(err)
	}

	currentCell, ok := player.Progress().CurrentCell()
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

func createTeleportCell(worldId string) (*core.Record, error) {
	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	firstCell, ok := adventuria.GameCells.GetByOrder(worldId, 0)
	if !ok {
		return nil, errors.New("createTeleportCell(): could not find cell 0")
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionCells))
	record.Set(schema.CellSchema.Name, "Cell Teleport")
	record.Set(schema.CellSchema.Icon, icon)
	record.Set(schema.CellSchema.Sort, 500)
	record.Set(schema.CellSchema.Type, "teleport")
	record.Set(schema.CellSchema.Value, fmt.Sprintf("{\"cell_id\": \"%s\"}", firstCell.ID()))
	record.Set(schema.CellSchema.World, worldId)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
