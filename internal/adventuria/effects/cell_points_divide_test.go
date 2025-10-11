package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func TestCellPointsDivide(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createCellPointsDivideItem()
	if err != nil {
		t.Fatal(err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	const points = 100
	user.SetPoints(points)

	invItemId, err := user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId)
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(adventuria.ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(adventuria.ActionTypeDone, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(user.Points())

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Current cell not found")
	}

	if user.Points() != points+cell.Points()/2 {
		t.Fatal("Points not divided")
	}
}

func createCellPointsDivideItem() (*core.Record, error) {
	collection, err := adventuria.GameCollections.Get(adventuria.TableItems)
	if err != nil {
		return nil, err
	}

	effectRecord, err := createCellPointsDivideEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(adventuria.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Set("name", "Cell Points Divide")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCellPointsDivideEffect() (*core.Record, error) {
	collection, err := adventuria.GameCollections.Get(adventuria.TableEffects)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	record.Set("name", "Cell Points Divide")
	record.Set("type", "cellPointsDivide")
	record.Set("value", 2)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
