package model

type OutboxStatus string

const (
	OutboxStatusPending    OutboxStatus = "pending"
	OutboxStatusProcessing OutboxStatus = "processing"
	OutboxStatusFailed     OutboxStatus = "failed"
	OutboxStatusCompleted  OutboxStatus = "completed"
)
