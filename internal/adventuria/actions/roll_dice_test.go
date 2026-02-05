package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/tests"
	"testing"
)

func Test_RollDice(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := tests.NewGameTest()
	if err != nil {
		t.Fatalf("Test_RollDice(): Error creating game: %s", err)
	}

	ctx := adventuria.AppContext{
		App: adventuria.PocketBase,
	}
	user, err := adventuria.GameUsers.GetByName(ctx, "user1")
	if err != nil {
		t.Fatalf("Test_RollDice(): Error getting user: %s", err)
	}

	initialCellsPassed := user.CellsPassed()

	_, err = game.DoAction(ctx.App, user.ID(), ActionTypeRollDice, adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_RollDice(): Error action roll dice: %s", err)
	}

	if user.CellsPassed() <= initialCellsPassed {
		t.Fatalf("Test_RollDice(): Expected that cells passed increased, got %d", user.CellsPassed())
	}
}
