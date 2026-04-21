package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/schema"
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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := player.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	genre, err := ctx.App.FindFirstRecordByFilter(
		adventuria.GameCollections.Get(schema.CollectionGenres),
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{
		InvItemId: invItemId,
		Data: map[string]any{
			"genre_id": genre.Id,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	filter, err := player.LastAction().CustomActivityFilter()
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(filter.Genres, genre.Id) {
		t.Fatalf("Test_AddGameGenre(): Genre not added to player, want = %s, got = %s", genre.Id, filter.Genres)
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

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Add Game Genre")
	record.Set(schema.ItemSchema.Effects, []string{effectRecord.Id})
	record.Set(schema.ItemSchema.Icon, icon)
	record.Set(schema.ItemSchema.IsUsingSlot, true)
	record.Set(schema.ItemSchema.CanDrop, true)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createAddGameGenreEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Add Game Genre")
	record.Set(schema.EffectSchema.Type, "addGameGenre")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
