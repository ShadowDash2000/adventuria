package npcs

import (
	"adventuria/internal/adventuria/action_events"
	"adventuria/internal/adventuria/board/actors"
	"adventuria/internal/adventuria/board/actors/player"
	"adventuria/internal/adventuria/effects"
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/event"
	"adventuria/pkg/locker"
	"context"
	"sync/atomic"
	"time"
)

type repository interface {
	UpdateShiftAndPositionByID(ctx context.Context, id string, nextShiftChangeAt time.Time, position actors.Position) error
	GetAllByFilter(ctx context.Context, filter Filter) ([]*NPCDetailed, error)
}

type behaviorsService interface {
	Move(ctx context.Context, npc *model.NPC, behavior *model.NPCBehavior) (actors.Position, error)
}

type effectsService interface {
	CreateFromEffectInfos(infos []*model.EffectInfo) ([]effects.Effect, error)
}

type actionEvents interface {
	CreateFromActionEventInfo(info *model.ActionEventInfo) (action_events.ActionEvent, error)
}

type players interface {
	GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.Player, error)
	Save(ctx context.Context, player *model.Player) error
}

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
}

type seasons interface {
	IsSeasonEnded(ctx context.Context, id string) (bool, error)
}

type NPCs struct {
	running       atomic.Bool
	playersLocker *locker.Locker[string]
	repository    repository
	behaviors     behaviorsService
	effects       effectsService
	actionEvents  actionEvents
	players       players
	settings      settings
	seasons       seasons

	onUpdateStart *event.Hook[*UpdateStartEvent]
	onUpdateEnd   *event.Hook[*UpdateEndEvent]
}

func NewNPCs(
	playersLocker *locker.Locker[string],
	repository repository,
	behaviors behaviorsService,
	effects effectsService,
	actionEvents actionEvents,
	players players,
	settings settings,
	seasons seasons,
) *NPCs {
	return &NPCs{
		playersLocker: playersLocker,
		repository:    repository,
		behaviors:     behaviors,
		effects:       effects,
		actionEvents:  actionEvents,
		players:       players,
		settings:      settings,
		seasons:       seasons,

		onUpdateStart: &event.Hook[*UpdateStartEvent]{},
		onUpdateEnd:   &event.Hook[*UpdateEndEvent]{},
	}
}

func (n *NPCs) IsRunning() bool {
	return n.running.Load()
}

func (n *NPCs) CheckNPCSchedules(ctx context.Context) (err error) {
	currentSeason, err := n.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	isSeasonEnded, err := n.seasons.IsSeasonEnded(ctx, currentSeason)
	if err != nil {
		return err
	}

	if isSeasonEnded {
		return
	}

	isNotRunning := n.running.CompareAndSwap(false, true)
	if !isNotRunning {
		return errs.ErrCellEventSchedulerAlreadyRunning
	}

	shouldTriggerUpdateEnd := false
	defer func() {
		n.running.Store(false)

		if shouldTriggerUpdateEnd {
			eventErr := n.OnUpdateEnd().Trigger(ctx, &UpdateEndEvent{})
			if eventErr != nil {
				//err = eventErr
			}
		}
	}()

	err = n.waitForPlayersUnlock(ctx)
	if err != nil {
		return err
	}

	npcs, err := n.GetAll(ctx)
	if err != nil {
		return err
	}

	var npcsToUpdate []*NPC
	for _, npc := range npcs {
		if npc.NPC().NextShiftChangeAt().Compare(time.Now()) <= 0 {
			npcsToUpdate = append(npcsToUpdate, npc)
		}
	}

	if len(npcsToUpdate) == 0 {
		return nil
	}

	err = n.OnUpdateStart().Trigger(ctx, &UpdateStartEvent{})
	if err != nil {
		return err
	}

	shouldTriggerUpdateEnd = true

	players, err := n.players.GetAllBySeasonID(ctx, currentSeason)
	if err != nil {
		return err
	}

	playersByCellId := make(map[string][]*model.Player)
	for _, player := range players {
		cellId := player.LastAction().Cell()
		playersByCellId[cellId] = append(playersByCellId[cellId], player)
	}

	err = n.disposeActionEvents(ctx, npcsToUpdate, playersByCellId)
	if err != nil {
		return err
	}

	for _, npc := range npcsToUpdate {
		position, err := n.behaviors.Move(ctx, npc.NPC(), npc.Behavior())
		if err != nil {
			return err
		}

		npc.NPC().SetNextShiftChangeAt(
			time.Now().Add(time.Duration(npc.NPC().ShiftInterval()) * time.Second),
		)

		err = n.repository.UpdateShiftAndPositionByID(ctx, npc.NPC().ID(), npc.NPC().NextShiftChangeAt(), position)
		if err != nil {
			return err
		}
	}

	err = n.initActionEvents(ctx, npcsToUpdate, playersByCellId)
	if err != nil {
		return err
	}

	return nil
}

func (n *NPCs) waitForPlayersUnlock(ctx context.Context) error {
	if !n.playersLocker.HasAny() {
		return nil
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if !n.playersLocker.HasAny() {
				return nil
			}
		}
	}
}

