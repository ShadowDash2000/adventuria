package actions

import (
	"adventuria/internal/adventuria"
	"testing"
)

func Test_WithBaseActions(t *testing.T) {
	WithBaseActions()

	actionTypes := []adventuria.ActionType{
		ActionTypeRollDice,
		ActionTypeDone,
		ActionTypeReroll,
		ActionTypeDrop,
		ActionTypeRollWheel,
		ActionTypeRollItem,
		ActionTypeBuyItem,
		ActionTypeUpdateComment,
		ActionTypeRerollFilter,
	}

	got := 0
	for _, actionType := range actionTypes {
		action, err := adventuria.NewActionFromType(actionType)
		if err != nil {
			t.Fatalf("Test_WithBaseActions(): Error creating action: %s", err)
		}

		if action.Type() == actionType {
			got++
		}
	}

	expected := len(actionTypes)
	if got != expected {
		t.Fatalf("Test_WithBaseActions(): Expected %d actions, got %d", expected, got)
	}
}
