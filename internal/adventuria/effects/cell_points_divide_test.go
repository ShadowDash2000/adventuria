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

func Test_CellPointsDivide(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createCellPointsDivideItem()
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

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	const points = 100
	player.Progress().AddPoints(points)

	if err = ctx.App.Save(player.Progress().ProxyRecord()); err != nil {
		t.Fatalf("Test_Buy(): Error saving player's progress: %s", err)
	}

	invItemId, err := player.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), actions.ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), actions.ActionTypeDone, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := player.Progress().CurrentCell()
	if !ok {
		t.Fatal("Test_CellPointsDivide(): Current cell not found")
	}

	wantPoints := points + cell.Points()/2
	if player.Progress().Points() != wantPoints {
		t.Fatalf("Test_CellPointsDivide(): Points not divided, want = %d, got = %d", wantPoints, player.Progress().Points())
	}
}

func createCellPointsDivideItem() (*core.Record, error) {
	effectRecord, err := createCellPointsDivideEffect()
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
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCellPointsDivideEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Cell Points Divide")
	record.Set(schema.EffectSchema.Type, "cellPointsDivide")
	record.Set(schema.EffectSchema.Value, 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
