package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"reflect"
	"testing"
)

func Test_Drop(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatalf("Test_Drop(): Error creating game: %s", err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatalf("Test_Drop(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Drop(): Error moving: %s", err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_Drop(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	err = user.LastAction().Save()
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

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error creating game: %s", err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error moving: %s", err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_Drop_inJail(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	user.SetIsInJail(true)

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ActionTypeDrop, user.ID(), adventuria.ActionRequest{})
	if err == nil {
		t.Fatalf("Test_Drop_inJail(): Expected that you can't drop in jail: %s", err)
	}
}
