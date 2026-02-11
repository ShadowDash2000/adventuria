package adventuria

import (
	"adventuria/pkg/result"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type Action interface {
	Type() ActionType
	Categories() []string
	InCategory(string) bool
	InCategories(categories []string) bool
	CanDo(ActionContext) bool
	Do(ActionContext, ActionRequest) (*result.Result, error)
	GetVariants(ActionContext) any

	setType(ActionType)
}

type ActionContext struct {
	AppContext
	User User
}

type ActionRecord interface {
	core.RecordProxy

	ID() string
	User() string
	SetUser(string)
	CellId() string
	setCell(string)
	Comment() string
	SetComment(string)
	Activity() string
	SetActivity(string)
	Type() ActionType
	SetType(ActionType)
	DiceRoll() int
	SetDiceRoll(int)
	UsedItemAppend(string)
	ItemsList() ([]string, error)
	SetItemsList([]string)
	CanMove() bool
	SetCanMove(bool)
	CustomActivityFilter() (*CustomActivityFilter, error)
	SetCustomActivityFilter(CustomActivityFilter)
	ClearCustomActivityFilter()
	MarkAsNew()
}

type CustomActivityFilter struct {
	Platforms       []string
	Developers      []string
	Publishers      []string
	Genres          []string
	Tags            []string
	Themes          []string
	MinPrice        int
	MaxPrice        int
	ReleaseDateFrom types.DateTime
	ReleaseDateTo   types.DateTime
	MinCampaignTime float64
	MaxCampaignTime float64
}

type ActionType string

const (
	ActionTypeNone ActionType = "none"
	ActionTypeMove ActionType = "move"
)

type ActionRequest map[string]any

var actionsList = map[ActionType]ActionDef{
	ActionTypeNone: NewAction(ActionTypeNone, &NoneAction{}),
}

type ActionDef struct {
	Type       ActionType
	Categories []string
	New        func() Action
}

func RegisterActions(actions []ActionDef) {
	for _, actionDef := range actions {
		actionsList[actionDef.Type] = actionDef
	}
}

func NewAction(t ActionType, a Action, categories ...string) ActionDef {
	return ActionDef{
		Type:       t,
		Categories: categories,
		New: func() Action {
			a.setType(t)
			return a
		},
	}
}
