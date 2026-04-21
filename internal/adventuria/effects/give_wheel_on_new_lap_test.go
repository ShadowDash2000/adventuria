package effects

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/actions"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/tests"
	"testing"
)

func Test_GiveWheelOnNewLap(t *testing.T) {
	actions.WithBaseActions()
	cells.WithBaseCells()
	WithBaseEffects()

	_, err := tests.NewGameTest()
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

	_, err = player.Move(ctx, 8)
	if err != nil {
		t.Fatal(err)
	}

	want := 2
	if player.Progress().ItemWheelsCount() != want {
		t.Fatalf("Test_GiveWheelOnNewLap(): Wheels count is %d, expected %d", player.Progress().ItemWheelsCount(), want)
	}
}
