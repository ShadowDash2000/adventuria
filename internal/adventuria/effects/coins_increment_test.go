package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/schema"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_CoinsIncrement(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	_, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}

	item, err := createCoinsIncrementItem(ctx)
	if err != nil {
		t.Fatal(err)
	}

	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	wantBalance := 2
	if user.Balance() != wantBalance {
		t.Fatalf("Test_CoinsIncrement(): Balance = %d, want = %d", user.Balance(), wantBalance)
	}
}

func createCoinsIncrementItem(ctx adventuria.AppContext) (*core.Record, error) {
	effectRecord, err := createCoinsIncrementEffect(ctx)
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set("name", "Coins Increment")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("isUsingSlot", false)
	record.Set("isActiveByDefault", true)
	record.Set("canDrop", true)
	err = ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCoinsIncrementEffect(ctx adventuria.AppContext) (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Coins Increment")
	record.Set("type", "coinsIncrement")
	record.Set("value", "2;onAfterItemSave")
	err := ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
