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

func Test_AddGameGenre(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createAddGameGenreItem()
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

	genre, err := adventuria.PocketBase.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(adventuria.CollectionGenres),
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId, adventuria.UseItemRequest{
		"genre_id": genre.Id,
	})
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(user.LastAction().CustomActivityFilter().Genres, genre.Id) {
		t.Fatalf("Test_AddGameGenre(): Genre not added to user, want = %s, got = %s", genre.Id, user.LastAction().CustomActivityFilter().Genres)
	}
}

func createAddGameGenreItem() (*core.Record, error) {
	effectRecord, err := createAddGameGenreEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Add Game Genre")
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

func createAddGameGenreEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Add Game Genre")
	record.Set("type", "addGameGenre")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
