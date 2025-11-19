package adventuria

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Timer interface {
	core.RecordProxy
	Start() error
	Stop() error
	GetTimeLeft() time.Duration
	IsTimeExceeded() bool
	UserId() string
	IsActive() bool
	SetIsActive(active bool)
	TimePassed() time.Duration
	SetTimePassed(tp time.Duration)
	TimeLimit() time.Duration
	SetTimeLimit(tp time.Duration)
	StartTime() types.DateTime
	SetStartTime(time types.DateTime)
	AddSecondsTimeLimit(secs int) error
}
