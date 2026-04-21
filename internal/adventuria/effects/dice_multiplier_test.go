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

func Test_DiceMultiplier(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createDiceMultiplierItem()
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

	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	res, err := game.DoAction(ctx.App, player.ID(), actions.ActionTypeRollDice, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	rollDiceRes, ok := res.Data.(actions.RollDiceResult)
	if !ok {
		t.Fatal("Test_DiceMultiplier(): Result data is not RollDiceResult")
	}

	t.Log("Test_DiceMultiplier(): Roll dice result:", rollDiceRes)

	dicesSum := 0
	for _, roll := range rollDiceRes.DiceRolls {
		dicesSum += roll.Roll
	}

	wantRoll := dicesSum * 2
	if wantRoll != rollDiceRes.Roll {
		t.Fatalf("Test_DiceMultiplier(): Roll not incremented, want = %d, got = %d", wantRoll, rollDiceRes.Roll)
	}
}

func createDiceMultiplierItem() (*core.Record, error) {
	effectRecord, err := createDiceMultiplierEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Dice Multiplier")
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

func createDiceMultiplierEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Dice Multiplier")
	record.Set(schema.EffectSchema.Type, "diceMultiplier")
	record.Set(schema.EffectSchema.Value, 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
