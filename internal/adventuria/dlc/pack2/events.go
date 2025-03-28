package pack2

import (
	"adventuria/internal/adventuria"
	"adventuria/pkg/helper"
)

func WithBaseEvents(g adventuria.Game) adventuria.Game {
	g.Event().On(adventuria.OnAfterItemRoll, OnAfterItemRollEffects)
	//g.Event().On(adventuria.OnAfterAction, ClearNextStepTypeItems)
	g.Event().On(adventuria.OnAfterRoll, OnAfterRollEffects)
	g.Event().On(adventuria.OnAfterAction, TeleportToRandomCellByIds)
	return g
}

func OnAfterItemRollEffects(user *adventuria.User, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnRollItem)
	if err != nil {
		return err
	}

	nextStepType := effects.Effect(EffectTypeChangeNextStepType).String()
	if nextStepType != "" {

	}

	return nil
}

func ClearNextStepTypeItems(user *adventuria.User, event string, gc *adventuria.GameComponents) error {
	var invItemsEffectsIds map[string][]string
	for invItemId, invItem := range user.Inventory.Items() {
		effects := invItem.GetEffectsByEvent(adventuria.OnBeforeNextStepType)
		for _, effect := range effects {
			if effect.Type() != EffectTypeChangeNextStepType {
				continue
			}

			nextStepType := effect.String()
			if nextStepType != user.LastAction.Type() {
				continue
			}

			invItemsEffectsIds[invItemId] = append(invItemsEffectsIds[invItemId], invItemId)
		}
	}

	if len(invItemsEffectsIds) > 0 {
		err := user.Inventory.ApplyEffects(invItemsEffectsIds)
		if err != nil {
			return err
		}
	}

	return nil
}

func OnAfterRollEffects(user *adventuria.User, rollResult *adventuria.RollResult, gc *adventuria.GameComponents) error {
	_, _, err := user.Inventory.GetEffects(adventuria.EffectUseOnRoll)
	if err != nil {
		return err
	}

	return nil
}

func TeleportToRandomCellByIds(user *adventuria.User, event string, gc *adventuria.GameComponents) error {
	effects, _, err := user.Inventory.GetEffects(event)
	if err != nil {
		return err
	}

	effect := effects.Effect(EffectTypeTeleportToRandomCellByIds)
	cellIds := effect.Slice()
	if len(cellIds) > 0 {
		cellId := helper.RandomItemFromSlice(cellIds)
		err = user.MoveToCellId(cellId)
		if err != nil {
			return err
		}
	}

	return nil
}
