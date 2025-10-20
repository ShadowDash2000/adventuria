package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_RollReverse(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createRollReverseItem()
	if err != nil {
		t.Fatal(err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId)
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_RollReverse(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	// TODO in some cases new cell still unreached
	res, err := game.DoAction(actions.ActionTypeRollDice, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	rollDiceRes, ok := res.Data.(actions.RollDiceResult)
	if !ok {
		t.Fatal("Test_RollReverse(): Result data is not RollDiceResult")
	}

	t.Log("Test_RollReverse(): Roll dice result:", rollDiceRes)

	dicesSum := 0
	for _, n := range rollDiceRes.DiceRolls {
		dicesSum += n
	}

	wantRoll := dicesSum * -1
	if wantRoll != rollDiceRes.Roll {
		t.Fatalf("Test_DiceIncrement(): Roll not incremented, want = %d, got = %d", wantRoll, rollDiceRes.Roll)
	}
}

func createRollReverseItem() (*core.Record, error) {
	effectRecord, err := createRollReverseEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(adventuria.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Roll Reverse")
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

func createRollReverseEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Roll Reverse")
	record.Set("type", "rollReverse")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
