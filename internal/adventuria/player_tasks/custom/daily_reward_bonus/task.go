package daily_reward_bonus

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_tasks"
	"encoding/json"
	"fmt"
)

var _ model.PlayerTask = (*DailyRewardBonus)(nil)

const Type model.PlayerTaskType = "daily_reward_bonus"

type DailyRewardBonus struct {
	player_tasks.PlayerTaskBase
	payload Payload
}

func NewDef() player_tasks.PlayerTaskDef {
	return player_tasks.NewPlayerTask(
		Type,
		func(payload string) (model.PlayerTask, error) {
			task, err := NewFromJSON(payload)
			if err != nil {
				return nil, err
			}

			return task, nil
		},
	)
}

func New(payload Payload) (*DailyRewardBonus, error) {
	payload = payload.clone()

	err := payload.validate()
	if err != nil {
		return nil, err
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("daily_reward_bonus: failed to marshal payload: %w", err)
	}

	return &DailyRewardBonus{
		PlayerTaskBase: player_tasks.NewPlayerTaskBase(Type, "daily_reward", string(encoded)),
		payload:        payload,
	}, nil
}

func NewFromJSON(payload string) (*DailyRewardBonus, error) {
	var decoded Payload
	err := json.Unmarshal([]byte(payload), &decoded)
	if err != nil {
		return nil, fmt.Errorf("daily_reward_bonus: failed to unmarshal payload: %w", err)
	}

	err = decoded.validate()
	if err != nil {
		return nil, err
	}

	return &DailyRewardBonus{
		PlayerTaskBase: player_tasks.NewPlayerTaskBase(Type, "daily_reward", payload),
		payload:        decoded,
	}, nil
}
