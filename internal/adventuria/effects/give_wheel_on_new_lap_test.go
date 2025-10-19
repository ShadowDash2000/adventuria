package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"testing"
)

func Test_GiveWheelOnNewLap(t *testing.T) {
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

	_, err = user.Move(7)
	if err != nil {
		t.Fatal(err)
	}

	want := 2
	if user.ItemWheelsCount() != want {
		t.Fatalf("Test_GiveWheelOnNewLap(): Wheels count is %d, expected %d", user.ItemWheelsCount(), want)
	}
}
