package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/tests"
	"testing"
)

func Test_Reroll(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_Reroll(): Error creating game: %s", err)
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	player, err := adventuria.GamePlayers.GetByName(ctx, "player1")
	if err != nil {
		t.Fatalf("Test_Reroll(): Error getting player: %s", err)
	}

	_, err = player.Move(ctx, 1)
	if err != nil {
		t.Fatalf("Test_Reroll(): Error moving: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeReroll, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action done: %s", err)
	}

	_, err = game.DoAction(ctx.App, player.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}
}
