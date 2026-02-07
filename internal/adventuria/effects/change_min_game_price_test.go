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

func Test_ChangeMinGamePrice(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createChangeMinGamePriceItem()
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

	_, err = user.Move(ctx, 1)
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

	filter, err := user.LastAction().CustomActivityFilter()
	if err != nil {
		t.Fatal(err)
	}

	if filter.MinPrice != 20 {
		t.Fatalf("Test_ChangeMaxGamePrice(): Min price is %d, expected 20", filter.MinPrice)
	}
}

func createChangeMinGamePriceItem() (*core.Record, error) {
	effectRecord, err := createChangeMinGamePriceEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set("name", "Change Min Activity Price")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", false)
	record.Set("isActiveByDefault", false)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeMinGamePriceEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Change Min Activity Price")
	record.Set("type", "changeMinGamePrice")
	record.Set("value", 20)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