func (n *NPCs) ListenToCellEvents(ctx context.Context, actor *player.Actor) error {
	npcs, err := n.GetAll(ctx)
	if err != nil {
		return err
	}

	if len(npcs) == 0 {
		return nil
	}

	npcsByCellId := make(map[string][]*NPC)
	for _, npc := range npcs {
		currentCell := npc.NPC().CurrentCell()
		if currentCell != "" {
			npcsByCellId[currentCell] = append(npcsByCellId[currentCell], npc)
		}
	}

	if len(npcsByCellId) == 0 {
		return nil
	}

	err = n.listenToActionEvents(npcsByCellId, actor)
	if err != nil {
		return err
	}

	err = n.listenToCellEventEffects(ctx, npcsByCellId, actor)
	if err != nil {
		return err
	}

	return nil
}

func (n *NPCs) listenToActionEvents(npcsByCellId map[string][]*NPC, actor *player.Actor) error {
	actor.Events().OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
		npcs, ok := npcsByCellId[e.CurrentCell.ID()]
		if !ok {
			return e.Next()
		}

		for _, npc := range npcs {
			err := npc.ActionEvent().Init(ctx, actor.Player())
			if err != nil {
				return err
			}
		}

		return e.Next()
	})

	return nil
}

func (n *NPCs) listenToCellEventEffects(ctx context.Context, npcsByCellId map[string][]*NPC, actor *player.Actor) error {
	unsubKeys, err := n.subscribeCellEventEffects(ctx, npcsByCellId, actor)
	if err != nil {
		return err
	}

	actor.Events().OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
		if len(unsubKeys) > 0 {
			for _, unsubKey := range unsubKeys {
				actor.Events().Unsubscribe(unsubKey)
			}
			unsubKeys = nil
		}

		unsubKeys, err = n.subscribeCellEventEffects(ctx, npcsByCellId, actor)
		if err != nil {
			return err
		}

		return e.Next()
	})

	return nil
}

func (n *NPCs) subscribeCellEventEffects(ctx context.Context, npcsByCellId map[string][]*NPC, actor *player.Actor) ([]string, error) {
	npcs, ok := npcsByCellId[actor.Player().LastAction().Cell()]
	if !ok {
		return nil, nil
	}

	var npcEffects []effects.Effect
	var npcIds []string
	for _, npc := range npcs {
		npcEffects = append(npcEffects, npc.Effects()...)
		npcIds = append(npcIds, npc.npc.ID())
	}

	if len(npcEffects) == 0 {
		return nil, nil
	}

	unsubKeys := make([]string, len(npcEffects))
	for i, effect := range npcEffects {
		unsubKey := "npc_effect:" +
			actor.Player().ID() + ":" +
			npcIds[i] + ":" +
			effect.Data().ID() + ":" +
			actor.Player().LastAction().Cell()

		unsubs, err := effect.Subscribe(
			ctx,
			actor,
			effects.Context{
				Priority: 50,
			},
			func(ctx context.Context) {
				actor.Events().Unsubscribe(unsubKey)
			},
		)
		if err != nil {
			return nil, err
		}

		actor.Events().AddUnsubs(unsubKey, unsubs...)
		unsubKeys[i] = unsubKey
	}

	return unsubKeys, nil
}

func (n *NPCs) initActionEvents(ctx context.Context, npcs []*NPC, playersByCellId map[string][]*model.Player) error {
	for _, npc := range npcs {
		actionEvent := npc.ActionEvent()
		if actionEvent == nil {
			continue
		}

		players, ok := playersByCellId[npc.NPC().CurrentCell()]
		if !ok {
			continue
		}

		for _, player := range players {
			err := actionEvent.Init(ctx, player)
			if err != nil {
				return err
			}

			err = n.players.Save(ctx, player)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *NPCs) disposeActionEvents(ctx context.Context, npcs []*NPC, playersByCellId map[string][]*model.Player) error {
	for _, npc := range npcs {
		actionEvent := npc.ActionEvent()
		if actionEvent == nil {
			continue
		}

		players, ok := playersByCellId[npc.NPC().CurrentCell()]
		if !ok {
			continue
		}

		for _, player := range players {
			err := actionEvent.Dispose(ctx, player)
			if err != nil {
				return err
			}

			err = n.players.Save(ctx, player)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *NPCs) GetAll(ctx context.Context) ([]*NPC, error) {
	npcsDetailed, err := n.repository.GetAllByFilter(ctx, Filter{})
	if err != nil {
		return nil, err
	}

	if len(npcsDetailed) == 0 {
		return nil, nil
	}

	npcs := make([]*NPC, len(npcsDetailed))
	for i, npcDetailed := range npcsDetailed {
		var npcActionEvent action_events.ActionEvent
		if npcDetailed.ActionEvent != nil {
			npcActionEvent, err = n.actionEvents.CreateFromActionEventInfo(npcDetailed.ActionEvent)
			if err != nil {
				return nil, err
			}
		}

		npcEffects, err := n.effects.CreateFromEffectInfos(npcDetailed.Effects)
		if err != nil {
			return nil, err
		}

		npcs[i] = RestoreNPC(
			npcDetailed.NPC,
			npcDetailed.Behavior,
			npcActionEvent,
			npcEffects,
		)
	}

	return npcs, nil
}
