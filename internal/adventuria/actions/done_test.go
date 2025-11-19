package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/tests"
	"reflect"
	"testing"
)

func Test_Done(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Done(): Error creating game: %s", err)
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatalf("Test_Done(): Error getting user: %s", err)
	}

	user.SetIsInJail(true)
	user.SetDropsInARow(2)

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Done(): Error moving: %s", err)
	}

	err = adventuria.PocketBase.Save(user.LastAction().ProxyRecord())
	if err != nil {
		t.Fatalf("Test_Done(): Error saving action: %s", err)
	}

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Done(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ActionTypeDone, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Done(): Error action done: %s", err)
	}

	type testCompare struct {
		IsInJail    bool
		DropsInARow int
		Points      int
	}

	want := &testCompare{
		false,
		0,
		20,
	}
	got := &testCompare{
		user.IsInJail(),
		user.DropsInARow(),
		user.Points(),
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Test_Done(): Want %v, got %v", want, got)
	}
}
