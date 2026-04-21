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

	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatalf("Test_Buy(): Error getting player: %s", err)
	}

	_, err = player.Move(ctx, 3)
	if err != nil {
		t.Fatalf("Test_Buy(): Error moving: %s", err)
	}

	player.Progress().AddBalance(2)

	res, err := game.DoAction(ctx.App, player.ID(), ActionTypeBuyItem, adventuria.ActionRequest{
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

	if player.Progress().Balance() != 0 {
		t.Fatalf("Test_Buy(): Player balance = %d, want = 0", player.Progress().Balance())
	}

	if _, _, err = player.Inventory().UseItem(ctx, invItemId); err != nil {
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
	record.Set(schema.ItemSchema.Name, "Cell Points Divide")
	record.Set(schema.ItemSchema.Effects, []string{effectRecord.Id})
	record.Set(schema.ItemSchema.Icon, icon)
	record.Set(schema.ItemSchema.IsUsingSlot, true)
	record.Set(schema.ItemSchema.CanDrop, true)
	record.Set(schema.ItemSchema.Type, "buff")
	record.Set(schema.ItemSchema.IsRollable, true)
	record.Set(schema.ItemSchema.Price, 2)
	err = ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCellPointsDivideEffect(ctx adventuria.AppContext) (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Cell Points Divide")
	record.Set(schema.EffectSchema.Type, "cellPointsDivide")
	record.Set(schema.EffectSchema.Value, 2)
	err := ctx.App.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
