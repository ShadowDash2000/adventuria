package adventuria

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Timer interface {
	core.RecordProxy
	Start(ctx AppContext) error
	Stop(ctx AppContext) error
	GetTimeLeft() int64
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
	AddSecondsTimeLimit(ctx AppContext, secs int) error
}
