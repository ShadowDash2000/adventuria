package cell_events_schedules

import (
	"adventuria/internal/adventuria/errs"
	"adventuria/internal/adventuria/model"
	"adventuria/pkg/helper"
	"context"
	"errors"
	"maps"
	"slices"
	"time"
)

type repository interface {
	UpdateActiveCellByID(ctx context.Context, id, cellId string) error
	GetByActiveCellID(ctx context.Context, activeCellId string) (*model.CellEventSchedule, error)
	GetIDByActiveCellID(ctx context.Context, activeCellId string) (string, error)
	GetAll(ctx context.Context) ([]*model.CellEventSchedule, error)
}

type cells interface {
	GetAllByWorldID(ctx context.Context, worldId string) ([]*model.CellInfo, error)
}

type effects interface {
	GetByIDs(ctx context.Context, ids []string) ([]model.Effect, error)
}

type actionEvents interface {
	GetByIDWrapped(ctx context.Context, id string) (model.ActionEvent, error)
	ListenToActionEvents(events *model.Events, player *model.Player) error
}

type players interface {
	GetAllBySeasonID(ctx context.Context, seasonId string) ([]*model.Player, error)
	Save(ctx context.Context, player *model.Player) error
}

type settings interface {
	CurrentSeason(ctx context.Context) (string, error)
}

type CellEventsSchedules struct {
	repository   repository
	cells        cells
	effects      effects
	actionEvents actionEvents
	players      players
	settings     settings
}

func NewCellEventsSchedules(
	repository repository,
	cells cells,
	effects effects,
	actionEvents actionEvents,
	players players,
	settings settings,
) *CellEventsSchedules {
	return &CellEventsSchedules{
		repository:   repository,
		cells:        cells,
		effects:      effects,
		actionEvents: actionEvents,
		players:      players,
		settings:     settings,
	}
}

