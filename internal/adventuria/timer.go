package adventuria

import (
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Timer struct {
	core.BaseRecordProxy
}

func NewBaseTimerFromRecord(record *core.Record) *Timer {
	timer := &Timer{}
	timer.SetProxyRecord(record)
	return timer
}

func NewTimer(userId string) (*Timer, error) {
	record, err := PocketBase.FindFirstRecordByFilter(
		TableTimers,
		"user.id = {:userId}",
		dbx.Params{"userId": userId},
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	timer := &Timer{}
	if record != nil {
		timer.SetProxyRecord(record)
	} else {
		timer, err = CreateTimer(userId, GameSettings.TimerTimeLimit())
		if err != nil {
			return nil, err
		}
	}

	timer.bindHooks()

	return timer, nil
}

func (t *Timer) bindHooks() {
	PocketBase.OnRecordAfterUpdateSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == t.UserId() {
			t.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	PocketBase.OnRecordAfterDeleteSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == t.Id {
			timersCollection, _ := GameCollections.Get(TableTimers)
			t.SetProxyRecord(core.NewRecord(timersCollection))
		}
		return e.Next()
	})
}

func (t *Timer) Start() error {
	if t.IsActive() {
		return nil
	}
	if t.TimeLimit() < 0 {
		return errors.New("time limit is less than 0")
	}

	t.SetIsActive(true)
	t.SetStartTime(types.NowDateTime())

	return PocketBase.Save(t)
}

func (t *Timer) Stop() error {
	if !t.IsActive() {
		return nil
	}

	t.SetIsActive(false)
	timePassed := t.TimePassed() + time.Now().Sub(t.StartTime().Time())
	t.SetTimePassed(timePassed)

	return PocketBase.Save(t)
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

func (t *Timer) UserId() string {
	return t.GetString("user")
}

func (t *Timer) setUserId(userId string) {
	t.Set("user", userId)
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

func (t *Timer) SetTimeLimit(tp time.Duration) {
	t.Set("timeLimit", int(tp/time.Second))
}

func (t *Timer) StartTime() types.DateTime {
	return t.GetDateTime("startTime")
}

func (t *Timer) SetStartTime(time types.DateTime) {
	t.Set("startTime", time)
}

func (t *Timer) AddSecondsTimeLimit(secs int) error {
	t.SetTimeLimit(t.TimeLimit() + (time.Duration(secs) * time.Second))
	return PocketBase.Save(t)
}

func CreateTimer(userId string, timeLimit int) (*Timer, error) {
	collection, err := GameCollections.Get(TableTimers)
	if err != nil {
		return nil, err
	}

	timer := &Timer{}
	timer.SetProxyRecord(core.NewRecord(collection))
	timer.Set("user", userId)
	timer.Set("timeLimit", timeLimit)
	timer.Set("timePassed", 0)
	timer.Set("isActive", false)
	err = PocketBase.Save(timer)
	if err != nil {
		return nil, err
	}

	return timer, nil
}

func ResetAllTimers(timeLimit int, limitExceedPenalty int) error {
	records, err := PocketBase.FindAllRecords(TableTimers)
	if err != nil {
		return err
	}

	for _, record := range records {
		timer := NewBaseTimerFromRecord(record)

		timePassed := timer.TimePassed()
		if timer.IsActive() {
			timePassed += time.Now().Sub(timer.StartTime().Time())

			timer.SetStartTime(types.NowDateTime())
		} else {
			timer.SetStartTime(types.DateTime{})
		}

		newTimeLimit := time.Duration(timeLimit) * time.Second
		if timePassed > timer.TimeLimit() {
			newTimeLimit -= (timePassed - timer.TimeLimit()) * time.Duration(limitExceedPenalty)
		}

		if newTimeLimit < 0 {
			timer.SetIsActive(false)
		}

		timer.SetTimeLimit(newTimeLimit)
		timer.SetTimePassed(0)
		err = PocketBase.Save(timer)
		if err != nil {
			return err
		}
	}

	return nil
}
