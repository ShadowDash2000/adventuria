package v1

import (
	actionsRepo "adventuria/internal/adventuria/actions/repository"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/reviews"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
	"github.com/spf13/cobra"
)

const (
	actionsPath    = "/api/collections/actions/records"
	actionsPerPage = 500
)

func migrateActionsCommand(pb core.App, activities activities, cells cells, reviews reviewsService) *cobra.Command {
	migrator := newActionsMigrator(pb, activities, cells, reviews)

	return &cobra.Command{
		Use:          "actions <base-url>",
		Example:      "migrate-data v1 actions http://127.0.0.1:8080",
		Short:        "Migrates actions from Adventuria v1 through Pocketbase HTTP API",
		SilenceUsage: true,
		Args:         cobra.ExactArgs(1),
		RunE: func(command *cobra.Command, args []string) error {
			return migrator.migrate(command.Context(), args[0])
		},
	}
}

type getActionsQuery struct {
	Page    int    `url:"page"`
	PerPage int    `url:"perPage"`
	Sort    string `url:"sort,omitempty"`
	Filter  string `url:"filter,omitempty"`
	Expand  string `url:"expand,omitempty"`
	Fields  string `url:"fields,omitempty"`
}

type action struct {
	CollectionID string         `json:"collectionId"`
	ID           string         `json:"id"`
	Created      types.DateTime `json:"created"`
	Updated      types.DateTime `json:"updated"`
	User         string         `json:"user"`
	Cell         string         `json:"cell"`
	Type         string         `json:"type"`
	Comment      string         `json:"comment"`
	Icon         string         `json:"icon"`
	Value        string         `json:"value"`
	Expand       actionExpand   `json:"expand"`
}

type actionExpand struct {
	User expandedActionUser `json:"user"`
	Cell expandedActionCell `json:"cell"`
}

type expandedActionUser struct {
	Name string `json:"name"`
}

type expandedActionCell struct {
	Name string `json:"name"`
}

type getActionsResponse struct {
	Items      []action `json:"items"`
	TotalPages int      `json:"totalPages"`
}

type actionsMigrator struct {
	pb         core.App
	activities activities
	cells      cells
	reviews    reviewsService
}

func newActionsMigrator(pb core.App, activities activities, cells cells, reviews reviewsService) *actionsMigrator {
	return &actionsMigrator{
		pb:         pb,
		activities: activities,
		cells:      cells,
		reviews:    reviews,
	}
}

func (a *actionsMigrator) migrate(ctx context.Context, baseUrl string) error {
	res, err := getActions(ctx, baseUrl, 1, actionsPerPage)
	if err != nil {
		return err
	}

	err = a.saveActions(ctx, baseUrl, res.Items)
	if err != nil {
		return err
	}

	for page := 2; page <= res.TotalPages; page++ {
		res, err := getActions(ctx, baseUrl, page, actionsPerPage)
		if err != nil {
			return err
		}

		err = a.saveActions(ctx, baseUrl, res.Items)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *actionsMigrator) saveActions(ctx context.Context, baseUrl string, actions []action) error {
	for _, oldAction := range actions {
		err := pbtransaction.RunInTransaction(ctx, a.pb, func(ctx context.Context, txApp core.App) error {
			return a.saveAction(ctx, baseUrl, oldAction)
		})
		if err != nil {
			return fmt.Errorf("migrate action %s: %w", oldAction.ID, err)
		}
	}
	return nil
}

func (a *actionsMigrator) saveAction(ctx context.Context, baseUrl string, oldAction action) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, a.pb)

	playerID, err := getPlayerIDByName(pb, oldAction.Expand.User.Name)
	if err != nil {
		if errors.Is(err, errs.ErrPlayerNotFound) {
			return nil
		}
		return err
	}

	currentCellName, ok := cellAliases[oldAction.Expand.Cell.Name]
	if !ok {
		return nil
	}
	cell, err := a.cells.GetByName(ctx, currentCellName, true)
	if err != nil {
		return fmt.Errorf("find cell %q: %w", currentCellName, err)
	}

	status, err := actionTypeToStatus(oldAction.Type)
	if err != nil {
		return err
	}
	newAction, err := model.NewAction(model.ActionCreate{Player: playerID, Cell: cell.ID(), Status: status})
	if err != nil {
		return err
	}

	activity, err := a.activities.GetByName(ctx, oldAction.Value)
	if err == nil {
		newAction.SetActivity(activity.ID())
	} else if !errors.Is(err, errs.ErrActivityNotFound) {
		return fmt.Errorf("find activity %q: %w", oldAction.Value, err)
	}

	comment, err := buildComment(ctx, baseUrl, oldAction)
	if err != nil {
		return err
	}

	review, err := a.reviews.Create(ctx, reviews.CreateInput{
		Comment: comment,
		Score:   0,
	})
	if err != nil {
		return err
	}

	newAction.SetReview(review.ID())

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionActions)
	if err != nil {
		return err
	}

	record := core.NewRecord(collection)
	err = actionsRepo.ActionToRecord(newAction, record)
	if err != nil {
		return err
	}

	record.SetRaw("created", oldAction.Created)
	record.SetRaw("updated", oldAction.Updated)

	return pb.Save(record)
}

