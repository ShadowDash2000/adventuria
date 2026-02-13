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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 50; i++ {
		invItemId, err := user.Inventory().AddItemById(ctx, item.Id)
		if err != nil {
			t.Fatal(err)
		}

		_, err = game.UseItem(ctx.App, user.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
		if err != nil {
			t.Fatal(err)
		}

		res, err := game.DoAction(ctx.App, user.ID(), actions.ActionTypeRollDice, adventuria.ActionRequest{})
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
		if err = ctx.App.Save(user.LastAction().ProxyRecord()); err != nil {
			t.Fatalf("Test_Buy(): Error saving user: %s", err)
		}
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

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set("name", "Roll Reverse")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", true)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createRollReverseEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set("name", "Roll Reverse")
	record.Set("type", "rollReverse")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
