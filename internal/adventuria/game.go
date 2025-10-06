package adventuria

import (
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game interface {
	ServiceLocator

	OnServe(fn func(se *core.ServeEvent) error)
	Start() error

	Init()
	GetUser(userId string) (User, error)
	GetNextStepType(userId string) (string, error)
	DoAction(actionType, userId string, req ActionRequest) (*ActionResult, error)
	UpdateAction(actionId string, comment string, userId string) error
	GetLastAction(userId string) (bool, Action, error)
	UseItem(userId, itemId string) error
	DropItem(userId, itemId string) error
	StartTimer(userId string) error
	StopTimer(userId string) error
	GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error)
}
