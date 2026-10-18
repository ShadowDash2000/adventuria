package model

import (
	"adventuria/internal/adventuria/model"

	"github.com/pocketbase/pocketbase/tools/types"
)

type ActionImport struct {
	Action  *model.ActionInfo
	ID      string
	Created types.DateTime
	Updated types.DateTime
}
