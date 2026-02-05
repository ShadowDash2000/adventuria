package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"adventuria/pkg/helper"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_ChangeGameById(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
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

	_, err = game.DoAction(ctx.App, user.ID(), "rollWheel", adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	user, err = adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	itemsList, err := user.LastAction().ItemsList()
	if err != nil {
		t.Fatal(err)
	}

	if index := slices.Index(itemsList, user.LastAction().Activity()); index != -1 {
		itemsList = slices.Delete(itemsList, index, index+1)
	}

	gameId := helper.RandomItemFromSlice(itemsList)
	item, err := createChangeGameByIdItem(gameId)
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

	user, err = adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	if user.LastAction().Activity() != gameId {
		t.Fatalf("Test_ChangeGameById(): Activity id = %s, want = %s", user.LastAction().Activity(), gameId)
	}
}

func createChangeGameByIdItem(gameId string) (*core.Record, error) {
	effectRecord, err := createChangeGameByIdEffect(gameId)
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Change Activity By Id")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", true)
	record.Set("isActiveByDefault", false)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeGameByIdEffect(gameId string) (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Change Activity By Id")
	record.Set("type", "changeGameById")
	record.Set("value", gameId)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
