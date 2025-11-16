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

	user, err := adventuria.GameUsers.GetByName("user1")
	if err != nil {
		t.Fatalf("Test_Reroll(): Error getting user: %s", err)
	}

	_, err = user.Move(1)
	if err != nil {
		t.Fatalf("Test_Reroll(): Error moving: %s", err)
	}

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
