package v2

import (
	actionsRepo "adventuria/internal/adventuria/actions/repository"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	reviewsRepo "adventuria/internal/adventuria/reviews/repository"
	"adventuria/internal/adventuria/schema"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/google/go-querystring/query"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
)

func migrateActionsCommand(pb core.App) *cobra.Command {
	return &cobra.Command{
		Use:          "actions <base-url>",
		Example:      "migrate-data v2 actions http://127.0.0.1:8080",
		Short:        "Migrates actions from Adventuria v2 through Pocketbase HTTP API",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return migrateActions(pb, args[0])
		},
	}
}

const (
	actionsPath    = "/api/collections/actions/records"
	actionsPerPage = 500
)

func migrateActions(pb core.App, baseUrl string) error {
	res, err := getActions(baseUrl, 1, actionsPerPage)
	if err != nil {
		return err
	}

	err = saveActions(pb, res.Items)
	if err != nil {
		return err
	}

	for page := 2; page <= res.TotalPages; page++ {
		res, err := getActions(baseUrl, page, actionsPerPage)
		if err != nil {
			return err
		}

		err = saveActions(pb, res.Items)
		if err != nil {
			return err
		}
	}

	return nil
}

type getActionsQuery struct {
	Page    int    `url:"page"`
	PerPage int    `url:"perPage"`
	Sort    string `url:"sort,omitempty"`
	Expand  string `url:"expand,omitempty"`
	Fields  string `url:"fields,omitempty"`
}

type action struct {
	Id       string         `json:"id"`
	Created  types.DateTime `json:"created"`
	Updated  types.DateTime `json:"updated"`
	User     string         `json:"user"`
	Cell     string         `json:"cell"`
	Type     string         `json:"type"`
	Activity string         `json:"activity"`
	Comment  string         `json:"comment"`
	DiceRoll int            `json:"diceRoll"`

	Expand actionExpand `json:"expand"`
}

type actionExpand struct {
	User expandedActionUser `json:"user"`
}

type expandedActionUser struct {
	Name string `json:"name"`
}

type getActionsResponse struct {
	Items      []action `json:"items"`
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	TotalItems int      `json:"totalItems"`
	TotalPages int      `json:"totalPages"`
}

func getActions(baseUrl string, page, perPage int) (*getActionsResponse, error) {
	q, err := query.Values(getActionsQuery{
		Page:    page,
		PerPage: perPage,
		Sort:    "created",
		Expand:  "user",
		Fields:  "*,expand.user.name",
	})
	if err != nil {
		return nil, fmt.Errorf("encode actions query: %w", err)
	}

	return requestActions(baseUrl, q)
}

func requestActions(baseUrl string, q url.Values) (*getActionsResponse, error) {
	requestUrl, err := url.JoinPath(baseUrl, actionsPath)
	if err != nil {
		return nil, fmt.Errorf("build actions URL: %w", err)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("parse actions URL: %w", err)
	}
	parsedUrl.RawQuery = q.Encode()

	res, err := http.Get(parsedUrl.String())
	if err != nil {
		return nil, fmt.Errorf("request actions: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("request actions: unexpected HTTP status %s", res.Status)
	}

	var actionsResponse getActionsResponse
	if err = json.NewDecoder(res.Body).Decode(&actionsResponse); err != nil {
		return nil, fmt.Errorf("decode actions response: %w", err)
	}

	return &actionsResponse, nil
}

func saveActions(pb core.App, actions []action) error {
	for _, action := range actions {
		err := pb.RunInTransaction(func(txApp core.App) error {
			return saveAction(txApp, action)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func saveAction(pb core.App, action action) error {
	playerId, err := getPlayerIdByName(pb, action.Expand.User.Name)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return nil
		}
		return err
	}

	status, err := actionTypeToStatus(action.Type)
	if err != nil {
		return err
	}

	newAction, err := model.NewAction(model.ActionCreate{
		Player: playerId,
		Cell:   action.Cell,
		Status: status,
	})
	if err != nil {
		return err
	}

	// we must save review if action type is one of those types to preserve edit ability, even if comment is empty
	mustSaveReview := action.Comment != "" || slices.Contains([]string{"done", "drop", "reroll"}, action.Type)
	if mustSaveReview {
		reviewRecord, err := saveReview(pb, action.Comment)
		if err != nil {
			return err
		}

		newAction.SetReview(reviewRecord.Id)
	}

	newAction.SetActivity(action.Activity)
	newAction.SetCellsPassed(action.DiceRoll)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionActions)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	err = actionsRepo.ActionToRecord(newAction, record)
	if err != nil {
		return err
	}

	record.SetRaw("created", action.Created)
	record.SetRaw("updated", action.Updated)

	return pb.Save(record)
}

var actionTypeToActionStatus = map[string]model.ActionStatus{
	"move":           model.ActionStatusMove,
	"rollDice":       model.ActionStatusRollDice,
	"done":           model.ActionStatusDone,
	"reroll":         model.ActionStatusReroll,
	"drop":           model.ActionStatusDrop,
	"rollWheel":      model.ActionStatusRollWheel,
	"rollItemOnCell": model.ActionStatusRollItemOnCell,
	"teleport":       model.ActionStatusTeleport,
}

func actionTypeToStatus(actionType string) (model.ActionStatus, error) {
	status, ok := actionTypeToActionStatus[actionType]
	if !ok {
		return "", fmt.Errorf("unknown action type: %s", actionType)
	}

	return status, nil
}

func saveReview(pb core.App, comment string) (*core.Record, error) {
	collection, err := pb.FindCollectionByNameOrId(schema.CollectionReviews)
	if err != nil {
		return nil, err
	}

	review, err := model.NewReview(comment, 0)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	reviewsRepo.ReviewToRecord(review, record)

	err = pb.Save(record)
	if err != nil {
		return nil, err
	}

	return record, nil
}

func getPlayerIdByName(pb core.App, name string) (string, error) {
	var id string
	err := pb.DB().
		Select(schema.PlayerSchema.Id).
		From(schema.CollectionPlayers).
		Where(dbx.HashExp{
			schema.PlayerSchema.Name: name,
		}).
		Limit(1).
		Row(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errs.ErrPlayerNotFound
		}
		return "", err
	}

	return id, nil
}
