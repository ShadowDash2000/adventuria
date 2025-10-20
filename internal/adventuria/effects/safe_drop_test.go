package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_SafeDrop(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createSafeDropItem()
	if err != nil {
		t.Fatal(err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_SafeDrop(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	err = user.LastAction().Save()
	if err != nil {
		t.Fatalf("Test_SafeDrop(): Error saving action: %s", err)
	}

	_, err = game.DoAction(actions.ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(actions.ActionTypeDrop, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	wantDropsInARow := 0
	if user.DropsInARow() != wantDropsInARow {
		t.Fatalf("Test_SafeDrop(): Drops in a row is %d, expected %d", user.DropsInARow(), wantDropsInARow)
	}
}

func createSafeDropItem() (*core.Record, error) {
	effectRecord, err := createSafeDropEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(adventuria.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Safe Drop")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", true)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createSafeDropEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Safe Drop")
	record.Set("type", "safeDrop")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
