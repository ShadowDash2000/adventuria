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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatalf("Test_Done(): Error getting player: %s", err)
	}

	player.Progress().SetIsInJail(true)
	player.Progress().SetDropsInARow(2)

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatalf("Test_Done(): Error moving: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Done(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeDone, adventuria.ActionRequest{})
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
		player.Progress().IsInJail(),
		player.Progress().DropsInARow(),
		player.Progress().Points(),
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Test_Done(): Want %v, got %v", want, got)
	}
}
