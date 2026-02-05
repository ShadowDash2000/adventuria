package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_SafeDrop(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createSafeDropItem()
	if err != nil {
		t.Fatal(err)
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := user.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(ctx.App, user.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), actions.ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), actions.ActionTypeDrop, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	user, err = adventuria.GameUsers.GetByName(ctx, "user1")
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

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
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
