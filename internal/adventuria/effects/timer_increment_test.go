package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

func Test_TimerIncrement(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	item, err := createTimerIncrementItem()
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

	wantTime := time.Duration(adventuria.GameSettings.TimerTimeLimit()+1000) * time.Second
	if user.Timer().TimeLimit() != wantTime {
		t.Fatalf("Test_TimerIncrement(): Time limit is %d, expected %d", user.Timer().TimeLimit(), wantTime)
	}
}

func createTimerIncrementItem() (*core.Record, error) {
	effectRecord, err := createTimerIncrementEffect()
	if err != nil {
		return nil, err
	}

	icon, err := filesystem.NewFileFromBytes(adventuria.Placeholder, "icon")
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionItems))
	record.Set("name", "Timer Increment")
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

func createTimerIncrementEffect() (*core.Record, error) {
	record := core.NewRecord(adventuria.GameCollections.Get(adventuria.CollectionEffects))
	record.Set("name", "Timer Increment")
	record.Set("type", "timerIncrement")
	record.Set("value", 1000)
	err := adventuria.PocketBase.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}