// CheckEventsSchedules TODO implement players actions block while updating
func (c *CellEventsSchedules) CheckEventsSchedules(ctx context.Context) error {
	events, err := c.repository.GetAll(ctx)
	if err != nil {
		return err
	}

	var eventsToUpdate []*model.CellEventSchedule
	for _, event := range events {
		if event.LastShiftChange().Add(time.Duration(event.ShiftInterval()) * time.Second).Before(time.Now()) {
			eventsToUpdate = append(eventsToUpdate, event)
		}
	}

	if len(eventsToUpdate) == 0 {
		return nil
	}

	err = c.pickCellsForEvents(ctx, eventsToUpdate)
	if err != nil {
		return err
	}

	for _, event := range eventsToUpdate {
		err = c.repository.UpdateActiveCellByID(ctx, event.ID(), event.ActiveCell())
		if err != nil {
			return err
		}
	}

	err = c.initActionEvents(ctx, eventsToUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (c *CellEventsSchedules) pickCellsForEvents(ctx context.Context, events []*model.CellEventSchedule) error {
	cellsByWorldId := make(map[string]map[model.CellType][]*model.CellInfo)
	for _, event := range events {
		for _, worldId := range event.Worlds() {
			cellsByWorldId[worldId] = nil
		}
	}

	cellTypesByWorldId := make(map[string][]model.CellType, len(cellsByWorldId))
	for worldId := range cellsByWorldId {
		cells, err := c.cells.GetAllByWorldID(ctx, worldId)
		if err != nil {
			return err
		}

		cellsByTypes := make(map[model.CellType][]*model.CellInfo)
		for _, cell := range cells {
			cellsByTypes[cell.Type()] = append(cellsByTypes[cell.Type()], cell)
		}

		for cellType := range maps.Keys(cellsByTypes) {
			cellTypesByWorldId[worldId] = append(cellTypesByWorldId[worldId], cellType)
		}

		cellsByWorldId[worldId] = cellsByTypes
	}

	for _, event := range events {
		worldId := helper.RandomItemFromSlice(event.Worlds())
		availableCellTypes := cellTypesByWorldId[worldId]
		if len(availableCellTypes) == 0 {
			event.SetActiveCell("")
			continue
		}

		cellsForEvent := cellsByWorldId[worldId]
		eventCellTypes := event.CellTypes()

		var eventCellType model.CellType
		if len(eventCellTypes) == 0 {
			eventCellType = helper.RandomItemFromSlice(availableCellTypes)
		} else {
			eventCellType = helper.RandomItemFromSlice(helper.SlicesIntersection(availableCellTypes, eventCellTypes))
		}

		availableCells := cellsForEvent[eventCellType]
		newActiveCell, index := helper.RandomItemFromSliceWithIndex(availableCells)

		event.SetActiveCell(newActiveCell.ID())

		lastIndex := len(availableCells) - 1
		availableCells[index] = availableCells[lastIndex]
		availableCells = availableCells[:lastIndex]

		if len(availableCells) == 0 {
			delete(cellsForEvent, eventCellType)
			slices.DeleteFunc(availableCellTypes, func(cellType model.CellType) bool {
				return cellType == eventCellType
			})
		}
	}

	return nil
}

func (c *CellEventsSchedules) ListenToCellEvents(ctx context.Context, events *model.Events, player *model.Player) error {
	err := c.actionEvents.ListenToActionEvents(events, player)
	if err != nil {
		return err
	}

	err = c.ListenToCellEventEffects(ctx, events, player)
	if err != nil {
		return err
	}

	return nil
}

func (c *CellEventsSchedules) ListenToCellEventEffects(ctx context.Context, events *model.Events, player *model.Player) error {
	unsubKeys, err := c.subscribeCellEventEffects(ctx, events, player)
	if err != nil {
		return err
	}

	events.OnAfterMove().BindFunc(func(ctx context.Context, e *model.OnAfterMoveEvent) error {
		if len(unsubKeys) > 0 {
			for _, unsubKey := range unsubKeys {
				events.Unsubscribe(unsubKey)
			}
			unsubKeys = nil
		}

		unsubKeys, err = c.subscribeCellEventEffects(ctx, events, player)
		if err != nil {
			return err
		}

		return e.Next()
	})

	return nil
}

func (c *CellEventsSchedules) subscribeCellEventEffects(ctx context.Context, events *model.Events, player *model.Player) ([]string, error) {
	cellEvent, err := c.repository.GetByActiveCellID(ctx, player.LastAction().Cell())
	if err != nil {
		if errors.Is(err, errs.ErrCellEventScheduleNotFound) {
			return nil, nil
		}
		return nil, err
	}

	effects, err := c.effects.GetByIDs(ctx, cellEvent.Effects())
	if err != nil {
		return nil, err
	}

	if len(effects) == 0 {
		return nil, nil
	}

	unsubKeys := make([]string, len(effects))
	for i, effect := range effects {
		unsubKey := "cell_event_effect:" + player.ID() + ":" + effect.Data().ID() + ":" + player.LastAction().Cell()
		unsubs, err := effect.Subscribe(
			ctx,
			events,
			player,
			model.EffectContext{
				Priority: 50,
			},
			func(ctx context.Context) {
				events.Unsubscribe(unsubKey)
			},
		)
		if err != nil {
			return nil, err
		}

		events.AddUnsubs(unsubKey, unsubs...)
		unsubKeys[i] = unsubKey
	}

	return unsubKeys, nil
}

func (c *CellEventsSchedules) initActionEvents(ctx context.Context, events []*model.CellEventSchedule) error {
	currentSeason, err := c.settings.CurrentSeason(ctx)
	if err != nil {
		return err
	}

	players, err := c.players.GetAllBySeasonID(ctx, currentSeason)
	if err != nil {
		return err
	}

	if len(players) == 0 {
		return nil
	}

	for _, event := range events {
		actionEvent, err := c.actionEvents.GetByIDWrapped(ctx, event.ActionEvent())
		if err != nil {
			return err
		}

		for _, player := range players {
			if event.ActiveCell() != player.LastAction().Cell() {
				continue
			}

			err = actionEvent.Init(ctx, player)
			if err != nil {
				return err
			}

			err = c.players.Save(ctx, player)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *CellEventsSchedules) GetIDByActiveCellID(ctx context.Context, activeCellId string) (string, error) {
	return c.repository.GetIDByActiveCellID(ctx, activeCellId)
}

func (c *CellEventsSchedules) GetByActiveCellID(ctx context.Context, activeCellId string) (*model.CellEventSchedule, error) {
	return c.repository.GetByActiveCellID(ctx, activeCellId)
}
