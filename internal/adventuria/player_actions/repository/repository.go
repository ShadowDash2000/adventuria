package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"database/sql"
	"errors"
	"time"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

type Repository struct {
	pb core.App
}

func NewRepository(pb core.App) *Repository {
	return &Repository{pb: pb}
}

func (r *Repository) Create(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	collection, err := pb.FindCollectionByNameOrId(schema.CollectionPlayerActions)
	if err != nil {
		return nil, err
	}

	record := core.NewRecord(collection)
	PlayerActionToRecord(action, record)

	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerAction(record), nil
}

func (r *Repository) Update(ctx context.Context, action *model.PlayerActionInfo) (*model.PlayerActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	record, err := pb.FindRecordById(schema.CollectionPlayerActions, action.ID())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrPlayerActionNotFound
		}
		return nil, err
	}

	PlayerActionToRecord(action, record)
	err = pb.SaveWithContext(ctx, record)
	if err != nil {
		return nil, err
	}

	return RecordToPlayerAction(record), nil
}

func (r *Repository) UpdatePayloadByID(ctx context.Context, id string, payload []byte) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Update(
			schema.CollectionPlayerActions,
			dbx.Params{
				schema.PlayerActionsSchema.Payload: payload,
			},
			dbx.HashExp{
				schema.PlayerActionsSchema.Id: id,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrPlayerActionNotFound
	}

	return nil
}

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	res, err := pb.DB().
		Delete(
			schema.CollectionPlayerActions,
			dbx.HashExp{
				schema.PlayerActionsSchema.Id: id,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrPlayerActionNotFound
	}

	return nil
}

func (r *Repository) DeleteByNPC(ctx context.Context, playerId, actionId, npcId, actionEventId string) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	_, err := pb.DB().
		Delete(
			schema.CollectionPlayerActions,
			dbx.HashExp{
				schema.PlayerActionsSchema.Player:      playerId,
				schema.PlayerActionsSchema.Action:      actionId,
				schema.PlayerActionsSchema.NPC:         npcId,
				schema.PlayerActionsSchema.ActionEvent: actionEventId,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteOlderThanExceptActionIDs(ctx context.Context, createdBefore time.Time, actionIds []string) error {
	if len(actionIds) == 0 {
		return nil
	}

	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	createdBeforeDate, _ := types.ParseDateTime(createdBefore)

	_, err := pb.DB().
		Delete(
			schema.CollectionPlayerActions,
			dbx.And(
				dbx.Not(dbx.In(
					schema.PlayerActionsSchema.Action,
					pbhelper.SliceToAny(actionIds)...,
				)),
				dbx.NewExp(
					"created < {:createdBefore}",
					dbx.Params{
						"createdBefore": createdBeforeDate,
					},
				),
			),
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllByPlayerAndActionID(ctx context.Context, playerId, actionId string) ([]*model.PlayerActionInfo, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var records []*core.Record
	err := pb.RecordQuery(schema.CollectionPlayerActions).
		WithContext(ctx).
		Where(dbx.HashExp{
			schema.PlayerActionsSchema.Player: playerId,
			schema.PlayerActionsSchema.Action: actionId,
		}).
		OrderBy("created ASC, id ASC").
		All(&records)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	return RecordsToPlayerActions(records), nil
}
