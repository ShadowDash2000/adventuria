package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"
)

func Test_GiveWheelOnDone(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatal(err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatal(err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_GiveWheelOnDone(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	_, err = game.DoAction(actions.ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(actions.ActionTypeDone, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	want := 1
	if user.ItemWheelsCount() != want {
		t.Fatalf("Test_GiveWheelOnDone(): Wheels count is %d, expected %d", user.ItemWheelsCount(), want)
	}
}
