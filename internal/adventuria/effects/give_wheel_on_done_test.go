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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatal(err)
	}

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), actions.ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), actions.ActionTypeDone, adventuria.ActionRequest{})
	if err != nil {
		t.Fatal(err)
	}

	want := 1
	if player.Progress().ItemWheelsCount() != want {
		t.Fatalf("Test_GiveWheelOnDone(): Wheels count is %d, expected %d", player.Progress().ItemWheelsCount(), want)
	}
}
