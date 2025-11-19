package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/tests"
	"reflect"
	"testing"
)

func Test_Drop(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Drop(): Error creating game: %s", err)
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatalf("Test_Drop(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Drop(): Error moving: %s", err)
	}

	err = adventuria.PocketBase.Save(user.LastAction().ProxyRecord())
	if err != nil {
		t.Fatalf("Test_Drop(): Error saving action: %s", err)
	}

	type testCompare struct {
		IsInJail    bool
		DropsInARow int
		Points      int
	}

	want := &testCompare{
		false,
		1,
		user.Points() + adventuria.GameSettings.PointsForDrop(),
	}

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ActionTypeDrop, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop(): Error action drop: %s", err)
	}

	got := &testCompare{
		user.IsInJail(),
		user.DropsInARow(),
		user.Points(),
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Test_Drop(): Want %v, got %v", want, got)
	}
}

func Test_Drop_inJail(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error creating game: %s", err)
	}

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error moving: %s", err)
	}

	user.SetIsInJail(true)

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error action roll wheel: %s", err)
	}

	canDo := adventuria.GameActions.CanDo(user, ActionTypeDrop)
	if canDo {
		t.Fatalf("Test_Drop_inJail(): Expected that you can't drop in jail: %s", err)
	}
}
