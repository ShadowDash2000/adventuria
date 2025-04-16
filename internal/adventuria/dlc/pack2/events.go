package pack2

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

func WithBaseEvents(g adventuria.Game) adventuria.Game {
	adventuria.GameEvent.On(adventuria.OnBeforeNextStepType, ChangeNextStepType)
	adventuria.GameEvent.On(adventuria.OnAfterAction, ClearNextStepTypeEffects)
	adventuria.GameEvent.On(adventuria.OnAfterAction, TeleportToRandomCellByIds)
	adventuria.GameEvent.On(adventuria.OnBeforeNextStepType, ChangeCellType)
	adventuria.GameEvent.On(adventuria.OnAfterMove, ClearChangeCellTypeEffects)
	return g
}

func ChangeNextStepType(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnRollItem)
	fields := e.Fields().(*adventuria.OnBeforeNextStepFields)

	fields.NextStepType = effects.Effect(EffectTypeChangeNextStepType).String()

	return nil
}

func ClearNextStepTypeEffects(e adventuria.EventFields) error {
	var invItemsEffectsIds map[string][]string
	for invItemId, invItem := range e.User().Inventory.Items() {
		effects := invItem.EffectsByEvent(adventuria.OnBeforeNextStepType)
		for _, effect := range effects {
			if effect.Type() != EffectTypeChangeNextStepType {
				continue
			}

			nextStepType := effect.String()
			if nextStepType != e.User().LastAction.Type() {
				continue
			}

			invItemsEffectsIds[invItemId] = append(invItemsEffectsIds[invItemId], effect.ID())
		}
	}

	if len(invItemsEffectsIds) > 0 {
		err := e.User().Inventory.ApplyEffects(invItemsEffectsIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func TeleportToRandomCellByIds(e adventuria.EventFields) error {
	fields := e.Fields().(*adventuria.OnAfterActionFields)
	effects := e.Effects(fields.Event)

	effect := effects.Effect(EffectTypeTeleportToRandomCellByIds)
	cellIds := effect.Slice()
	if len(cellIds) > 0 {
		cellId := helper.RandomItemFromSlice(cellIds)
		err := e.User().MoveToCellId(cellId)
		if err != nil {
			return err
		}
	}

	return nil
}

func ChangeCellType(e adventuria.EventFields) error {
	effects := e.Effects(adventuria.EffectUseOnRollItem)
	fields := e.Fields().(*adventuria.OnBeforeNextStepFields)

	cellType := adventuria.CellType(effects.Effect(EffectChangeCellType).String())
	if cellType != "" {
		fields.CurrentCell.SetType(cellType)
	}

	return nil
}

func ClearChangeCellTypeEffects(e adventuria.EventFields) error {
	return e.User().Inventory.ApplyEffectsByTypes([]string{EffectChangeCellType})
}
