package adventuria

import (
	"adventuria/pkg/collections"
	"time"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Game interface {
	OnServe(fn func(se *core.ServeEvent) error)
	Start() error

	Init()
	GetUser(userId string) (User, error)
	GetUserByName(name string) (User, error)
	DoAction(actionType ActionType, userId string, req ActionRequest) (*ActionResult, error)
	UpdateAction(actionId string, comment string, userId string) error
	UseItem(userId, itemId string) error
	DropItem(userId, itemId string) error
	StartTimer(userId string) error
	StopTimer(userId string) error
	GetTimeLeft(userId string) (time.Duration, bool, types.DateTime, error)
	GetAvailableActions(userId string) ([]ActionType, error)
}

var (
	PocketBase      core.App
	GameCells       *Cells
	GameItems       *Items
	GameCollections *collections.Collections
	GameSettings    *Settings
)
