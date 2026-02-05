package actions

import (
	"adventuria/internal/adventuria"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

type UpdateCommentAction struct {
	adventuria.ActionBase
}

func (a *UpdateCommentAction) CanDo(_ adventuria.ActionContext) bool {
	return true
}

func (a *UpdateCommentAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*adventuria.ActionResult, error) {
	requiredFields := []string{"action_id", "comment"}
	for _, field := range requiredFields {
		if _, ok := req[field]; !ok {
			return &adventuria.ActionResult{
				Success: false,
				Error:   fmt.Sprintf("request error: %s not specified", field),
			}, nil
		}
	}

	actionId, ok := req["action_id"].(string)
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: action_id is not string",
		}, nil
	}

	comment, ok := req["comment"].(string)
	if !ok {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: comment is not string",
		}, nil
	}

	record := &core.Record{}
	err := ctx.AppContext.App.
		RecordQuery(adventuria.GameCollections.Get(adventuria.CollectionActions)).
		AndWhere(
			dbx.HashExp{
				"user": ctx.User.ID(),
				"id":   actionId,
			},
		).
		AndWhere(
			dbx.Or(
				// TODO get rid of hard coded types
				dbx.HashExp{"type": "done"},
				dbx.HashExp{"type": "drop"},
				dbx.HashExp{"type": "reroll"},
			),
		).
		Limit(1).
		One(record)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "request error: can't find action with id",
		}, nil
	}

	action := adventuria.NewActionRecordFromRecord(record)
	action.SetComment(comment)

	err = ctx.AppContext.App.Save(record)
	if err != nil {
		return &adventuria.ActionResult{
			Success: false,
			Error:   "internal error",
		}, fmt.Errorf("update_comment.do(): %w", err)
	}

	return &adventuria.ActionResult{
		Success: true,
	}, nil
}

func (a *UpdateCommentAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}
