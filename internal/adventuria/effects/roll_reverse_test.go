package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_RollReverse(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createRollReverseItem()
	if err != nil {
		t.Fatal(err)
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 100; i++ {
		invItemId, err := user.Inventory().AddItemById(item.Id)
		if err != nil {
			t.Fatal(err)
		}

		err = game.UseItem(user.ID(), invItemId, adventuria.UseItemRequest{})
		if err != nil {
			t.Fatal(err)
		}

		res, err := game.DoAction(actions.ActionTypeRollDice, user.ID(), adventuria.ActionRequest{})
		if err != nil {
			t.Fatalf("Test_RollReverse(): Error rolling dice: %s", err)
		}

		rollDiceRes, ok := res.Data.(actions.RollDiceResult)
		if !ok {
			t.Fatal("Test_RollReverse(): Result data is not RollDiceResult")
		}

		dicesSum := 0
		for _, roll := range rollDiceRes.DiceRolls {
			dicesSum += roll.Roll
		}

		wantRoll := dicesSum * -1
		if wantRoll != rollDiceRes.Roll {
			t.Fatalf("Test_RollReverse(): Roll not reversed, want = %d, got = %d", wantRoll, rollDiceRes.Roll)
		}

		user.LastAction().SetCanMove(true)
	}
}

func createRollReverseItem() (*core.Record, error) {
	effectRecord, err := createRollReverseEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
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
