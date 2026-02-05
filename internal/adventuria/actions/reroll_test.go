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
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatalf("Test_Reroll(): Error getting user: %s", err)
	}

	_, err = user.Move(ctx, 1)
	if err != nil {
		t.Fatalf("Test_Reroll(): Error moving: %s", err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), ActionTypeReroll, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action done: %s", err)
	}

	_, err = game.DoAction(ctx.App, user.ID(), ActionTypeRollWheel, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}
}
