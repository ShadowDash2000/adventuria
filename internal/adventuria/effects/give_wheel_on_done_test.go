package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"testing"
)

func Test_GiveWheelOnDone(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	game, err := tests.NewGameTest()
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

	err = user.LastAction().Save()
	if err != nil {
		t.Fatalf("Test_GiveWheelOnDone(): Error saving action: %s", err)
	}

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
