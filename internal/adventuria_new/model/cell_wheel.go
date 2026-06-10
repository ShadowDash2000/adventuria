package model

type RollWheelRequest map[string]any

type WheelRollResult struct {
	WinnerId string `json:"winnerId"`
}
