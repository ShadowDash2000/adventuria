package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_DiceIncrement(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createDiceIncrementItem()
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
		t.Fatal("Test_DiceIncrement(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	res, err := game.DoAction(actions.ActionTypeRollDice, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	rollDiceRes, ok := res.Data.(actions.RollDiceResult)
	if !ok {
		t.Fatal("Test_DiceIncrement(): Result data is not RollDiceResult")
	}

	t.Log("Test_DiceIncrement(): Roll dice result:", rollDiceRes)

	dicesSum := 0
	for _, n := range rollDiceRes.DiceRolls {
		dicesSum += n
	}

	wantRoll := dicesSum + 2
	if wantRoll != rollDiceRes.Roll {
		t.Fatalf("Test_DiceIncrement(): Roll not incremented, want = %d, got = %d", wantRoll, rollDiceRes.Roll)
	}
}

func createDiceIncrementItem() (*core.Record, error) {
	effectRecord, err := createDiceIncrementEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(adventuria.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Dice Increment")
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

func createDiceIncrementEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Dice Increment")
	record.Set("type", "diceIncrement")
	record.Set("value", 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
