package model

import (
	"adventuria/pkg/event"
	"context"
	"errors"
	"time"
)

type PlayerTaskType string
type PlayerTaskTrigger string

type PlayerTaskCompleteFunc func(ctx context.Context) error

type PlayerTask interface {
	Type() PlayerTaskType
	Trigger() PlayerTaskTrigger
	Payload() string
}

type PlayerTaskSubscriber interface {
	Subscribe(ctx context.Context, events *Events, player *Player, complete PlayerTaskCompleteFunc, payload string) ([]event.Unsubscribe, error)
}

type PlayerTaskData struct {
	Id        string
	Completed time.Time
	Progress  string
	Type      PlayerTaskType
	Trigger   PlayerTaskTrigger
	Payload   string
}

type PlayerTaskInfo struct {
	data  PlayerTaskData
	isNew bool
}

type PlayerTaskCreate struct {
	ProgressId string
	Task       PlayerTask
}

func NewPlayerTask(data PlayerTaskCreate) (*PlayerTaskInfo, error) {
	if data.ProgressId == "" {
		return nil, errors.New("player_task: progress_id is empty")
	}
	if data.Task == nil {
		return nil, errors.New("player_task: task is empty")
	}
	if data.Task.Type() == "" {
		return nil, errors.New("player_task: type is empty")
	}

	return &PlayerTaskInfo{
		data: PlayerTaskData{
			Progress: data.ProgressId,
			Type:     data.Task.Type(),
			Trigger:  data.Task.Trigger(),
			Payload:  data.Task.Payload(),
		},
		isNew: true,
	}, nil
}

func RestorePlayerTaskInfo(data PlayerTaskData) *PlayerTaskInfo {
	return &PlayerTaskInfo{
		data:  data,
		isNew: false,
	}
}

func (p *PlayerTaskInfo) IsNew() bool {
	return p.isNew
}

func (p *PlayerTaskInfo) ID() string {
	return p.data.Id
}

func (p *PlayerTaskInfo) Completed() time.Time {
	return p.data.Completed
}

func (p *PlayerTaskInfo) Progress() string {
	return p.data.Progress
}

func (p *PlayerTaskInfo) Type() PlayerTaskType {
	return p.data.Type
}

func (p *PlayerTaskInfo) Trigger() PlayerTaskTrigger {
	return p.data.Trigger
}

func (p *PlayerTaskInfo) Payload() string {
	return p.data.Payload
}
