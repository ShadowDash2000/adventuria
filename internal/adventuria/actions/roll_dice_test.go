package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"testing"
)

func Test_RollDice(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatalf("Test_RollDice(): Error creating game: %s", err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatalf("Test_RollDice(): Error getting user: %s", err)
	}

	initialCellsPassed := user.CellsPassed()

	_, err = game.DoAction(ActionTypeRollDice, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_RollDice(): Error action roll dice: %s", err)
	}

	if user.CellsPassed() <= initialCellsPassed {
		t.Fatalf("Test_RollDice(): Expected that cells passed increased, got %d", user.CellsPassed())
	}
}
