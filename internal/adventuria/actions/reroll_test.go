package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/cells"
	"adventuria/internal/adventuria/effects"
	"testing"
)

func Test_Reroll(t *testing.T) {
	WithBaseActions()
	cells.WithBaseCells()
	effects.WithBaseEffects()

	game, err := adventuria.NewTestGame()
	if err != nil {
		t.Fatalf("Test_Reroll(): Error creating game: %s", err)
	}

	user, err := game.GetUserByName("user1")
	if err != nil {
		t.Fatalf("Test_Reroll(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Reroll(): Error moving: %s", err)
	}

	cell, ok := user.CurrentCell()
	if !ok {
		t.Fatal("Test_Reroll(): Current cell not found")
	}

	user.LastAction().SetCell(cell.ID())

	err = user.LastAction().Save()
	if err != nil {
		t.Fatalf("Test_Reroll(): Error saving action: %s", err)
	}

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}

	_, err = game.DoAction(ActionTypeReroll, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action done: %s", err)
	}

	_, err = game.DoAction(ActionTypeRollWheel, user.ID(), adventuria.ActionRequest{})
	if err != nil {
		t.Fatalf("Test_Reroll(): Error action roll wheel: %s", err)
	}
}
