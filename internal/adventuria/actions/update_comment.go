package actions

import (
	"adventuria/internal/adventuria"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/result"
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

func (a *UpdateCommentAction) Do(ctx adventuria.ActionContext, req adventuria.ActionRequest) (*result.Result, error) {
	requiredFields := []string{"action_id", "comment"}
	for _, field := range requiredFields {
		if _, ok := req[field]; !ok {
			return result.Err(fmt.Sprintf("request error: %s not specified", field)), nil
		}
	}

	actionId, ok := req["action_id"].(string)
	if !ok {
		return result.Err("action_id is not string"), nil
	}

	comment, ok := req["comment"].(string)
	if !ok {
		return result.Err("comment is not string"), nil
	}

	var record core.Record
	err := ctx.AppContext.App.
		RecordQuery(adventuria.GameCollections.Get(schema.CollectionActions)).
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
		One(&record)
	if err != nil {
		return result.Err("can't find action with id"), nil
	}

	record.Set(schema.ActionSchema.Comment, comment)
	err = ctx.AppContext.App.Save(&record)
	if err != nil {
		return result.Err("internal error: failed to update action's comment"),
			fmt.Errorf("update_comment.do(): %w", err)
	}

	return result.Ok(), nil
}

func (a *UpdateCommentAction) GetVariants(_ adventuria.ActionContext) any {
	return nil
}
