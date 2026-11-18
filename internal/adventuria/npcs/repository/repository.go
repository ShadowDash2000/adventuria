package repository

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/npcs"
	"adventuria/internal/adventuria/schema"
	"adventuria/pkg/pbhelper"
	"adventuria/pkg/pbtransaction"
	"context"
	"fmt"
	"time"

	"ganigeorgiev.com/pocketbase/dbx"
	"ganigeorgiev.com/pocketbase/pocketbase/core"
	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

type effectsRepo interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.EffectInfo, error)
}

type actionEventsRepo interface {
	GetByIDs(ctx context.Context, ids []string) ([]*model.ActionEventInfo, error)
}

type Repository struct {
	pb           core.App
	effects      effectsRepo
	actionEvents actionEventsRepo
}

func NewRepository(pb core.App, effectsRepo effectsRepo, actionEventsRepo actionEventsRepo) *Repository {
	return &Repository{
		pb:           pb,
		effects:      effectsRepo,
		actionEvents: actionEventsRepo,
	}
}

func (r *Repository) GetActiveNPCs(ctx context.Context) ([]*npcs.ActiveNPC, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	var activeNPCs []*activeNPCRow
	err := pb.DB().
		Select(
			schema.NPCsSchema.Id,
			schema.NPCsSchema.Name,
			schema.NPCsSchema.Avatar,
			schema.NPCsSchema.Color,
		).
		From(schema.CollectionNPCs).
		WithContext(ctx).
		All(&activeNPCs)
	if err != nil {
		return nil, err
	}

	if len(activeNPCs) == 0 {
		return nil, nil
	}

	res := make([]*npcs.ActiveNPC, len(activeNPCs))
	for i, activeNPC := range activeNPCs {
		res[i] = activeNPC.ToActiveNPC()
	}

	return res, nil
}

func (r *Repository) UpdateShiftAndPositionByID(ctx context.Context, id string, nextShiftChangeAt time.Time, position model.ActorPosition) error {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	nextShiftChangeDate, _ := types.ParseDateTime(nextShiftChangeAt)

	res, err := pb.DB().
		Update(
			schema.CollectionNPCs,
			dbx.Params{
				schema.NPCsSchema.NextShiftChangeAt: nextShiftChangeDate,
				schema.NPCsSchema.CurrentWorld:      position.WorldId,
				schema.NPCsSchema.CurrentCell:       position.CellId,
				schema.NPCsSchema.CellsPassed:       position.CellsPassed,
			},
			dbx.HashExp{
				schema.NPCsSchema.Id: id,
			},
		).
		WithContext(ctx).
		Execute()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errs.ErrNPCNotFound
	}

	return nil
}

