package model

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type OutboxType string

type Outbox interface {
	Verifiable
	Process(ctx context.Context, outbox *OutboxInfo) error
}

type OutboxData struct {
	Id      string
	Type    OutboxType
	Payload string
	Status  OutboxStatus
}

type OutboxInfo struct {
	data  OutboxData
	isNew bool
}

type OutBoxCreate struct {
	Type    OutboxType
	Payload string
}

func NewOutbox(id uuid.UUID, data OutBoxCreate) (*OutboxInfo, error) {
	if id == uuid.Nil {
		return nil, errors.New("outbox: id cannot be nil")
	}
	if data.Type == "" {
		return nil, errors.New("outbox: type is empty")
	}

	return &OutboxInfo{
		data: OutboxData{
			Id:      id.String(),
			Type:    data.Type,
			Payload: data.Payload,
			Status:  OutboxStatusPending,
		},
		isNew: true,
	}, nil
}

func RestoreOutbox(data OutboxData) *OutboxInfo {
	return &OutboxInfo{
		data:  data,
		isNew: false,
	}
}

func (o *OutboxInfo) IsNew() bool {
	return o.isNew
}

func (o *OutboxInfo) ID() string {
	return o.data.Id
}

func (o *OutboxInfo) Type() OutboxType {
	return o.data.Type
}

func (o *OutboxInfo) Payload() string {
	return o.data.Payload
}

func (o *OutboxInfo) Status() OutboxStatus {
	return o.data.Status
}

func (o *OutboxInfo) SetStatus(status OutboxStatus) {
	o.data.Status = status
}
