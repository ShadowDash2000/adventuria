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

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	invItemId, err := user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	err = game.UseItem(user.ID(), invItemId, adventuria.UseItemRequest{})
	if err != nil {
		t.Fatal(err)
	}

	if user.LastAction().CustomGameFilter().MaxPrice != 20 {
		t.Fatalf("Test_ChangeMaxGamePrice(): Max price is %d, expected 20", user.LastAction().CustomGameFilter().MaxPrice)
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

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Change Max Game Price Usable")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", false)
	record.Set("isActiveByDefault", false)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeMaxGamePriceUsableEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Change Max Game Price Usable")
	record.Set("type", "changeMaxGamePrice")
	record.Set("value", "20;usable")
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

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Inventory().AddItemById(item.Id)
	if err != nil {
		t.Fatal(err)
	}

	if user.LastAction().CustomGameFilter().MaxPrice != 20 {
		t.Fatalf("Test_ChangeMaxGamePrice(): Max price is %d, expected 20", user.LastAction().CustomGameFilter().MaxPrice)
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

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Change Max Game Price Unusable")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", false)
	record.Set("canDrop", false)
	record.Set("isActiveByDefault", true)

	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createChangeMaxGamePriceUnusableEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Change Max Game Price Unusable")
	record.Set("type", "changeMaxGamePrice")
	record.Set("value", "20;unusable")
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