func (r *Repository) GetAllDetailedByFilter(ctx context.Context, filter npcs.Filter) ([]*npcs.NPCDetailed, error) {
	pb := pbtransaction.GetCtxTransactionOrApp(ctx, r.pb)

	npcsTable := schema.CollectionNPCs
	behaviorsTable := schema.CollectionNPCBehaviors

	npcId := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.Id)
	npcBehavior := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.Behavior)
	npcActionEvent := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.ActionEvent)
	npcEffects := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.Effects)
	npcCurrentWorld := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.CurrentWorld)
	npcCurrentCell := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.CurrentCell)
	npcCellsPassed := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.CellsPassed)
	npcShiftInterval := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.ShiftInterval)
	npcNextShiftChangeAt := pbhelper.DotExpand(npcsTable, schema.NPCsSchema.NextShiftChangeAt)

	behaviorId := pbhelper.DotExpand(behaviorsTable, schema.NPCBehaviorsSchema.Id)
	behaviorType := pbhelper.DotExpand(behaviorsTable, schema.NPCBehaviorsSchema.Type)
	behaviorPayload := pbhelper.DotExpand(behaviorsTable, schema.NPCBehaviorsSchema.Payload)
	behaviorData := pbhelper.DotExpand(behaviorsTable, schema.NPCBehaviorsSchema.Data)

	var npcsWithBehavior []*npcWithBehaviorRow
	q := pb.DB().
		Select(
			fmt.Sprintf("%s as npc_id", npcId),
			fmt.Sprintf("%s as npc_action_event", npcActionEvent),
			fmt.Sprintf("%s as npc_effects", npcEffects),
			fmt.Sprintf("%s as npc_current_world", npcCurrentWorld),
			fmt.Sprintf("%s as npc_current_cell", npcCurrentCell),
			fmt.Sprintf("%s as npc_cells_passed", npcCellsPassed),
			fmt.Sprintf("%s as npc_shift_interval", npcShiftInterval),
			fmt.Sprintf("%s as npc_next_shift_change_at", npcNextShiftChangeAt),

			fmt.Sprintf("%s as behavior_id", behaviorId),
			fmt.Sprintf("%s as behavior_type", behaviorType),
			fmt.Sprintf("%s as behavior_payload", behaviorPayload),
			fmt.Sprintf("%s as behavior_data", behaviorData),
		).
		From(npcsTable).
		InnerJoin(
			behaviorsTable,
			dbx.NewExp(pbhelper.Eq(behaviorId, npcBehavior)),
		)

	buildQueryFromFilter(filter, q, filterColumns{
		CurrentCell: npcCurrentCell,
	})

	err := q.WithContext(ctx).All(&npcsWithBehavior)
	if err != nil {
		return nil, err
	}

	if len(npcsWithBehavior) == 0 {
		return nil, nil
	}

	actionEventsByNPC, err := r.getNPCsActionEvents(ctx, npcsWithBehavior)
	if err != nil {
		return nil, err
	}

	effectsByNPC, err := r.getNPCsEffects(ctx, npcsWithBehavior)
	if err != nil {
		return nil, err
	}

	res := make([]*npcs.NPCDetailed, len(npcsWithBehavior))
	for i, row := range npcsWithBehavior {
		res[i] = row.toNPCWithBehavior(
			actionEventsByNPC[row.NPCId],
			effectsByNPC[row.NPCId],
		)
	}

	return res, nil
}

func (r *Repository) getNPCsEffects(ctx context.Context, npcsWithBehavior []*npcWithBehaviorRow) (map[string][]*model.EffectInfo, error) {
	var ids []string
	seen := make(map[string]struct{})
	for _, row := range npcsWithBehavior {
		for _, id := range row.NPCEffects {
			if _, ok := seen[id]; ok {
				continue
			}

			seen[id] = struct{}{}
			ids = append(ids, id)
		}
	}

	if len(ids) == 0 {
		return make(map[string][]*model.EffectInfo, len(npcsWithBehavior)), nil
	}

	effects, err := r.effects.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	effectsMap := make(map[string]*model.EffectInfo, len(effects))
	for _, effect := range effects {
		effectsMap[effect.ID()] = effect
	}

	res := make(map[string][]*model.EffectInfo, len(npcsWithBehavior))
	for _, row := range npcsWithBehavior {
		res[row.NPCId] = make([]*model.EffectInfo, len(row.NPCEffects))
		for i, id := range row.NPCEffects {
			res[row.NPCId][i] = effectsMap[id].Clone()
		}
	}

	return res, nil
}

func (r *Repository) getNPCsActionEvents(ctx context.Context, npcsWithBehavior []*npcWithBehaviorRow) (map[string]*model.ActionEventInfo, error) {
	var ids []string
	seen := make(map[string]struct{})
	for _, row := range npcsWithBehavior {
		if _, ok := seen[row.NPCActionEvent]; ok {
			continue
		}

		seen[row.NPCActionEvent] = struct{}{}
		ids = append(ids, row.NPCActionEvent)
	}

	if len(ids) == 0 {
		return make(map[string]*model.ActionEventInfo, len(npcsWithBehavior)), nil
	}

	actionEvents, err := r.actionEvents.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	actionEventsMap := make(map[string]*model.ActionEventInfo, len(actionEvents))
	for _, actionEvent := range actionEvents {
		actionEventsMap[actionEvent.ID()] = actionEvent
	}

	res := make(map[string]*model.ActionEventInfo, len(npcsWithBehavior))
	for _, row := range npcsWithBehavior {
		res[row.NPCId] = actionEventsMap[row.NPCActionEvent].Clone()
	}

	return res, nil
}
