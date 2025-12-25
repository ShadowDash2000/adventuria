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

func Test_ChooseGame(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createChooseGameItem()
	if err != nil {
		t.Fatal(err)
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction("rollWheel", user.ID(), adventuria.ActionRequest{})
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
	err = game.UseItem(user.ID(), invItemId, adventuria.UseItemRequest{"game_id": gameId})
	if err != nil {
		t.Fatal(err)
	}

	if user.LastAction().Activity() != gameId {
		t.Fatalf("Test_ChooseGame(): Activity id = %s, want = %s", user.LastAction().Activity(), gameId)
	}
}

func createChooseGameItem() (*core.Record, error) {
	effectRecord, err := createChooseGameEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Choose Activity")
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

func createChooseGameEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Choose Activity")
	record.Set("type", "chooseGame")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
