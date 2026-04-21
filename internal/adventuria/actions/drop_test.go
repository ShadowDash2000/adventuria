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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatalf("Test_Drop(): Error getting player: %s", err)
	}

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatalf("Test_Drop(): Error moving: %s", err)
	}

	player.Progress().AddPoints(2)
	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop(): Error action roll wheel: %s", err)
	}

	type testCompare struct {
		IsInJail    bool
		DropsInARow int
		Points      int
	}

	want := &testCompare{
		false,
		1,
		player.Progress().Points() + adventuria.GameSettings.PointsForDrop(),
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeDrop, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop(): Error action drop: %s", err)
	}

	got := &testCompare{
		player.Progress().IsInJail(),
		player.Progress().DropsInARow(),
		player.Progress().Points(),
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

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error getting player: %s", err)
	}

	player.Progress().SetIsInJail(true)

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error moving: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Drop_inJail(): Error action roll wheel: %s", err)
	}

	canDo := adventuria.GameActions.CanDo(ctx, player, ActionTypeDrop)
	if canDo {
		t.Fatalf("Test_Drop_inJail(): Expected that you can't drop in jail: %s", err)
	}
}
