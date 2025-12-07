package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/tests"
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_Buy(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Buy(): Error creating game: %s", err)
	}

	var itemId string
	for i := 0; i < 3; i++ {
		item, err := createCellPointsDivideItem()
		if err != nil {
			t.Fatal(err)
		}
		itemId = item.Id
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatalf("Test_Buy(): Error getting user: %s", err)
	}

	_, err = user.Move(3)
	if err != nil {
		t.Fatalf("Test_Buy(): Error moving: %s", err)
	}

	user.SetBalance(2)

	err = adventuria.PocketBase.Save(user.LastAction().ProxyRecord())
	if err != nil {
		t.Fatalf("Test_Buy(): Error saving action: %s", err)
	}

	res, err := game.DoAction(ActionTypeBuyItem, user.ID(), adventuria.ActionRequest{
		"item_id": itemId,
	})
	if err != nil {
		t.Fatalf("Test_Buy(): Error performing action: %s", err)
	}

	var invItemId string
	invItemId, ok := res.Data.(string)
	if !ok {
		t.Fatalf("Test_Buy(): Error getting item id from response: %s", err)
	}

	if user.Balance() != 0 {
		t.Fatalf("Test_Buy(): User balance = %d, want = 0", user.Balance())
	}

	if _, _, err = user.Inventory().UseItem(invItemId); err != nil {
		t.Fatalf("Test_Buy(): Error using item: %s", err)
	}
}

func createCellPointsDivideItem() (*core.Record, error) {
	effectRecord, err := createCellPointsDivideEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(tests.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Cell Points Divide")
	record.Set("effects", []string{effectRecord.Id})
	record.Set("icon", icon)
	record.Set("order", 1)
	record.Set("isUsingSlot", true)
	record.Set("canDrop", true)
	record.Set("type", "buff")
	record.Set("isRollable", true)
	record.Set("price", 2)
	err = adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func createCellPointsDivideEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Cell Points Divide")
	record.Set("type", "cellPointsDivide")
	record.Set("value", 2)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
