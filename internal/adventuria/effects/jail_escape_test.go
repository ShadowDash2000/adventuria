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

func Test_JailEscape(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createJailEscapeItem()
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

	player.Progress().SetIsInJail(true)
	if err = ctx.App.Save(player.Progress().ProxyRecord()); err != nil {
		t.Fatalf("Test_Buy(): Error saving player's progress: %s", err)
	}

	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	if player.Progress().IsInJail() {
		t.Fatalf("Test_JailEscape(): Player is in jail, expected not in jail")
	}
}

func createJailEscapeItem() (*core.Record, error) {
	effectRecord, err := createJailEscapeEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Jail Escape")
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

func createJailEscapeEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Jail Escape")
	record.Set(schema.EffectSchema.Type, "jailEscape")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
