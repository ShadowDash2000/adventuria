package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/schema"
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

	_, err = game.DoAction(ctx.App, player.ID(), "rollWheel", adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	itemsList, err := player.LastAction().ItemsList()
	if err != nil {
		t.Fatal(err)
	}

	if index := slices.Index(itemsList, player.LastAction().Activity()); index != -1 {
		itemsList = slices.Delete(itemsList, index, index+1)
	}

	gameId := helper.RandomItemFromSlice(itemsList)
	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{
		InvItemId: invItemId,
		Data: map[string]any{
			"game_id": gameId,
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	if player.LastAction().Activity() != gameId {
		t.Fatalf("Test_ChooseGame(): Activity id = %s, want = %s", player.LastAction().Activity(), gameId)
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

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Choose Activity")
	record.Set(schema.ItemSchema.Effects, []string{effectRecord.Id})
	record.Set(schema.ItemSchema.Icon, icon)
	record.Set(schema.ItemSchema.IsUsingSlot, true)
	record.Set(schema.ItemSchema.CanDrop, true)
	record.Set(schema.ItemSchema.IsActiveByDefault, false)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChooseGameEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Choose Activity")
	record.Set(schema.EffectSchema.Type, "chooseGame")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
