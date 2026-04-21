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

func Test_ChangeMaxGamePriceUsable(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createChangeMaxGamePriceUsableItem()
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

	invItemId, err := player.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.UseItem(ctx.App, player.ID(), adventuria.UseItemRequest{InvItemId: invItemId})
	if err != nil {
		t.Fatal(err)
	}

	filter, err := player.LastAction().CustomActivityFilter()
	if err != nil {
		t.Fatal(err)
	}

	if filter.MaxPrice != 20 {
		t.Fatalf("Test_ChangeMaxGamePrice(): Max price is %d, expected 20", filter.MaxPrice)
	}
}

func createChangeMaxGamePriceUsableItem() (*core.Record, error) {
	effectRecord, err := createChangeMaxGamePriceUsableEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Change Max Activity Price Usable")
	record.Set(schema.ItemSchema.Effects, []string{effectRecord.Id})
	record.Set(schema.ItemSchema.Icon, icon)
	record.Set(schema.ItemSchema.IsUsingSlot, true)
	record.Set(schema.ItemSchema.CanDrop, false)
	record.Set(schema.ItemSchema.IsActiveByDefault, false)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeMaxGamePriceUsableEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Change Max Activity Price Usable")
	record.Set(schema.EffectSchema.Type, "changeMaxGamePrice")
	record.Set(schema.EffectSchema.Value, "20;usable")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func Test_ChangeMaxGamePriceUnusable(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	_, err := tests.NewGameTest()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createChangeMaxGamePriceUnusableItem()
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

	_, err = player.Inventory().AddItemById(ctx, item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = ctx.App.Save(player.LastAction().ProxyRecord())
	if err != nil {
		t.Fatal(err)
	}

	_ = player.Refetch(ctx)

	filter, err := player.LastAction().CustomActivityFilter()
	if err != nil {
		t.Fatal(err)
	}

	if filter.MaxPrice != 20 {
		t.Fatalf("Test_ChangeMaxGamePrice(): Max price is %d, expected 20", filter.MaxPrice)
	}
}

func createChangeMaxGamePriceUnusableItem() (*core.Record, error) {
	effectRecord, err := createChangeMaxGamePriceUnusableEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionItems))
	record.Set(schema.ItemSchema.Name, "Change Max Activity Price Unusable")
	record.Set(schema.ItemSchema.Effects, []string{effectRecord.Id})
	record.Set(schema.ItemSchema.Icon, icon)
	record.Set(schema.ItemSchema.IsUsingSlot, false)
	record.Set(schema.ItemSchema.CanDrop, false)
	record.Set(schema.ItemSchema.IsActiveByDefault, true)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeMaxGamePriceUnusableEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(schema.CollectionEffects))
	record.Set(schema.EffectSchema.Name, "Change Max Activity Price Unusable")
	record.Set(schema.EffectSchema.Type, "changeMaxGamePrice")
	record.Set(schema.EffectSchema.Value, "20;unusable")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
