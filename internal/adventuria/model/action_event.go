package model

import (
	"context"
)

type ActionEventType string

type ActionEvent interface {
	Data() *ActionEventInfo
	Init(ctx context.Context, player *Player) error
}

type ActionEventData struct {
	Id         string
	Name       string
	Type       ActionEventType
	ActionType ActionType
	Value      string
}

type ActionEventInfo struct {
	data ActionEventData
}

func RestoreActionEvent(data ActionEventData) *ActionEventInfo {
	return &ActionEventInfo{data: data}
}

func (c *ActionEventInfo) ID() string {
	return c.data.Id
}

func (c *ActionEventInfo) Name() string {
	return c.data.Name
}

func (c *ActionEventInfo) Type() ActionEventType {
	return c.data.Type
}

func (c *ActionEventInfo) ActionType() ActionType {
	return c.data.ActionType
}

func (c *ActionEventInfo) Value() string {
	return c.data.Value
}
