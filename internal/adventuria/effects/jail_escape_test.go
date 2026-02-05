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

func Test_JailEscape(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createJailEscapeItem()
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

	user.SetIsInJail(true)
	if err = ctx.App.Save(user.ProxyRecord()); err != nil {
		t.Fatalf("Test_Buy(): Error saving user: %s", err)
	}

	err = game.UseItem(ctx.App, user.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	user, err = adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	if user.IsInJail() {
		t.Fatalf("Test_JailEscape(): User is in jail, expected not in jail")
	}
}

func createJailEscapeItem() (*core.Record, error) {
	effectRecord, err := createJailEscapeEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Jail Escape")
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

func createJailEscapeEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Jail Escape")
	record.Set("type", "jailEscape")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
