package v2

import (
	actionsRepo "adventuria/internal/adventuria/actions/repository"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
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

func migrateActionsCommand(pb core.App, reviews reviewsService, items itemsService) *cobra.Command {
	migrator := newActionsMigrator(pb, reviews, items)

	return &cobra.Command{
		Use:          "actions <base-url>",
		Example:      "migrate-data v2 actions http://127.0.0.1:8080",
		Short:        "Migrates actions from Adventuria v2 through Pocketbase HTTP API",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return migrator.migrate(command.Context(), args[0])
		},
	}
}

const (
	actionsPath    = "/api/collections/actions/records"
	actionsPerPage = 500
	itemsPath      = "/api/collections/items/records"
	itemsPerPage   = 100
)

type actionsMigrator struct {
	pb       core.App
	itemsMap map[string]item
	reviews  reviewsService
	items    itemsService
}

func newActionsMigrator(pb core.App, reviews reviewsService, items itemsService) *actionsMigrator {
	return &actionsMigrator{
		pb:      pb,
		reviews: reviews,
		items:   items,
	}
}

func (a *actionsMigrator) migrate(ctx context.Context, baseUrl string) error {
	items, err := getAllItems(ctx, baseUrl)
	if err != nil {
		return err
	}

	a.itemsMap = make(map[string]item, len(items))
	for _, item := range items {
		a.itemsMap[item.Id] = item
	}

	res, err := getActions(ctx, baseUrl, 1, actionsPerPage)
	if err != nil {
		return err
	}

	err = a.saveActions(ctx, res.Items)
	if err != nil {
		return err
	}

	for page := 2; page <= res.TotalPages; page++ {
		res, err := getActions(ctx, baseUrl, page, actionsPerPage)
		if err != nil {
			return err
		}

		err = a.saveActions(ctx, res.Items)
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
	Id        string         `json:"id"`
	Created   types.DateTime `json:"created"`
	Updated   types.DateTime `json:"updated"`
	User      string         `json:"user"`
	Cell      string         `json:"cell"`
	Type      string         `json:"type"`
	Activity  string         `json:"activity"`
	Comment   string         `json:"comment"`
	DiceRoll  int            `json:"diceRoll"`
	UsedItems []string       `json:"used_items"`

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

func getActions(ctx context.Context, baseUrl string, page, perPage int) (*getActionsResponse, error) {
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

	requestUrl, err := url.JoinPath(baseUrl, actionsPath)
	if err != nil {
		return nil, fmt.Errorf("build actions URL: %w", err)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("parse actions URL: %w", err)
	}
	parsedUrl.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
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

func (a *actionsMigrator) saveActions(ctx context.Context, actions []action) error {
	for _, action := range actions {
		err := pbtransaction.RunInTransaction(ctx, a.pb, func(ctx context.Context, txApp core.App) error {
			return a.saveAction(ctx, action)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *actionsMigrator) saveAction(ctx context.Context, action action) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, a.pb)

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
		review, err := a.reviews.Create(ctx, reviews.CreateInput{
			Comment: action.Comment,
			Score:   0,
		})
		if err != nil {
			return err
		}

		newAction.SetReview(review.ID())
	}

	usedItemsState, err := a.oldItemsIDsToState(ctx, action.UsedItems)
	if err != nil {
		return err
	}
	newAction.State().UsedItems = usedItemsState

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

func (a *actionsMigrator) oldItemsIDsToState(ctx context.Context, oldIDs []string) (model.ActionUsedItemsState, error) {
	var res model.ActionUsedItemsState
	for _, oldID := range oldIDs {
		oldItem, ok := a.itemsMap[oldID]
		if !ok {
			return nil, fmt.Errorf("item with id %s not found", oldID)
		}

		item, err := a.items.GetByName(ctx, oldItem.Name)
		if err != nil {
			return nil, err
		}

		res = append(res, model.ActionUsedItemState{
			Id: item.ID(),
		})
	}

	return res, nil
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

type getItemsQuery struct {
	Page    int    `url:"page"`
	PerPage int    `url:"perPage"`
	Fields  string `url:"fields,omitempty"`
}

type item struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type getItemsResponse struct {
	Items      []item `json:"items"`
	Page       int    `json:"page"`
	PerPage    int    `json:"perPage"`
	TotalItems int    `json:"totalItems"`
	TotalPages int    `json:"totalPages"`
}

func getItems(ctx context.Context, baseUrl string, page, perPage int) (*getItemsResponse, error) {
	q, err := query.Values(getItemsQuery{
		Page:    page,
		PerPage: perPage,
		Fields:  "id,name",
	})
	if err != nil {
		return nil, fmt.Errorf("encode items query: %w", err)
	}

	requestUrl, err := url.JoinPath(baseUrl, itemsPath)
	if err != nil {
		return nil, fmt.Errorf("build items URL: %w", err)
	}

	parsedUrl, err := url.Parse(requestUrl)
	if err != nil {
		return nil, fmt.Errorf("parse items URL: %w", err)
	}
	parsedUrl.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request items: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("request items: unexpected HTTP status %s", res.Status)
	}

	var itemsResponse getItemsResponse
	if err = json.NewDecoder(res.Body).Decode(&itemsResponse); err != nil {
		return nil, fmt.Errorf("decode items response: %w", err)
	}

	return &itemsResponse, nil
}

func getAllItems(ctx context.Context, baseUrl string) ([]item, error) {
	res, err := getItems(ctx, baseUrl, 1, itemsPerPage)
	if err != nil {
		return nil, err
	}

	var items []item
	items = append(items, res.Items...)

	for page := 2; page <= res.TotalPages; page++ {
		res, err := getItems(ctx, baseUrl, page, itemsPerPage)
		if err != nil {
			return nil, err
		}

		items = append(items, res.Items...)
	}

	return items, nil
}
