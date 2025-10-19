package adventuria

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type GameFilterRecord interface {
	core.RecordProxy

	ID() string
	Name() string
	Platforms() []string
	Developers() []string
	Publishers() []string
	Genres() []string
	Tags() []string
	MinPrice() int
	MaxPrice() int
	ReleaseDateFrom() types.DateTime
	ReleaseDateTo() types.DateTime
	MinCampaignTime() float64
	MaxCampaignTime() float64
	Games() []string
}