func getActions(ctx context.Context, baseUrl string, page, perPage int) (*getActionsResponse, error) {
	q, err := query.Values(getActionsQuery{
		Page:    page,
		PerPage: perPage,
		Sort:    "created",
		Filter:  `type = "chooseResult" || type = "drop"`,
		Expand:  "user,cell",
		Fields:  "*,expand.user.name,expand.cell.name",
	})
	if err != nil {
		return nil, fmt.Errorf("encode actions query: %w", err)
	}

	return requestActions(ctx, baseUrl, q)
}

func requestActions(ctx context.Context, baseUrl string, q url.Values) (*getActionsResponse, error) {
	requestURL, err := url.JoinPath(baseUrl, actionsPath)
	if err != nil {
		return nil, fmt.Errorf("build actions URL: %w", err)
	}

	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, fmt.Errorf("parse actions URL: %w", err)
	}
	parsedURL.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
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
	err = json.NewDecoder(res.Body).Decode(&actionsResponse)
	if err != nil {
		return nil, fmt.Errorf("decode actions response: %w", err)
	}

	return &actionsResponse, nil
}

var actionTypeToActionStatus = map[string]model.ActionStatus{
	"chooseResult": model.ActionStatusDone,
	"drop":         model.ActionStatusDrop,
}

func actionTypeToStatus(actionType string) (model.ActionStatus, error) {
	status, ok := actionTypeToActionStatus[actionType]
	if !ok {
		return "", fmt.Errorf("unknown action type: %s", actionType)
	}

	return status, nil
}

func buildComment(ctx context.Context, baseUrl string, oldAction action) (string, error) {
	if oldAction.Icon == "" {
		return oldAction.Comment, nil
	}

	fileURL, err := url.JoinPath(baseUrl, "/api/files", oldAction.CollectionID, oldAction.ID, oldAction.Icon)
	if err != nil {
		return "", fmt.Errorf("build icon URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fileURL, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request icon: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("request icon: unexpected HTTP status %s", res.Status)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read icon: %w", err)
	}

	contentType, _, err := mime.ParseMediaType(res.Header.Get("Content-Type"))
	if err != nil {
		return "", fmt.Errorf("parse content type: %w", err)
	}

	image := fmt.Sprintf(
		`<img src="data:%s;base64,%s" width="%d" height="%d">`,
		contentType,
		base64.StdEncoding.EncodeToString(data),
		350,
		250,
	)
	if oldAction.Comment == "" {
		return image, nil
	}

	if len(oldAction.Comment)+len(image) > model.MaxCommentLength {
		return oldAction.Comment, nil
	}

	return "<p>" + oldAction.Comment + "</p><p>" + image + "</p>", nil
}

func getPlayerIDByName(pb core.App, name string) (string, error) {
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
