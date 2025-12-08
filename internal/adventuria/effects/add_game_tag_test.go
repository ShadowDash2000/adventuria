package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_AddGameTag(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createAddGameTagItem()
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

	_, err = game.DoAction(actions.ActionTypeRollDice, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	tag, err := adventuria.PocketBase.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionTags),
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId, adventuria.UseItemRequest{
		"tag_id": tag.Id,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(user.LastAction().CustomGameFilter().Tags, tag.Id) {
		t.Fatalf("Test_AddGameTag(): Tag not added to user, want = %s, got = %s", tag.Id, user.LastAction().CustomGameFilter().Tags)
	}
}

func createAddGameTagItem() (*core.Record, error) {
	effectRecord, err := createAddGameTagEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Add Game Tag")
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

func createAddGameTagEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Add Game Tag")
	record.Set("type", "addGameTag")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
