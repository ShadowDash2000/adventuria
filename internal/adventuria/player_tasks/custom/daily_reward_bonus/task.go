package daily_reward_bonus

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/player_tasks"
	"bytes"
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
		func(data []byte) (model.PlayerTask, error) {
			task, err := NewFromJSON(data)
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

	encoded, err := payload.Marshal()
	if err != nil {
		return nil, err
	}

	return &DailyRewardBonus{
		PlayerTaskBase: player_tasks.NewPlayerTaskBase(Type, "daily_reward", encoded),
		payload:        payload,
	}, nil
}

func NewFromJSON(data []byte) (*DailyRewardBonus, error) {
	payload, err := NewPayloadFromJSON(data)
	if err != nil {
		return nil, err
	}

	return &DailyRewardBonus{
		PlayerTaskBase: player_tasks.NewPlayerTaskBase(Type, "daily_reward", bytes.Clone(data)),
		payload:        payload,
	}, nil
}
