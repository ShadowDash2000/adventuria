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
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	const points = 100
	user.SetPoints(points)

	if err = ctx.App.Save(user.ProxyRecord()); err != nil {
		t.Fatalf("Test_Buy(): Error saving user: %s", err)
	}

	invItemId, err := user.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(ctx.App, user.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_CellPointsDivide(): Current cell not found")
	}

	_, err = game.DoAction(ctx.App, user.ID(), actions.ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), actions.ActionTypeDone, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Test_CellPointsDivide(): Points:", user.Points())

	user, err = adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	wantPoints := points + cell.Points()/2
	if user.Points() != wantPoints {
		t.Fatalf("Test_CellPointsDivide(): Points not divided, want = %d, got = %d", wantPoints, user.Points())
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
	record.Set("name", "Cell Points Divide")
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

func createCellPointsDivideEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Cell Points Divide")
	record.Set("type", "cellPointsDivide")
	record.Set("value", 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
