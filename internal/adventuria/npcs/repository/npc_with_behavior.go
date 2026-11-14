package repository

import (
	"adventuria/internal/adventuria/model"
	"adventuria/internal/adventuria/npcs"

	"ganigeorgiev.com/pocketbase/pocketbase/tools/types"
)

type npcWithBehaviorRow struct {
	NPCId                string                  `db:"npc_id"`
	NPCActionEvent       string                  `db:"npc_action_event"`
	NPCEffects           types.JSONArray[string] `db:"npc_effects"`
	NPCCurrentWorld      string                  `db:"npc_current_world"`
	NPCCurrentCell       string                  `db:"npc_current_cell"`
	NPCCellsPassed       int                     `db:"npc_cells_passed"`
	NPCShiftInterval     int                     `db:"npc_shift_interval"`
	NPCNextShiftChangeAt types.DateTime          `db:"npc_next_shift_change_at"`

	BehaviorId      string                `db:"behavior_id"`
	BehaviorType    model.NPCBehaviorType `db:"behavior_type"`
	BehaviorPayload types.JSONRaw         `db:"behavior_payload"`
	BehaviorData    types.JSONRaw         `db:"behavior_data"`
}

func (n *npcWithBehaviorRow) toNPCWithBehavior(actionEvent *model.ActionEventInfo, effects []*model.EffectInfo) *npcs.NPCDetailed {
	if n == nil {
		return nil
	}

	return &npcs.NPCDetailed{
		NPC: model.RestoreNPC(model.NPCData{
			Id:                n.NPCId,
			Behavior:          n.BehaviorId,
			ActionEvent:       n.NPCActionEvent,
			Effects:           n.NPCEffects,
			CurrentWorld:      n.NPCCurrentWorld,
			CurrentCell:       n.NPCCurrentCell,
			CellsPassed:       n.NPCCellsPassed,
			ShiftInterval:     n.NPCShiftInterval,
			NextShiftChangeAt: n.NPCNextShiftChangeAt.Time(),
		}),
		Behavior: model.RestoreNPCBehavior(model.NPCBehaviorData{
			Id:      n.BehaviorId,
			Type:    n.BehaviorType,
			Payload: n.BehaviorPayload,
			Data:    n.BehaviorData,
		}),
		ActionEvent: actionEvent,
		Effects:     effects,
	}
}
