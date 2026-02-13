package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_Buy(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Buy(): Error creating game: %s", err)
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}

	item, err := createCellPointsDivideItem(ctx)
	if err != nil {
		t.Fatal(err)
	}

	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatalf("Test_Buy(): Error getting user: %s", err)
	}

	_, err = user.Move(ctx, 3)
	if err != nil {
		t.Fatalf("Test_Buy(): Error moving: %s", err)
	}

	user.AddBalance(2)

	res, err := game.DoAction(ctx.App, user.ID(), ActionTypeBuyItem, adventuria.ActionRequest{
		"item_id": item.Id,
	})
	if err != nil {
		t.Fatalf("Test_Buy(): Error performing action: %s", err)
	}

	var invItemId string
	invItemId, ok := res.Data.(string)
	if !ok {
		t.Fatalf("Test_Buy(): Error getting item id from response: %s", err)
	}

	if user.Balance() != 0 {
		t.Fatalf("Test_Buy(): User balance = %d, want = 0", user.Balance())
	}

	if _, _, err = user.Inventory().UseItem(ctx, invItemId); err != nil {
		t.Fatalf("Test_Buy(): Error using item: %s", err)
	}
}

func createCellPointsDivideItem(ctx adventuria.AppContext) (*core.Record, error) {
	effectRecord, err := createCellPointsDivideEffect(ctx)
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set("name", "Cell Points Divide")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", true)
	record.Set("type", "buff")
	record.Set("isRollable", true)
	record.Set("price", 2)
	err = ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCellPointsDivideEffect(ctx adventuria.AppContext) (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Cell Points Divide")
	record.Set("type", "cellPointsDivide")
	record.Set("value", 2)
	err := ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
