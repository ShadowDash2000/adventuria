package adventuria

import (
	"github.com/pocketbase/pocketbase/tools/filesystem"
	"github.com/pocketbase/pocketbase/tools/types"
	"time"
)

type Game interface {
	Init()
	Settings() *Settings
	Event() Event
	GetUser(userId string) (*User, error)
	ChooseGame(game string, userId string) error
	GetNextStepType(userId string) (string, error)
	UpdateAction(actionId string, comment string, file *filesystem.File, userId string) error
	Reroll(comment string, file *filesystem.File, userId string) error
	Roll(userId string) (int, []int, Cell, error)
	Drop(comment string, file *filesystem.File, userId string) error
	Done(comment string, file *filesystem.File, userId string) error
	RollWheel(userId string) (*WheelRollResult, error)
	GetLastAction(userId string) (bool, Action, error)
	GetItemsEffects(userId string, event EffectUse) (*Effects, error)
	UseItem(userId, itemId string) error
	DropItem(userId, itemId string) error
	StartTimer(userId string) error
	StopTimer(userId string) error
	GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error)
}
