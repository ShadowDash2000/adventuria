package adventuria

import (
	"errors"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Timer struct {
	core.BaseRecordProxy
	app    core.App
	userId string
}

func NewTimer(userId string, app core.App) (*Timer, error) {
	t := &Timer{
		app:    app,
		userId: userId,
	}

	records, err := app.FindRecordsByFilter(
		TableTimers,
		"user.id = {:userId}",
		"",
		1,
		0,
		dbx.Params{"userId": userId},
	)
	if err != nil {
		return nil, err
	}

	if len(records) != 0 {
		t.SetProxyRecord(records[0])
	}

	t.bindHooks()

	return t, nil
}

func (t *Timer) bindHooks() {
	t.app.OnRecordAfterCreateSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == t.userId {
			t.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	t.app.OnRecordAfterUpdateSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == t.userId {
			t.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
}

func (t *Timer) Start() error {
	if t.IsActive() {
		return nil
	}
	if t.IsTimeExceeded() {
		return errors.New("timer exceeds limit")
	}

	t.SetIsActive(true)
	t.SetStartTime(types.NowDateTime())

	return t.app.Save(t)
}

func (t *Timer) Stop() error {
	if !t.IsActive() {
		return nil
	}

	t.SetIsActive(false)
	timePassed := t.TimePassed() + time.Now().Sub(t.StartTime().Time())
	t.SetTimePassed(timePassed)

	return t.app.Save(t)
}

// GetTimeLeft returns time.Duration in seconds
func (t *Timer) GetTimeLeft() time.Duration {
	timeLeft := t.TimeLimit() - t.TimePassed()
	if t.IsActive() {
		timeLeft -= time.Now().Sub(t.StartTime().Time())
	}
	return timeLeft / time.Second
}

func (t *Timer) IsTimeExceeded() bool {
	return t.TimePassed() >= t.TimeLimit()
}

func (t *Timer) IsActive() bool {
	return t.GetBool("isActive")
}

func (t *Timer) SetIsActive(active bool) {
	t.Set("isActive", active)
}

func (t *Timer) TimePassed() time.Duration {
	return time.Duration(t.GetInt("timePassed")) * time.Second
}

func (t *Timer) SetTimePassed(tp time.Duration) {
	t.Set("timePassed", int(tp/time.Second))
}

func (t *Timer) TimeLimit() time.Duration {
	return time.Duration(t.GetInt("timeLimit")) * time.Second
}

func (t *Timer) StartTime() types.DateTime {
	return t.GetDateTime("startTime")
}

func (t *Timer) SetStartTime(time types.DateTime) {
	t.Set("startTime", time)
}

func ClearAllTimers() error {
	return nil
}
