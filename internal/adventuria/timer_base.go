package adventuria

import (
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type TimerBase struct {
	core.BaseRecordProxy
	locator SettingsLocator
}

func NewBaseTimerFromRecord(locator SettingsLocator, record *core.Record) Timer {
	timer := &TimerBase{locator: locator}
	timer.SetProxyRecord(record)
	return timer
}

func NewTimer(locator SettingsLocator, userId string) (Timer, error) {
	record, err := locator.PocketBase().FindFirstRecordByFilter(
		TableTimers,
		"user.id = {:userId}",
		dbx.Params{"userId": userId},
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	timer := &TimerBase{locator: locator}
	if record != nil {
		timer.SetProxyRecord(record)
	} else {
		timer, err = CreateTimer(locator, userId, locator.Settings().TimerTimeLimit())
		if err != nil {
			return nil, err
		}
	}

	timer.bindHooks()

	return timer, nil
}

func (t *TimerBase) bindHooks() {
	t.locator.PocketBase().OnRecordAfterUpdateSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("user") == t.UserId() {
			t.SetProxyRecord(e.Record)
		}
		return e.Next()
	})
	t.locator.PocketBase().OnRecordAfterDeleteSuccess(TableTimers).BindFunc(func(e *core.RecordEvent) error {
		if e.Record.Id == t.Id {
			timersCollection, _ := t.locator.Collections().Get(TableTimers)
			t.SetProxyRecord(core.NewRecord(timersCollection))
		}
		return e.Next()
	})
}

func (t *TimerBase) Start() error {
	if t.IsActive() {
		return nil
	}
	if t.TimeLimit() < 0 {
		return errors.New("time limit is less than 0")
	}

	t.SetIsActive(true)
	t.SetStartTime(types.NowDateTime())

	return t.locator.PocketBase().Save(t)
}

func (t *TimerBase) Stop() error {
	if !t.IsActive() {
		return nil
	}

	t.SetIsActive(false)
	timePassed := t.TimePassed() + time.Now().Sub(t.StartTime().Time())
	t.SetTimePassed(timePassed)

	return t.locator.PocketBase().Save(t)
}

// GetTimeLeft returns time.Duration in seconds
func (t *TimerBase) GetTimeLeft() time.Duration {
	timeLeft := t.TimeLimit() - t.TimePassed()
	if t.IsActive() {
		timeLeft -= time.Now().Sub(t.StartTime().Time())
	}
	return timeLeft / time.Second
}

func (t *TimerBase) IsTimeExceeded() bool {
	return t.TimePassed() >= t.TimeLimit()
}

func (t *TimerBase) UserId() string {
	return t.GetString("user")
}

func (t *TimerBase) setUserId(userId string) {
	t.Set("user", userId)
}

func (t *TimerBase) IsActive() bool {
	return t.GetBool("isActive")
}

func (t *TimerBase) SetIsActive(active bool) {
	t.Set("isActive", active)
}

func (t *TimerBase) TimePassed() time.Duration {
	return time.Duration(t.GetInt("timePassed")) * time.Second
}

func (t *TimerBase) SetTimePassed(tp time.Duration) {
	t.Set("timePassed", int(tp/time.Second))
}

func (t *TimerBase) TimeLimit() time.Duration {
	return time.Duration(t.GetInt("timeLimit")) * time.Second
}

func (t *TimerBase) SetTimeLimit(tp time.Duration) {
	t.Set("timeLimit", int(tp/time.Second))
}

func (t *TimerBase) StartTime() types.DateTime {
	return t.GetDateTime("startTime")
}

func (t *TimerBase) SetStartTime(time types.DateTime) {
	t.Set("startTime", time)
}

func (t *TimerBase) AddSecondsTimeLimit(secs int) error {
	t.SetTimeLimit(t.TimeLimit() + (time.Duration(secs) * time.Second))
	return t.locator.PocketBase().Save(t)
}

func (t *TimerBase) Save() error {
	return t.locator.PocketBase().Save(t)
}

func CreateTimer(locator SettingsLocator, userId string, timeLimit int) (*TimerBase, error) {
	collection, err := locator.Collections().Get(TableTimers)
	if err != nil {
		return nil, err
	}

	timer := &TimerBase{locator: locator}
	timer.SetProxyRecord(core.NewRecord(collection))
	timer.Set("user", userId)
	timer.Set("timeLimit", timeLimit)
	timer.Set("timePassed", 0)
	timer.Set("isActive", false)
	err = locator.PocketBase().Save(timer)
	if err != nil {
		return nil, err
	}

	return timer, nil
}

func ResetAllTimers(locator SettingsLocator, timeLimit int, limitExceedPenalty int) error {
	records, err := locator.PocketBase().FindAllRecords(TableTimers)
	if err != nil {
		return err
	}

	for _, record := range records {
		timer := NewBaseTimerFromRecord(locator, record)

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
		err = timer.Save()
		if err != nil {
			return err
		}
	}

	return nil
}
