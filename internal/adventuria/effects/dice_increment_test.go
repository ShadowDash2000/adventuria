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

func Test_DiceIncrement(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createDiceIncrementItem()
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

	invItemId, err := user.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(ctx.App, user.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	res, err := game.DoAction(ctx.App, user.ID(), actions.ActionTypeRollDice, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	rollDiceRes, ok := res.Data.(actions.RollDiceResult)
	if !ok {
		t.Fatal("Test_DiceIncrement(): Result data is not RollDiceResult")
	}

	t.Log("Test_DiceIncrement(): Roll dice result:", rollDiceRes)

	dicesSum := 0
	for _, roll := range rollDiceRes.DiceRolls {
		dicesSum += roll.Roll
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

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
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
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Dice Increment")
	record.Set("type", "diceIncrement")
	record.Set("value", 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
