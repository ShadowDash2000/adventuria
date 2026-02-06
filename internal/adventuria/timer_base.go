package adventuria

import (
	"adventuria/internal/adventuria/schema"
	"database/sql"
	"errors"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type TimerBase struct {
	core.BaseRecordProxy
}

func NewBaseTimerFromRecord(record *core.Record) Timer {
	timer := &TimerBase{}
	timer.SetProxyRecord(record)
	return timer
}

func NewTimer(ctx AppContext, userId string) (Timer, error) {
	var record core.Record
	err := ctx.App.
		RecordQuery(schema.CollectionTimers).
		Where(dbx.HashExp{schema.TimerSchema.User: userId}).
		One(&record)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	timer := &TimerBase{}
	if errors.Is(err, sql.ErrNoRows) {
		timer, err = CreateTimer(ctx, userId, GameSettings.TimerTimeLimit())
		if err != nil {
			return nil, err
		}
	} else {
		timer.SetProxyRecord(&record)
	}

	return timer, nil
}

func (t *TimerBase) Start(ctx AppContext) error {
	if t.IsActive() {
		return nil
	}
	if t.TimeLimit() < 0 {
		return errors.New("time limit is less than 0")
	}

	t.SetIsActive(true)
	t.SetStartTime(types.NowDateTime())

	return ctx.App.Save(t)
}

func (t *TimerBase) Stop(ctx AppContext) error {
	if !t.IsActive() {
		return nil
	}

	t.SetIsActive(false)
	timePassed := t.TimePassed() + time.Now().Sub(t.StartTime().Time())
	t.SetTimePassed(timePassed)

	return ctx.App.Save(t)
}

// GetTimeLeft returns the time left in seconds
func (t *TimerBase) GetTimeLeft() int64 {
	timeLeft := t.TimeLimit() - t.TimePassed()
	if t.IsActive() {
		timeLeft -= time.Now().Sub(t.StartTime().Time())
	}
	return int64(timeLeft / time.Second)
}

func (t *TimerBase) IsTimeExceeded() bool {
	return t.TimePassed() >= t.TimeLimit()
}

func (t *TimerBase) UserId() string {
	return t.GetString(schema.TimerSchema.User)
}

func (t *TimerBase) setUserId(userId string) {
	t.Set(schema.TimerSchema.User, userId)
}

func (t *TimerBase) IsActive() bool {
	return t.GetBool(schema.TimerSchema.IsActive)
}

func (t *TimerBase) SetIsActive(active bool) {
	t.Set(schema.TimerSchema.IsActive, active)
}

func (t *TimerBase) TimePassed() time.Duration {
	return time.Duration(t.GetInt(schema.TimerSchema.TimePassed)) * time.Second
}

func (t *TimerBase) SetTimePassed(tp time.Duration) {
	t.Set(schema.TimerSchema.TimePassed, int(tp/time.Second))
}

// TimeLimit returns time.Duration in seconds
func (t *TimerBase) TimeLimit() time.Duration {
	return time.Duration(t.GetInt(schema.TimerSchema.TimeLimit)) * time.Second
}

func (t *TimerBase) SetTimeLimit(tp time.Duration) {
	t.Set(schema.TimerSchema.TimeLimit, int(tp/time.Second))
}

func (t *TimerBase) StartTime() types.DateTime {
	return t.GetDateTime(schema.TimerSchema.StartTime)
}

func (t *TimerBase) SetStartTime(time types.DateTime) {
	t.Set(schema.TimerSchema.StartTime, time)
}

func (t *TimerBase) AddSecondsTimeLimit(ctx AppContext, secs int) error {
	t.SetTimeLimit(t.TimeLimit() + (time.Duration(secs) * time.Second))
	return ctx.App.Save(t)
}

func CreateTimer(ctx AppContext, userId string, timeLimit int) (*TimerBase, error) {
	timer := &TimerBase{}
	timer.SetProxyRecord(core.NewRecord(GameCollections.Get(schema.CollectionTimers)))
	timer.Set(schema.TimerSchema.User, userId)
	timer.Set(schema.TimerSchema.TimeLimit, timeLimit)
	timer.Set(schema.TimerSchema.TimePassed, 0)
	timer.Set(schema.TimerSchema.IsActive, false)
	err := ctx.App.Save(timer)
	if err != nil {
		return nil, err
	}

	return timer, nil
}

func ResetAllTimers(ctx AppContext, timeLimit int, limitExceedPenalty int) error {
	records, err := ctx.App.FindAllRecords(schema.CollectionTimers)
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
		err = ctx.App.Save(timer.ProxyRecord())
		if err != nil {
			return err
		}
	}

	return nil
}
